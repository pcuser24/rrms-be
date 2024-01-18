// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: property.sql

package database

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const changePropertyVisibility = `-- name: ChangePropertyVisibility :exec
UPDATE properties SET
  is_public = $2
WHERE id = $1
`

type ChangePropertyVisibilityParams struct {
	ID       uuid.UUID `json:"id"`
	IsPublic bool      `json:"is_public"`
}

func (q *Queries) ChangePropertyVisibility(ctx context.Context, arg ChangePropertyVisibilityParams) error {
	_, err := q.db.Exec(ctx, changePropertyVisibility, arg.ID, arg.IsPublic)
	return err
}

const createProperty = `-- name: CreateProperty :one
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
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7,
  $8,
  $9,
  $10,
  $11,
  $12,
  $13,
  $14,
  $15,
  $16,
  $17,
  $18,
  $19,
  NOW(),
  NOW()
) RETURNING id, creator_id, name, building, project, area, number_of_floors, year_built, orientation, entrance_width, facade, full_address, city, district, ward, lat, lng, place_url, description, type, is_public, created_at, updated_at
`

type CreatePropertyParams struct {
	CreatorID      uuid.UUID     `json:"creator_id"`
	Name           pgtype.Text   `json:"name"`
	Building       pgtype.Text   `json:"building"`
	Project        pgtype.Text   `json:"project"`
	Area           float32       `json:"area"`
	NumberOfFloors pgtype.Int4   `json:"number_of_floors"`
	YearBuilt      pgtype.Int4   `json:"year_built"`
	Orientation    pgtype.Text   `json:"orientation"`
	EntranceWidth  pgtype.Float4 `json:"entrance_width"`
	Facade         pgtype.Float4 `json:"facade"`
	FullAddress    string        `json:"full_address"`
	District       string        `json:"district"`
	City           string        `json:"city"`
	Ward           pgtype.Text   `json:"ward"`
	Lat            pgtype.Float8 `json:"lat"`
	Lng            pgtype.Float8 `json:"lng"`
	PlaceUrl       string        `json:"place_url"`
	Description    pgtype.Text   `json:"description"`
	Type           PROPERTYTYPE  `json:"type"`
}

