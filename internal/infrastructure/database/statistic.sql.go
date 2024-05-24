// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: statistic.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const getLeastRentedProperties = `-- name: GetLeastRentedProperties :many
SELECT r.id, COALESCE(c.count, 0) AS count
FROM 
  (SELECT properties.id FROM properties WHERE 
    EXISTS (
      SELECT 1 FROM property_managers WHERE property_managers.manager_id = $1 AND property_managers.property_id = properties.id 
    )
  ) AS r LEFT JOIN (
    SELECT property_id, COUNT(property_id) AS count FROM rentals WHERE 
    EXISTS (
      SELECT 1 FROM property_managers WHERE property_managers.manager_id = $1 AND property_managers.property_id = rentals.property_id 
    ) 
    GROUP BY property_id
  ) AS c ON r.id = c.property_id
ORDER BY count ASC
LIMIT $2
OFFSET $3
`

type GetLeastRentedPropertiesParams struct {
	ManagerID uuid.UUID `json:"manager_id"`
	Limit     int32     `json:"limit"`
	Offset    int32     `json:"offset"`
}

type GetLeastRentedPropertiesRow struct {
	ID    uuid.UUID `json:"id"`
	Count int64     `json:"count"`
}

func (q *Queries) GetLeastRentedProperties(ctx context.Context, arg GetLeastRentedPropertiesParams) ([]GetLeastRentedPropertiesRow, error) {
	rows, err := q.db.Query(ctx, getLeastRentedProperties, arg.ManagerID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetLeastRentedPropertiesRow
	for rows.Next() {
		var i GetLeastRentedPropertiesRow
		if err := rows.Scan(&i.ID, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getLeastRentedUnits = `-- name: GetLeastRentedUnits :many
SELECT r.id, r.property_id, COALESCE(c.count, 0) AS count
FROM 
  (SELECT units.id, units.property_id FROM units WHERE 
    EXISTS (
      SELECT 1 FROM property_managers WHERE property_managers.manager_id = $1 AND property_managers.property_id = units.property_id 
    )
  ) AS r LEFT JOIN (
    SELECT unit_id, COUNT(unit_id) AS count FROM rentals WHERE 
    EXISTS (
      SELECT 1 FROM property_managers WHERE property_managers.manager_id = $1 AND property_managers.property_id = rentals.property_id 
    ) 
    GROUP BY unit_id
  ) AS c ON r.id = c.unit_id
ORDER BY count ASC
LIMIT $2
OFFSET $3
`

type GetLeastRentedUnitsParams struct {
	ManagerID uuid.UUID `json:"manager_id"`
	Limit     int32     `json:"limit"`
	Offset    int32     `json:"offset"`
}

type GetLeastRentedUnitsRow struct {
	ID         uuid.UUID `json:"id"`
	PropertyID uuid.UUID `json:"property_id"`
	Count      int64     `json:"count"`
}

func (q *Queries) GetLeastRentedUnits(ctx context.Context, arg GetLeastRentedUnitsParams) ([]GetLeastRentedUnitsRow, error) {
	rows, err := q.db.Query(ctx, getLeastRentedUnits, arg.ManagerID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetLeastRentedUnitsRow
	for rows.Next() {
		var i GetLeastRentedUnitsRow
		if err := rows.Scan(&i.ID, &i.PropertyID, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMaintenanceRequests = `-- name: GetMaintenanceRequests :many
SELECT id FROM rental_complaints WHERE 
  EXISTS (
    SELECT 1 FROM rentals WHERE 
      rental_complaints.rental_id = rentals.id AND
      EXISTS (
        SELECT 1 FROM property_managers WHERE manager_id = $1 AND property_managers.property_id = rentals.property_id 
      )
  ) AND
  DATE_TRUNC('month', created_at) = DATE_TRUNC('month', $2)
