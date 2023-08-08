BEGIN;

CREATE TYPE "PROPERTYTYPE" AS ENUM (
  'APARTMENT',
  'SINGLE_RESIDENCE',
  'ROOM',
  'BLOCK'
);

CREATE TABLE IF NOT EXISTS "properties" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "owner_id" UUID NOT NULL,
  "name" TEXT NOT NULL DEFAULT '',
  "area" REAL NOT NULL DEFAULT 0,
  "number_of_floors" INTEGER,
  "year_built" INTEGER,
  "orientation" VARCHAR(4),
  "full_address" TEXT NOT NULL,
  "district" TEXT NOT NULL,
  "city" TEXT NOT NULL,
  "lat" DOUBLE PRECISION NOT NULL,
  "lng" DOUBLE PRECISION NOT NULL,
  "type" "PROPERTYTYPE" NOT NULL,
  "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
ALTER TABLE "properties" ADD CONSTRAINT "property_user_fkey" FOREIGN KEY ("owner_id") REFERENCES "User"("id") ON DELETE CASCADE;
COMMENT ON COLUMN "properties"."orientation" IS 'n,s,w,e,nw,ne,sw,se';

CREATE TABLE IF NOT EXISTS "property_amenity" (
  "property_id" UUID NOT NULL,
  "amenity" TEXT NOT NULL,
  "description" TEXT,

  PRIMARY KEY("property_id", "amenity")
);
ALTER TABLE "property_amenity" ADD CONSTRAINT "property_id_amenity_fkey" FOREIGN KEY ("property_id") REFERENCES "properties"("id") ON DELETE CASCADE;
COMMENT ON TABLE "property_amenity" IS 'Elevator, Security camera, Pool, Yard, ...';

CREATE TABLE IF NOT EXISTS "property_feature" (
  "property_id" UUID NOT NULL,
  "feature" TEXT NOT NULL,
  "description" TEXT,

  PRIMARY KEY("property_id", "feature")
);
ALTER TABLE "property_feature" ADD CONSTRAINT "property_id_feature_fkey" FOREIGN KEY ("property_id") REFERENCES "properties"("id") ON DELETE CASCADE;
COMMENT ON TABLE "property_feature" IS 'Security guard, Parking, Gym, ...';

CREATE TABLE IF NOT EXISTS "property_tag" (
  "id" SERIAL PRIMARY KEY,
  "property_id" UUID NOT NULL,
  "tag" TEXT NOT NULL
);
ALTER TABLE "property_tag" ADD CONSTRAINT "property_id_tag_fkey" FOREIGN KEY ("property_id") REFERENCES "properties"("id") ON DELETE CASCADE;
COMMENT ON TABLE "property_tag" IS '';

CREATE TYPE "MEDIATYPE" AS ENUM (
  'IMAGE',
  'VIDEO'
);
CREATE TABLE IF NOT EXISTS "property_media" (
  "id" SERIAL PRIMARY KEY,
  "property_id" UUID NOT NULL,
  "url" TEXT NOT NULL,
  "type" "MEDIATYPE" NOT NULL
);
ALTER TABLE "property_media" ADD CONSTRAINT "property_id_media_fkey" FOREIGN KEY ("property_id") REFERENCES "properties"("id") ON DELETE CASCADE;

END;