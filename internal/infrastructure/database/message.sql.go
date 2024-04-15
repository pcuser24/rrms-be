// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: message.sql

package database

import (
	"context"

	"github.com/google/uuid"
)

const checkMsgGroupMembership = `-- name: CheckMsgGroupMembership :one
SELECT COUNT(*) > 0
FROM msg_group_members 
WHERE 
  group_id = $1 AND 
  user_id = $2
`

type CheckMsgGroupMembershipParams struct {
	GroupID int64     `json:"group_id"`
	UserID  uuid.UUID `json:"user_id"`
}

func (q *Queries) CheckMsgGroupMembership(ctx context.Context, arg CheckMsgGroupMembershipParams) (bool, error) {
	row := q.db.QueryRow(ctx, checkMsgGroupMembership, arg.GroupID, arg.UserID)
	var column_1 bool
	err := row.Scan(&column_1)
	return column_1, err
}

const createMessage = `-- name: CreateMessage :one
INSERT INTO messages (
  group_id,
  from_user,
  content
) VALUES (
  $1,
  $2,
  $3
) RETURNING id, group_id, from_user, content, status, type, created_at, updated_at
`

type CreateMessageParams struct {
	GroupID  int64     `json:"group_id"`
	FromUser uuid.UUID `json:"from_user"`
	Content  string    `json:"content"`
}