`

type GetMaintenanceRequestsParams struct {
	ManagerID uuid.UUID       `json:"manager_id"`
	Month     pgtype.Interval `json:"month"`
}

func (q *Queries) GetMaintenanceRequests(ctx context.Context, arg GetMaintenanceRequestsParams) ([]int64, error) {
	rows, err := q.db.Query(ctx, getMaintenanceRequests, arg.ManagerID, arg.Month)
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

const getManagedPropertiesByRole = `-- name: GetManagedPropertiesByRole :many
SELECT property_id FROM property_managers WHERE manager_id = $1 AND role = $2
`

type GetManagedPropertiesByRoleParams struct {
	ManagerID uuid.UUID `json:"manager_id"`
	Role      string    `json:"role"`
}

func (q *Queries) GetManagedPropertiesByRole(ctx context.Context, arg GetManagedPropertiesByRoleParams) ([]uuid.UUID, error) {
	rows, err := q.db.Query(ctx, getManagedPropertiesByRole, arg.ManagerID, arg.Role)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []uuid.UUID
	for rows.Next() {
		var property_id uuid.UUID
		if err := rows.Scan(&property_id); err != nil {
			return nil, err
		}
		items = append(items, property_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getManagedUnits = `-- name: GetManagedUnits :many
SELECT id FROM units WHERE property_id IN (SELECT property_id FROM property_managers WHERE manager_id = $1)
`

func (q *Queries) GetManagedUnits(ctx context.Context, managerID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := q.db.Query(ctx, getManagedUnits, managerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
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

const getMostRentedProperties = `-- name: GetMostRentedProperties :many
SELECT r.id, COALESCE(c.count, 0) AS count
FROM 
  (SELECT properties.id FROM properties WHERE 
    EXISTS (
      SELECT 1 FROM property_managers WHERE property_managers.manager_id = $1 AND property_managers.property_id = properties.id 
    )
  ) AS r LEFT JOIN (
    SELECT property_id, COUNT(property_id) AS count FROM rentals WHERE 
    EXISTS (
      SELECT 1 FROM property_managers WHERE property_managers.manager_id = $1 AND property_managers.property_id = rentals.property_id 
    ) 
    GROUP BY property_id
  ) AS c ON r.id = c.property_id
ORDER BY count DESC
LIMIT $2
OFFSET $3
`

type GetMostRentedPropertiesParams struct {
	ManagerID uuid.UUID `json:"manager_id"`
	Limit     int32     `json:"limit"`
	Offset    int32     `json:"offset"`
}

type GetMostRentedPropertiesRow struct {
	ID    uuid.UUID `json:"id"`
	Count int64     `json:"count"`
}

func (q *Queries) GetMostRentedProperties(ctx context.Context, arg GetMostRentedPropertiesParams) ([]GetMostRentedPropertiesRow, error) {
	rows, err := q.db.Query(ctx, getMostRentedProperties, arg.ManagerID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetMostRentedPropertiesRow
	for rows.Next() {
		var i GetMostRentedPropertiesRow
		if err := rows.Scan(&i.ID, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getMostRentedUnits = `-- name: GetMostRentedUnits :many
SELECT r.id, r.property_id, COALESCE(c.count, 0) AS count
FROM 
  (SELECT units.id, units.property_id FROM units WHERE 
    EXISTS (
      SELECT 1 FROM property_managers WHERE property_managers.manager_id = $1 AND property_managers.property_id = units.property_id 
    )
  ) AS r LEFT JOIN (
    SELECT unit_id, COUNT(unit_id) AS count FROM rentals WHERE 
    EXISTS (
      SELECT 1 FROM property_managers WHERE property_managers.manager_id = $1 AND property_managers.property_id = rentals.property_id 
    ) 
    GROUP BY unit_id
  ) AS c ON r.id = c.unit_id
ORDER BY count DESC
LIMIT $2
OFFSET $3
`

type GetMostRentedUnitsParams struct {
	ManagerID uuid.UUID `json:"manager_id"`
	Limit     int32     `json:"limit"`
	Offset    int32     `json:"offset"`
}

type GetMostRentedUnitsRow struct {
	ID         uuid.UUID `json:"id"`
	PropertyID uuid.UUID `json:"property_id"`
	Count      int64     `json:"count"`
}

func (q *Queries) GetMostRentedUnits(ctx context.Context, arg GetMostRentedUnitsParams) ([]GetMostRentedUnitsRow, error) {
	rows, err := q.db.Query(ctx, getMostRentedUnits, arg.ManagerID, arg.Limit, arg.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetMostRentedUnitsRow
	for rows.Next() {
		var i GetMostRentedUnitsRow
		if err := rows.Scan(&i.ID, &i.PropertyID, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNewApplications = `-- name: GetNewApplications :many
SELECT id FROM applications WHERE 
  EXISTS (
    SELECT 1 FROM property_managers WHERE manager_id = $1 AND property_managers.property_id = applications.property_id 
  ) AND
  DATE_TRUNC('month', created_at) = DATE_TRUNC('month', $2)
