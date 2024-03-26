BEGIN;

CREATE TYPE "TENANTTYPE" AS ENUM ('INDIVIDUAL', 'ORGANIZATION');
CREATE TYPE "CONTRACTTYPE" AS ENUM ('DIGITAL', 'FILE', 'IMAGE');
CREATE TYPE "PRERENTALSTATUS" AS ENUM ('PENDING', 'FINISHED', 'CANCELED');

CREATE TABLE IF NOT EXISTS "prerentals" (
  "id" BIGSERIAL PRIMARY KEY,
  "creator_id" UUID NOT NULL,
  "property_id" UUID NOT NULL,
  "unit_id" UUID NOT NULL,
  "application_id" BIGINT,
  "tenant_id" UUID,
  "profile_image" TEXT NOT NULL,

  "tenant_type" "TENANTTYPE" NOT NULL,
  "tenant_name" VARCHAR(100) NOT NULL,
  "tenant_identity" VARCHAR(20) NOT NULL,
  "tenant_dob" DATE NOT NULL,
  "tenant_phone" VARCHAR(20) NOT NULL,
  "tenant_email" VARCHAR(100) NOT NULL,
  "tenant_address" VARCHAR(255),

  -- contract
  "contract_type" "CONTRACTTYPE",
  "contract_content" TEXT,
  "contract_last_update_at" TIMESTAMPTZ DEFAULT NOW(),
  "contract_last_update_by" UUID,

  "land_area" REAL NOT NULL CHECK (land_area >= 0),
  "unit_area" REAL NOT NULL CHECK (unit_area >= 0),

  "start_date" DATE,
  "movein_date" DATE NOT NULL,
  "rental_period" INTEGER NOT NULL CHECK (rental_period >= 0),
  "rental_price" REAL NOT NULL CHECK (rental_price >= 0),

  "status" "PRERENTALSTATUS" NOT NULL DEFAULT 'PENDING',

  "note" TEXT,

  UNIQUE ("application_id")
);
ALTER TABLE "prerentals" ADD CONSTRAINT "prerental_application_id_fkey" FOREIGN KEY ("application_id") REFERENCES "applications"("id") ON DELETE SET NULL;
ALTER TABLE "prerentals" ADD CONSTRAINT "prerental_creator_id_fkey" FOREIGN KEY ("creator_id") REFERENCES "User"("id") ON DELETE CASCADE;
ALTER TABLE "prerentals" ADD CONSTRAINT "prerental_tenant_id_fkey" FOREIGN KEY ("tenant_id") REFERENCES "User"("id") ON DELETE CASCADE;
ALTER TABLE "prerentals" ADD CONSTRAINT "prerental_property_id_fkey" FOREIGN KEY ("property_id") REFERENCES "properties"("id") ON DELETE CASCADE;
ALTER TABLE "prerentals" ADD CONSTRAINT "prerental_unit_id_fkey" FOREIGN KEY ("unit_id") REFERENCES "units"("id") ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS "prerental_coaps" (
  "prerental_id" BIGINT NOT NULL,
  "full_name" TEXT,
  "dob" DATE,
  "job" TEXT,
  "income" INTEGER,
  "email" TEXT,
  "phone" TEXT,
  "description" TEXT
);
ALTER TABLE "prerental_coaps" ADD CONSTRAINT "prerental_coaps_prerental_id_fkey" FOREIGN KEY ("prerental_id") REFERENCES "prerentals"("id") ON DELETE CASCADE;

END;
