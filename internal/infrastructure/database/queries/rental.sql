-- name: CreatePreRental :one
INSERT INTO prerentals (
  application_id,
  creator_id,
  tenant_id,
  profile_image,
  property_id,
  unit_id,
  tenant_type,
  tenant_name,
  tenant_dob,
  tenant_identity,
  tenant_phone,
  tenant_email,
  tenant_address,
  contract_type,
  contract_content,
  contract_last_update_by,
  land_area,
  unit_area,
  start_date,
  movein_date,
  rental_period,
  rental_price,
  note
) VALUES (
  sqlc.narg(application_id),
  sqlc.arg(creator_id),
  sqlc.narg(tenant_id),
  sqlc.arg(profile_image),
  sqlc.arg(property_id),
  sqlc.arg(unit_id),
  sqlc.arg(tenant_type),
  sqlc.arg(tenant_name),
  sqlc.narg(tenant_dob),
  sqlc.arg(tenant_identity),
  sqlc.arg(tenant_phone),
  sqlc.arg(tenant_email),
  sqlc.arg(tenant_address),
  sqlc.narg(contract_type),
  sqlc.narg(contract_content),
  sqlc.arg(creator_id),
  sqlc.arg(land_area),
  sqlc.arg(unit_area),
  sqlc.narg(start_date),
  sqlc.arg(movein_date),
  sqlc.arg(rental_period),
  sqlc.arg(rental_price),
  sqlc.narg(note)
) RETURNING *;

-- name: CreatePreRentalCoap :one
INSERT INTO prerental_coaps (
  prerental_id,
  full_name,
  dob,
  job,
  income,
  email,
  phone,
  description
) VALUES (
  sqlc.arg(prerental_id),
  sqlc.arg(full_name),
  sqlc.narg(dob),
  sqlc.narg(job),
  sqlc.narg(income),
  sqlc.narg(email),
  sqlc.narg(phone),
  sqlc.narg(description)
) RETURNING *;

-- name: GetPreRental :one
SELECT * FROM prerentals WHERE id = $1;

-- name: GetPreRentalContract :one
SELECT id, contract_type, contract_content, contract_last_update_at, contract_last_update_by FROM prerentals WHERE id = $1;

-- name: GetPreRentalCoapByPreRentalID :many
SELECT * FROM prerental_coaps WHERE prerental_id = $1;

-- name: UpdatePreRental :exec
UPDATE prerentals SET
  tenant_id = coalesce(sqlc.narg(tenant_id), tenant_id),
  profile_image = coalesce(sqlc.narg(profile_image), profile_image),
  tenant_type = coalesce(sqlc.narg(tenant_type), tenant_type),
  tenant_name = coalesce(sqlc.narg(tenant_name), tenant_name),
  tenant_dob = coalesce(sqlc.narg(tenant_dob), tenant_dob),
  tenant_identity = coalesce(sqlc.narg(tenant_identity), tenant_identity),
  tenant_phone = coalesce(sqlc.narg(tenant_phone), tenant_phone),
  tenant_email = coalesce(sqlc.narg(tenant_email), tenant_email),
  tenant_address = coalesce(sqlc.narg(tenant_address), tenant_address),
  start_date = coalesce(sqlc.narg(start_date), start_date),
  movein_date = coalesce(sqlc.narg(movein_date), movein_date),
  rental_period = coalesce(sqlc.narg(rental_period), rental_period),
  rental_price = coalesce(sqlc.narg(rental_price), rental_price),
  note = coalesce(sqlc.narg(note), note),
  status = coalesce(sqlc.narg(status), status)
WHERE id = $1;

-- name: UpdatePreRentalContract :exec
UPDATE prerentals SET
  contract_type = coalesce(sqlc.narg(contract_type), contract_type),
  contract_content = coalesce(sqlc.narg(contract_content), contract_content),
  contract_last_update_at = NOW(),
  contract_last_update_by = sqlc.arg(contract_last_update_by)
WHERE id = $1;

-- name: DeletePreRental :exec
DELETE FROM prerentals WHERE id = $1;
