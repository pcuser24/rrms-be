-- name: CreateProperty :one
INSERT INTO properties (
  creator_id,
  name,
  building,
  project,
  area,
  number_of_floors,
  year_built,
  orientation,
  entrance_width,
  facade,
  full_address,
  district,
  city,
  ward,
  lat,
  lng,
  place_url,
  description,
  type,
  created_at,
  updated_at
) VALUES (
  sqlc.arg(creator_id),
  sqlc.arg(name),
  sqlc.narg(building),
  sqlc.narg(project),
  sqlc.arg(area),
  sqlc.narg(number_of_floors),
  sqlc.narg(year_built),
  sqlc.narg(orientation),
  sqlc.narg(entrance_width),
  sqlc.narg(facade),
  sqlc.arg(full_address),
  sqlc.arg(district),
  sqlc.arg(city),
  sqlc.narg(ward),
  sqlc.narg(lat),
  sqlc.narg(lng),
  sqlc.arg(place_url),
  sqlc.narg(description),
  sqlc.arg(type),
  NOW(),
  NOW()
) RETURNING *;

-- name: CreatePropertyManager :one
INSERT INTO property_managers (
  property_id,
  manager_id,
  role
) VALUES (
  sqlc.arg(property_id),
  sqlc.arg(manager_id),
  sqlc.arg(role)
) RETURNING *;

-- name: CreatePropertyMedia :one
INSERT INTO property_media (
  property_id,
  url,
  type,
  description
) VALUES (
  sqlc.arg(property_id),
  sqlc.arg(url),
  sqlc.arg(type),
  sqlc.narg(description)
) RETURNING *;

-- name: CreatePropertyFeature :one
INSERT INTO property_features (
  property_id,
  feature_id,
  description
) VALUES (
  sqlc.arg(property_id),
  sqlc.arg(feature_id),
  sqlc.narg(description)
) RETURNING *;

-- name: CreatePropertyTag :one
INSERT INTO property_tags (
  property_id,
  tag
) VALUES (
  sqlc.arg(property_id),
  sqlc.arg(tag)
) RETURNING *;

-- name: GetPropertyById :one
SELECT * FROM properties WHERE id = $1 LIMIT 1;

-- name: GetAllPropertyFeatures :many
SELECT * FROM p_features;

-- name: GetPropertyFeatures :many
SELECT * FROM property_features WHERE property_id = $1;

-- name: GetPropertyTags :many
SELECT * FROM property_tags WHERE property_id = $1;

-- name: GetPropertyMedia :many
SELECT * FROM property_media WHERE property_id = $1;

-- name: GetPropertyManagers :many
SELECT * FROM property_managers WHERE property_id = $1;

-- name: GetManagedProperties :many
SELECT property_id, role FROM property_managers WHERE manager_id = $1;

-- name: IsPropertyPublic :one
SELECT is_public FROM properties WHERE id = $1 LIMIT 1;

-- name: UpdateProperty :exec
UPDATE properties SET
  name = coalesce(sqlc.narg(name), name),
  building = coalesce(sqlc.narg(building), building),
  project = coalesce(sqlc.narg(project), project),
  area = coalesce(sqlc.narg(area), area),
  number_of_floors = coalesce(sqlc.narg(number_of_floors), number_of_floors),
  year_built = coalesce(sqlc.narg(year_built), year_built),
  orientation = coalesce(sqlc.narg(orientation), orientation),
  entrance_width = coalesce(sqlc.narg(entrance_width), entrance_width),
  facade = coalesce(sqlc.narg(facade), facade),
  full_address = coalesce(sqlc.narg(full_address), full_address),
  district = coalesce(sqlc.narg(district), district),
  city = coalesce(sqlc.narg(city), city),
  ward = coalesce(sqlc.narg(ward), ward),
  lat = coalesce(sqlc.narg(lat), lat),
  lng = coalesce(sqlc.narg(lng), lng),
  place_url = coalesce(sqlc.narg(place_url), place_url),
  description = coalesce(sqlc.narg(description), description),
  is_public = coalesce(sqlc.narg(is_public), is_public),
  updated_at = NOW()
WHERE id = $1;

-- name: UpdatePropertyManager :exec
UPDATE property_managers SET
  role = sqlc.arg(role)
WHERE property_id = $1 AND manager_id = $2;

-- name: ChangePropertyVisibility :exec
UPDATE properties SET
  is_public = sqlc.arg(is_public)
WHERE id = $1;

-- name: DeletePropertyManager :exec
DELETE FROM property_managers WHERE property_id = $1 AND manager_id = $2;

-- name: DeletePropertyMedia :exec
DELETE FROM property_media WHERE property_id = $1 AND id = $2;

-- name: DeletePropertyFeature :exec
DELETE FROM property_features WHERE property_id = $1 AND feature_id = $2;

-- name: DeletePropertyTag :exec
DELETE FROM property_tags WHERE property_id = $1 AND id = $2;

-- name: DeleteProperty :exec
DELETE FROM properties WHERE id = $1;
