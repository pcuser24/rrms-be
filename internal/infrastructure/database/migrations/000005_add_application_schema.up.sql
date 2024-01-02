BEGIN;

CREATE TYPE "APPLICATION_STATUS" AS ENUM ('PENDING', 'APPROVED', 'CONDITIONALLY_APPROVED', 'REJECTED');
CREATE TABLE IF NOT EXISTS "applications" (
  id BIGSERIAL PRIMARY KEY,
  creator_id UUID NOT NULL,
  listing_id UUID NOT NULL,
  property_id UUID NOT NULL,
  unit_ids UUID[] NOT NULL,
  status "APPLICATION_STATUS" NOT NULL DEFAULT 'PENDING',
  created_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  updated_at TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  -- basic info
  full_name TEXT NOT NULL,
  email TEXT NOT NULL,
  phone TEXT NOT NULL,
  dob TIMESTAMPTZ NOT NULL,
  profile_image TEXT NOT NULL,
  movein_date TIMESTAMPTZ NOT NULL,
  preferred_term INTEGER NOT NULL CHECK (preferred_term >= 0),
  -- rental history
  rh_address TEXT,
  rh_city TEXT,
  rh_district TEXT,
  rh_ward TEXT,
  rh_rental_duration INTEGER CHECK (rh_rental_duration >= 0),
  rh_monthly_payment REAL CHECK (rh_monthly_payment >= 0),
  rh_reason_for_leaving TEXT,
  -- employment
  employment_status VARCHAR(10) NOT NULL,
  employment_company_name TEXT,
  employment_position TEXT,
  employment_monthly_income REAL CHECK (employment_monthly_income >= 0),
  employment_comment TEXT,
  employment_proofs_of_income TEXT[],
  -- identity
  identity_type VARCHAR(10) NOT NULL,
  identity_number TEXT NOT NULL,
  identity_issued_date TIMESTAMPTZ NOT NULL,
  identity_issued_by TEXT NOT NULL
);
ALTER TABLE "applications" ADD CONSTRAINT "applications_creator_id_fkey" FOREIGN KEY ("creator_id") REFERENCES "User"("id") ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS "application_minors" (
  application_id BIGINT NOT NULL,
  full_name TEXT NOT NULL,
  dob TIMESTAMPTZ NOT NULL,
  email TEXT,
  phone TEXT,
  description TEXT
);
ALTER TABLE "application_minors" ADD CONSTRAINT "application_minors_application_id_fkey" FOREIGN KEY ("application_id") REFERENCES "applications"("id") ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS "application_coaps" (
  application_id BIGINT NOT NULL,
  full_name TEXT NOT NULL,
  dob TIMESTAMPTZ NOT NULL,
  job TEXT NOT NULL,
  income INTEGER NOT NULL,
  email TEXT,
  phone TEXT,
  description TEXT
);
ALTER TABLE "application_coaps" ADD CONSTRAINT "application_coaps_application_id_fkey" FOREIGN KEY ("application_id") REFERENCES "applications"("id") ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS "application_pets" (
  application_id BIGINT NOT NULL,
  type VARCHAR(10) NOT NULL,
  weight REAL,
  description TEXT
);
ALTER TABLE "application_pets" ADD CONSTRAINT "application_pets_application_id_fkey" FOREIGN KEY ("application_id") REFERENCES "applications"("id") ON DELETE CASCADE;


CREATE TABLE IF NOT EXISTS "application_vehicles" (
  application_id BIGINT NOT NULL,
  type VARCHAR(10) NOT NULL,
  model TEXT,
  code TEXT NOT NULL,
  description TEXT
);
ALTER TABLE "application_vehicles" ADD CONSTRAINT "application_vehicles_application_id_fkey" FOREIGN KEY ("application_id") REFERENCES "applications"("id") ON DELETE CASCADE;

END;
