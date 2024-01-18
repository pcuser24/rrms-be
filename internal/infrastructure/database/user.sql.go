// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: user.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createSession = `-- name: CreateSession :one
INSERT INTO "Session" ("id", "userId", "sessionToken", "expires", "user_agent", "client_ip", "created_at")
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, "sessionToken", "userId", expires, user_agent, client_ip, is_blocked, created_at
`

type CreateSessionParams struct {
	ID           uuid.UUID   `json:"id"`
	Userid       uuid.UUID   `json:"userid"`
	Sessiontoken string      `json:"sessiontoken"`
	Expires      time.Time   `json:"expires"`
	UserAgent    pgtype.Text `json:"user_agent"`
	ClientIp     pgtype.Text `json:"client_ip"`
	CreatedAt    time.Time   `json:"created_at"`
}

func (q *Queries) CreateSession(ctx context.Context, arg CreateSessionParams) (Session, error) {
	row := q.db.QueryRow(ctx, createSession,
		arg.ID,
		arg.Userid,
		arg.Sessiontoken,
		arg.Expires,
		arg.UserAgent,
		arg.ClientIp,
		arg.CreatedAt,
	)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.SessionToken,
		&i.UserId,
		&i.Expires,
		&i.UserAgent,
		&i.ClientIp,
		&i.IsBlocked,
		&i.CreatedAt,
	)
	return i, err
}

const getSessionById = `-- name: GetSessionById :one
SELECT id, "sessionToken", "userId", expires, user_agent, client_ip, is_blocked, created_at FROM "Session" WHERE id = $1 LIMIT 1
`

func (q *Queries) GetSessionById(ctx context.Context, id uuid.UUID) (Session, error) {
	row := q.db.QueryRow(ctx, getSessionById, id)
	var i Session
	err := row.Scan(
		&i.ID,
		&i.SessionToken,
		&i.UserId,
		&i.Expires,
		&i.UserAgent,
		&i.ClientIp,
		&i.IsBlocked,
		&i.CreatedAt,
	)
	return i, err
}

const getUserByEmail = `-- name: GetUserByEmail :one
SELECT id, email, password, group_id, created_at, updated_at, created_by, updated_by, deleted_f FROM "User" WHERE email = $1 LIMIT 1
`

func (q *Queries) GetUserByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, getUserByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.GroupID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.DeletedF,
	)
	return i, err
}

const getUserById = `-- name: GetUserById :one
SELECT id, email, password, group_id, created_at, updated_at, created_by, updated_by, deleted_f FROM "User" WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUserById(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUserById, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.GroupID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.DeletedF,
	)
	return i, err
}

const insertUser = `-- name: InsertUser :one
INSERT INTO "User" (email, password, created_at, updated_at)
VALUES ($1, $2, NOW(), NOW())
RETURNING id, email, password, group_id, created_at, updated_at, created_by, updated_by, deleted_f
`

type InsertUserParams struct {
	Email    string      `json:"email"`
	Password pgtype.Text `json:"password"`
}

func (q *Queries) InsertUser(ctx context.Context, arg InsertUserParams) (User, error) {
	row := q.db.QueryRow(ctx, insertUser, arg.Email, arg.Password)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Password,
		&i.GroupID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.CreatedBy,
		&i.UpdatedBy,
		&i.DeletedF,
	)
	return i, err
}

const updateSessionBlockingStatus = `-- name: UpdateSessionBlockingStatus :exec
UPDATE "Session" SET is_blocked = $1 WHERE id = $2
`

type UpdateSessionBlockingStatusParams struct {
	IsBlocked bool      `json:"is_blocked"`
	ID        uuid.UUID `json:"id"`
}

func (q *Queries) UpdateSessionBlockingStatus(ctx context.Context, arg UpdateSessionBlockingStatusParams) error {
	_, err := q.db.Exec(ctx, updateSessionBlockingStatus, arg.IsBlocked, arg.ID)
	return err
}
