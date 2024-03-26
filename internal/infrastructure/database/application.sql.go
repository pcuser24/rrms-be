// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: application.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const checkApplicationUpdatabilty = `-- name: CheckApplicationUpdatabilty :one
SELECT count(*) > 0 FROM applications WHERE 
  id = $1 
  AND (
    property_id IN (SELECT property_id FROM property_managers WHERE manager_id = $2)
  )
`

type CheckApplicationUpdatabiltyParams struct {
	ID        int64     `json:"id"`
	ManagerID uuid.UUID `json:"manager_id"`
}

func (q *Queries) CheckApplicationUpdatabilty(ctx context.Context, arg CheckApplicationUpdatabiltyParams) (bool, error) {
	row := q.db.QueryRow(ctx, checkApplicationUpdatabilty, arg.ID, arg.ManagerID)
	var column_1 bool
	err := row.Scan(&column_1)
	return column_1, err
}

const checkApplicationVisibility = `-- name: CheckApplicationVisibility :one
SELECT count(*) > 0 FROM applications WHERE 
  id = $1 
  AND (
    property_id IN (SELECT property_id FROM property_managers WHERE manager_id = $2)
    OR creator_id = $2
  )
`

type CheckApplicationVisibilityParams struct {
	ID        int64     `json:"id"`
	ManagerID uuid.UUID `json:"manager_id"`
}

func (q *Queries) CheckApplicationVisibility(ctx context.Context, arg CheckApplicationVisibilityParams) (bool, error) {
	row := q.db.QueryRow(ctx, checkApplicationVisibility, arg.ID, arg.ManagerID)
	var column_1 bool
	err := row.Scan(&column_1)
	return column_1, err
}

const createApplication = `-- name: CreateApplication :one
INSERT INTO applications (
  creator_id,
  listing_id,
  property_id,
  unit_id,
  listing_price,
  offered_price,
  -- basic info
  tenant_type,
  full_name,
  dob,
  email,
  phone,
  profile_image,
  movein_date,
  preferred_term,
  rental_intention,
  organization_name,
  organization_hq_address,
  organization_scale,
  -- rental history
  rh_address,
  rh_city,
  rh_district,
  rh_ward,
  rh_rental_duration,
  rh_monthly_payment,
  rh_reason_for_leaving,
  -- employment
  employment_status,
  employment_company_name,
  employment_position,
  employment_monthly_income,
  employment_comment
  -- identity
  -- identity_type,
  -- identity_number
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  -- basic info
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
  -- rental history
  $19,
  $20,
  $21,
  $22,
  $23,
  $24,
  $25,
  -- employment
  $26,
  $27,
  $28,
  $29,
  $30
  -- identity
  -- sqlc.arg(identity_type),
  -- sqlc.arg(identity_number)
) RETURNING id, creator_id, listing_id, property_id, unit_id, listing_price, offered_price, status, created_at, updated_at, tenant_type, full_name, email, phone, dob, profile_image, movein_date, preferred_term, rental_intention, organization_name, organization_hq_address, organization_scale, rh_address, rh_city, rh_district, rh_ward, rh_rental_duration, rh_monthly_payment, rh_reason_for_leaving, employment_status, employment_company_name, employment_position, employment_monthly_income, employment_comment
`

type CreateApplicationParams struct {
	CreatorID               pgtype.UUID `json:"creator_id"`
	ListingID               uuid.UUID   `json:"listing_id"`
	PropertyID              uuid.UUID   `json:"property_id"`
	UnitID                  uuid.UUID   `json:"unit_id"`
	ListingPrice            int64       `json:"listing_price"`
	OfferedPrice            int64       `json:"offered_price"`
	TenantType              TENANTTYPE  `json:"tenant_type"`
	FullName                string      `json:"full_name"`
	Dob                     pgtype.Date `json:"dob"`
	Email                   string      `json:"email"`
	Phone                   string      `json:"phone"`
	ProfileImage            string      `json:"profile_image"`
	MoveinDate              pgtype.Date `json:"movein_date"`
	PreferredTerm           int32       `json:"preferred_term"`
	RentalIntention         string      `json:"rental_intention"`
	OrganizationName        pgtype.Text `json:"organization_name"`
	OrganizationHqAddress   pgtype.Text `json:"organization_hq_address"`
	OrganizationScale       pgtype.Text `json:"organization_scale"`
	RhAddress               pgtype.Text `json:"rh_address"`
	RhCity                  pgtype.Text `json:"rh_city"`
	RhDistrict              pgtype.Text `json:"rh_district"`
	RhWard                  pgtype.Text `json:"rh_ward"`
	RhRentalDuration        pgtype.Int4 `json:"rh_rental_duration"`
	RhMonthlyPayment        pgtype.Int8 `json:"rh_monthly_payment"`
	RhReasonForLeaving      pgtype.Text `json:"rh_reason_for_leaving"`
	EmploymentStatus        string      `json:"employment_status"`
	EmploymentCompanyName   pgtype.Text `json:"employment_company_name"`
	EmploymentPosition      pgtype.Text `json:"employment_position"`
	EmploymentMonthlyIncome pgtype.Int8 `json:"employment_monthly_income"`
	EmploymentComment       pgtype.Text `json:"employment_comment"`
}

func (q *Queries) CreateApplication(ctx context.Context, arg CreateApplicationParams) (Application, error) {
	row := q.db.QueryRow(ctx, createApplication,
		arg.CreatorID,
		arg.ListingID,
		arg.PropertyID,
		arg.UnitID,
		arg.ListingPrice,
		arg.OfferedPrice,
		arg.TenantType,
		arg.FullName,
		arg.Dob,
		arg.Email,
		arg.Phone,
		arg.ProfileImage,
		arg.MoveinDate,
		arg.PreferredTerm,
		arg.RentalIntention,
		arg.OrganizationName,
		arg.OrganizationHqAddress,
		arg.OrganizationScale,
		arg.RhAddress,
		arg.RhCity,
		arg.RhDistrict,
		arg.RhWard,
		arg.RhRentalDuration,
		arg.RhMonthlyPayment,
		arg.RhReasonForLeaving,
		arg.EmploymentStatus,
		arg.EmploymentCompanyName,
		arg.EmploymentPosition,
		arg.EmploymentMonthlyIncome,
		arg.EmploymentComment,
	)
	var i Application
	err := row.Scan(
		&i.ID,
		&i.CreatorID,
		&i.ListingID,
		&i.PropertyID,
		&i.UnitID,
		&i.ListingPrice,
		&i.OfferedPrice,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.TenantType,
		&i.FullName,
		&i.Email,
		&i.Phone,
		&i.Dob,
		&i.ProfileImage,
		&i.MoveinDate,
		&i.PreferredTerm,
		&i.RentalIntention,
		&i.OrganizationName,
		&i.OrganizationHqAddress,
		&i.OrganizationScale,
		&i.RhAddress,
		&i.RhCity,
		&i.RhDistrict,
		&i.RhWard,
		&i.RhRentalDuration,
		&i.RhMonthlyPayment,
		&i.RhReasonForLeaving,
		&i.EmploymentStatus,
		&i.EmploymentCompanyName,
		&i.EmploymentPosition,
		&i.EmploymentMonthlyIncome,
		&i.EmploymentComment,
	)
	return i, err
}

const createApplicationCoap = `-- name: CreateApplicationCoap :one
INSERT INTO application_coaps (
  application_id,
  full_name,
  dob,
  job,
  income,
  email,
  phone,
  description
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6,
  $7,
  $8
) RETURNING application_id, full_name, dob, job, income, email, phone, description
`

type CreateApplicationCoapParams struct {
	ApplicationID int64       `json:"application_id"`
	FullName      string      `json:"full_name"`
	Dob           time.Time   `json:"dob"`
	Job           string      `json:"job"`
	Income        int32       `json:"income"`
	Email         pgtype.Text `json:"email"`
	Phone         pgtype.Text `json:"phone"`
	Description   pgtype.Text `json:"description"`
}

func (q *Queries) CreateApplicationCoap(ctx context.Context, arg CreateApplicationCoapParams) (ApplicationCoap, error) {
	row := q.db.QueryRow(ctx, createApplicationCoap,
		arg.ApplicationID,
		arg.FullName,
		arg.Dob,
		arg.Job,
		arg.Income,
		arg.Email,
		arg.Phone,
		arg.Description,
	)
	var i ApplicationCoap
	err := row.Scan(
		&i.ApplicationID,
		&i.FullName,
		&i.Dob,
		&i.Job,
		&i.Income,
		&i.Email,
		&i.Phone,
		&i.Description,
	)
	return i, err
}

const createApplicationMinor = `-- name: CreateApplicationMinor :one
INSERT INTO application_minors (
  application_id,
  full_name,
  dob,
  email,
  phone,
  description
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5,
  $6
) RETURNING application_id, full_name, dob, email, phone, description
`

type CreateApplicationMinorParams struct {
	ApplicationID int64       `json:"application_id"`
	FullName      string      `json:"full_name"`
	Dob           time.Time   `json:"dob"`
	Email         pgtype.Text `json:"email"`
	Phone         pgtype.Text `json:"phone"`
	Description   pgtype.Text `json:"description"`
}

func (q *Queries) CreateApplicationMinor(ctx context.Context, arg CreateApplicationMinorParams) (ApplicationMinor, error) {
	row := q.db.QueryRow(ctx, createApplicationMinor,
		arg.ApplicationID,
		arg.FullName,
		arg.Dob,
		arg.Email,
		arg.Phone,
		arg.Description,
	)
	var i ApplicationMinor
	err := row.Scan(
		&i.ApplicationID,
		&i.FullName,
		&i.Dob,
		&i.Email,
		&i.Phone,
		&i.Description,
	)
	return i, err
}

const createApplicationPet = `-- name: CreateApplicationPet :one
INSERT INTO application_pets (
  application_id,
  type,
  weight,
  description
) VALUES (
  $1,
  $2,
  $3,
  $4
) RETURNING application_id, type, weight, description
`

type CreateApplicationPetParams struct {
	ApplicationID int64         `json:"application_id"`
	Type          string        `json:"type"`
	Weight        pgtype.Float4 `json:"weight"`
	Description   pgtype.Text   `json:"description"`
}

func (q *Queries) CreateApplicationPet(ctx context.Context, arg CreateApplicationPetParams) (ApplicationPet, error) {
	row := q.db.QueryRow(ctx, createApplicationPet,
		arg.ApplicationID,
		arg.Type,
		arg.Weight,
		arg.Description,
	)
	var i ApplicationPet
	err := row.Scan(
		&i.ApplicationID,
		&i.Type,
		&i.Weight,
		&i.Description,
	)
	return i, err
}

const createApplicationVehicle = `-- name: CreateApplicationVehicle :one
INSERT INTO application_vehicles (
  application_id,
  type,
  model,
  code,
  description
) VALUES (
  $1,
  $2,
  $3,
  $4,
  $5
) RETURNING application_id, type, model, code, description
`

type CreateApplicationVehicleParams struct {
	ApplicationID int64       `json:"application_id"`
	Type          string      `json:"type"`
	Model         pgtype.Text `json:"model"`
	Code          string      `json:"code"`
	Description   pgtype.Text `json:"description"`
}

func (q *Queries) CreateApplicationVehicle(ctx context.Context, arg CreateApplicationVehicleParams) (ApplicationVehicle, error) {
	row := q.db.QueryRow(ctx, createApplicationVehicle,
		arg.ApplicationID,
		arg.Type,
		arg.Model,
		arg.Code,
		arg.Description,
	)
	var i ApplicationVehicle
	err := row.Scan(
		&i.ApplicationID,
		&i.Type,
		&i.Model,
		&i.Code,
		&i.Description,
	)
	return i, err
}

const deleteApplication = `-- name: DeleteApplication :exec
DELETE FROM applications WHERE id = $1
`

func (q *Queries) DeleteApplication(ctx context.Context, id int64) error {
	_, err := q.db.Exec(ctx, deleteApplication, id)
	return err
}

const getApplicationByID = `-- name: GetApplicationByID :one
SELECT id, creator_id, listing_id, property_id, unit_id, listing_price, offered_price, status, created_at, updated_at, tenant_type, full_name, email, phone, dob, profile_image, movein_date, preferred_term, rental_intention, organization_name, organization_hq_address, organization_scale, rh_address, rh_city, rh_district, rh_ward, rh_rental_duration, rh_monthly_payment, rh_reason_for_leaving, employment_status, employment_company_name, employment_position, employment_monthly_income, employment_comment FROM applications WHERE id = $1 LIMIT 1
`

func (q *Queries) GetApplicationByID(ctx context.Context, id int64) (Application, error) {
	row := q.db.QueryRow(ctx, getApplicationByID, id)
	var i Application
	err := row.Scan(
		&i.ID,
		&i.CreatorID,
		&i.ListingID,
		&i.PropertyID,
		&i.UnitID,
		&i.ListingPrice,
		&i.OfferedPrice,
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.TenantType,
		&i.FullName,
		&i.Email,
		&i.Phone,
		&i.Dob,
		&i.ProfileImage,
		&i.MoveinDate,
		&i.PreferredTerm,
		&i.RentalIntention,
		&i.OrganizationName,
		&i.OrganizationHqAddress,
		&i.OrganizationScale,
		&i.RhAddress,
		&i.RhCity,
		&i.RhDistrict,
		&i.RhWard,
		&i.RhRentalDuration,
		&i.RhMonthlyPayment,
		&i.RhReasonForLeaving,
		&i.EmploymentStatus,
		&i.EmploymentCompanyName,
		&i.EmploymentPosition,
		&i.EmploymentMonthlyIncome,
		&i.EmploymentComment,
	)
	return i, err
}

const getApplicationCoaps = `-- name: GetApplicationCoaps :many
SELECT application_id, full_name, dob, job, income, email, phone, description FROM application_coaps WHERE application_id = $1
`

func (q *Queries) GetApplicationCoaps(ctx context.Context, applicationID int64) ([]ApplicationCoap, error) {
	rows, err := q.db.Query(ctx, getApplicationCoaps, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ApplicationCoap
	for rows.Next() {
		var i ApplicationCoap
		if err := rows.Scan(
			&i.ApplicationID,
			&i.FullName,
			&i.Dob,
			&i.Job,
			&i.Income,
			&i.Email,
			&i.Phone,
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

const getApplicationMinors = `-- name: GetApplicationMinors :many
SELECT application_id, full_name, dob, email, phone, description FROM application_minors WHERE application_id = $1
`

func (q *Queries) GetApplicationMinors(ctx context.Context, applicationID int64) ([]ApplicationMinor, error) {
	rows, err := q.db.Query(ctx, getApplicationMinors, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ApplicationMinor
	for rows.Next() {
		var i ApplicationMinor
		if err := rows.Scan(
			&i.ApplicationID,
			&i.FullName,
			&i.Dob,
			&i.Email,
			&i.Phone,
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

const getApplicationPets = `-- name: GetApplicationPets :many
SELECT application_id, type, weight, description FROM application_pets WHERE application_id = $1
`

func (q *Queries) GetApplicationPets(ctx context.Context, applicationID int64) ([]ApplicationPet, error) {
	rows, err := q.db.Query(ctx, getApplicationPets, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ApplicationPet
	for rows.Next() {
		var i ApplicationPet
		if err := rows.Scan(
			&i.ApplicationID,
			&i.Type,
			&i.Weight,
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

const getApplicationVehicles = `-- name: GetApplicationVehicles :many
SELECT application_id, type, model, code, description FROM application_vehicles WHERE application_id = $1
`

func (q *Queries) GetApplicationVehicles(ctx context.Context, applicationID int64) ([]ApplicationVehicle, error) {
	rows, err := q.db.Query(ctx, getApplicationVehicles, applicationID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ApplicationVehicle
	for rows.Next() {
		var i ApplicationVehicle
		if err := rows.Scan(
			&i.ApplicationID,
			&i.Type,
			&i.Model,
			&i.Code,
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

const getApplicationsByUserId = `-- name: GetApplicationsByUserId :many
SELECT 
  id 
