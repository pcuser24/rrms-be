BEGIN;

CREATE TABLE IF NOT EXISTS "listings" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4() ,
  "creator_id" UUID NOT NULL,
  "property_id" UUID NOT NULL,
  "title" TEXT NOT NULL DEFAULT '',
  "description" TEXT NOT NULL DEFAULT '',
  
  -- contact info
  "full_name" TEXT NOT NULL,
  "email" TEXT NOT NULL,
  "phone" TEXT NOT NULL,
  "contact_type" TEXT NOT NULL,

  "price" REAL NOT NULL,
  "price_negotiable" BOOL NOT NULL DEFAULT FALSE,
  "security_deposit" REAL,
  "lease_term" INTEGER,

  "pets_allowed" BOOL,
  "number_of_residents" INTEGER CHECK (number_of_residents >= 0),

  "priority" INTEGER NOT NULL DEFAULT 1,
  "active" BOOL NOT NULL DEFAULT FALSE,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  -- "post_at" TIMESTAMPTZ NOT NULL,
  "expired_at" TIMESTAMPTZ NOT NULL
);
ALTER TABLE "listings" ADD CONSTRAINT "listings_creator_id_fkey" FOREIGN KEY ("creator_id") REFERENCES "User"("id") ON DELETE CASCADE;
ALTER TABLE "listings" ADD CONSTRAINT "listings_property_id_fkey" FOREIGN KEY ("property_id") REFERENCES "properties"("id") ON DELETE CASCADE;
COMMENT ON COLUMN "listings"."price" IS 'Rental price per month in vietnamese dong';
COMMENT ON COLUMN "listings"."priority" IS 'Priority of the listing, range from 1 to 5, 1 is the lowest';
COMMENT ON COLUMN "listings"."lease_term" IS 'Lease term in months';
-- COMMENT ON COLUMN "listings"."post_at" IS 'The time when the listing goes public';
COMMENT ON COLUMN "listings"."expired_at" IS 'The time when the listing is expired. The listing is expired if the current time is greater than this time.';

CREATE TABLE IF NOT EXISTS "listing_units" (
  "listing_id" UUID NOT NULL,
  "unit_id" UUID NOT NULL,
  "price" BIGINT NOT NULL CHECK (price >= 0),
  PRIMARY KEY("listing_id", "unit_id")
);
ALTER TABLE "listing_units" ADD CONSTRAINT "listing_units_listing_id_fkey" FOREIGN KEY ("listing_id") REFERENCES "listings"("id") ON DELETE CASCADE;
ALTER TABLE "listing_units" ADD CONSTRAINT "listing_units_unit_id_fkey" FOREIGN KEY ("unit_id") REFERENCES "units"("id") ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS "l_policies" (
  "id" BIGSERIAL PRIMARY KEY,
  "policy" TEXT NOT NULL
);
INSERT INTO l_policies (policy) VALUES
('l_policy-payment'),
('l_policy-maintenance'),
('l_policy-insurance'),
('l_policy-noise'),
('l_policy-lease_renewal'),
('l_policy-change_to_property'),
('l_policy-parking'),
('l_policy-pets'),
('l_policy-subletting'),
('l_policy-business'),
('l_policy-consequences'),
('l_policy-other');

CREATE TABLE IF NOT EXISTS "listing_policies" (
  "listing_id" UUID NOT NULL,
  "policy_id" BIGINT NOT NULL,
  "note" TEXT,

  PRIMARY KEY("listing_id", "policy_id")
);
ALTER TABLE "listing_policies" ADD CONSTRAINT "listing_policies_listing_id_fkey" FOREIGN KEY ("listing_id") REFERENCES "listings"("id") ON DELETE CASCADE;
ALTER TABLE "listing_policies" ADD CONSTRAINT "listing_policies_policy_id_fkey" FOREIGN KEY ("policy_id") REFERENCES "l_policies"("id") ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS "listing_tags" (
  "id" BIGSERIAL PRIMARY KEY,
  "listing_id" UUID NOT NULL,
  "tag" VARCHAR(32) NOT NULL
);
ALTER TABLE "listing_tags" ADD CONSTRAINT "listing_tags_listing_id_fkey" FOREIGN KEY ("listing_id") REFERENCES "listings"("id") ON DELETE CASCADE;

END;
