-- name: CreateRental :one
INSERT INTO rentals (
  application_id,
  creator_id,
  property_id,
  unit_id,
  profile_image,
  
  tenant_id,
  tenant_type,
  tenant_name,
  tenant_phone,
  tenant_email,
  organization_name,
  organization_hq_address,

  start_date,
  movein_date,
  rental_period,
  rental_price,
  rental_intention,
  deposit,
  deposit_paid,

  electricity_payment_type,
  electricity_price,
  water_payment_type,
  water_price,

  note
) VALUES (
  sqlc.narg(application_id),
  sqlc.arg(creator_id),
  sqlc.arg(property_id),
  sqlc.arg(unit_id),
  sqlc.arg(profile_image),
  
  sqlc.narg(tenant_id),
  sqlc.arg(tenant_type),
  sqlc.arg(tenant_name),
  sqlc.arg(tenant_phone),
  sqlc.arg(tenant_email),
  sqlc.narg(organization_name),
  sqlc.narg(organization_hq_address),

  sqlc.arg(start_date),
  sqlc.arg(movein_date),
  sqlc.arg(rental_period),
  sqlc.arg(rental_price),
  sqlc.arg(rental_intention),
  sqlc.arg(deposit),
  sqlc.arg(deposit_paid),

  sqlc.arg(electricity_payment_type),
  sqlc.narg(electricity_price),
  sqlc.arg(water_payment_type),
  sqlc.narg(water_price),

  sqlc.narg(note)
) RETURNING *;

-- name: CreateRentalCoap :one
INSERT INTO rental_coaps (
  rental_id,
  full_name,
  dob,
  job,
  income,
  email,
  phone,
  description
) VALUES (
  sqlc.arg(rental_id),
  sqlc.arg(full_name),
  sqlc.narg(dob),
  sqlc.narg(job),
  sqlc.narg(income),
  sqlc.narg(email),
  sqlc.narg(phone),
  sqlc.narg(description)
) RETURNING *;

-- name: CreateRentalMinor :one
INSERT INTO "rental_minors" (
  "rental_id",
  "full_name",
  "dob",
  "email",
  "phone",
  "description"
) VALUES (
  sqlc.arg(rental_id),
  sqlc.arg(full_name),
  sqlc.arg(dob),
  sqlc.narg(email),
  sqlc.narg(phone),
  sqlc.narg(description)
) RETURNING *;

-- name: CreateRentalPet :one
INSERT INTO "rental_pets" (
  "rental_id",
  "type",
  "weight",
  "description"
) VALUES (
  sqlc.arg(rental_id),
  sqlc.arg(type),
  sqlc.narg(weight),
  sqlc.narg(description)
) RETURNING *;

-- name: CreateRentalService :one
INSERT INTO "rental_services" (
  "rental_id",
  "name",
  "setupBy",
  "provider",
  "price"
) VALUES (
  sqlc.arg(rental_id),
  sqlc.arg(name),
  sqlc.arg(setupBy),
  sqlc.narg(provider),
  sqlc.narg(price)
) RETURNING *;

-- name: GetRental :one
SELECT * FROM rentals WHERE id = $1 LIMIT 1;

-- name: GetRentalByApplicationId :one
SELECT * FROM rentals WHERE application_id = $1 LIMIT 1;

-- name: GetRentalCoapsByRentalID :many
SELECT * FROM rental_coaps WHERE rental_id = $1;

-- name: GetRentalMinorsByRentalID :many
SELECT * FROM rental_minors WHERE rental_id = $1;

-- name: GetRentalPetsByRentalID :many
SELECT * FROM rental_pets WHERE rental_id = $1;

-- name: GetRentalServicesByRentalID :many
SELECT * FROM rental_services WHERE rental_id = $1;

-- name: CheckRentalVisibility :one
SELECT count(*) > 0 FROM rentals 
WHERE 
  id = $1 
  AND (
    tenant_id = sqlc.arg(user_id) 
    OR EXISTS (
      SELECT 1 FROM property_managers 
      WHERE 
        property_managers.property_id = rentals.property_id 
        AND property_managers.manager_id = sqlc.arg(user_id)
      )
  );

-- name: UpdateRental :exec
UPDATE rentals SET
  tenant_id = coalesce(sqlc.narg(tenant_id), tenant_id),
  profile_image = coalesce(sqlc.narg(profile_image), profile_image),
  tenant_type = coalesce(sqlc.narg(tenant_type), tenant_type),
  tenant_name = coalesce(sqlc.narg(tenant_name), tenant_name),
  tenant_phone = coalesce(sqlc.narg(tenant_phone), tenant_phone),
  tenant_email = coalesce(sqlc.narg(tenant_email), tenant_email),
  organization_name = coalesce(sqlc.narg(organization_name), organization_name),
  organization_hq_address = coalesce(sqlc.narg(organization_hq_address), organization_hq_address),
  start_date = coalesce(sqlc.narg(start_date), start_date),
  movein_date = coalesce(sqlc.narg(movein_date), movein_date),
  rental_period = coalesce(sqlc.narg(rental_period), rental_period),
  rental_price = coalesce(sqlc.narg(rental_price), rental_price),
  rental_intention = coalesce(sqlc.narg(rental_intention), rental_intention),
  deposit = coalesce(sqlc.narg(deposit), deposit),
  deposit_paid = coalesce(sqlc.narg(deposit_paid), deposit_paid),
  electricity_payment_type = coalesce(sqlc.narg(electricity_payment_type), electricity_payment_type),
  electricity_price = coalesce(sqlc.narg(electricity_price), electricity_price),
  water_payment_type = coalesce(sqlc.narg(water_payment_type), water_payment_type),
  water_price = coalesce(sqlc.narg(water_price), water_price),
  note = coalesce(sqlc.narg(note), note),
  updated_at = NOW()
WHERE id = $1;

-- name: DeleteRental :exec
DELETE FROM rentals WHERE id = $1;
