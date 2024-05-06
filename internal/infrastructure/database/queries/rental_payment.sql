-- name: CreateRentalPayment :one
INSERT INTO "rental_payments" (
  "code",
  "rental_id",
  "payment_date",
  "updated_by",
  "status",
  "amount",
  "discount",
  "penalty",
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
  sqlc.narg(penalty),
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
  expiry_date = coalesce(sqlc.narg(expiry_date), expiry_date),
  payment_date = coalesce(sqlc.narg(payment_date), payment_date),
  discount = coalesce(sqlc.narg(discount), discount),
  penalty = coalesce(sqlc.narg(penalty), penalty),
  updated_by = sqlc.arg(user_id),
  updated_at = NOW()
WHERE "id" = $1;
