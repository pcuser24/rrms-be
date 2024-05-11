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

-- name: CreateReminderMember :one
INSERT INTO "reminder_members" (
  "reminder_id",
  "user_id"
) VALUES (
  sqlc.arg(reminder_id),
  sqlc.arg(user_id)
) RETURNING *;

-- name: GetReminderById :one
SELECT * FROM "reminders" WHERE "id" = $1;

-- name: GetReminderMembers :many
SELECT "user_id" FROM "reminder_members" WHERE "reminder_id" = $1;

-- name: GetRemindersByCreator :many
SELECT * FROM "reminders" WHERE "creator_id" = $1;

-- name: GetRemindersOfUser :many
SELECT * FROM "reminders" WHERE "id" IN (SELECT "reminder_id" FROM "reminder_members" WHERE "user_id" = $1);

-- name: GetRemindersOfUserWithResourceTag :many
SELECT * 
FROM "reminders" 
WHERE 
  "id" IN (SELECT "reminder_id" FROM "reminder_members" WHERE "user_id" = $1)
  AND "resource_tag" = $2;

-- name: GetRemindersInDate :many
SELECT * FROM "reminders" WHERE DATE_TRUNC('month', start_at) = DATE_TRUNC('month', $1);

-- name: CheckReminderVisibility :one
SELECT EXISTS(SELECT 1 FROM "reminder_members" WHERE "reminder_id" = $1 AND "user_id" = $2);

-- name: CheckOverlappingReminder :one
SELECT EXISTS(
  SELECT 1 
  FROM "reminders" 
  WHERE 
    status IN ('INPROGRESS', 'COMPLETED') AND
    EXISTS (
      SELECT 1 FROM reminder_members WHERE reminders.id = reminder_members.reminder_id AND reminder_members.user_id = $1
    ) AND (
      (start_at, end_at) OVERLAPS (sqlc.arg(start_time), sqlc.arg(end_time))
      OR (start_at >= sqlc.arg(start_time) AND start_at < sqlc.arg(end_time)) 
      OR (end_at > sqlc.arg(start_time) AND end_at <= sqlc.arg(end_time))
    )
);

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
  "status" = coalesce(sqlc.narg(status), status),
  "updated_at" = NOW()
WHERE 
  "id" = $1 
RETURNING *;

-- name: DeleteReminder :exec
DELETE FROM "reminders" WHERE "id" = $1;
