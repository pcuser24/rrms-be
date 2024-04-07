-- name: CreateUser :one
INSERT INTO "User" (
  email, 
  password, 
  created_at, 
  updated_at,
  first_name,
  last_name,
  role
) VALUES (
  sqlc.arg(email), 
  sqlc.arg(password), 
  NOW(), 
  NOW(),
  sqlc.arg(first_name),
  sqlc.arg(last_name),
  sqlc.arg(role)
) RETURNING *;

-- name: CreateSession :one
INSERT INTO "Session" ("id", "userId", "sessionToken", "expires", "user_agent", "client_ip", "created_at")
VALUES (sqlc.arg(id), sqlc.arg(userId), sqlc.arg(sessionToken), sqlc.arg(expires), sqlc.arg(user_agent), sqlc.arg(client_ip), sqlc.arg(created_at))
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM "User" WHERE email = $1 LIMIT 1;

-- name: GetUserById :one
SELECT * FROM "User" WHERE id = $1 LIMIT 1;

-- name: GetSessionById :one
SELECT * FROM "Session" WHERE id = $1 LIMIT 1;

-- name: UpdateUser :exec
UPDATE "User" SET 
  email = coalesce(sqlc.narg(email), email), 
  password = coalesce(sqlc.narg(password), password), 
  first_name = coalesce(sqlc.narg(first_name), first_name),
  last_name = coalesce(sqlc.narg(last_name), last_name),
  phone = coalesce(sqlc.narg(phone), phone),
  avatar = coalesce(sqlc.narg(avatar), avatar),
  address = coalesce(sqlc.narg(address), address),
  city = coalesce(sqlc.narg(city), city),
  district = coalesce(sqlc.narg(district), district),
  ward = coalesce(sqlc.narg(ward), ward),
  role = coalesce(sqlc.narg(role), role),
  updated_at = NOW(),
  updated_by = $1
WHERE id = $2;

-- name: UpdateSessionBlockingStatus :exec
UPDATE "Session" SET is_blocked = $1 WHERE id = $2;
