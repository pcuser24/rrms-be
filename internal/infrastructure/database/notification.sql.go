// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: notification.sql

package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const createNotification = `-- name: CreateNotification :one
INSERT INTO notifications (
  user_id,
  title,
  content,
  data,
  target,
  channel
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
) RETURNING id, user_id, title, content, data, seen, target, channel, created_at, updated_at
`

type CreateNotificationParams struct {
	UserID  pgtype.UUID         `json:"user_id"`
	Title   string              `json:"title"`
	Content string              `json:"content"`
	Data    []byte              `json:"data"`
	Target  string              `json:"target"`
	Channel NOTIFICATIONCHANNEL `json:"channel"`
}

func (q *Queries) CreateNotification(ctx context.Context, arg CreateNotificationParams) (Notification, error) {
	row := q.db.QueryRow(ctx, createNotification,
		arg.UserID,
		arg.Title,
		arg.Content,
		arg.Data,
		arg.Target,
		arg.Channel,
	)
	var i Notification
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.Content,
		&i.Data,
		&i.Seen,
		&i.Target,
		&i.Channel,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createNotificationDevice = `-- name: CreateNotificationDevice :one
INSERT INTO user_notification_devices (
  user_id,
  session_id,
  token,
  platform
) VALUES (
  $1,
  $2,
  $3,
  $4
) ON CONFLICT ("user_id", "session_id") DO UPDATE 
SET "token" = EXCLUDED."token","last_accessed" = EXCLUDED."last_accessed"
RETURNING user_id, session_id, token, platform, last_accessed, created_at
`

type CreateNotificationDeviceParams struct {
	UserID    uuid.UUID `json:"user_id"`
	SessionID uuid.UUID `json:"session_id"`
	Token     string    `json:"token"`
	Platform  PLATFORM  `json:"platform"`
}

func (q *Queries) CreateNotificationDevice(ctx context.Context, arg CreateNotificationDeviceParams) (UserNotificationDevice, error) {
	row := q.db.QueryRow(ctx, createNotificationDevice,
		arg.UserID,
		arg.SessionID,
		arg.Token,
		arg.Platform,
	)
	var i UserNotificationDevice
	err := row.Scan(
		&i.UserID,
		&i.SessionID,
		&i.Token,
		&i.Platform,
		&i.LastAccessed,
		&i.CreatedAt,
	)
	return i, err
}

const deleteExpiredTokens = `-- name: DeleteExpiredTokens :exec
DELETE FROM user_notification_devices 
WHERE 
  "last_accessed" < NOW() - $1::INTEGER * INTERVAL '1 day'
`

func (q *Queries) DeleteExpiredTokens(ctx context.Context, interval int32) error {
	_, err := q.db.Exec(ctx, deleteExpiredTokens, interval)
	return err
}

const deleteNotificationDeviceToken = `-- name: DeleteNotificationDeviceToken :exec
DELETE FROM user_notification_devices 
WHERE 
  "user_id" = $1 AND 
  "session_id" = $2 AND 
  "token" = $3
`

type DeleteNotificationDeviceTokenParams struct {
	UserID       uuid.UUID `json:"user_id"`
	SessionID    uuid.UUID `json:"session_id"`
	CurrentToken string    `json:"current_token"`
}

func (q *Queries) DeleteNotificationDeviceToken(ctx context.Context, arg DeleteNotificationDeviceTokenParams) error {
	_, err := q.db.Exec(ctx, deleteNotificationDeviceToken, arg.UserID, arg.SessionID, arg.CurrentToken)
	return err
}

const getNotification = `-- name: GetNotification :one
SELECT id, user_id, title, content, data, seen, target, channel, created_at, updated_at
FROM notifications
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetNotification(ctx context.Context, id int64) (Notification, error) {
	row := q.db.QueryRow(ctx, getNotification, id)
	var i Notification
	err := row.Scan(
		&i.ID,
		&i.UserID,
		&i.Title,
		&i.Content,
		&i.Data,
		&i.Seen,
		&i.Target,
		&i.Channel,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getNotificationDevice = `-- name: GetNotificationDevice :one
SELECT user_id, session_id, token, platform, last_accessed, created_at 
FROM user_notification_devices 
WHERE 
  user_id = $1 
  AND session_id = $2
  AND platform = $3::"PLATFORM"
  AND CASE
    WHEN $4::TEXT <> '' THEN token = $4::TEXT
    ELSE TRUE
  END
LIMIT 1
`

type GetNotificationDeviceParams struct {
	UserID    uuid.UUID `json:"user_id"`
	SessionID uuid.UUID `json:"session_id"`
	Platform  PLATFORM  `json:"platform"`
	Token     string    `json:"token"`
}

func (q *Queries) GetNotificationDevice(ctx context.Context, arg GetNotificationDeviceParams) (UserNotificationDevice, error) {
	row := q.db.QueryRow(ctx, getNotificationDevice,
		arg.UserID,
		arg.SessionID,
		arg.Platform,
		arg.Token,
	)
	var i UserNotificationDevice
	err := row.Scan(
		&i.UserID,
		&i.SessionID,
		&i.Token,
		&i.Platform,
		&i.LastAccessed,
		&i.CreatedAt,
	)
	return i, err
}

const getNotificationsOfUser = `-- name: GetNotificationsOfUser :many
SELECT id, user_id, title, content, data, seen, target, channel, created_at, updated_at
FROM notifications
WHERE user_id = $3
ORDER BY created_at DESC
LIMIT $1 OFFSET $2
`

type GetNotificationsOfUserParams struct {
	Limit  int32       `json:"limit"`
	Offset int32       `json:"offset"`
	UserID pgtype.UUID `json:"user_id"`
}

func (q *Queries) GetNotificationsOfUser(ctx context.Context, arg GetNotificationsOfUserParams) ([]Notification, error) {
	rows, err := q.db.Query(ctx, getNotificationsOfUser, arg.Limit, arg.Offset, arg.UserID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Notification
	for rows.Next() {
		var i Notification
		if err := rows.Scan(
			&i.ID,
			&i.UserID,
			&i.Title,
			&i.Content,
			&i.Data,
			&i.Seen,
			&i.Target,
			&i.Channel,
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

const updateNotificationDeviceTokenTimestamp = `-- name: UpdateNotificationDeviceTokenTimestamp :exec
UPDATE "user_notification_devices"
SET "last_accessed" = NOW()
WHERE 
  "user_id" = $1 AND 
  "session_id" = $2
`

type UpdateNotificationDeviceTokenTimestampParams struct {
	UserID    uuid.UUID `json:"user_id"`
	SessionID uuid.UUID `json:"session_id"`
}

func (q *Queries) UpdateNotificationDeviceTokenTimestamp(ctx context.Context, arg UpdateNotificationDeviceTokenTimestampParams) error {
	_, err := q.db.Exec(ctx, updateNotificationDeviceTokenTimestamp, arg.UserID, arg.SessionID)
	return err
}

const updatedNotification = `-- name: UpdatedNotification :exec
UPDATE notifications
SET
  title = coalesce($1, title),
  content = coalesce($2, content),
  data = coalesce($3, data),
  seen = coalesce($4, seen),
  updated_at = NOW()
WHERE id = $5
`

type UpdatedNotificationParams struct {
	Title   pgtype.Text `json:"title"`
	Content pgtype.Text `json:"content"`
	Data    []byte      `json:"data"`
	Seen    pgtype.Bool `json:"seen"`
	ID      int64       `json:"id"`
}

func (q *Queries) UpdatedNotification(ctx context.Context, arg UpdatedNotificationParams) error {
	_, err := q.db.Exec(ctx, updatedNotification,
		arg.Title,
		arg.Content,
		arg.Data,
		arg.Seen,
		arg.ID,
	)
	return err
}