func (q *Queries) CreateProperty(ctx context.Context, arg CreatePropertyParams) (Property, error) {
	row := q.db.QueryRow(ctx, createProperty,
		arg.CreatorID,
		arg.Name,
		arg.Building,
		arg.Project,
		arg.Area,
		arg.NumberOfFloors,
		arg.YearBuilt,
		arg.Orientation,
		arg.EntranceWidth,
		arg.Facade,
		arg.FullAddress,
		arg.District,
		arg.City,
		arg.Ward,
		arg.Lat,
		arg.Lng,
		arg.PlaceUrl,
		arg.Description,
		arg.Type,
	)
	var i Property
	err := row.Scan(
		&i.ID,
		&i.CreatorID,
		&i.Name,
		&i.Building,
		&i.Project,
		&i.Area,
		&i.NumberOfFloors,
		&i.YearBuilt,
		&i.Orientation,
		&i.EntranceWidth,
		&i.Facade,
		&i.FullAddress,
		&i.City,
		&i.District,
		&i.Ward,
		&i.Lat,
		&i.Lng,
		&i.PlaceUrl,
		&i.Description,
		&i.Type,
		&i.IsPublic,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const createPropertyFeature = `-- name: CreatePropertyFeature :one
INSERT INTO property_features (
  property_id,
  feature_id,
  description
) VALUES (
  $1,
  $2,
  $3
) RETURNING property_id, feature_id, description
`

type CreatePropertyFeatureParams struct {
	PropertyID  uuid.UUID   `json:"property_id"`
	FeatureID   int64       `json:"feature_id"`
	Description pgtype.Text `json:"description"`
}

func (q *Queries) CreatePropertyFeature(ctx context.Context, arg CreatePropertyFeatureParams) (PropertyFeature, error) {
	row := q.db.QueryRow(ctx, createPropertyFeature, arg.PropertyID, arg.FeatureID, arg.Description)
	var i PropertyFeature
	err := row.Scan(&i.PropertyID, &i.FeatureID, &i.Description)
	return i, err
}

const createPropertyManager = `-- name: CreatePropertyManager :one
INSERT INTO property_managers (
  property_id,
  manager_id,
  role
) VALUES (
  $1,
  $2,
  $3
) RETURNING property_id, manager_id, role
`

type CreatePropertyManagerParams struct {
	PropertyID uuid.UUID `json:"property_id"`
	ManagerID  uuid.UUID `json:"manager_id"`
	Role       string    `json:"role"`
}

func (q *Queries) CreatePropertyManager(ctx context.Context, arg CreatePropertyManagerParams) (PropertyManager, error) {
	row := q.db.QueryRow(ctx, createPropertyManager, arg.PropertyID, arg.ManagerID, arg.Role)
	var i PropertyManager
	err := row.Scan(&i.PropertyID, &i.ManagerID, &i.Role)
	return i, err
}

const createPropertyMedia = `-- name: CreatePropertyMedia :one
INSERT INTO property_media (
  property_id,
  url,
  type,
  description
) VALUES (
  $1,
  $2,
  $3,
  $4
) RETURNING id, property_id, url, type, description
`

type CreatePropertyMediaParams struct {
	PropertyID  uuid.UUID   `json:"property_id"`
	Url         string      `json:"url"`
	Type        MEDIATYPE   `json:"type"`
	Description pgtype.Text `json:"description"`
}

func (q *Queries) CreatePropertyMedia(ctx context.Context, arg CreatePropertyMediaParams) (PropertyMedium, error) {
	row := q.db.QueryRow(ctx, createPropertyMedia,
		arg.PropertyID,
		arg.Url,
		arg.Type,
		arg.Description,
	)
	var i PropertyMedium
	err := row.Scan(
		&i.ID,
		&i.PropertyID,
		&i.Url,
		&i.Type,
		&i.Description,
	)
	return i, err
}

const createPropertyTag = `-- name: CreatePropertyTag :one
INSERT INTO property_tags (
  property_id,
  tag
) VALUES (
  $1,
  $2
) RETURNING id, property_id, tag
`

type CreatePropertyTagParams struct {
	PropertyID uuid.UUID `json:"property_id"`
	Tag        string    `json:"tag"`
}

func (q *Queries) CreatePropertyTag(ctx context.Context, arg CreatePropertyTagParams) (PropertyTag, error) {
	row := q.db.QueryRow(ctx, createPropertyTag, arg.PropertyID, arg.Tag)
	var i PropertyTag
	err := row.Scan(&i.ID, &i.PropertyID, &i.Tag)
	return i, err
}

const deleteProperty = `-- name: DeleteProperty :exec
DELETE FROM properties WHERE id = $1
`

func (q *Queries) DeleteProperty(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.Exec(ctx, deleteProperty, id)
	return err
}

const deletePropertyFeature = `-- name: DeletePropertyFeature :exec
DELETE FROM property_features WHERE property_id = $1 AND feature_id = $2
`

type DeletePropertyFeatureParams struct {
	PropertyID uuid.UUID `json:"property_id"`
	FeatureID  int64     `json:"feature_id"`
}

func (q *Queries) DeletePropertyFeature(ctx context.Context, arg DeletePropertyFeatureParams) error {
	_, err := q.db.Exec(ctx, deletePropertyFeature, arg.PropertyID, arg.FeatureID)
	return err
}

const deletePropertyManager = `-- name: DeletePropertyManager :exec
DELETE FROM property_managers WHERE property_id = $1 AND manager_id = $2
`

type DeletePropertyManagerParams struct {
	PropertyID uuid.UUID `json:"property_id"`
	ManagerID  uuid.UUID `json:"manager_id"`
}

func (q *Queries) DeletePropertyManager(ctx context.Context, arg DeletePropertyManagerParams) error {
	_, err := q.db.Exec(ctx, deletePropertyManager, arg.PropertyID, arg.ManagerID)
	return err
}

const deletePropertyMedia = `-- name: DeletePropertyMedia :exec
DELETE FROM property_media WHERE property_id = $1 AND id = $2
`

type DeletePropertyMediaParams struct {
	PropertyID uuid.UUID `json:"property_id"`
	ID         int64     `json:"id"`
}

func (q *Queries) DeletePropertyMedia(ctx context.Context, arg DeletePropertyMediaParams) error {
	_, err := q.db.Exec(ctx, deletePropertyMedia, arg.PropertyID, arg.ID)
	return err
}

const deletePropertyTag = `-- name: DeletePropertyTag :exec
DELETE FROM property_tags WHERE property_id = $1 AND id = $2
`

type DeletePropertyTagParams struct {
	PropertyID uuid.UUID `json:"property_id"`
	ID         int64     `json:"id"`
}

func (q *Queries) DeletePropertyTag(ctx context.Context, arg DeletePropertyTagParams) error {
	_, err := q.db.Exec(ctx, deletePropertyTag, arg.PropertyID, arg.ID)
	return err
}

const getAllPropertyFeatures = `-- name: GetAllPropertyFeatures :many
SELECT id, feature FROM p_features
`

func (q *Queries) GetAllPropertyFeatures(ctx context.Context) ([]PFeature, error) {
	rows, err := q.db.Query(ctx, getAllPropertyFeatures)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PFeature
	for rows.Next() {
		var i PFeature
		if err := rows.Scan(&i.ID, &i.Feature); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getManagedProperties = `-- name: GetManagedProperties :many
SELECT property_id, role FROM property_managers WHERE manager_id = $1
`

type GetManagedPropertiesRow struct {
	PropertyID uuid.UUID `json:"property_id"`
	Role       string    `json:"role"`
}

func (q *Queries) GetManagedProperties(ctx context.Context, managerID uuid.UUID) ([]GetManagedPropertiesRow, error) {
	rows, err := q.db.Query(ctx, getManagedProperties, managerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetManagedPropertiesRow
	for rows.Next() {
		var i GetManagedPropertiesRow
		if err := rows.Scan(&i.PropertyID, &i.Role); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPropertyById = `-- name: GetPropertyById :one
SELECT id, creator_id, name, building, project, area, number_of_floors, year_built, orientation, entrance_width, facade, full_address, city, district, ward, lat, lng, place_url, description, type, is_public, created_at, updated_at FROM properties WHERE id = $1 LIMIT 1
`

func (q *Queries) GetPropertyById(ctx context.Context, id uuid.UUID) (Property, error) {
	row := q.db.QueryRow(ctx, getPropertyById, id)
	var i Property
	err := row.Scan(
		&i.ID,
		&i.CreatorID,
		&i.Name,
		&i.Building,
		&i.Project,
		&i.Area,
		&i.NumberOfFloors,
		&i.YearBuilt,
		&i.Orientation,
		&i.EntranceWidth,
		&i.Facade,
		&i.FullAddress,
		&i.City,
		&i.District,
		&i.Ward,
		&i.Lat,
		&i.Lng,
		&i.PlaceUrl,
		&i.Description,
		&i.Type,
		&i.IsPublic,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getPropertyFeatures = `-- name: GetPropertyFeatures :many
SELECT property_id, feature_id, description FROM property_features WHERE property_id = $1
`

func (q *Queries) GetPropertyFeatures(ctx context.Context, propertyID uuid.UUID) ([]PropertyFeature, error) {
	rows, err := q.db.Query(ctx, getPropertyFeatures, propertyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PropertyFeature
	for rows.Next() {
		var i PropertyFeature
		if err := rows.Scan(&i.PropertyID, &i.FeatureID, &i.Description); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPropertyManagers = `-- name: GetPropertyManagers :many
SELECT property_id, manager_id, role FROM property_managers WHERE property_id = $1
`

func (q *Queries) GetPropertyManagers(ctx context.Context, propertyID uuid.UUID) ([]PropertyManager, error) {
	rows, err := q.db.Query(ctx, getPropertyManagers, propertyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PropertyManager
	for rows.Next() {
		var i PropertyManager
		if err := rows.Scan(&i.PropertyID, &i.ManagerID, &i.Role); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPropertyMedia = `-- name: GetPropertyMedia :many
SELECT id, property_id, url, type, description FROM property_media WHERE property_id = $1
`

func (q *Queries) GetPropertyMedia(ctx context.Context, propertyID uuid.UUID) ([]PropertyMedium, error) {
	rows, err := q.db.Query(ctx, getPropertyMedia, propertyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PropertyMedium
	for rows.Next() {
		var i PropertyMedium
		if err := rows.Scan(
			&i.ID,
			&i.PropertyID,
			&i.Url,
			&i.Type,
			&i.Description,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPropertyTags = `-- name: GetPropertyTags :many
SELECT id, property_id, tag FROM property_tags WHERE property_id = $1
`

func (q *Queries) GetPropertyTags(ctx context.Context, propertyID uuid.UUID) ([]PropertyTag, error) {
	rows, err := q.db.Query(ctx, getPropertyTags, propertyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []PropertyTag
	for rows.Next() {
		var i PropertyTag
		if err := rows.Scan(&i.ID, &i.PropertyID, &i.Tag); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const isPropertyPublic = `-- name: IsPropertyPublic :one
SELECT is_public FROM properties WHERE id = $1 LIMIT 1
`

func (q *Queries) IsPropertyPublic(ctx context.Context, id uuid.UUID) (bool, error) {
	row := q.db.QueryRow(ctx, isPropertyPublic, id)
	var is_public bool
	err := row.Scan(&is_public)
	return is_public, err
}

const updateProperty = `-- name: UpdateProperty :exec
UPDATE properties SET
  name = coalesce($2, name),
  building = coalesce($3, building),
  project = coalesce($4, project),
  area = coalesce($5, area),
  number_of_floors = coalesce($6, number_of_floors),
  year_built = coalesce($7, year_built),
  orientation = coalesce($8, orientation),
  entrance_width = coalesce($9, entrance_width),
  facade = coalesce($10, facade),
  full_address = coalesce($11, full_address),
  district = coalesce($12, district),
  city = coalesce($13, city),
  ward = coalesce($14, ward),
  lat = coalesce($15, lat),
  lng = coalesce($16, lng),
  place_url = coalesce($17, place_url),
  description = coalesce($18, description),
  updated_at = NOW()
WHERE id = $1
`

type UpdatePropertyParams struct {
	ID             uuid.UUID     `json:"id"`
	Name           pgtype.Text   `json:"name"`
	Building       pgtype.Text   `json:"building"`
	Project        pgtype.Text   `json:"project"`
	Area           pgtype.Float4 `json:"area"`
	NumberOfFloors pgtype.Int4   `json:"number_of_floors"`
	YearBuilt      pgtype.Int4   `json:"year_built"`
	Orientation    pgtype.Text   `json:"orientation"`
	EntranceWidth  pgtype.Float4 `json:"entrance_width"`
	Facade         pgtype.Float4 `json:"facade"`
	FullAddress    pgtype.Text   `json:"full_address"`
	District       pgtype.Text   `json:"district"`
	City           pgtype.Text   `json:"city"`
	Ward           pgtype.Text   `json:"ward"`
	Lat            pgtype.Float8 `json:"lat"`
	Lng            pgtype.Float8 `json:"lng"`
	PlaceUrl       pgtype.Text   `json:"place_url"`
	Description    pgtype.Text   `json:"description"`
}

func (q *Queries) UpdateProperty(ctx context.Context, arg UpdatePropertyParams) error {
	_, err := q.db.Exec(ctx, updateProperty,
		arg.ID,
		arg.Name,
		arg.Building,
		arg.Project,
		arg.Area,
		arg.NumberOfFloors,
		arg.YearBuilt,
		arg.Orientation,
		arg.EntranceWidth,
		arg.Facade,
		arg.FullAddress,
		arg.District,
		arg.City,
		arg.Ward,
		arg.Lat,
		arg.Lng,
		arg.PlaceUrl,
		arg.Description,
	)
	return err
}

const updatePropertyManager = `-- name: UpdatePropertyManager :exec
UPDATE property_managers SET
  role = $3
WHERE property_id = $1 AND manager_id = $2
`

type UpdatePropertyManagerParams struct {
	PropertyID uuid.UUID `json:"property_id"`
	ManagerID  uuid.UUID `json:"manager_id"`
	Role       string    `json:"role"`
}

func (q *Queries) UpdatePropertyManager(ctx context.Context, arg UpdatePropertyManagerParams) error {
	_, err := q.db.Exec(ctx, updatePropertyManager, arg.PropertyID, arg.ManagerID, arg.Role)
	return err
}
