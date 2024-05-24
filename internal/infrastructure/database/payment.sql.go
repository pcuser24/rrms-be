// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: payment.sql

package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const checkPaymentAccessible = `-- name: CheckPaymentAccessible :one
SELECT EXISTS (SELECT 1 FROM "payments" WHERE "id" = $1 AND "user_id" = $2)
`

type CheckPaymentAccessibleParams struct {
	ID     int64     `json:"id"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) CheckPaymentAccessible(ctx context.Context, arg CheckPaymentAccessibleParams) (bool, error) {
	row := q.db.QueryRow(ctx, checkPaymentAccessible, arg.ID, arg.UserID)
	var exists bool
	err := row.Scan(&exists)
	return exists, err
}

const createPayment = `-- name: CreatePayment :one
INSERT INTO "payments" (
  "user_id",
  "order_id",
  "order_info",
  "amount"
) VALUES (
  $1,
  $2,
  $3,
  $4
) RETURNING id, user_id, order_id, order_info, amount, status, created_at, updated_at
`

type CreatePaymentParams struct {
	UserID    uuid.UUID `json:"user_id"`
	OrderID   string    `json:"order_id"`
	OrderInfo string    `json:"order_info"`
	Amount    int64     `json:"amount"`
}

func (q *Queries) CreatePayment(ctx context.Context, arg CreatePaymentParams) (Payment, error) {
	row := q.db.QueryRow(ctx, createPayment,
		arg.UserID,
		arg.OrderID,
		arg.OrderInfo,
		arg.Amount,
	)
	var i Payment
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.OrderID,
		&i.OrderInfo,
		&i.Amount,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createPaymentItem = `-- name: CreatePaymentItem :one
INSERT INTO "payment_items" (
  "payment_id",
  "name",
  "price",
  "quantity",
  "discount"
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5
) RETURNING payment_id, name, price, quantity, discount
`

type CreatePaymentItemParams struct {
	PaymentID int64  `json:"payment_id"`
	Name      string `json:"name"`
	Price     int64  `json:"price"`
	Quantity  int32  `json:"quantity"`
	Discount  int32  `json:"discount"`
}

func (q *Queries) CreatePaymentItem(ctx context.Context, arg CreatePaymentItemParams) (PaymentItem, error) {
	row := q.db.QueryRow(ctx, createPaymentItem,
		arg.PaymentID,
		arg.Name,
		arg.Price,
		arg.Quantity,
		arg.Discount,
	)
	var i PaymentItem
	err := row.Scan(
		&i.PaymentID,
		&i.Name,
		&i.Price,
		&i.Quantity,
		&i.Discount,
	)
	return i, err
}

const deletePayment = `-- name: DeletePayment :exec
DELETE FROM "payments" WHERE "id" = $1
`

func (q *Queries) DeletePayment(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deletePayment, id)
	return err
}

const getPaymentById = `-- name: GetPaymentById :one
SELECT id, user_id, order_id, order_info, amount, status, created_at, updated_at FROM "payments" WHERE "id" = $1
`

func (q *Queries) GetPaymentById(ctx context.Context, id int64) (Payment, error) {
	row := q.db.QueryRow(ctx, getPaymentById, id)
	var i Payment
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.OrderID,
		&i.OrderInfo,
		&i.Amount,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPaymentItemsByPaymentId = `-- name: GetPaymentItemsByPaymentId :many
SELECT payment_id, name, price, quantity, discount FROM "payment_items" WHERE "payment_id" = $1
`

func (q *Queries) GetPaymentItemsByPaymentId(ctx context.Context, paymentID int64) ([]PaymentItem, error) {
	rows, err := q.db.Query(ctx, getPaymentItemsByPaymentId, paymentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PaymentItem
	for rows.Next() {
		var i PaymentItem
		if err := rows.Scan(
			&i.PaymentID,
			&i.Name,
			&i.Price,
			&i.Quantity,
			&i.Discount,
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

const getPaymentsOfUser = `-- name: GetPaymentsOfUser :many
SELECT id, user_id, order_id, order_info, amount, status, created_at, updated_at 
FROM "payments" 
WHERE "user_id" = $3
ORDER BY "created_at" DESC
LIMIT $1 OFFSET $2
`

type GetPaymentsOfUserParams struct {
	Limit  int32     `json:"limit"`
	Offset int32     `json:"offset"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) GetPaymentsOfUser(ctx context.Context, arg GetPaymentsOfUserParams) ([]Payment, error) {
	rows, err := q.db.Query(ctx, getPaymentsOfUser, arg.Limit, arg.Offset, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Payment
	for rows.Next() {
		var i Payment
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.OrderID,
			&i.OrderInfo,
			&i.Amount,
			&i.Status,
			&i.CreatedAt,
			&i.UpdatedAt,
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

const updatePayment = `-- name: UpdatePayment :exec
UPDATE "payments" SET 
  order_id = coalesce($2, order_id),
  order_info = coalesce($3, order_info),
  amount = coalesce($4, amount),
  status = coalesce($5, status),
  updated_at = NOW()
WHERE "id" = $1
`

type UpdatePaymentParams struct {
	ID        int64             `json:"id"`
	OrderID   pgtype.Text       `json:"order_id"`
	OrderInfo pgtype.Text       `json:"order_info"`
	Amount    pgtype.Int8       `json:"amount"`
	Status    NullPAYMENTSTATUS `json:"status"`
}

func (q *Queries) UpdatePayment(ctx context.Context, arg UpdatePaymentParams) error {
	_, err := q.db.Exec(ctx, updatePayment,
		arg.ID,
		arg.OrderID,
		arg.OrderInfo,
		arg.Amount,
		arg.Status,
	)
	return err
}
