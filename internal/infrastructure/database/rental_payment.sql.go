// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: rental_payment.sql

package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"
)

const createRentalPayment = `-- name: CreateRentalPayment :one
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
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7,
  $8,
  $9,
  $10
) RETURNING id, code, rental_id, created_at, updated_at, start_date, end_date, expiry_date, payment_date, updated_by, status, amount, discount, paid, payamount, fine, note
`

type CreateRentalPaymentParams struct {
	Code        string                  `json:"code"`
	RentalID    int64                   `json:"rental_id"`
	PaymentDate pgtype.Date             `json:"payment_date"`
	UserID      pgtype.UUID             `json:"user_id"`
	Status      NullRENTALPAYMENTSTATUS `json:"status"`
	Amount      float32                 `json:"amount"`
	Discount    pgtype.Float4           `json:"discount"`
	Note        pgtype.Text             `json:"note"`
	StartDate   pgtype.Date             `json:"start_date"`
	EndDate     pgtype.Date             `json:"end_date"`
}

func (q *Queries) CreateRentalPayment(ctx context.Context, arg CreateRentalPaymentParams) (RentalPayment, error) {
	row := q.db.QueryRow(ctx, createRentalPayment,
		arg.Code,
		arg.RentalID,
		arg.PaymentDate,
		arg.UserID,
		arg.Status,
		arg.Amount,
		arg.Discount,
		arg.Note,
		arg.StartDate,
		arg.EndDate,
	)
	var i RentalPayment
	err := row.Scan(
		&i.ID,
		&i.Code,
		&i.RentalID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.StartDate,
		&i.EndDate,
		&i.ExpiryDate,
		&i.PaymentDate,
		&i.UpdatedBy,
		&i.Status,
		&i.Amount,
		&i.Discount,
		&i.Paid,
		&i.Payamount,
		&i.Fine,
		&i.Note,
	)
	return i, err
}

const getPaymentsOfRental = `-- name: GetPaymentsOfRental :many
SELECT id, code, rental_id, created_at, updated_at, start_date, end_date, expiry_date, payment_date, updated_by, status, amount, discount, paid, payamount, fine, note FROM "rental_payments" WHERE "rental_id" = $1
`

func (q *Queries) GetPaymentsOfRental(ctx context.Context, rentalID int64) ([]RentalPayment, error) {
	rows, err := q.db.Query(ctx, getPaymentsOfRental, rentalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RentalPayment
	for rows.Next() {
		var i RentalPayment
		if err := rows.Scan(
			&i.ID,
			&i.Code,
			&i.RentalID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.StartDate,
			&i.EndDate,
			&i.ExpiryDate,
			&i.PaymentDate,
			&i.UpdatedBy,
			&i.Status,
			&i.Amount,
			&i.Discount,
			&i.Paid,
			&i.Payamount,
			&i.Fine,
			&i.Note,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRentalPayment = `-- name: GetRentalPayment :one
SELECT id, code, rental_id, created_at, updated_at, start_date, end_date, expiry_date, payment_date, updated_by, status, amount, discount, paid, payamount, fine, note FROM "rental_payments" WHERE "id" = $1 LIMIT 1
`

func (q *Queries) GetRentalPayment(ctx context.Context, id int64) (RentalPayment, error) {
	row := q.db.QueryRow(ctx, getRentalPayment, id)
	var i RentalPayment
	err := row.Scan(
		&i.ID,
		&i.Code,
		&i.RentalID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.StartDate,
		&i.EndDate,
		&i.ExpiryDate,
		&i.PaymentDate,
		&i.UpdatedBy,
		&i.Status,
		&i.Amount,
		&i.Discount,
		&i.Paid,
		&i.Payamount,
		&i.Fine,
		&i.Note,
	)
	return i, err
}

const planRentalPayment = `-- name: PlanRentalPayment :many
SELECT plan_rental_payment($1)
`

func (q *Queries) PlanRentalPayment(ctx context.Context, rentalID int64) ([]int64, error) {
	rows, err := q.db.Query(ctx, planRentalPayment, rentalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var plan_rental_payment int64
		if err := rows.Scan(&plan_rental_payment); err != nil {
			return nil, err
		}
		items = append(items, plan_rental_payment)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const planRentalPayments = `-- name: PlanRentalPayments :many
SELECT plan_rental_payments()
`

func (q *Queries) PlanRentalPayments(ctx context.Context) ([]int64, error) {
	rows, err := q.db.Query(ctx, planRentalPayments)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var plan_rental_payments int64
		if err := rows.Scan(&plan_rental_payments); err != nil {
			return nil, err
		}
		items = append(items, plan_rental_payments)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateFinePayments = `-- name: UpdateFinePayments :exec
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
WHERE rp.id = up.id
`

func (q *Queries) UpdateFinePayments(ctx context.Context) error {
	_, err := q.db.Exec(ctx, updateFinePayments)
	return err
}

const updateFinePaymentsOfRental = `-- name: UpdateFinePaymentsOfRental :exec
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
        r.id = $1  AND
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
WHERE rp.id = up.id
`

func (q *Queries) UpdateFinePaymentsOfRental(ctx context.Context, rentalID int64) error {
	_, err := q.db.Exec(ctx, updateFinePaymentsOfRental, rentalID)
	return err
}

const updateRentalPayment = `-- name: UpdateRentalPayment :exec
UPDATE "rental_payments" SET
  status = coalesce($2, status),
  note = coalesce($3, note),
  amount = coalesce($4, amount),
  paid = coalesce($5, paid),
  payamount = coalesce($6, payamount),
  fine = coalesce($7, fine),
  expiry_date = coalesce($8, expiry_date),
  payment_date = coalesce($9, payment_date),
  discount = coalesce($10, discount),
  updated_by = $11,
  updated_at = NOW()
WHERE "id" = $1
`

type UpdateRentalPaymentParams struct {
	ID          int64                   `json:"id"`
	Status      NullRENTALPAYMENTSTATUS `json:"status"`
	Note        pgtype.Text             `json:"note"`
	Amount      pgtype.Float4           `json:"amount"`
	Paid        pgtype.Float4           `json:"paid"`
	Payamount   pgtype.Float4           `json:"payamount"`
	Fine        pgtype.Float4           `json:"fine"`
	ExpiryDate  pgtype.Date             `json:"expiry_date"`
	PaymentDate pgtype.Date             `json:"payment_date"`
	Discount    pgtype.Float4           `json:"discount"`
	UserID      pgtype.UUID             `json:"user_id"`
}

func (q *Queries) UpdateRentalPayment(ctx context.Context, arg UpdateRentalPaymentParams) error {
	_, err := q.db.Exec(ctx, updateRentalPayment,
		arg.ID,
		arg.Status,
		arg.Note,
		arg.Amount,
		arg.Paid,
		arg.Payamount,
		arg.Fine,
		arg.ExpiryDate,
		arg.PaymentDate,
		arg.Discount,
		arg.UserID,
	)
	return err
}