`

type GetNewApplicationsParams struct {
	ManagerID uuid.UUID       `json:"manager_id"`
	Month     pgtype.Interval `json:"month"`
}

func (q *Queries) GetNewApplications(ctx context.Context, arg GetNewApplicationsParams) ([]int64, error) {
	rows, err := q.db.Query(ctx, getNewApplications, arg.ManagerID, arg.Month)
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

const getOccupiedProperties = `-- name: GetOccupiedProperties :many
SELECT property_id FROM property_managers 
WHERE 
  manager_id = $1 AND
  property_id IN (
    SELECT property_id FROM rentals WHERE start_date + INTERVAL '1 month' * rental_period >= CURRENT_DATE AND status = 'INPROGRESS'
  )
`

func (q *Queries) GetOccupiedProperties(ctx context.Context, managerID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := q.db.Query(ctx, getOccupiedProperties, managerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []uuid.UUID
	for rows.Next() {
		var property_id uuid.UUID
		if err := rows.Scan(&property_id); err != nil {
			return nil, err
		}
		items = append(items, property_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getOccupiedUnits = `-- name: GetOccupiedUnits :many
SELECT DISTINCT unit_id FROM rentals WHERE
  EXISTS (
    SELECT 1 FROM property_managers WHERE manager_id = $1 AND property_managers.property_id = rentals.property_id 
  ) AND 
  start_date + INTERVAL '1 month' * rental_period >= CURRENT_DATE AND
  status = 'INPROGRESS'
`

func (q *Queries) GetOccupiedUnits(ctx context.Context, managerID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := q.db.Query(ctx, getOccupiedUnits, managerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []uuid.UUID
	for rows.Next() {
		var unit_id uuid.UUID
		if err := rows.Scan(&unit_id); err != nil {
			return nil, err
		}
		items = append(items, unit_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getPaymentsStatistic = `-- name: GetPaymentsStatistic :one
SELECT coalesce(SUM(amount), 0)::REAL 
FROM payments 
WHERE 
  status = 'SUCCESS' AND 
  user_id = $1 AND
  created_at >= $2 AND
  created_at <= $3
`

type GetPaymentsStatisticParams struct {
	UserID    uuid.UUID `json:"user_id"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

func (q *Queries) GetPaymentsStatistic(ctx context.Context, arg GetPaymentsStatisticParams) (float32, error) {
	row := q.db.QueryRow(ctx, getPaymentsStatistic, arg.UserID, arg.StartDate, arg.EndDate)
	var column_1 float32
	err := row.Scan(&column_1)
	return column_1, err
}

const getPropertiesWithActiveListing = `-- name: GetPropertiesWithActiveListing :many
SELECT DISTINCT property_id FROM listings WHERE 
  EXISTS (SELECT 1 FROM property_managers WHERE property_managers.property_id = listings.property_id AND property_managers.manager_id = $1) AND
  active = TRUE AND
  expired_at > NOW()
`

