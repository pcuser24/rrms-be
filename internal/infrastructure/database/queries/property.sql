-- name: CreateProperty :one
INSERT INTO properties (
  owner_id,
  name,
  area,
  number_of_floors,
  year_built,
  orientation,
  full_address,
  district,
  city,
  lat,
  lng,
  type,
  created_at,
  updated_at
) VALUES (
  sqlc.arg(owner_id),
  sqlc.narg(name),
  sqlc.arg(area),
  sqlc.narg(number_of_floors),
  sqlc.narg(year_built),
  sqlc.narg(orientation),
  sqlc.arg(full_address),
  sqlc.arg(district),
  sqlc.arg(city),
  sqlc.arg(lat),
  sqlc.arg(lng),
  sqlc.arg(type),
  NOW(),
  NOW()
) RETURNING *;

-- name: CreatePropertyMedia :one
INSERT INTO property_media (
  property_id,
  url,
  type
) VALUES (
  sqlc.arg(property_id),
  sqlc.arg(url),
  sqlc.arg(type)
) RETURNING *;

-- name: CreatePropertyFeature :one
INSERT INTO property_feature (
  property_id,
  feature_id,
  description
) VALUES (
  sqlc.arg(property_id),
  sqlc.arg(feature_id),
  sqlc.narg(description)
) RETURNING *;

-- name: GetPropertyById :one
SELECT * FROM properties WHERE id = $1 LIMIT 1;

-- name: GetAllPropertyFeatures :many
SELECT * FROM p_features;

-- name: GetPropertyFeatures :many
SELECT * FROM property_feature WHERE property_id = $1;

-- name: GetPropertyTags :many
SELECT * FROM property_tag WHERE property_id = $1;

-- name: GetPropertyMedium :many
SELECT * FROM property_media WHERE property_id = $1;

-- name: GetPropertiesByOwnerId :many
SELECT * FROM properties WHERE owner_id = $1 LIMIT $2 OFFSET $3;

-- name: CheckPropertyOwnerShip :one
SELECT count(*) FROM properties WHERE id = $1 AND owner_id = $2 LIMIT 1;

-- name: UpdateProperty :exec
UPDATE properties SET
  name = coalesce(sqlc.narg(name), name),
  area = coalesce(sqlc.narg(area), area),
  number_of_floors = coalesce(sqlc.narg(number_of_floors), number_of_floors),
  year_built = coalesce(sqlc.narg(year_built), year_built),
  orientation = coalesce(sqlc.narg(orientation), orientation),
  full_address = coalesce(sqlc.narg(full_address), full_address),
  district = coalesce(sqlc.narg(district), district),
  city = coalesce(sqlc.narg(city), city),
  lat = coalesce(sqlc.narg(lat), lat),
  lng = coalesce(sqlc.narg(lng), lng),
  updated_at = NOW()
WHERE id = $1;

-- name: DeleteProperty :exec
DELETE FROM properties WHERE id = $1;
