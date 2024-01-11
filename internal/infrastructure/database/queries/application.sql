-- name: CreateApplication :one
INSERT INTO applications (
  creator_id,
  listing_id,
  property_id,
  unit_ids,
  -- basic info
  full_name,
  dob,
  email,
  phone,
  profile_image,
  movein_date,
  preferred_term,
  -- rental history
  rh_address,
  rh_city,
  rh_district,
  rh_ward,
  rh_rental_duration,
  rh_monthly_payment,
  rh_reason_for_leaving,
  -- employment
  employment_status,
  employment_company_name,
  employment_position,
  employment_monthly_income,
  employment_comment,
  employment_proofs_of_income,
  -- identity
  identity_type,
  identity_number,
  identity_issued_date,
  identity_issued_by
) VALUES (
  sqlc.arg(creator_id),
  sqlc.arg(listing_id),
  sqlc.arg(property_id),
  sqlc.arg(unit_ids),
  -- basic info
  sqlc.arg(full_name),
  sqlc.arg(dob),
  sqlc.arg(email),
  sqlc.arg(phone),
  sqlc.arg(profile_image),
  sqlc.arg(movein_date),
  sqlc.arg(preferred_term),
  -- rental history
  sqlc.narg(rh_address),
  sqlc.narg(rh_city),
  sqlc.narg(rh_district),
  sqlc.narg(rh_ward),
  sqlc.narg(rh_rental_duration),
  sqlc.narg(rh_monthly_payment),
  sqlc.narg(rh_reason_for_leaving),
  -- employment
  sqlc.arg(employment_status),
  sqlc.narg(employment_company_name),
  sqlc.narg(employment_position),
  sqlc.narg(employment_monthly_income),
  sqlc.narg(employment_comment),
  sqlc.narg(employment_proofs_of_income),
  -- identity
  sqlc.arg(identity_type),
  sqlc.arg(identity_number),
  sqlc.arg(identity_issued_date),
  sqlc.arg(identity_issued_by)
) RETURNING *;

-- name: CreateApplicationCoap :one
INSERT INTO application_coaps (
  application_id,
  full_name,
  dob,
  job,
  income,
  email,
  phone,
  description
) VALUES (
  sqlc.arg(application_id),
  sqlc.arg(full_name),
  sqlc.arg(dob),
  sqlc.arg(job),
  sqlc.arg(income),
  sqlc.narg(email),
  sqlc.narg(phone),
  sqlc.narg(description)
) RETURNING *;

-- name: CreateApplicationMinor :one
INSERT INTO application_minors (
  application_id,
  full_name,
  dob,
  email,
  phone,
  description
) VALUES (
  sqlc.arg(application_id),
  sqlc.arg(full_name),
  sqlc.arg(dob),
  sqlc.narg(email),
  sqlc.narg(phone),
  sqlc.narg(description)
) RETURNING *;

-- name: CreateApplicationPet :one
INSERT INTO application_pets (
  application_id,
  type,
  weight,
  description
) VALUES (
  sqlc.arg(application_id),
  sqlc.arg(type),
  sqlc.narg(weight),
  sqlc.narg(description)
) RETURNING *;

-- name: CreateApplicationVehicle :one
INSERT INTO application_vehicles (
  application_id,
  type,
  model,
  code,
  description
) VALUES (
  sqlc.arg(application_id),
  sqlc.arg(type),
  sqlc.narg(model),
  sqlc.arg(code),
  sqlc.narg(description)
) RETURNING *;

-- name: GetApplicationByID :one
SELECT * FROM applications WHERE id = $1 LIMIT 1;

-- name: GetApplicationMinors :many
SELECT * FROM application_minors WHERE application_id = $1;

-- name: GetApplicationCoaps :many
SELECT * FROM application_coaps WHERE application_id = $1;

-- name: GetApplicationPets :many
SELECT * FROM application_pets WHERE application_id = $1;

-- name: GetApplicationVehicles :many
SELECT * FROM application_vehicles WHERE application_id = $1;

-- name: GetApplicationsByUserId :many
SELECT * FROM applications WHERE creator_id = $1;

-- name: GetApplicationsToUser :many
SELECT * FROM applications WHERE property_id IN (
  SELECT property_id FROM property_managers WHERE manager_id = $1
);

-- name: UpdateApplicationStatus :exec
UPDATE applications SET status = $1, updated_at = NOW() WHERE id = $2;

-- name: DeleteApplication :exec
DELETE FROM applications WHERE id = $1;
