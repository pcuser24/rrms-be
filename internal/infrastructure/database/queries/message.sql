-- name: CreateMsgGroup :one
INSERT INTO msg_groups (
  name,
  created_by
) VALUES (
  sqlc.arg(name),
  sqlc.arg(created_by)
) RETURNING *;

-- name: CreateMsgGroupMember :one
INSERT INTO msg_group_members (
  group_id,
  user_id
) VALUES (
  sqlc.arg(group_id),
  sqlc.arg(user_id)
) RETURNING *;

-- name: CreateMessage :one
INSERT INTO messages (
  group_id,
  from_user,
  content
) VALUES (
  sqlc.arg(group_id),
  sqlc.arg(from_user),
  sqlc.arg(content)
) RETURNING *;

-- name: GetMsgGroup :one
SELECT * FROM msg_groups WHERE group_id = sqlc.arg(group_id);

-- name: GetMsgGroupMembers :many
SELECT user_id, users.first_name, users.last_name, users.email
FROM msg_group_members 
INNER JOIN "User" as users ON msg_group_members.user_id = users.id 
WHERE msg_group_members.group_id = sqlc.arg(group_id);

-- name: GetMessagesOfGroup :many
SELECT * 
FROM messages 
WHERE 
  group_id = sqlc.arg(group_id) 
ORDER BY
  created_at DESC 
LIMIT $1 OFFSET $2;

-- name: GetMsgGroupByName :one 
SELECT * 
FROM msg_groups 
WHERE 
  name = sqlc.arg(name)
  AND EXISTS (
    SELECT 1 
    FROM msg_group_members 
    WHERE 
      group_id = msg_groups.group_id 
      AND user_id = sqlc.arg(user_id)
  )
LIMIT 1;

-- name: CheckMsgGroupMembership :one
SELECT COUNT(*) > 0
FROM msg_group_members 
WHERE 
  group_id = sqlc.arg(group_id) AND 
  user_id = sqlc.arg(user_id);

-- name: UpdateMessage :many
UPDATE messages
SET
  content = coalesce(sqlc.arg(content), content),
  type = coalesce(sqlc.arg(type), type),
  status = coalesce(sqlc.arg(status), status),
  updated_at = NOW()
WHERE
  id = sqlc.arg(id)
  AND from_user = sqlc.arg(from_user)
  AND group_id = sqlc.arg(group_id)
RETURNING id;

-- name: DeleteMsgGroupMember :exec
DELETE FROM msg_group_members WHERE group_id = sqlc.arg(group_id) AND user_id = sqlc.arg(user_id);

-- name: DeleteMsgGroup :exec
DELETE FROM msg_groups WHERE group_id = sqlc.arg(group_id);