func (q *Queries) GetPropertiesWithActiveListing(ctx context.Context, managerID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := q.db.Query(ctx, getPropertiesWithActiveListing, managerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []uuid.UUID
	for rows.Next() {
		var property_id uuid.UUID
		if err := rows.Scan(&property_id); err != nil {
			return nil, err
		}
		items = append(items, property_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRentalPaymentArrears = `-- name: GetRentalPaymentArrears :many
SELECT rental_payments.id, rental_payments.code, rental_payments.rental_id, rental_payments.created_at, rental_payments.updated_at, rental_payments.start_date, rental_payments.end_date, rental_payments.expiry_date, rental_payments.payment_date, rental_payments.updated_by, rental_payments.status, rental_payments.amount, rental_payments.discount, rental_payments.penalty, rental_payments.note, (rental_payments.expiry_date - CURRENT_DATE) AS expiry_duration, rentals.tenant_id, rentals.tenant_name, rentals.property_id, rentals.unit_id 
FROM rental_payments INNER JOIN rentals ON rentals.id = rental_payments.rental_id
WHERE 
  rental_payments.status IN ('ISSUED', 'PENDING', 'REQUEST2PAY') AND 
  EXISTS (
    SELECT 1 FROM rentals WHERE 
      rental_payments.rental_id = rentals.id AND
      EXISTS (
        SELECT 1 FROM property_managers WHERE manager_id = $1 AND property_managers.property_id = rentals.property_id 
      )
  ) AND
  rental_payments.expiry_date >= $4 AND
  rental_payments.expiry_date <= $5
ORDER BY
  (rental_payments.expiry_date - CURRENT_DATE) ASC
LIMIT $2
OFFSET $3
`

type GetRentalPaymentArrearsParams struct {
	ManagerID uuid.UUID   `json:"manager_id"`
	Limit     int32       `json:"limit"`
	Offset    int32       `json:"offset"`
	StartDate pgtype.Date `json:"start_date"`
	EndDate   pgtype.Date `json:"end_date"`
}

type GetRentalPaymentArrearsRow struct {
	ID             int64               `json:"id"`
	Code           string              `json:"code"`
	RentalID       int64               `json:"rental_id"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
	StartDate      pgtype.Date         `json:"start_date"`
	EndDate        pgtype.Date         `json:"end_date"`
	ExpiryDate     pgtype.Date         `json:"expiry_date"`
	PaymentDate    pgtype.Date         `json:"payment_date"`
	UpdatedBy      pgtype.UUID         `json:"updated_by"`
	Status         RENTALPAYMENTSTATUS `json:"status"`
	Amount         float32             `json:"amount"`
	Discount       pgtype.Float4       `json:"discount"`
	Penalty        pgtype.Float4       `json:"penalty"`
	Note           pgtype.Text         `json:"note"`
	ExpiryDuration int32               `json:"expiry_duration"`
	TenantID       pgtype.UUID         `json:"tenant_id"`
	TenantName     string              `json:"tenant_name"`
	PropertyID     uuid.UUID           `json:"property_id"`
	UnitID         uuid.UUID           `json:"unit_id"`
}

func (q *Queries) GetRentalPaymentArrears(ctx context.Context, arg GetRentalPaymentArrearsParams) ([]GetRentalPaymentArrearsRow, error) {
	rows, err := q.db.Query(ctx, getRentalPaymentArrears,
		arg.ManagerID,
		arg.Limit,
		arg.Offset,
		arg.StartDate,
		arg.EndDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRentalPaymentArrearsRow
	for rows.Next() {
		var i GetRentalPaymentArrearsRow
		if err := rows.Scan(
			&i.ID,
			&i.Code,
			&i.RentalID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.StartDate,
			&i.EndDate,
			&i.ExpiryDate,
			&i.PaymentDate,
			&i.UpdatedBy,
			&i.Status,
			&i.Amount,
			&i.Discount,
			&i.Penalty,
			&i.Note,
			&i.ExpiryDuration,
			&i.TenantID,
			&i.TenantName,
			&i.PropertyID,
			&i.UnitID,
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

const getRentalPaymentIncomes = `-- name: GetRentalPaymentIncomes :one
SELECT coalesce(SUM(amount), 0)::REAL 
FROM rental_payments 
WHERE 
  status = 'PAID' AND 
  EXISTS (
    SELECT 1 FROM rentals WHERE 
      rental_payments.rental_id = rentals.id AND
      EXISTS (
        SELECT 1 FROM property_managers WHERE manager_id = $1 AND property_managers.property_id = rentals.property_id 
      )
  ) AND
  payment_date >= $2 AND
  payment_date <= $3
`

type GetRentalPaymentIncomesParams struct {
	ManagerID uuid.UUID   `json:"manager_id"`
	StartDate pgtype.Date `json:"start_date"`
	EndDate   pgtype.Date `json:"end_date"`
}

func (q *Queries) GetRentalPaymentIncomes(ctx context.Context, arg GetRentalPaymentIncomesParams) (float32, error) {
	row := q.db.QueryRow(ctx, getRentalPaymentIncomes, arg.ManagerID, arg.StartDate, arg.EndDate)
	var column_1 float32
	err := row.Scan(&column_1)
	return column_1, err
}

const getRentedProperties = `-- name: GetRentedProperties :many
SELECT property_id FROM rentals WHERE tenant_id = $1
`

func (q *Queries) GetRentedProperties(ctx context.Context, tenantID pgtype.UUID) ([]uuid.UUID, error) {
	rows, err := q.db.Query(ctx, getRentedProperties, tenantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []uuid.UUID
	for rows.Next() {
		var property_id uuid.UUID
		if err := rows.Scan(&property_id); err != nil {
			return nil, err
		}
		items = append(items, property_id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
