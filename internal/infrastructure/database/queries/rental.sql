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

  payment_type,

  rental_price,
  rental_payment_basis,
  rental_intention,
  deposit,
  deposit_paid,
  notice_period,

  electricity_setup_by,
  electricity_payment_type,
  electricity_price,
  electricity_customer_code,
  electricity_provider,
  water_setup_by,
  water_payment_type,
  water_price,
  water_customer_code,
  water_provider,

  -- rental_payment_grace_period,
  -- rental_payment_late_fee_percentage,

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

  sqlc.narg(payment_type),

  sqlc.arg(rental_price),
  sqlc.arg(rental_payment_basis),
  sqlc.arg(rental_intention),
  sqlc.arg(deposit),
  sqlc.arg(deposit_paid),
  sqlc.narg(notice_period),

  sqlc.arg(electricity_setup_by),
  sqlc.narg(electricity_payment_type),
  sqlc.narg(electricity_price),
  sqlc.narg(electricity_customer_code),
  sqlc.narg(electricity_provider),
  sqlc.arg(water_setup_by),
  sqlc.narg(water_payment_type),
  sqlc.narg(water_price),
  sqlc.narg(water_customer_code),
  sqlc.narg(water_provider),

  -- sqlc.arg(rental_payment_grace_period),
  -- sqlc.narg(rental_payment_late_fee_percentage),
  
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
  "setup_by",
  "provider",
  "price"
) VALUES (
  sqlc.arg(rental_id),
  sqlc.arg(name),
  sqlc.arg(setup_by),
  sqlc.narg(provider),
  sqlc.narg(price)
) RETURNING *;

-- name: CreateRentalPolicy :one
INSERT INTO "rental_policies" (
  "rental_id",
  "title",
  "content"
) VALUES (
  sqlc.arg(rental_id),
  sqlc.arg(title),
  sqlc.arg(content)
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

-- name: GetRentalPoliciesByRentalID :many
SELECT * FROM rental_policies WHERE rental_id = $1;

-- name: GetRentalServicesByRentalID :many
SELECT * FROM rental_services WHERE rental_id = $1;

-- Get rental side: Side A (lanlord and managers) and Side B (tenant). Otherwise return C
-- name: GetRentalSide :one
SELECT 
  CASE 
    WHEN EXISTS (SELECT 1 FROM property_managers WHERE property_id = rentals.property_id AND manager_id = sqlc.arg(user_id)) THEN 'A'
    WHEN rentals.tenant_id = sqlc.arg(user_id) THEN 'B'
    ELSE 'C'
  END AS side
FROM rentals
WHERE id = $1;

-- name: GetManagedRentals :many
SELECT id FROM rentals WHERE property_id IN (SELECT property_id FROM property_managers WHERE manager_id = sqlc.arg(user_id));

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
  rental_payment_basis = coalesce(sqlc.narg(rental_payment_basis), rental_payment_basis),
  rental_intention = coalesce(sqlc.narg(rental_intention), rental_intention),
  notice_period = coalesce(sqlc.narg(notice_period), notice_period),
  deposit = coalesce(sqlc.narg(deposit), deposit),
  deposit_paid = coalesce(sqlc.narg(deposit_paid), deposit_paid),
  electricity_setup_by = coalesce(sqlc.narg(electricity_setup_by), electricity_setup_by),
  electricity_payment_type = coalesce(sqlc.narg(electricity_payment_type), electricity_payment_type),
  electricity_price = coalesce(sqlc.narg(electricity_price), electricity_price),
  electricity_customer_code = coalesce(sqlc.narg(electricity_customer_code), electricity_customer_code),
  electricity_provider = coalesce(sqlc.narg(electricity_provider), electricity_provider),
  water_setup_by = coalesce(sqlc.narg(water_setup_by), water_setup_by),
  water_payment_type = coalesce(sqlc.narg(water_payment_type), water_payment_type),
  water_price = coalesce(sqlc.narg(water_price), water_price),
  water_customer_code = coalesce(sqlc.narg(water_customer_code), water_customer_code),
  water_provider = coalesce(sqlc.narg(water_provider), water_provider),
  -- rental_payment_grace_period = coalesce(sqlc.narg(rental_payment_grace_period), rental_payment_grace_period),
  -- rental_payment_late_fee_percentage = coalesce(sqlc.narg(rental_payment_late_fee_percentage), rental_payment_late_fee_percentage),
  note = coalesce(sqlc.narg(note), note),
  updated_at = NOW()
WHERE id = $1;

-- name: DeleteRental :exec
DELETE FROM rentals WHERE id = $1;