func (q *Queries) CreateMessage(ctx context.Context, arg CreateMessageParams) (Message, error) {
	row := q.db.QueryRow(ctx, createMessage, arg.GroupID, arg.FromUser, arg.Content)
	var i Message
	err := row.Scan(
		&i.ID,
		&i.GroupID,
		&i.FromUser,
		&i.Content,
		&i.Status,
		&i.Type,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createMsgGroup = `-- name: CreateMsgGroup :one
INSERT INTO msg_groups (
  name,
  created_by
) VALUES (
  $1,
  $2
) RETURNING group_id, name, created_at, created_by
`

type CreateMsgGroupParams struct {
	Name      string    `json:"name"`
	CreatedBy uuid.UUID `json:"created_by"`
}

func (q *Queries) CreateMsgGroup(ctx context.Context, arg CreateMsgGroupParams) (MsgGroup, error) {
	row := q.db.QueryRow(ctx, createMsgGroup, arg.Name, arg.CreatedBy)
	var i MsgGroup
	err := row.Scan(
		&i.GroupID,
		&i.Name,
		&i.CreatedAt,
		&i.CreatedBy,
	)
	return i, err
}

const createMsgGroupMember = `-- name: CreateMsgGroupMember :one
INSERT INTO msg_group_members (
  group_id,
  user_id
) VALUES (
  $1,
  $2
) RETURNING group_id, user_id
`

type CreateMsgGroupMemberParams struct {
	GroupID int64     `json:"group_id"`
	UserID  uuid.UUID `json:"user_id"`
}

func (q *Queries) CreateMsgGroupMember(ctx context.Context, arg CreateMsgGroupMemberParams) (MsgGroupMember, error) {
	row := q.db.QueryRow(ctx, createMsgGroupMember, arg.GroupID, arg.UserID)
	var i MsgGroupMember
	err := row.Scan(&i.GroupID, &i.UserID)
	return i, err
}

const deleteMsgGroup = `-- name: DeleteMsgGroup :exec
DELETE FROM msg_groups WHERE group_id = $1
`

func (q *Queries) DeleteMsgGroup(ctx context.Context, groupID int64) error {
	_, err := q.db.Exec(ctx, deleteMsgGroup, groupID)
	return err
}

const deleteMsgGroupMember = `-- name: DeleteMsgGroupMember :exec
DELETE FROM msg_group_members WHERE group_id = $1 AND user_id = $2
`

type DeleteMsgGroupMemberParams struct {
	GroupID int64     `json:"group_id"`
	UserID  uuid.UUID `json:"user_id"`
}

func (q *Queries) DeleteMsgGroupMember(ctx context.Context, arg DeleteMsgGroupMemberParams) error {
	_, err := q.db.Exec(ctx, deleteMsgGroupMember, arg.GroupID, arg.UserID)
	return err
}

const getMessagesOfGroup = `-- name: GetMessagesOfGroup :many
SELECT id, group_id, from_user, content, status, type, created_at, updated_at 
FROM messages 
WHERE 
  group_id = $3 
ORDER BY
  created_at DESC 
LIMIT $1 OFFSET $2
`

type GetMessagesOfGroupParams struct {
	Limit   int32 `json:"limit"`
	Offset  int32 `json:"offset"`
	GroupID int64 `json:"group_id"`
}

func (q *Queries) GetMessagesOfGroup(ctx context.Context, arg GetMessagesOfGroupParams) ([]Message, error) {
	rows, err := q.db.Query(ctx, getMessagesOfGroup, arg.Limit, arg.Offset, arg.GroupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Message
	for rows.Next() {
		var i Message
		if err := rows.Scan(
			&i.ID,
			&i.GroupID,
			&i.FromUser,
			&i.Content,
			&i.Status,
			&i.Type,
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

const getMsgGroup = `-- name: GetMsgGroup :one
SELECT group_id, name, created_at, created_by FROM msg_groups WHERE group_id = $1
`

func (q *Queries) GetMsgGroup(ctx context.Context, groupID int64) (MsgGroup, error) {
	row := q.db.QueryRow(ctx, getMsgGroup, groupID)
	var i MsgGroup
	err := row.Scan(
		&i.GroupID,
		&i.Name,
		&i.CreatedAt,
		&i.CreatedBy,
	)
	return i, err
}

const getMsgGroupByName = `-- name: GetMsgGroupByName :one
SELECT group_id, name, created_at, created_by 
FROM msg_groups 
WHERE 
  name = $1
  AND EXISTS (
    SELECT 1 
    FROM msg_group_members 
    WHERE 
      group_id = msg_groups.group_id 
      AND user_id = $2
  )
LIMIT 1
`

type GetMsgGroupByNameParams struct {
	Name   string    `json:"name"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) GetMsgGroupByName(ctx context.Context, arg GetMsgGroupByNameParams) (MsgGroup, error) {
	row := q.db.QueryRow(ctx, getMsgGroupByName, arg.Name, arg.UserID)
	var i MsgGroup
	err := row.Scan(
		&i.GroupID,
		&i.Name,
		&i.CreatedAt,
		&i.CreatedBy,
	)
	return i, err
}

const getMsgGroupMembers = `-- name: GetMsgGroupMembers :many
SELECT user_id, users.first_name, users.last_name, users.email
FROM msg_group_members 
INNER JOIN "User" as users ON msg_group_members.user_id = users.id 
WHERE msg_group_members.group_id = $1
`

type GetMsgGroupMembersRow struct {
	UserID    uuid.UUID `json:"user_id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
}

func (q *Queries) GetMsgGroupMembers(ctx context.Context, groupID int64) ([]GetMsgGroupMembersRow, error) {
	rows, err := q.db.Query(ctx, getMsgGroupMembers, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetMsgGroupMembersRow
	for rows.Next() {
		var i GetMsgGroupMembersRow
		if err := rows.Scan(
			&i.UserID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
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

const updateMessage = `-- name: UpdateMessage :many
UPDATE messages
SET
  content = coalesce($1, content),
  type = coalesce($2, type),
  status = coalesce($3, status),
  updated_at = NOW()
WHERE
  id = $4
  AND from_user = $5
  AND group_id = $6
RETURNING id
`

type UpdateMessageParams struct {
	Content  string        `json:"content"`
	Type     MESSAGETYPE   `json:"type"`
	Status   MESSAGESTATUS `json:"status"`
	ID       int64         `json:"id"`
	FromUser uuid.UUID     `json:"from_user"`
	GroupID  int64         `json:"group_id"`
}

func (q *Queries) UpdateMessage(ctx context.Context, arg UpdateMessageParams) ([]int64, error) {
	rows, err := q.db.Query(ctx, updateMessage,
		arg.Content,
		arg.Type,
		arg.Status,
		arg.ID,
		arg.FromUser,
		arg.GroupID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
