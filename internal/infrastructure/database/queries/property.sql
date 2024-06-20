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

-- name: CreateNewPropertyManagerRequest :one
INSERT INTO "new_property_manager_requests" (
  "creator_id",
  "property_id",
  "user_id",
  "email",
  "created_at",
  "updated_at"
) VALUES (
  sqlc.arg(creator_id), 
  sqlc.arg(property_id),
  sqlc.narg(user_id),
  sqlc.arg(email),
  NOW(), NOW()
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

-- name: GetNewPropertyManagerRequest :one
SELECT * FROM "new_property_manager_requests" WHERE "id" = $1 LIMIT 1;

-- name: GetNewPropertyManagerRequestsToUser :many
SELECT * 
FROM "new_property_manager_requests" 
WHERE "user_id" = $1
ORDER BY "created_at" DESC
LIMIT $2
OFFSET $3;

-- name: IsPropertyVisible :one
SELECT (
  SELECT is_public FROM "properties" WHERE properties.id = sqlc.arg(property_id) LIMIT 1
) OR (
  SELECT EXISTS (SELECT 1 FROM "User" WHERE "User".id = sqlc.arg(user_id) AND "User".role = 'ADMIN' LIMIT 1)
) OR (
  SELECT EXISTS (SELECT 1 FROM "property_managers" WHERE property_managers.property_id = sqlc.arg(property_id) AND property_managers.manager_id = sqlc.arg(user_id) LIMIT 1)
) OR (
  SELECT EXISTS (SELECT 1 FROM "new_property_manager_requests" WHERE new_property_manager_requests.property_id = sqlc.arg(property_id) AND new_property_manager_requests.user_id = sqlc.arg(user_id) LIMIT 1)
);

-- name: UpdateNewPropertyManagerRequest :exec
UPDATE "new_property_manager_requests" SET
  "approved" = sqlc.arg(approved),
  "updated_at" = NOW()
WHERE "id" = $1;

-- name: AddPropertyManager :exec
INSERT INTO property_managers (
  property_id,
  manager_id,
  role
) VALUES (
  (SELECT property_id FROM new_property_manager_requests WHERE id = sqlc.arg(request_id) LIMIT 1),
  sqlc.arg(user_id),
  'MANAGER'
);

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
  primary_image = coalesce(sqlc.narg(primary_image), primary_image),
  description = coalesce(sqlc.narg(description), description),
  is_public = coalesce(sqlc.narg(is_public), is_public),
  updated_at = NOW()
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

-- name: CreatePropertyVerificationRequest :one
INSERT INTO "property_verification_requests" (
  "creator_id",
  "property_id",
  "video_url",
  "house_ownership_certificate",
  "certificate_of_landuse_right",
  "front_idcard",
  "back_idcard",
  "note",
  "created_at",
  "updated_at"
) VALUES (
  sqlc.arg(creator_id),
  sqlc.arg(property_id),
  sqlc.arg(video_url),
  sqlc.narg(house_ownership_certificate),
  sqlc.narg(certificate_of_landuse_right),
  sqlc.arg(front_idcard),
  sqlc.arg(back_idcard),
  sqlc.narg(note),
  NOW(),
  NOW()
) RETURNING *;

-- name: GetPropertyVerificationRequest :one
SELECT * FROM "property_verification_requests" WHERE "id" = $1 LIMIT 1;

-- name: GetPropertyVerificationRequestsOfProperty :many
SELECT * FROM "property_verification_requests" WHERE "property_id" = sqlc.arg(property_id) ORDER BY "updated_at" DESC LIMIT $1 OFFSET $2;

-- name: UpdatePropertyVerificationRequest :exec
UPDATE "property_verification_requests" SET
  "video_url" = coalesce(sqlc.narg(video_url), video_url),
  "house_ownership_certificate" = coalesce(sqlc.narg(house_ownership_certificate), house_ownership_certificate),
  "certificate_of_landuse_right" = coalesce(sqlc.narg(certificate_of_landuse_right), certificate_of_landuse_right),
  "front_idcard" = coalesce(sqlc.narg(front_idcard), front_idcard),
  "back_idcard" = coalesce(sqlc.narg(back_idcard), back_idcard),
  "note" = coalesce(sqlc.narg(note), note),
  "feedback" = coalesce(sqlc.narg(feedback), feedback),
  "status" = coalesce(sqlc.narg(status), status),
  "updated_at" = NOW()
WHERE "id" = $1;