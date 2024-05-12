-- name: CreateReminder :one
INSERT INTO "reminders" (
  "creator_id",
  "title",
  "start_at",
  "end_at",
  "note",
  "location",
  "priority",
  "recurrence_day",
  "recurrence_month",
  "recurrence_mode",
  "resource_tag"
) VALUES (
  sqlc.arg(creator_id),
  sqlc.arg(title),
  sqlc.arg(start_at),
  sqlc.arg(end_at),
  sqlc.narg(note),
  sqlc.arg(location),
  sqlc.narg(priority),
  sqlc.narg(recurrence_day),
  sqlc.narg(recurrence_month),
  sqlc.arg(recurrence_mode),
  sqlc.arg(resource_tag)
) RETURNING *;

-- name: GetReminderById :one
SELECT * FROM "reminders" WHERE "id" = $1;

-- name: GetRemindersByCreator :many
SELECT * FROM "reminders" WHERE "creator_id" = $1;

-- name: GetRemindersOfUserWithResourceTag :many
SELECT * 
FROM "reminders" 
WHERE 
  "creator_id" = $1
  AND "resource_tag" = $2;

-- name: GetRemindersInDate :many
SELECT * FROM "reminders" WHERE DATE_TRUNC('month', start_at) = DATE_TRUNC('month', $1);

-- name: CheckReminderVisibility :one
SELECT EXISTS(SELECT 1 FROM "reminders" WHERE "id" = $1 AND "creator_id" = $2);

-- name: UpdateReminder :many
UPDATE "reminders" SET
  "title" = coalesce(sqlc.narg(title), title),
  "start_at" = coalesce(sqlc.narg(start_at), start_at),
  "end_at" = coalesce(sqlc.narg(end_at), end_at),
  "note" = coalesce(sqlc.narg(note), note),
  "location" = coalesce(sqlc.narg(location), location),
  "priority" = coalesce(sqlc.narg(priority), priority),
  "recurrence_day" = coalesce(sqlc.narg(recurrence_day), recurrence_day),
  "recurrence_month" = coalesce(sqlc.narg(recurrence_month), recurrence_month),
  "recurrence_mode" = coalesce(sqlc.narg(recurrence_mode), recurrence_mode),
  "updated_at" = NOW()
WHERE 
  "id" = $1 
RETURNING *;

-- name: DeleteReminder :exec
DELETE FROM "reminders" WHERE "id" = $1;
