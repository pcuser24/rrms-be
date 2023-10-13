// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: listing.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const checkListingOwnership = `-- name: CheckListingOwnership :one
SELECT count(*) FROM listings WHERE id = $1 AND creator_id = $2 LIMIT 1
`

type CheckListingOwnershipParams struct {
	ID        uuid.UUID `json:"id"`
	CreatorID uuid.UUID `json:"creator_id"`
}

func (q *Queries) CheckListingOwnership(ctx context.Context, arg CheckListingOwnershipParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, checkListingOwnership, arg.ID, arg.CreatorID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const checkValidUnitForListing = `-- name: CheckValidUnitForListing :one
SELECT count(*) FROM units WHERE units.id = $1 AND units.property_id IN (SELECT listings.property_id FROM listings WHERE listings.id = $2) LIMIT 1
`

type CheckValidUnitForListingParams struct {
	ID   uuid.UUID `json:"id"`
	ID_2 uuid.UUID `json:"id_2"`
}

func (q *Queries) CheckValidUnitForListing(ctx context.Context, arg CheckValidUnitForListingParams) (int64, error) {
	row := q.db.QueryRowContext(ctx, checkValidUnitForListing, arg.ID, arg.ID_2)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const createListing = `-- name: CreateListing :one
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
  NOW(), NOW(), 
  $16,
  $17,
  $18
) RETURNING id, creator_id, property_id, title, description, full_name, email, phone, contact_type, price, price_negotiable, security_deposit, lease_term, pets_allowed, number_of_residents, priority, active, created_at, updated_at, post_at, expired_at
`

type CreateListingParams struct {
	CreatorID         uuid.UUID     `json:"creator_id"`
	PropertyID        uuid.UUID     `json:"property_id"`
	Title             string        `json:"title"`
	Description       string        `json:"description"`
	FullName          string        `json:"full_name"`
	Email             string        `json:"email"`
	Phone             string        `json:"phone"`
	ContactType       string        `json:"contact_type"`
	Price             int64         `json:"price"`
	PriceNegotiable   sql.NullBool  `json:"price_negotiable"`
	SecurityDeposit   sql.NullInt64 `json:"security_deposit"`
	LeaseTerm         int32         `json:"lease_term"`
	PetsAllowed       sql.NullBool  `json:"pets_allowed"`
	NumberOfResidents sql.NullInt32 `json:"number_of_residents"`
	Priority          int32         `json:"priority"`
	PostAt            time.Time     `json:"post_at"`
	Active            bool          `json:"active"`
	ExpiredAt         time.Time     `json:"expired_at"`
}

func (q *Queries) CreateListing(ctx context.Context, arg CreateListingParams) (Listing, error) {
	row := q.db.QueryRowContext(ctx, createListing,
		arg.CreatorID,
		arg.PropertyID,
		arg.Title,
		arg.Description,
		arg.FullName,
		arg.Email,
		arg.Phone,
		arg.ContactType,
		arg.Price,
		arg.PriceNegotiable,
		arg.SecurityDeposit,
		arg.LeaseTerm,
		arg.PetsAllowed,
		arg.NumberOfResidents,
		arg.Priority,
		arg.PostAt,
		arg.Active,
		arg.ExpiredAt,
	)
	var i Listing
	err := row.Scan(
		&i.ID,
		&i.CreatorID,
		&i.PropertyID,
		&i.Title,
		&i.Description,
		&i.FullName,
		&i.Email,
		&i.Phone,
		&i.ContactType,
		&i.Price,
		&i.PriceNegotiable,
		&i.SecurityDeposit,
		&i.LeaseTerm,
		&i.PetsAllowed,
		&i.NumberOfResidents,
		&i.Priority,
		&i.Active,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PostAt,
		&i.ExpiredAt,
	)
	return i, err
}

const createListingPolicy = `-- name: CreateListingPolicy :one
INSERT INTO listing_policy (
  listing_id,
  policy_id,
  note
) VALUES (
  $1,
  $2,
  $3
) RETURNING listing_id, policy_id, note
`

type CreateListingPolicyParams struct {
	ListingID uuid.UUID      `json:"listing_id"`
	PolicyID  int64          `json:"policy_id"`
	Note      sql.NullString `json:"note"`
}

func (q *Queries) CreateListingPolicy(ctx context.Context, arg CreateListingPolicyParams) (ListingPolicy, error) {
	row := q.db.QueryRowContext(ctx, createListingPolicy, arg.ListingID, arg.PolicyID, arg.Note)
	var i ListingPolicy
	err := row.Scan(&i.ListingID, &i.PolicyID, &i.Note)
	return i, err
}

const createListingUnit = `-- name: CreateListingUnit :one
INSERT INTO listing_unit (
  listing_id,
  unit_id
) VALUES (
  $1,
  $2
) RETURNING listing_id, unit_id
`

type CreateListingUnitParams struct {
	ListingID uuid.UUID `json:"listing_id"`
	UnitID    uuid.UUID `json:"unit_id"`
}

func (q *Queries) CreateListingUnit(ctx context.Context, arg CreateListingUnitParams) (ListingUnit, error) {
	row := q.db.QueryRowContext(ctx, createListingUnit, arg.ListingID, arg.UnitID)
	var i ListingUnit
	err := row.Scan(&i.ListingID, &i.UnitID)
	return i, err
}

const deleteListing = `-- name: DeleteListing :exec
DELETE FROM listings WHERE id = $1
`

func (q *Queries) DeleteListing(ctx context.Context, id uuid.UUID) error {
	_, err := q.db.ExecContext(ctx, deleteListing, id)
	return err
}

const getAllRentalPolicies = `-- name: GetAllRentalPolicies :many
SELECT id, policy FROM rental_policies
`

func (q *Queries) GetAllRentalPolicies(ctx context.Context) ([]RentalPolicy, error) {
	rows, err := q.db.QueryContext(ctx, getAllRentalPolicies)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []RentalPolicy
	for rows.Next() {
		var i RentalPolicy
		if err := rows.Scan(&i.ID, &i.Policy); err != nil {
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

const getListingByID = `-- name: GetListingByID :one
SELECT id, creator_id, property_id, title, description, full_name, email, phone, contact_type, price, price_negotiable, security_deposit, lease_term, pets_allowed, number_of_residents, priority, active, created_at, updated_at, post_at, expired_at FROM listings WHERE id = $1 LIMIT 1
`

func (q *Queries) GetListingByID(ctx context.Context, id uuid.UUID) (Listing, error) {
	row := q.db.QueryRowContext(ctx, getListingByID, id)
	var i Listing
	err := row.Scan(
		&i.ID,
		&i.CreatorID,
		&i.PropertyID,
		&i.Title,
		&i.Description,
		&i.FullName,
		&i.Email,
		&i.Phone,
		&i.ContactType,
		&i.Price,
		&i.PriceNegotiable,
		&i.SecurityDeposit,
		&i.LeaseTerm,
		&i.PetsAllowed,
		&i.NumberOfResidents,
		&i.Priority,
		&i.Active,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.PostAt,
		&i.ExpiredAt,
	)
	return i, err
}

const getListingPolicies = `-- name: GetListingPolicies :many
SELECT listing_id, policy_id, note FROM listing_policy WHERE listing_id = $1
`

func (q *Queries) GetListingPolicies(ctx context.Context, listingID uuid.UUID) ([]ListingPolicy, error) {
	rows, err := q.db.QueryContext(ctx, getListingPolicies, listingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListingPolicy
	for rows.Next() {
		var i ListingPolicy
		if err := rows.Scan(&i.ListingID, &i.PolicyID, &i.Note); err != nil {
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

const getListingUnits = `-- name: GetListingUnits :many
SELECT listing_id, unit_id FROM listing_unit WHERE listing_id = $1
`

func (q *Queries) GetListingUnits(ctx context.Context, listingID uuid.UUID) ([]ListingUnit, error) {
	rows, err := q.db.QueryContext(ctx, getListingUnits, listingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListingUnit
	for rows.Next() {
		var i ListingUnit
		if err := rows.Scan(&i.ListingID, &i.UnitID); err != nil {
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

const updateListing = `-- name: UpdateListing :exec
UPDATE listings SET
  title = coalesce($1, title),
  description = coalesce($2, description),
  full_name = coalesce($3, full_name),
  email = coalesce($4, email),
  phone = coalesce($5, phone),
  contact_type = coalesce($6, contact_type),
  price = coalesce($7, price),
  price_negotiable = coalesce($8, price_negotiable),
  security_deposit = coalesce($9, security_deposit),
  lease_term = coalesce($10, lease_term),
  pets_allowed = coalesce($11, pets_allowed),
  number_of_residents = coalesce($12, number_of_residents),
  updated_at = NOW(),
  post_at = coalesce($13, post_at)
