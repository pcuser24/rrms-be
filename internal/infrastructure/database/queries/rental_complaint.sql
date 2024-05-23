-- name: CreateRentalComplaint :one
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
  sqlc.arg(rental_id),
  sqlc.arg(creator_id),
  sqlc.arg(title),
  sqlc.arg(content),
  sqlc.narg(suggestion),
  sqlc.arg(media),
  sqlc.arg(occurred_at),
  sqlc.arg(type),
  sqlc.arg(creator_id)
) RETURNING *;

-- name: UpdateRentalComplaint :exec
UPDATE rental_complaints
SET
  title = coalesce(sqlc.narg(title), title),
  content = coalesce(sqlc.narg(content), content),
  suggestion = coalesce(sqlc.narg(suggestion), suggestion),
  media = coalesce(sqlc.narg(media), media),
  occurred_at = coalesce(sqlc.narg(occurred_at), occurred_at),
  status = coalesce(sqlc.narg(status), status),
  updated_at = NOW(),
  updated_by = sqlc.arg(user_id)
WHERE id = sqlc.arg(id);

-- name: GetRentalComplaint :one
SELECT * FROM rental_complaints WHERE id = $1 LIMIT 1;

-- name: GetRentalComplaintsByRentalId :many
SELECT * FROM rental_complaints WHERE rental_id = $1;

-- name: CreateRentalComplaintReply :one
INSERT INTO rental_complaint_replies (
  complaint_id, 
  replier_id, 
  reply,
  media
) VALUES (
  sqlc.arg(complaint_id),
  sqlc.arg(replier_id),
  sqlc.arg(reply),
  sqlc.narg(media)
) RETURNING *;

-- name: GetRentalComplaintReplies :many
SELECT * FROM rental_complaint_replies WHERE complaint_id = $1;

-- SELECT 
--   rental_complaint_replies.*, "User".first_name as replier_firstname, "User".last_name as replier_lastname 
-- FROM 
--   rental_complaint_replies LEFT JOIN "User" ON "User".id = rental_complaint_replies.replier_id 
-- WHERE 
--   rental_complaint_replies.complaint_id = $1;