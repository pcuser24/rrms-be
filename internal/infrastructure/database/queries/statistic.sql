-- name: GetManagedPropertiesByRole :many
SELECT property_id FROM property_managers WHERE manager_id = $1 AND role = $2;

-- name: GetRentedProperties :many
SELECT property_id FROM rentals WHERE tenant_id = $1;

-- name: GetPropertiesWithActiveListing :many
SELECT DISTINCT property_id FROM listings WHERE 
  EXISTS (SELECT 1 FROM property_managers WHERE property_managers.property_id = listings.property_id AND property_managers.manager_id = $1) AND
  active = TRUE AND
  expired_at > NOW()
;

-- name: GetOccupiedProperties :many
SELECT property_id FROM property_managers 
WHERE 
  manager_id = $1 AND
  property_id IN (
    SELECT property_id FROM rentals WHERE start_date + INTERVAL '1 month' * rental_period >= CURRENT_DATE AND status = 'INPROGRESS'
  );

-- name: GetMostRentedProperties :many
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
;

-- name: GetLeastRentedProperties :many
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
;

-- name: GetMostRentedUnits :many
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
;

-- name: GetLeastRentedUnits :many
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
;

-- name: GetManagedUnits :many
SELECT id FROM units WHERE property_id IN (SELECT property_id FROM property_managers WHERE manager_id = $1);

-- name: GetOccupiedUnits :many
SELECT DISTINCT unit_id FROM rentals WHERE
  EXISTS (
    SELECT 1 FROM property_managers WHERE manager_id = $1 AND property_managers.property_id = rentals.property_id 
  ) AND 
  start_date + INTERVAL '1 month' * rental_period >= CURRENT_DATE AND
  status = 'INPROGRESS';

-- name: GetApplicationsInMonth :many
SELECT id FROM applications WHERE 
  EXISTS (
    SELECT 1 FROM property_managers WHERE manager_id = $1 AND property_managers.property_id = applications.property_id 
  ) AND
  DATE_TRUNC('month', created_at) = DATE_TRUNC('month', sqlc.arg(month)::TIMESTAMP)
;

-- name: GetMaintenanceRequests :many
SELECT id FROM rental_complaints WHERE 
  EXISTS (
    SELECT 1 FROM rentals WHERE 
      rental_complaints.rental_id = rentals.id AND
      EXISTS (
        SELECT 1 FROM property_managers WHERE manager_id = $1 AND property_managers.property_id = rentals.property_id 
      )
  ) AND
  DATE_TRUNC('month', created_at) = DATE_TRUNC('month', sqlc.arg(month))
;

-- name: GetRentalPaymentArrears :many
SELECT rental_payments.*, (rental_payments.expiry_date - CURRENT_DATE) AS expiry_duration, rentals.tenant_id, rentals.tenant_name, rentals.property_id, rentals.unit_id 
FROM rental_payments INNER JOIN rentals ON rentals.id = rental_payments.rental_id
WHERE 
  rental_payments.status IN ('ISSUED', 'PENDING', 'REQUEST2PAY', 'PARTIALLYPAID', 'PAYFINE') AND 
  EXISTS (
    SELECT 1 FROM rentals WHERE 
      rental_payments.rental_id = rentals.id AND
      EXISTS (
        SELECT 1 FROM property_managers WHERE manager_id = $1 AND property_managers.property_id = rentals.property_id 
      )
  ) AND
  rental_payments.expiry_date >= sqlc.arg(start_date) AND
  rental_payments.expiry_date <= sqlc.arg(end_date)
ORDER BY
  (rental_payments.expiry_date - CURRENT_DATE) ASC
LIMIT $2
OFFSET $3
;

-- name: GetRentalPaymentIncomes :one
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
  payment_date >= sqlc.arg(start_date) AND
  payment_date <= sqlc.arg(end_date)
  ;

-- name: GetPaymentsStatistic :one
SELECT coalesce(SUM(amount), 0)::REAL 
FROM payments 
WHERE 
  status = 'SUCCESS' AND 
  user_id = $1 AND
  created_at >= sqlc.arg(start_date) AND
  created_at <= sqlc.arg(end_date)
  ;



-- name: GetListingsCountByCity :one
SELECT COUNT(*) 
FROM listings 
WHERE
  EXISTS (
    SELECT 1 FROM properties WHERE city = sqlc.arg(city) AND properties.id = listings.property_id
  );

