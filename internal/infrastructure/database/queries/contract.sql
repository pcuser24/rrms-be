-- name: CreateContract :one
INSERT INTO "contracts" (
  rental_id,
  a_fullname,
  a_dob,
  a_phone,
  a_address,
  a_household_registration,
  a_identity,
  a_identity_issued_by,
  a_identity_issued_at,
  a_documents,
  a_bank_account,
  a_bank,
  a_registration_number,
  b_fullname,
  b_phone,
  payment_method,
  payment_day,
  n_copies,
  created_at_place,
  content,
  created_at,
  updated_at,
  created_by,
  updated_by
) VALUES (
  sqlc.arg(rental_id),
  sqlc.arg(a_fullname),
  sqlc.arg(a_dob),
  sqlc.arg(a_phone),
  sqlc.arg(a_address),
  sqlc.arg(a_household_registration),
  sqlc.arg(a_identity),
  sqlc.arg(a_identity_issued_by),
  sqlc.arg(a_identity_issued_at),
  sqlc.narg(a_documents),
  sqlc.narg(a_bank_account),
  sqlc.narg(a_bank),
  sqlc.arg(a_registration_number),
  sqlc.arg(b_fullname),
  sqlc.arg(b_phone),
  sqlc.arg(payment_method),
  sqlc.arg(payment_day),
  sqlc.arg(n_copies),
  sqlc.arg(created_at_place),
  sqlc.narg(content),
  NOW(),
  NOW(),
  sqlc.arg(user_id),
  sqlc.arg(user_id)
) RETURNING *;

-- name: GetContractByID :one
SELECT * FROM "contracts" WHERE "id" = $1;

-- name: GetContractByRentalID :one
SELECT * FROM "contracts" WHERE "rental_id" = $1;

-- name: PingContractByRentalID :one
SELECT id, rental_id, status, updated_by, updated_at FROM "contracts" WHERE "rental_id" = $1;

-- name: UpdateContract :exec
UPDATE "contracts" SET
  a_fullname = coalesce(sqlc.narg(a_fullname), a_fullname),
  a_dob = coalesce(sqlc.narg(a_dob), a_dob),
  a_phone = coalesce(sqlc.narg(a_phone), a_phone),
  a_address = coalesce(sqlc.narg(a_address), a_address),
  a_household_registration = coalesce(sqlc.narg(a_household_registration), a_household_registration),
  a_identity = coalesce(sqlc.narg(a_identity), a_identity),
  a_identity_issued_by = coalesce(sqlc.narg(a_identity_issued_by), a_identity_issued_by),
  a_identity_issued_at = coalesce(sqlc.narg(a_identity_issued_at), a_identity_issued_at),
  a_documents = coalesce(sqlc.narg(a_documents), a_documents),
  a_bank_account = coalesce(sqlc.narg(a_bank_account), a_bank_account),
  a_bank = coalesce(sqlc.narg(a_bank), a_bank),
  a_registration_number = coalesce(sqlc.narg(a_registration_number), a_registration_number),
  
  b_fullname = coalesce(sqlc.narg(b_fullname), b_fullname),
  b_organization_name = coalesce(sqlc.narg(b_organization_name), b_organization_name),
  b_organization_hq_address = coalesce(sqlc.narg(b_organization_hq_address), b_organization_hq_address),
  b_organization_code = coalesce(sqlc.narg(b_organization_code), b_organization_code),
  b_organization_code_issued_at = coalesce(sqlc.narg(b_organization_code_issued_at), b_organization_code_issued_at),
  b_organization_code_issued_by = coalesce(sqlc.narg(b_organization_code_issued_by), b_organization_code_issued_by),
  b_dob = coalesce(sqlc.narg(b_dob), b_dob),
  b_phone = coalesce(sqlc.narg(b_phone), b_phone),
  b_address = coalesce(sqlc.narg(b_address), b_address),
  b_household_registration = coalesce(sqlc.narg(b_household_registration), b_household_registration),
  b_identity = coalesce(sqlc.narg(b_identity), b_identity),
  b_identity_issued_by = coalesce(sqlc.narg(b_identity_issued_by), b_identity_issued_by),
  b_identity_issued_at = coalesce(sqlc.narg(b_identity_issued_at), b_identity_issued_at),
  b_bank_account = coalesce(sqlc.narg(b_bank_account), b_bank_account),
  b_bank = coalesce(sqlc.narg(b_bank), b_bank),
  b_tax_code = coalesce(sqlc.narg(b_tax_code), b_tax_code),
  
  payment_method = coalesce(sqlc.narg(payment_method), payment_method),
  payment_day = coalesce(sqlc.narg(payment_day), payment_day),
  n_copies = coalesce(sqlc.narg(n_copies), n_copies),
  created_at_place = coalesce(sqlc.narg(created_at_place), created_at_place),
  
  content = coalesce(sqlc.narg(content), content),
  updated_at = NOW(),
  updated_by = sqlc.arg(user_id)
WHERE id = $1;

-- name: UpdateContractContent :exec
UPDATE "contracts" SET
  content = coalesce(sqlc.narg(content), content),
  status = sqlc.arg(status),
  updated_at = NOW(),
  updated_by = sqlc.arg(user_id)
WHERE id = $1;
