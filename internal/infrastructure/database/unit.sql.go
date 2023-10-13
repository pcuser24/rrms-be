// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: unit.sql

package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const checkUnitManageability = `-- name: CheckUnitManageability :one
SELECT count(*) FROM units WHERE units.id = $1 AND units.property_id IN (SELECT property_id FROM property_managers WHERE property_managers.property_id = units.property_id AND manager_id=$2 LIMIT 1) LIMIT 1
`

type CheckUnitManageabilityParams struct {
	ID        uuid.UUID `json:"id"`
	ManagerID uuid.UUID `json:"manager_id"`
}

func (q *Queries) CheckUnitManageability(ctx context.Context, arg CheckUnitManageabilityParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, checkUnitManageability, arg.ID, arg.ManagerID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const checkUnitOfProperty = `-- name: CheckUnitOfProperty :one
SELECT count(*) FROM units WHERE id = $1 AND property_id = $2 LIMIT 1
`

type CheckUnitOfPropertyParams struct {
	ID         uuid.UUID `json:"id"`
	PropertyID uuid.UUID `json:"property_id"`
}

func (q *Queries) CheckUnitOfProperty(ctx context.Context, arg CheckUnitOfPropertyParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, checkUnitOfProperty, arg.ID, arg.PropertyID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createUnit = `-- name: CreateUnit :one
INSERT INTO units (
  property_id,
  name,
  area,
  floor,
  number_of_living_rooms,
  number_of_bedrooms,
  number_of_bathrooms,
  number_of_toilets,
  number_of_kitchens,
  number_of_balconies,
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
  NOW(),
  NOW()
) RETURNING id, property_id, name, area, floor, number_of_living_rooms, number_of_bedrooms, number_of_bathrooms, number_of_toilets, number_of_balconies, number_of_kitchens, type, created_at, updated_at
`

type CreateUnitParams struct {
	PropertyID          uuid.UUID      `json:"property_id"`
	Name                sql.NullString `json:"name"`
	Area                float32        `json:"area"`
	Floor               sql.NullInt32  `json:"floor"`
	NumberOfLivingRooms sql.NullInt32  `json:"number_of_living_rooms"`
	NumberOfBedrooms    sql.NullInt32  `json:"number_of_bedrooms"`
	NumberOfBathrooms   sql.NullInt32  `json:"number_of_bathrooms"`
	NumberOfToilets     sql.NullInt32  `json:"number_of_toilets"`
	NumberOfKitchens    sql.NullInt32  `json:"number_of_kitchens"`
	NumberOfBalconies   sql.NullInt32  `json:"number_of_balconies"`
	Type                UNITTYPE       `json:"type"`
}

func (q *Queries) CreateUnit(ctx context.Context, arg CreateUnitParams) (Unit, error) {
	row := q.db.QueryRowContext(ctx, createUnit,
		arg.PropertyID,
		arg.Name,
		arg.Area,
		arg.Floor,
		arg.NumberOfLivingRooms,
		arg.NumberOfBedrooms,
		arg.NumberOfBathrooms,
		arg.NumberOfToilets,
		arg.NumberOfKitchens,
		arg.NumberOfBalconies,
		arg.Type,
	)
	var i Unit
	err := row.Scan(
		&i.ID,
		&i.PropertyID,
		&i.Name,
		&i.Area,
		&i.Floor,
		&i.NumberOfLivingRooms,
		&i.NumberOfBedrooms,
		&i.NumberOfBathrooms,
		&i.NumberOfToilets,
		&i.NumberOfBalconies,
		&i.NumberOfKitchens,
		&i.Type,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const deleteUnit = `-- name: DeleteUnit :exec
DELETE FROM units WHERE id = $1
`

func (q *Queries) DeleteUnit(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteUnit, id)
	return err
}

const deleteUnitAmenity = `-- name: DeleteUnitAmenity :exec
DELETE FROM unit_amenities WHERE unit_id = $1 AND amenity_id = $2
`

type DeleteUnitAmenityParams struct {
	UnitID    uuid.UUID `json:"unit_id"`
	AmenityID int64     `json:"amenity_id"`
}

func (q *Queries) DeleteUnitAmenity(ctx context.Context, arg DeleteUnitAmenityParams) error {
	_, err := q.db.ExecContext(ctx, deleteUnitAmenity, arg.UnitID, arg.AmenityID)
	return err
}

const deleteUnitMedia = `-- name: DeleteUnitMedia :exec
DELETE FROM unit_media WHERE unit_id = $1 AND id = $2
`

type DeleteUnitMediaParams struct {
	UnitID uuid.UUID `json:"unit_id"`
	ID     int64     `json:"id"`
}

func (q *Queries) DeleteUnitMedia(ctx context.Context, arg DeleteUnitMediaParams) error {
	_, err := q.db.ExecContext(ctx, deleteUnitMedia, arg.UnitID, arg.ID)
	return err
}

const getAllUnitAmenities = `-- name: GetAllUnitAmenities :many
SELECT id, amenity FROM u_amenities
`

func (q *Queries) GetAllUnitAmenities(ctx context.Context) ([]UAmenity, error) {
	rows, err := q.db.QueryContext(ctx, getAllUnitAmenities)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UAmenity
	for rows.Next() {
		var i UAmenity
		if err := rows.Scan(&i.ID, &i.Amenity); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUnitAmenities = `-- name: GetUnitAmenities :many
SELECT unit_id, amenity_id, description FROM unit_amenities WHERE unit_id = $1
`

func (q *Queries) GetUnitAmenities(ctx context.Context, unitID uuid.UUID) ([]UnitAmenity, error) {
	rows, err := q.db.QueryContext(ctx, getUnitAmenities, unitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UnitAmenity
	for rows.Next() {
		var i UnitAmenity
		if err := rows.Scan(&i.UnitID, &i.AmenityID, &i.Description); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUnitById = `-- name: GetUnitById :one
SELECT id, property_id, name, area, floor, number_of_living_rooms, number_of_bedrooms, number_of_bathrooms, number_of_toilets, number_of_balconies, number_of_kitchens, type, created_at, updated_at FROM units WHERE id = $1 LIMIT 1
`

func (q *Queries) GetUnitById(ctx context.Context, id uuid.UUID) (Unit, error) {
	row := q.db.QueryRowContext(ctx, getUnitById, id)
	var i Unit
	err := row.Scan(
		&i.ID,
		&i.PropertyID,
		&i.Name,
		&i.Area,
		&i.Floor,
		&i.NumberOfLivingRooms,
		&i.NumberOfBedrooms,
		&i.NumberOfBathrooms,
		&i.NumberOfToilets,
		&i.NumberOfBalconies,
		&i.NumberOfKitchens,
		&i.Type,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUnitManagers = `-- name: GetUnitManagers :many
SELECT property_id, manager_id, role FROM property_managers WHERE property_id IN (SELECT property_id FROM units WHERE units.id = $1 LIMIT 1)
`

func (q *Queries) GetUnitManagers(ctx context.Context, id uuid.UUID) ([]PropertyManager, error) {
	rows, err := q.db.QueryContext(ctx, getUnitManagers, id)
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
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUnitMedia = `-- name: GetUnitMedia :many
SELECT id, unit_id, url, type, description FROM unit_media WHERE unit_id = $1
`

func (q *Queries) GetUnitMedia(ctx context.Context, unitID uuid.UUID) ([]UnitMedia, error) {
	rows, err := q.db.QueryContext(ctx, getUnitMedia, unitID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []UnitMedia
	for rows.Next() {
		var i UnitMedia
		if err := rows.Scan(
			&i.ID,
			&i.UnitID,
			&i.Url,
			&i.Type,
			&i.Description,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getUnitsOfProperty = `-- name: GetUnitsOfProperty :many
SELECT id, property_id, name, area, floor, number_of_living_rooms, number_of_bedrooms, number_of_bathrooms, number_of_toilets, number_of_balconies, number_of_kitchens, type, created_at, updated_at FROM units WHERE property_id = $1
`

func (q *Queries) GetUnitsOfProperty(ctx context.Context, propertyID uuid.UUID) ([]Unit, error) {
	rows, err := q.db.QueryContext(ctx, getUnitsOfProperty, propertyID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Unit
	for rows.Next() {
		var i Unit
		if err := rows.Scan(
			&i.ID,
			&i.PropertyID,
			&i.Name,
			&i.Area,
			&i.Floor,
			&i.NumberOfLivingRooms,
			&i.NumberOfBedrooms,
			&i.NumberOfBathrooms,
			&i.NumberOfToilets,
			&i.NumberOfBalconies,
			&i.NumberOfKitchens,
			&i.Type,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const insertUnitAmenity = `-- name: InsertUnitAmenity :one
INSERT INTO unit_amenities (
  unit_id,
  amenity_id,
  description
) VALUES (
  $1,
  $2,
  $3
) RETURNING unit_id, amenity_id, description
`

type InsertUnitAmenityParams struct {
	UnitID      uuid.UUID      `json:"unit_id"`
	AmenityID   int64          `json:"amenity_id"`
	Description sql.NullString `json:"description"`
}

func (q *Queries) InsertUnitAmenity(ctx context.Context, arg InsertUnitAmenityParams) (UnitAmenity, error) {
	row := q.db.QueryRowContext(ctx, insertUnitAmenity, arg.UnitID, arg.AmenityID, arg.Description)
	var i UnitAmenity
	err := row.Scan(&i.UnitID, &i.AmenityID, &i.Description)
	return i, err
}

const insertUnitMedia = `-- name: InsertUnitMedia :one
INSERT INTO unit_media (
  unit_id,
  url,
  type
) VALUES (
  $1,
  $2,
  $3
) RETURNING id, unit_id, url, type, description
`

type InsertUnitMediaParams struct {
	UnitID uuid.UUID `json:"unit_id"`
	Url    string    `json:"url"`
	Type   MEDIATYPE `json:"type"`
}

func (q *Queries) InsertUnitMedia(ctx context.Context, arg InsertUnitMediaParams) (UnitMedia, error) {
	row := q.db.QueryRowContext(ctx, insertUnitMedia, arg.UnitID, arg.Url, arg.Type)
	var i UnitMedia
	err := row.Scan(
		&i.ID,
		&i.UnitID,
		&i.Url,
		&i.Type,
		&i.Description,
	)
	return i, err
}

const isUnitPublic = `-- name: IsUnitPublic :one
SELECT is_public FROM properties WHERE properties.id IN (SELECT property_id from units WHERE units.id = $1 LIMIT 1) LIMIT 1
`

func (q *Queries) IsUnitPublic(ctx context.Context, id uuid.UUID) (bool, error) {
	row := q.db.QueryRowContext(ctx, isUnitPublic, id)
	var is_public bool
	err := row.Scan(&is_public)
	return is_public, err
}

const updateUnit = `-- name: UpdateUnit :exec
UPDATE units SET
  name = coalesce($2, name),
  area = coalesce($3, area),
  floor = coalesce($4, floor),
  number_of_living_rooms = coalesce($5, number_of_living_rooms),
  number_of_bedrooms = coalesce($6, number_of_bedrooms),
  number_of_bathrooms = coalesce($7, number_of_bathrooms),
  number_of_toilets = coalesce($8, number_of_toilets),
  number_of_kitchens = coalesce($9, number_of_kitchens),
  number_of_balconies = coalesce($10, number_of_balconies),
  updated_at = NOW()
WHERE id = $1
`

type UpdateUnitParams struct {
	ID                  uuid.UUID       `json:"id"`
	Name                sql.NullString  `json:"name"`
	Area                sql.NullFloat64 `json:"area"`
	Floor               sql.NullInt32   `json:"floor"`
	NumberOfLivingRooms sql.NullInt32   `json:"number_of_living_rooms"`
	NumberOfBedrooms    sql.NullInt32   `json:"number_of_bedrooms"`
	NumberOfBathrooms   sql.NullInt32   `json:"number_of_bathrooms"`
	NumberOfToilets     sql.NullInt32   `json:"number_of_toilets"`
	NumberOfKitchens    sql.NullInt32   `json:"number_of_kitchens"`
	NumberOfBalconies   sql.NullInt32   `json:"number_of_balconies"`
}

func (q *Queries) UpdateUnit(ctx context.Context, arg UpdateUnitParams) error {
	_, err := q.db.ExecContext(ctx, updateUnit,
		arg.ID,
		arg.Name,
		arg.Area,
		arg.Floor,
		arg.NumberOfLivingRooms,
		arg.NumberOfBedrooms,
		arg.NumberOfBathrooms,
		arg.NumberOfToilets,
		arg.NumberOfKitchens,
		arg.NumberOfBalconies,
	)
	return err
}