-- name: GetTenantPendingPayments :many
SELECT rental_payments.*, (rental_payments.expiry_date - CURRENT_DATE) AS expiry_duration, rentals.tenant_id, rentals.tenant_name, rentals.property_id, rentals.unit_id 
FROM rental_payments INNER JOIN rentals ON rentals.id = rental_payments.rental_id
WHERE 
  rental_payments.status IN ('ISSUED', 'PENDING', 'REQUEST2PAY') AND 
  EXISTS (
    SELECT 1 FROM rentals WHERE 
      rental_payments.rental_id = rentals.id AND
      rentals.tenant_id = sqlc.arg(user_id)
  )
ORDER BY
  (rental_payments.expiry_date - CURRENT_DATE) ASC
LIMIT $1 OFFSET $2;

-- name: GetTotalTenantPendingPayments :one
SELECT coalesce(SUM(rental_payments.amount), 0)::REAL
FROM rental_payments 
WHERE 
  status IN ('ISSUED', 'PENDING', 'REQUEST2PAY') AND 
  EXISTS (
    SELECT 1 FROM rentals WHERE 
      rental_payments.rental_id = rentals.id AND
      rentals.tenant_id = sqlc.arg(user_id)
  );

-- name: GetTenantExpenditure :one
SELECT coalesce(SUM(amount), 0)::REAL 
FROM rental_payments 
WHERE 
  status = 'PAID' AND 
  EXISTS (
    SELECT 1 FROM rentals WHERE 
      rental_payments.rental_id = rentals.id AND
      rentals.tenant_id = sqlc.arg(user_id)
  ) AND
  payment_date >= sqlc.arg(start_date) AND
  payment_date <= sqlc.arg(end_date)
  ;

-- name: GetRentalComplaintStatistics :one
SELECT COUNT(*) FROM rental_complaints
WHERE
  rental_complaints.status = sqlc.arg(status) AND
  EXISTS (
    SELECT 1 FROM rentals WHERE 
      rentals.id = rental_complaints.rental_id AND (
        rentals.tenant_id = sqlc.arg(user_id)
        OR EXISTS (
          SELECT 1 FROM property_managers WHERE property_managers.property_id = rentals.property_id AND manager_id = sqlc.arg(user_id)
        )
      )
  );

-- name: GetTotalTenantsManagedByUserStatistic :one
WITH rs AS (
  SELECT id FROM rentals WHERE 
  EXISTS (
    SELECT 1 FROM property_managers 
    WHERE property_managers.property_id = rentals.property_id 
    AND manager_id = sqlc.arg(user_id)
  ) 
  AND start_date >= sqlc.arg(start_time) 
  AND start_date <= sqlc.arg(end_time)
  AND start_date + INTERVAL '1 month' * rental_period >= CURRENT_DATE
  AND status = 'INPROGRESS'
)
SELECT 
  COALESCE(SUM(coaps_count + minors_count) + COUNT(DISTINCT counts.id), 0)::INTEGER AS total_associated_records
FROM (
  SELECT 
    rs.id, 
    COALESCE(rc.coaps_count, 0) AS coaps_count, 
    COALESCE(rm.minors_count, 0) AS minors_count
  FROM rs 
  LEFT JOIN (
    SELECT rental_coaps.rental_id, COUNT(*) AS coaps_count 
    FROM rental_coaps 
    GROUP BY rental_id
  ) rc ON rs.id = rc.rental_id
  LEFT JOIN (
    SELECT rental_id, COUNT(*) AS minors_count 
    FROM rental_minors
    GROUP BY rental_id
  ) rm ON rs.id = rm.rental_id
) AS counts;

-- name: GetTotalTenantsOfUnitStatistic :one
WITH rs AS (
  SELECT id FROM rentals WHERE 
  unit_id = sqlc.arg(unit_id)
  AND start_date + INTERVAL '1 month' * rental_period >= CURRENT_DATE
  AND status = 'INPROGRESS'
)
SELECT 
  COALESCE(SUM(coaps_count + minors_count) + COUNT(DISTINCT counts.id), 0)::INTEGER AS total_associated_records
FROM (
  SELECT 
    rs.id, 
    COALESCE(rc.coaps_count, 0) AS coaps_count, 
    COALESCE(rm.minors_count, 0) AS minors_count
  FROM rs 
  LEFT JOIN (
    SELECT rental_coaps.rental_id, COUNT(*) AS coaps_count 
    FROM rental_coaps 
    GROUP BY rental_id
  ) rc ON rs.id = rc.rental_id
  LEFT JOIN (
    SELECT rental_id, COUNT(*) AS minors_count 
    FROM rental_minors
    GROUP BY rental_id
  ) rm ON rs.id = rm.rental_id
) AS counts;