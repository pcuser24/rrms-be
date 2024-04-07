BEGIN;

CREATE TYPE "CONTRACTSTATUS" AS ENUM ('PENDING_A', 'PENDING_B', 'PENDING', 'SIGNED', 'REJECTED', 'CANCELED');
CREATE TABLE IF NOT EXISTS "contracts" (
  "id" BIGSERIAL PRIMARY KEY,
  "rental_id" BIGINT NOT NULL,

  "a_fullname" TEXT NOT NULL,
  "a_dob" DATE NOT NULL,
  "a_phone" TEXT NOT NULL,
  "a_address" TEXT NOT NULL,
  "a_household_registration" TEXT NOT NULL,
  "a_identity" TEXT NOT NULL,
  "a_identity_issued_by" TEXT NOT NULL,
  "a_identity_issued_at" DATE NOT NULL,
  "a_documents" TEXT[],
  "a_bank_account" TEXT,
  "a_bank" TEXT,
  "a_registration_number" TEXT NOT NULL,

  "b_fullname" TEXT NOT NULL,
  "b_organization_name" TEXT,
  "b_organization_hq_address" TEXT,
  "b_organization_code" TEXT,
  "b_organization_code_issued_at" DATE,
  "b_organization_code_issued_by" TEXT,
  "b_dob" TEXT,
  "b_phone" TEXT NOT NULL,
  "b_address" TEXT,
  "b_household_registration" TEXT,
  "b_identity" TEXT,
  "b_identity_issued_by" TEXT,
  "b_identity_issued_at" DATE,
  "b_bank_account" TEXT,
  "b_bank" TEXT,
  "b_tax_code" TEXT,

  "payment_method" TEXT NOT NULL,
  "payment_day" INTEGER NOT NULL CHECK (payment_day >= 1 AND payment_day <= 28),
  "n_copies" INTEGER NOT NULL CHECK (n_copies >= 1),
  "created_at_place" TEXT NOT NULL,

  "content" TEXT NOT NULL DEFAULT '',

  "status" "CONTRACTSTATUS" DEFAULT 'PENDING_A' NOT NULL,
  "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  "created_by" UUID NOT NULL,
  "updated_by" UUID NOT NULL,

  UNIQUE ("rental_id")
);
ALTER TABLE "contracts" ADD CONSTRAINT "fk_contracts_rental_id" FOREIGN KEY ("rental_id") REFERENCES "rentals" ("id") ON DELETE CASCADE;
ALTER TABLE "contracts" ADD CONSTRAINT "fk_contracts_created_by" FOREIGN KEY ("created_by") REFERENCES "User" ("id") ON DELETE CASCADE;

END;