FROM 
  applications 
WHERE 
  creator_id = $1 
  AND created_at >= $2
ORDER BY 
  created_at DESC 
LIMIT $3 OFFSET $4
`

type GetApplicationsByUserIdParams struct {
	CreatorID pgtype.UUID `json:"creator_id"`
	CreatedAt time.Time   `json:"created_at"`
	Limit     int32       `json:"limit"`
	Offset    int32       `json:"offset"`
}

func (q *Queries) GetApplicationsByUserId(ctx context.Context, arg GetApplicationsByUserIdParams) ([]int64, error) {
	rows, err := q.db.Query(ctx, getApplicationsByUserId,
		arg.CreatorID,
		arg.CreatedAt,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getApplicationsOfListing = `-- name: GetApplicationsOfListing :many
SELECT id FROM applications WHERE listing_id = $1
`

func (q *Queries) GetApplicationsOfListing(ctx context.Context, listingID uuid.UUID) ([]int64, error) {
	rows, err := q.db.Query(ctx, getApplicationsOfListing, listingID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getApplicationsToUser = `-- name: GetApplicationsToUser :many
SELECT 
  id 
FROM 
  applications 
WHERE 
  property_id IN (
    SELECT property_id FROM property_managers WHERE manager_id = $1
  ) AND created_at >= $2