WHERE id = $14
`

type UpdateListingParams struct {
	Title             sql.NullString `json:"title"`
	Description       sql.NullString `json:"description"`
	FullName          sql.NullString `json:"full_name"`
	Email             sql.NullString `json:"email"`
	Phone             sql.NullString `json:"phone"`
	ContactType       sql.NullString `json:"contact_type"`
	Price             sql.NullInt64  `json:"price"`
	PriceNegotiable   sql.NullBool   `json:"price_negotiable"`
	SecurityDeposit   sql.NullInt64  `json:"security_deposit"`
	LeaseTerm         sql.NullInt32  `json:"lease_term"`
	PetsAllowed       sql.NullBool   `json:"pets_allowed"`
	NumberOfResidents sql.NullInt32  `json:"number_of_residents"`
	PostAt            sql.NullTime   `json:"post_at"`
	ID                uuid.UUID      `json:"id"`
}

func (q *Queries) UpdateListing(ctx context.Context, arg UpdateListingParams) error {
	_, err := q.db.ExecContext(ctx, updateListing,
		arg.Title,
		arg.Description,
		arg.FullName,
		arg.Email,
		arg.Phone,
		arg.ContactType,
		arg.Price,
		arg.PriceNegotiable,
		arg.SecurityDeposit,
		arg.LeaseTerm,
		arg.PetsAllowed,
		arg.NumberOfResidents,
		arg.PostAt,
		arg.ID,
	)
	return err
}

const updateListingStatus = `-- name: UpdateListingStatus :exec
UPDATE listings SET active = $1 WHERE id = $2
`

type UpdateListingStatusParams struct {
	Active bool      `json:"active"`
	ID     uuid.UUID `json:"id"`
}

func (q *Queries) UpdateListingStatus(ctx context.Context, arg UpdateListingStatusParams) error {
	_, err := q.db.ExecContext(ctx, updateListingStatus, arg.Active, arg.ID)
	return err
}