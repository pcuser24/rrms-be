-- name: InsertUser :one
INSERT INTO "User" (email, password, created_at, updated_at)
VALUES (sqlc.arg(email), sqlc.arg(password), NOW(), NOW())
RETURNING *;

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

-- name: UpdateSessionBlockingStatus :exec
UPDATE "Session" SET is_blocked = $1 WHERE id = $2;
