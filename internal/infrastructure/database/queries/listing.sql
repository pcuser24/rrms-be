-- name: CreateListing :one
INSERT INTO listings (
  creator_id,
  property_id,
  title,
  description,
  full_name,
  email,
  phone,
  contact_type,
  price,
  price_negotiable,
  security_deposit,
  lease_term,
  pets_allowed,
  number_of_residents,
  priority,
  created_at,
  updated_at,
  post_at,
  active,
  expired_at
) VALUES (
  sqlc.arg(creator_id),
  sqlc.arg(property_id),
  sqlc.arg(title),
  sqlc.arg(description),
  sqlc.arg(full_name),
  sqlc.arg(email),
  sqlc.arg(phone),
  sqlc.arg(contact_type),
  sqlc.arg(price),
  sqlc.narg(price_negotiable),
  sqlc.narg(security_deposit),
  sqlc.arg(lease_term),
  sqlc.narg(pets_allowed),
  sqlc.narg(number_of_residents),
  sqlc.arg(priority),
  NOW(), NOW(), 
  sqlc.arg(post_at),
  sqlc.arg(active),
  sqlc.arg(expired_at)
) RETURNING *;

-- name: CreateListingPolicy :one
INSERT INTO listing_policies (
  listing_id,
  policy_id,
  note
) VALUES (
  sqlc.arg(listing_id),
  sqlc.arg(policy_id),
  sqlc.narg(note)
) RETURNING *;

-- name: CreateListingUnit :one
INSERT INTO listing_units (
  listing_id,
  unit_id,
  price
) VALUES (
  sqlc.arg(listing_id),
  sqlc.arg(unit_id),
  sqlc.arg(price)
) RETURNING *;

-- name: GetListingByID :one
SELECT * FROM listings WHERE id = $1 LIMIT 1;

-- name: GetListingPolicies :many
SELECT * FROM listing_policies WHERE listing_id = $1;

-- name: GetListingUnits :many
SELECT * FROM listing_units WHERE listing_id = $1;

-- name: GetAllRentalPolicies :many
SELECT * FROM rental_policies;

-- name: CheckListingOwnership :one
SELECT count(*) FROM listings WHERE id = $1 AND creator_id = $2 LIMIT 1;

-- name: CheckValidUnitForListing :one
SELECT count(*) FROM units WHERE units.id = $1 AND units.property_id IN (SELECT listings.property_id FROM listings WHERE listings.id = $2) LIMIT 1;

-- name: UpdateListing :exec
UPDATE listings SET
  title = coalesce(sqlc.narg(title), title),
  description = coalesce(sqlc.narg(description), description),
  full_name = coalesce(sqlc.narg(full_name), full_name),
  email = coalesce(sqlc.narg(email), email),
  phone = coalesce(sqlc.narg(phone), phone),
  contact_type = coalesce(sqlc.narg(contact_type), contact_type),
  price = coalesce(sqlc.narg(price), price),
  price_negotiable = coalesce(sqlc.narg(price_negotiable), price_negotiable),
  security_deposit = coalesce(sqlc.narg(security_deposit), security_deposit),
  lease_term = coalesce(sqlc.narg(lease_term), lease_term),
  pets_allowed = coalesce(sqlc.narg(pets_allowed), pets_allowed),
  number_of_residents = coalesce(sqlc.narg(number_of_residents), number_of_residents),
  updated_at = NOW(),
  post_at = coalesce(sqlc.narg(post_at), post_at)
WHERE id = sqlc.arg(id);

-- name: UpdateListingStatus :exec
UPDATE listings SET active = $1 WHERE id = $2;

-- name: DeleteListing :exec
DELETE FROM listings WHERE id = $1;
