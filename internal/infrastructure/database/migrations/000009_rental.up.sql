BEGIN;

CREATE TYPE "RENTALPAYMENTTYPE" AS ENUM ('PREPAID', 'POSTPAID');
CREATE TYPE "RENTALSTATUS" AS ENUM ('INPROGRESS', 'END');
CREATE TYPE "LATEPAYMENTPENALTYSCHEME" AS ENUM ('FIXED', 'PERCENT', 'NONE');
CREATE TABLE IF NOT EXISTS "rentals" (
  "id" BIGSERIAL PRIMARY KEY,
  "creator_id" UUID NOT NULL,
  "property_id" UUID NOT NULL,
  "unit_id" UUID NOT NULL,
  "application_id" BIGINT,

  "tenant_id" UUID,
  "profile_image" TEXT NOT NULL,
  "tenant_type" "TENANTTYPE" NOT NULL,
  "tenant_name" VARCHAR(100) NOT NULL,
  "tenant_phone" VARCHAR(20) NOT NULL,
  "tenant_email" VARCHAR(100) NOT NULL,
  "organization_name" TEXT,
  "organization_hq_address" TEXT,

  "start_date" DATE NOT NULL,
  "movein_date" DATE NOT NULL,
  "rental_period" INTEGER NOT NULL CHECK (rental_period >= 0),

  "payment_type" "RENTALPAYMENTTYPE" NOT NULL DEFAULT 'POSTPAID',

  "rental_price" REAL NOT NULL CHECK (rental_price >= 0),
  "rental_payment_basis" INTEGER NOT NULL CHECK(rental_payment_basis >= 1),
  CHECK(rental_payment_basis <= rental_period),
  "rental_intention" VARCHAR(20) NOT NULL,
  "notice_period" INTEGER DEFAULT 30,
  "grace_period" INTEGER DEFAULT 0 CHECK (grace_period >= 0),
  "late_payment_penalty_scheme" "LATEPAYMENTPENALTYSCHEME" DEFAULT 'FIXED',
  "late_payment_penalty_amount" REAL DEFAULT 0 CHECK (late_payment_penalty_amount >= 0),
  
  -- basic services
  "electricity_setup_by" VARCHAR(20) NOT NULL,
  "electricity_payment_type" VARCHAR(10),
  "electricity_customer_code" VARCHAR(50),
  "electricity_provider" TEXT,
  "electricity_price" REAL CHECK (electricity_price >= 0),
  "water_setup_by" VARCHAR(20) NOT NULL,
  "water_payment_type" VARCHAR(10),
  "water_customer_code" VARCHAR(50),
  "water_provider" TEXT,
  "water_price" REAL CHECK (water_price >= 0),

  -- policy
  -- "rental_payment_grace_period" INTEGER NOT NULL CHECK (rental_payment_grace_period >= 0),
  "note" TEXT,

  "status" "RENTALSTATUS" NOT NULL DEFAULT 'INPROGRESS',
  
  "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,

  UNIQUE ("application_id")
);
ALTER TABLE "rentals" ADD CONSTRAINT "rental_application_id_fkey" FOREIGN KEY ("application_id") REFERENCES "applications"("id") ON DELETE SET NULL;
ALTER TABLE "rentals" ADD CONSTRAINT "rental_creator_id_fkey" FOREIGN KEY ("creator_id") REFERENCES "User"("id") ON DELETE CASCADE;
ALTER TABLE "rentals" ADD CONSTRAINT "rental_tenant_id_fkey" FOREIGN KEY ("tenant_id") REFERENCES "User"("id") ON DELETE CASCADE;
ALTER TABLE "rentals" ADD CONSTRAINT "rental_property_id_fkey" FOREIGN KEY ("property_id") REFERENCES "properties"("id") ON DELETE CASCADE;
ALTER TABLE "rentals" ADD CONSTRAINT "rental_unit_id_fkey" FOREIGN KEY ("unit_id") REFERENCES "units"("id") ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS "rental_coaps" (
  "rental_id" BIGINT NOT NULL,
  "full_name" TEXT,
  "dob" DATE,
  "job" TEXT,
  "income" INTEGER,
  "email" TEXT,
  "phone" TEXT,
  "description" TEXT
);
ALTER TABLE "rental_coaps" ADD CONSTRAINT "rental_coaps_rental_id_fkey" FOREIGN KEY ("rental_id") REFERENCES "rentals"("id") ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS "rental_minors" (
  "rental_id" BIGINT NOT NULL,
  "full_name" TEXT NOT NULL,
  "dob" DATE NOT NULL,
  "email" TEXT,
  "phone" TEXT,
  description TEXT
);
ALTER TABLE "rental_minors" ADD CONSTRAINT "rental_minors_rental_id_fkey" FOREIGN KEY ("rental_id") REFERENCES "rentals"("id") ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS "rental_pets" (
  "rental_id" BIGINT NOT NULL,
  "type" VARCHAR(10) NOT NULL,
  "weight" REAL,
  "description" TEXT
);
ALTER TABLE "rental_pets" ADD CONSTRAINT "rental_pets_rental_id_fkey" FOREIGN KEY ("rental_id") REFERENCES "rentals"("id") ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS "rental_services" (
  "id" BIGSERIAL PRIMARY KEY,
  "rental_id" BIGINT NOT NULL,
  "name" TEXT NOT NULL,
  "setup_by" VARCHAR(20) NOT NULL,
  "provider" TEXT,
  "price" REAL CHECK (price >= 0)
);
ALTER TABLE "rental_services" ADD CONSTRAINT "rental_services_rental_id_fkey" FOREIGN KEY ("rental_id") REFERENCES "rentals"("id") ON DELETE CASCADE;
COMMENT ON COLUMN "rental_services"."setup_by" IS 'The party who set up the service, either "LANDLORD" or "TENANT"';

CREATE TABLE IF NOT EXISTS "rental_policies" (
  "rental_id" BIGINT NOT NULL,
  "title" TEXT NOT NULL,
  "content" TEXT NOT NULL
);
ALTER TABLE "rental_policies" ADD CONSTRAINT "rental_policies_rental_id_fkey" FOREIGN KEY ("rental_id") REFERENCES "rentals"("id") ON DELETE CASCADE;

END;
