-- name: CreateUnit :one
INSERT INTO units (
  property_id,
  name,
  area,
  floor,
  has_balcony,
  number_of_living_rooms,
  number_of_bedrooms,
  number_of_bathrooms,
  number_of_toilets,
  number_of_kitchens,
  type,
  created_at,
  updated_at
) VALUES (
  sqlc.arg(property_id),
  sqlc.narg(name),
  sqlc.arg(area),
  sqlc.narg(floor),
  sqlc.narg(has_balcony),
  sqlc.narg(number_of_living_rooms),
  sqlc.narg(number_of_bedrooms),
  sqlc.narg(number_of_bathrooms),
  sqlc.narg(number_of_toilets),
  sqlc.narg(number_of_kitchens),
  sqlc.arg(type),
  NOW(),
  NOW()
) RETURNING *;

-- name: CreateUnitMedia :one
INSERT INTO unit_media (
  unit_id,
  url,
  type
) VALUES (
  sqlc.arg(unit_id),
  sqlc.arg(url),
  sqlc.arg(type)
) RETURNING *;

-- name: CreateUnitAmenity :one
INSERT INTO unit_amenity (
  unit_id,
  amenity
) VALUES (
  sqlc.arg(unit_id),
  sqlc.arg(amenity)
) RETURNING *;

-- name: DeleteUnitMedia :exec
DELETE FROM unit_media WHERE unit_id = $1 AND url = $2;

-- name: DeleteUnitAmenity :exec
DELETE FROM unit_amenity WHERE unit_id = $1 AND amenity = $2;

-- name: GetUnitById :one
SELECT * FROM units WHERE id = $1 LIMIT 1;

-- name: GetUnitsOfProperty :many
SELECT * FROM units WHERE property_id = $1;

-- name: UpdateUnit :exec
UPDATE units SET
  name = coalesce(sqlc.narg(name), name),
  area = coalesce(sqlc.narg(area), area),
  floor = coalesce(sqlc.narg(floor), floor),
  has_balcony = coalesce(sqlc.narg(has_balcony), has_balcony),
  number_of_living_rooms = coalesce(sqlc.narg(number_of_living_rooms), number_of_living_rooms),
  number_of_bedrooms = coalesce(sqlc.narg(number_of_bedrooms), number_of_bedrooms),
  number_of_bathrooms = coalesce(sqlc.narg(number_of_bathrooms), number_of_bathrooms),
  number_of_toilets = coalesce(sqlc.narg(number_of_toilets), number_of_toilets),
  number_of_kitchens = coalesce(sqlc.narg(number_of_kitchens), number_of_kitchens),
  updated_at = NOW()
WHERE id = $1;

-- name: DeleteUnit :exec
DELETE FROM units WHERE id = $1;