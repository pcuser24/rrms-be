-- name: InsertUser :one
INSERT INTO "User" (email, password, created_at, updated_at)
VALUES (sqlc.arg(email), sqlc.arg(password), NOW(), NOW())
RETURNING *;

-- name: GetUserByEmail :one
SELECT * FROM "User" WHERE email = $1 LIMIT 1;

-- name: GetUserById :one
SELECT * FROM "User" WHERE id = $1 LIMIT 1;
