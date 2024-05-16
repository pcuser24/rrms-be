-- name: CreatePayment :one
INSERT INTO "payments" (
  "user_id",
  "order_id",
  "order_info",
  "amount"
) VALUES (
  sqlc.arg(user_id),
  sqlc.arg(order_id),
  sqlc.arg(order_info),
  sqlc.arg(amount)
) RETURNING *;

-- name: CreatePaymentItem :one
INSERT INTO "payment_items" (
  "payment_id",
  "name",
  "price",
  "quantity",
  "discount"
) VALUES (
  sqlc.arg(payment_id),
  sqlc.arg(name),
  sqlc.arg(price),
  sqlc.arg(quantity),
  sqlc.arg(discount)
) RETURNING *;

-- name: GetPaymentById :one
SELECT * FROM "payments" WHERE "id" = $1;

-- name: GetPaymentItemsByPaymentId :many
SELECT * FROM "payment_items" WHERE "payment_id" = $1;

-- name: CheckPaymentAccessible :one
SELECT EXISTS (SELECT 1 FROM "payments" WHERE "id" = $1 AND "user_id" = $2);

-- name: UpdatePayment :exec
UPDATE "payments" SET 
  order_id = coalesce(sqlc.narg(order_id), order_id),
  order_info = coalesce(sqlc.narg(order_info), order_info),
  amount = coalesce(sqlc.narg(amount), amount),
  status = coalesce(sqlc.narg(status), status),
  updated_at = NOW()
WHERE "id" = $1;

-- name: DeletePayment :exec
DELETE FROM "payments" WHERE "id" = $1;
