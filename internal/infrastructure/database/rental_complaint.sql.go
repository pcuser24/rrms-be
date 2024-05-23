// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: rental_complaint.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createRentalComplaint = `-- name: CreateRentalComplaint :one
INSERT INTO rental_complaints (
  rental_id, 
  creator_id, 
  title,
  content, 
  suggestion,
  media,
  occurred_at,
  type,
  updated_by
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7,
  $8,
  $2
) RETURNING id, rental_id, creator_id, title, content, suggestion, media, occurred_at, created_at, updated_at, updated_by, type, status
`

type CreateRentalComplaintParams struct {
	RentalID   int64               `json:"rental_id"`
	CreatorID  uuid.UUID           `json:"creator_id"`
	Title      string              `json:"title"`
	Content    string              `json:"content"`
	Suggestion pgtype.Text         `json:"suggestion"`
	Media      []string            `json:"media"`
	OccurredAt time.Time           `json:"occurred_at"`
	Type       RENTALCOMPLAINTTYPE `json:"type"`
}

func (q *Queries) CreateRentalComplaint(ctx context.Context, arg CreateRentalComplaintParams) (RentalComplaint, error) {
	row := q.db.QueryRow(ctx, createRentalComplaint,
		arg.RentalID,
		arg.CreatorID,
		arg.Title,
		arg.Content,
		arg.Suggestion,
		arg.Media,
		arg.OccurredAt,
		arg.Type,
	)
	var i RentalComplaint
	err := row.Scan(
		&i.ID,
		&i.RentalID,
		&i.CreatorID,
		&i.Title,
		&i.Content,
		&i.Suggestion,
		&i.Media,
		&i.OccurredAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UpdatedBy,
		&i.Type,
		&i.Status,
	)
	return i, err
}

const createRentalComplaintReply = `-- name: CreateRentalComplaintReply :one
INSERT INTO rental_complaint_replies (
  complaint_id, 
  replier_id, 
  reply,
  media
) VALUES (
  $1,
  $2,
  $3,
  $4
) RETURNING complaint_id, replier_id, reply, media, created_at
`

type CreateRentalComplaintReplyParams struct {
	ComplaintID int64     `json:"complaint_id"`
	ReplierID   uuid.UUID `json:"replier_id"`
	Reply       string    `json:"reply"`
	Media       []string  `json:"media"`
}

func (q *Queries) CreateRentalComplaintReply(ctx context.Context, arg CreateRentalComplaintReplyParams) (RentalComplaintReply, error) {
	row := q.db.QueryRow(ctx, createRentalComplaintReply,
		arg.ComplaintID,
		arg.ReplierID,
		arg.Reply,
		arg.Media,
	)
	var i RentalComplaintReply
	err := row.Scan(
		&i.ComplaintID,
		&i.ReplierID,
		&i.Reply,
		&i.Media,
		&i.CreatedAt,
	)
	return i, err
}

const getRentalComplaint = `-- name: GetRentalComplaint :one
SELECT id, rental_id, creator_id, title, content, suggestion, media, occurred_at, created_at, updated_at, updated_by, type, status FROM rental_complaints WHERE id = $1 LIMIT 1
`

func (q *Queries) GetRentalComplaint(ctx context.Context, id int64) (RentalComplaint, error) {
	row := q.db.QueryRow(ctx, getRentalComplaint, id)
	var i RentalComplaint
	err := row.Scan(
		&i.ID,
		&i.RentalID,
		&i.CreatorID,
		&i.Title,
		&i.Content,
		&i.Suggestion,
		&i.Media,
		&i.OccurredAt,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UpdatedBy,
		&i.Type,
		&i.Status,
	)
	return i, err
}

const getRentalComplaintReplies = `-- name: GetRentalComplaintReplies :many
SELECT complaint_id, replier_id, reply, media, created_at FROM rental_complaint_replies WHERE complaint_id = $1
`

func (q *Queries) GetRentalComplaintReplies(ctx context.Context, complaintID int64) ([]RentalComplaintReply, error) {
	rows, err := q.db.Query(ctx, getRentalComplaintReplies, complaintID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RentalComplaintReply
	for rows.Next() {
		var i RentalComplaintReply
		if err := rows.Scan(
			&i.ComplaintID,
			&i.ReplierID,
			&i.Reply,
			&i.Media,
			&i.CreatedAt,
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

const getRentalComplaintsByRentalId = `-- name: GetRentalComplaintsByRentalId :many
SELECT id, rental_id, creator_id, title, content, suggestion, media, occurred_at, created_at, updated_at, updated_by, type, status FROM rental_complaints WHERE rental_id = $1
`

func (q *Queries) GetRentalComplaintsByRentalId(ctx context.Context, rentalID int64) ([]RentalComplaint, error) {
	rows, err := q.db.Query(ctx, getRentalComplaintsByRentalId, rentalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RentalComplaint
	for rows.Next() {
		var i RentalComplaint
		if err := rows.Scan(
			&i.ID,
			&i.RentalID,
			&i.CreatorID,
			&i.Title,
			&i.Content,
			&i.Suggestion,
			&i.Media,
			&i.OccurredAt,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.UpdatedBy,
			&i.Type,
			&i.Status,
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

const updateRentalComplaint = `-- name: UpdateRentalComplaint :exec
UPDATE rental_complaints
SET
  title = coalesce($1, title),
  content = coalesce($2, content),
  suggestion = coalesce($3, suggestion),
  media = coalesce($4, media),
  occurred_at = coalesce($5, occurred_at),
  status = coalesce($6, status),
  updated_at = NOW(),
  updated_by = $7
WHERE id = $8
`

type UpdateRentalComplaintParams struct {
	Title      pgtype.Text               `json:"title"`
	Content    pgtype.Text               `json:"content"`
	Suggestion pgtype.Text               `json:"suggestion"`
	Media      []string                  `json:"media"`
	OccurredAt pgtype.Timestamptz        `json:"occurred_at"`
	Status     NullRENTALCOMPLAINTSTATUS `json:"status"`
	UserID     uuid.UUID                 `json:"user_id"`
	ID         int64                     `json:"id"`
}

func (q *Queries) UpdateRentalComplaint(ctx context.Context, arg UpdateRentalComplaintParams) error {
	_, err := q.db.Exec(ctx, updateRentalComplaint,
		arg.Title,
		arg.Content,
		arg.Suggestion,
		arg.Media,
		arg.OccurredAt,
		arg.Status,
		arg.UserID,
		arg.ID,
	)
	return err
}