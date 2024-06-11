-- name: CreateNotificationDevice :one
INSERT INTO user_notification_devices (
  user_id,
  session_id,
  token,
  platform
) VALUES (
  sqlc.arg(user_id),
  sqlc.arg(session_id),
  sqlc.arg(token),
  sqlc.arg(platform)
) ON CONFLICT ("user_id", "session_id") DO UPDATE 
SET "token" = EXCLUDED."token","last_accessed" = EXCLUDED."last_accessed"
RETURNING *;

-- name: GetNotificationDevice :one
SELECT * 
FROM user_notification_devices 
WHERE 
  user_id = sqlc.arg(user_id) 
  AND session_id = sqlc.arg(session_id)
  AND platform = sqlc.arg(platform)::"PLATFORM"
  AND CASE
    WHEN sqlc.arg(token)::TEXT <> '' THEN token = sqlc.arg(token)::TEXT
    ELSE TRUE
  END
LIMIT 1;

-- name: UpdateNotificationDeviceTokenTimestamp :exec
UPDATE "user_notification_devices"
SET "last_accessed" = NOW()
WHERE 
  "user_id" = sqlc.arg(user_id) AND 
  "session_id" = sqlc.arg(session_id);

-- name: DeleteNotificationDeviceToken :exec
DELETE FROM user_notification_devices 
WHERE 
  "user_id" = sqlc.arg(user_id) AND 
  "session_id" = sqlc.arg(session_id) AND 
  "token" = sqlc.arg(current_token);

-- name: DeleteExpiredTokens :exec
DELETE FROM user_notification_devices 
WHERE 
  "last_accessed" < NOW() - sqlc.arg(interval)::INTEGER * INTERVAL '1 day';


-- name: CreateNotification :one
INSERT INTO notifications (
  user_id,
  title,
  content,
  data,
  email,
  push,
  sms
) VALUES (
  sqlc.narg(user_id),
  sqlc.arg(title),
  sqlc.arg(content),
  sqlc.narg(data),
  sqlc.arg(email),
  sqlc.arg(push),
  sqlc.arg(sms)
) RETURNING *;

-- name: GetNotificationsOfUser :many
SELECT *
FROM notifications
WHERE user_id = sqlc.arg(user_id)
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: GetNotification :one
SELECT *
FROM notifications
WHERE id = sqlc.arg(id)
LIMIT 1;

-- name: UpdatedNotification :exec
UPDATE notifications
SET
  title = coalesce(sqlc.narg(title), title),
  content = coalesce(sqlc.narg(content), content),
  data = coalesce(sqlc.narg(data), data),
  email = coalesce(sqlc.narg(email), email),
  push = coalesce(sqlc.narg(push), push),
  sms = coalesce(sqlc.narg(sms), sms),
  seen = coalesce(sqlc.narg(seen), seen),
  updated_at = NOW()
WHERE id = sqlc.arg(id);

