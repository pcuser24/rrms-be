-- name: CreateRentalPayment :one
INSERT INTO "rental_payments" (
  "code",
  "rental_id",
  "payment_date",
  "updated_by",
  "status",
  "amount",
  "discount",
  "note",
  "start_date",
  "end_date"
) VALUES (
  sqlc.arg(code),
  sqlc.arg(rental_id),
  sqlc.narg(payment_date),
  sqlc.arg(user_id),
  sqlc.narg(status),
  sqlc.arg(amount),
  sqlc.narg(discount),
  sqlc.narg(note),
  sqlc.narg(start_date),
  sqlc.narg(end_date)
) RETURNING *;

-- name: GetRentalPayment :one
SELECT * FROM "rental_payments" WHERE "id" = $1 LIMIT 1;

-- name: GetPaymentsOfRental :many
SELECT * FROM "rental_payments" WHERE "rental_id" = $1;

-- name: PlanRentalPayments :many
SELECT plan_rental_payments();

-- name: PlanRentalPayment :many
SELECT plan_rental_payment($1);

-- name: UpdateRentalPayment :exec
UPDATE "rental_payments" SET
  status = coalesce(sqlc.narg(status), status),
  note = coalesce(sqlc.narg(note), note),
  amount = coalesce(sqlc.narg(amount), amount),
  paid = coalesce(sqlc.narg(paid), paid),
  payamount = coalesce(sqlc.narg(payamount), payamount),
  fine = coalesce(sqlc.narg(fine), fine),
  expiry_date = coalesce(sqlc.narg(expiry_date), expiry_date),
  payment_date = coalesce(sqlc.narg(payment_date), payment_date),
  discount = coalesce(sqlc.narg(discount), discount),
  updated_by = sqlc.arg(user_id),
  updated_at = NOW()
WHERE "id" = $1;

-- name: UpdateFinePayments :exec
WITH updated_payments AS (
    SELECT 
        rp.id,
        rp.amount,
        rp.discount,
        rp.paid,
        r.grace_period,
        r.late_payment_penalty_scheme,
        r.late_payment_penalty_amount,
        CASE 
            WHEN r.late_payment_penalty_scheme = 'FIXED' THEN (rp.amount - coalesce(rp.discount, 0) - rp.paid) + r.late_payment_penalty_amount
            WHEN r.late_payment_penalty_scheme = 'PERCENT' THEN (rp.amount - coalesce(rp.discount, 0) - rp.paid) * (1 + r.late_payment_penalty_amount / 100)
            WHEN r.late_payment_penalty_scheme = 'NONE' THEN (rp.amount - coalesce(rp.discount, 0) - rp.paid)
        END AS calculated_fine
    FROM rental_payments rp
    INNER JOIN rentals r ON rp.rental_id = r.id
    WHERE 
    	  rp.code LIKE '%_RENTAL_%' AND
        rp.status IN ('PENDING', 'REQUEST2PAY', 'PARTIALLYPAID') AND
        (rp.amount - coalesce(rp.discount, 0) - rp.paid) > 0 AND
        (rp.expiry_date + r.grace_period * INTERVAL '1 day') < CURRENT_DATE
)
UPDATE rental_payments rp
SET 
    status = 'PAYFINE',
    fine = up.calculated_fine,
    updated_at = NOW()
FROM updated_payments up
WHERE rp.id = up.id;

-- name: UpdateFinePaymentsOfRental :exec
WITH updated_payments AS (
    SELECT 
        rp.id,
        rp.amount,
        rp.discount,
        rp.paid,
        r.grace_period,
        r.late_payment_penalty_scheme,
        r.late_payment_penalty_amount,
        CASE 
            WHEN r.late_payment_penalty_scheme = 'FIXED' THEN (rp.amount - coalesce(rp.discount, 0) - rp.paid) + r.late_payment_penalty_amount
            WHEN r.late_payment_penalty_scheme = 'PERCENT' THEN (rp.amount - coalesce(rp.discount, 0) - rp.paid) * (1 + r.late_payment_penalty_amount / 100)
            WHEN r.late_payment_penalty_scheme = 'NONE' THEN (rp.amount - coalesce(rp.discount, 0) - rp.paid)
        END AS calculated_fine
    FROM rental_payments rp
    INNER JOIN rentals r ON rp.rental_id = r.id
    WHERE 
    	  rp.code LIKE '%_RENTAL_%' AND
        r.id = sqlc.arg(rental_id)  AND
        rp.status IN ('PENDING', 'REQUEST2PAY', 'PARTIALLYPAID') AND
        (rp.amount - coalesce(rp.discount, 0) - rp.paid) > 0 AND
        (rp.expiry_date + r.grace_period * INTERVAL '1 day') < CURRENT_DATE
)
UPDATE rental_payments rp
SET 
    status = 'PAYFINE',
    fine = up.calculated_fine,
    updated_at = NOW()
FROM updated_payments up
WHERE rp.id = up.id;
