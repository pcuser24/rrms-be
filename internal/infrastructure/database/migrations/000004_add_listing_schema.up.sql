BEGIN;

CREATE TABLE IF NOT EXISTS "listings" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4() ,
  "creator_id" UUID NOT NULL,
  "property_id" UUID NOT NULL,
  "title" TEXT NOT NULL DEFAULT '',
  "description" TEXT NOT NULL DEFAULT '',

  "price" BIGINT NOT NULL,
  "security_deposit" BIGINT NOT NULL DEFAULT 0,
  "lease_term" INTEGER NOT NULL DEFAULT 1,

  "pets_allowed" BOOL,
  "number_of_residents" INTEGER CHECK (number_of_residents >= 0),

  "priority" INTEGER NOT NULL DEFAULT 1,
  "active" BOOL NOT NULL DEFAULT TRUE,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT NOW(),
  "expired_at" TIMESTAMPTZ NOT NULL
);
ALTER TABLE "listings" ADD CONSTRAINT "listings_creator_id_fkey" FOREIGN KEY ("creator_id") REFERENCES "User"("id") ON DELETE CASCADE;
ALTER TABLE "listings" ADD CONSTRAINT "listings_property_id_fkey" FOREIGN KEY ("property_id") REFERENCES "properties"("id") ON DELETE CASCADE;
COMMENT ON COLUMN "listings"."price" IS 'Rental price per month in vietnamese dong';
COMMENT ON COLUMN "listings"."priority" IS 'Priority of the listing, range from 1 to 5, 1 is the lowest';
COMMENT ON COLUMN "listings"."lease_term" IS 'Lease term in months';

CREATE TABLE IF NOT EXISTS "listing_unit" (
  "listing_id" UUID NOT NULL,
  "unit_id" UUID NOT NULL,
  PRIMARY KEY("listing_id", "unit_id")
);
ALTER TABLE "listing_unit" ADD CONSTRAINT "listing_unit_listing_id_fkey" FOREIGN KEY ("listing_id") REFERENCES "listings"("id") ON DELETE CASCADE;
ALTER TABLE "listing_unit" ADD CONSTRAINT "listing_unit_unit_id_fkey" FOREIGN KEY ("unit_id") REFERENCES "units"("id") ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS "rental_policies" (
  "id" BIGSERIAL PRIMARY KEY,
  "policy" TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS "listing_policy" (
  "listing_id" UUID NOT NULL,
  "policy_id" BIGINT NOT NULL,
  "note" TEXT,

  PRIMARY KEY("listing_id", "policy_id")
);
ALTER TABLE "listing_policy" ADD CONSTRAINT "listing_policy_listing_id_fkey" FOREIGN KEY ("listing_id") REFERENCES "listings"("id") ON DELETE CASCADE;
ALTER TABLE "listing_policy" ADD CONSTRAINT "listing_policy_policy_id_fkey" FOREIGN KEY ("policy_id") REFERENCES "rental_policies"("id") ON DELETE CASCADE;

END;