// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.22.0
// source: application.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

const createApplication = `-- name: CreateApplication :one
INSERT INTO applications (
  creator_id,
  listing_id,
  property_id,
  unit_ids,
  -- basic info
  full_name,
  dob,
  email,
  phone,
  profile_image,
  movein_date,
  preferred_term,
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
  employment_comment,
  employment_proofs_of_income,
  -- identity
  identity_type,
  identity_number,
  identity_issued_date,
  identity_issued_by
) VALUES (
  $1,
  $2,
  $3,
  $4,
  -- basic info
  $5,
  $6,
  $7,
  $8,
  $9,
  $10,
  $11,
  -- rental history
  $12,
  $13,
  $14,
  $15,
  $16,
  $17,
  $18,
  -- employment
  $19,
  $20,
  $21,
  $22,
  $23,
  $24,
  -- identity
  $25,
  $26,
  $27,
  $28
) RETURNING id, creator_id, listing_id, property_id, unit_ids, status, created_at, updated_at, full_name, email, phone, dob, profile_image, movein_date, preferred_term, rh_address, rh_city, rh_district, rh_ward, rh_rental_duration, rh_monthly_payment, rh_reason_for_leaving, employment_status, employment_company_name, employment_position, employment_monthly_income, employment_comment, employment_proofs_of_income, identity_type, identity_number, identity_issued_date, identity_issued_by
`

type CreateApplicationParams struct {
	CreatorID                uuid.UUID       `json:"creator_id"`
	ListingID                uuid.UUID       `json:"listing_id"`
	PropertyID               uuid.UUID       `json:"property_id"`
	UnitIds                  []uuid.UUID     `json:"unit_ids"`
	FullName                 string          `json:"full_name"`
	Dob                      time.Time       `json:"dob"`
	Email                    string          `json:"email"`
	Phone                    string          `json:"phone"`
	ProfileImage             string          `json:"profile_image"`
	MoveinDate               time.Time       `json:"movein_date"`
	PreferredTerm            int32           `json:"preferred_term"`
	RhAddress                sql.NullString  `json:"rh_address"`
	RhCity                   sql.NullString  `json:"rh_city"`
	RhDistrict               sql.NullString  `json:"rh_district"`
	RhWard                   sql.NullString  `json:"rh_ward"`
	RhRentalDuration         sql.NullInt32   `json:"rh_rental_duration"`
	RhMonthlyPayment         sql.NullFloat64 `json:"rh_monthly_payment"`
	RhReasonForLeaving       sql.NullString  `json:"rh_reason_for_leaving"`
	EmploymentStatus         string          `json:"employment_status"`
	EmploymentCompanyName    sql.NullString  `json:"employment_company_name"`
	EmploymentPosition       sql.NullString  `json:"employment_position"`
	EmploymentMonthlyIncome  sql.NullFloat64 `json:"employment_monthly_income"`
	EmploymentComment        sql.NullString  `json:"employment_comment"`
	EmploymentProofsOfIncome []string        `json:"employment_proofs_of_income"`
	IdentityType             string          `json:"identity_type"`
	IdentityNumber           string          `json:"identity_number"`
	IdentityIssuedDate       time.Time       `json:"identity_issued_date"`
	IdentityIssuedBy         string          `json:"identity_issued_by"`
}