ORDER BY
  created_at DESC
LIMIT $3 OFFSET $4
`

type GetApplicationsToUserParams struct {
	ManagerID uuid.UUID `json:"manager_id"`
	CreatedAt time.Time `json:"created_at"`
	Limit     int32     `json:"limit"`
	Offset    int32     `json:"offset"`
}

func (q *Queries) GetApplicationsToUser(ctx context.Context, arg GetApplicationsToUserParams) ([]int64, error) {
	rows, err := q.db.Query(ctx, getApplicationsToUser,
		arg.ManagerID,
		arg.CreatedAt,
		arg.Limit,
		arg.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateApplicationStatus = `-- name: UpdateApplicationStatus :many
UPDATE applications 
SET 
  status = $1, 
  updated_at = NOW() 
WHERE 
  id = $2
  AND property_id IN (SELECT property_id FROM property_managers WHERE manager_id = $3)
RETURNING id
`

type UpdateApplicationStatusParams struct {
	Status    APPLICATIONSTATUS `json:"status"`
	ID        int64             `json:"id"`
	ManagerID uuid.UUID         `json:"manager_id"`
}

func (q *Queries) UpdateApplicationStatus(ctx context.Context, arg UpdateApplicationStatusParams) ([]int64, error) {
	rows, err := q.db.Query(ctx, updateApplicationStatus, arg.Status, arg.ID, arg.ManagerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