func (q *Queries) CreateApplication(ctx context.Context, arg CreateApplicationParams) (Application, error) {
	row := q.db.QueryRowContext(ctx, createApplication,
		arg.CreatorID,
		arg.ListingID,
		arg.PropertyID,
		pq.Array(arg.UnitIds),
		arg.FullName,
		arg.Dob,
		arg.Email,
		arg.Phone,
		arg.ProfileImage,
		arg.MoveinDate,
		arg.PreferredTerm,
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
		pq.Array(arg.EmploymentProofsOfIncome),
		arg.IdentityType,
		arg.IdentityNumber,
		arg.IdentityIssuedDate,
		arg.IdentityIssuedBy,
	)
	var i Application
	err := row.Scan(
		&i.ID,
		&i.CreatorID,
		&i.ListingID,
		&i.PropertyID,
		pq.Array(&i.UnitIds),
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.FullName,
		&i.Email,
		&i.Phone,
		&i.Dob,
		&i.ProfileImage,
		&i.MoveinDate,
		&i.PreferredTerm,
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
		pq.Array(&i.EmploymentProofsOfIncome),
		&i.IdentityType,
		&i.IdentityNumber,
		&i.IdentityIssuedDate,
		&i.IdentityIssuedBy,
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
	ApplicationID int64          `json:"application_id"`
	FullName      string         `json:"full_name"`
	Dob           time.Time      `json:"dob"`
	Job           string         `json:"job"`
	Income        int32          `json:"income"`
	Email         sql.NullString `json:"email"`
	Phone         sql.NullString `json:"phone"`
	Description   sql.NullString `json:"description"`
}

func (q *Queries) CreateApplicationCoap(ctx context.Context, arg CreateApplicationCoapParams) (ApplicationCoap, error) {
	row := q.db.QueryRowContext(ctx, createApplicationCoap,
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
	ApplicationID int64          `json:"application_id"`
	FullName      string         `json:"full_name"`
	Dob           time.Time      `json:"dob"`
	Email         sql.NullString `json:"email"`
	Phone         sql.NullString `json:"phone"`
	Description   sql.NullString `json:"description"`
}

func (q *Queries) CreateApplicationMinor(ctx context.Context, arg CreateApplicationMinorParams) (ApplicationMinor, error) {
	row := q.db.QueryRowContext(ctx, createApplicationMinor,
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
	ApplicationID int64           `json:"application_id"`
	Type          string          `json:"type"`
	Weight        sql.NullFloat64 `json:"weight"`
	Description   sql.NullString  `json:"description"`
}

func (q *Queries) CreateApplicationPet(ctx context.Context, arg CreateApplicationPetParams) (ApplicationPet, error) {
	row := q.db.QueryRowContext(ctx, createApplicationPet,
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
	ApplicationID int64          `json:"application_id"`
	Type          string         `json:"type"`
	Model         sql.NullString `json:"model"`
	Code          string         `json:"code"`
	Description   sql.NullString `json:"description"`
}

func (q *Queries) CreateApplicationVehicle(ctx context.Context, arg CreateApplicationVehicleParams) (ApplicationVehicle, error) {
	row := q.db.QueryRowContext(ctx, createApplicationVehicle,
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
	_, err := q.db.ExecContext(ctx, deleteApplication, id)
	return err
}

const getApplicationByID = `-- name: GetApplicationByID :one
SELECT id, creator_id, listing_id, property_id, unit_ids, status, created_at, updated_at, full_name, email, phone, dob, profile_image, movein_date, preferred_term, rh_address, rh_city, rh_district, rh_ward, rh_rental_duration, rh_monthly_payment, rh_reason_for_leaving, employment_status, employment_company_name, employment_position, employment_monthly_income, employment_comment, employment_proofs_of_income, identity_type, identity_number, identity_issued_date, identity_issued_by FROM applications WHERE id = $1 LIMIT 1
`

func (q *Queries) GetApplicationByID(ctx context.Context, id int64) (Application, error) {
	row := q.db.QueryRowContext(ctx, getApplicationByID, id)
	var i Application
	err := row.Scan(
		&i.ID,
		&i.CreatorID,
		&i.ListingID,
		&i.PropertyID,
		pq.Array(&i.UnitIds),
		&i.Status,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.FullName,
		&i.Email,
		&i.Phone,
		&i.Dob,
		&i.ProfileImage,
		&i.MoveinDate,
		&i.PreferredTerm,
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
		pq.Array(&i.EmploymentProofsOfIncome),
		&i.IdentityType,
		&i.IdentityNumber,
		&i.IdentityIssuedDate,
		&i.IdentityIssuedBy,
	)
	return i, err
}

const getApplicationCoaps = `-- name: GetApplicationCoaps :many
SELECT application_id, full_name, dob, job, income, email, phone, description FROM application_coaps WHERE application_id = $1
`

func (q *Queries) GetApplicationCoaps(ctx context.Context, applicationID int64) ([]ApplicationCoap, error) {
	rows, err := q.db.QueryContext(ctx, getApplicationCoaps, applicationID)
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
	if err := rows.Close(); err != nil {
		return nil, err
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
	rows, err := q.db.QueryContext(ctx, getApplicationMinors, applicationID)
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
	if err := rows.Close(); err != nil {
		return nil, err
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
	rows, err := q.db.QueryContext(ctx, getApplicationPets, applicationID)
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
	if err := rows.Close(); err != nil {
		return nil, err
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
	rows, err := q.db.QueryContext(ctx, getApplicationVehicles, applicationID)
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
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateApplicationStatus = `-- name: UpdateApplicationStatus :exec
UPDATE applications SET status = $1, updated_at = NOW() WHERE id = $2
`

type UpdateApplicationStatusParams struct {
	Status APPLICATIONSTATUS `json:"status"`
	ID     int64             `json:"id"`
}

func (q *Queries) UpdateApplicationStatus(ctx context.Context, arg UpdateApplicationStatusParams) error {
	_, err := q.db.ExecContext(ctx, updateApplicationStatus, arg.Status, arg.ID)
	return err
}
