BEGIN;

CREATE TYPE "UNITTYPE" AS ENUM ('ROOM', 'APARTMENT', 'STUDIO');
CREATE TABLE IF NOT EXISTS "units" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4() ,
  "property_id" UUID NOT NULL,
  "name" TEXT NOT NULL DEFAULT '',
  "area" REAL NOT NULL DEFAULT 0,
  "floor" INTEGER,
  "has_balcony" BOOL,
  "number_of_living_rooms" INTEGER CHECK (number_of_living_rooms >= 0),
  "number_of_bedrooms" INTEGER CHECK (number_of_bedrooms >= 0),
  "number_of_bathrooms" INTEGER CHECK (number_of_bathrooms >= 0),
  "number_of_toilets" INTEGER CHECK (number_of_bathrooms >= 0),
  "number_of_kitchens" INTEGER CHECK (number_of_kitchens >= 0),
  "type" "UNITTYPE" NOT NULL,
  "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
ALTER TABLE "units" ADD CONSTRAINT "property_unit_fkey" FOREIGN KEY ("property_id") REFERENCES "properties"("id") ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS "unit_media" (
  "id" SERIAL PRIMARY KEY,
  "unit_id" UUID NOT NULL,
  "url" TEXT NOT NULL,
  "type" "MEDIATYPE" NOT NULL
);
ALTER TABLE "unit_media" ADD CONSTRAINT "unit_id_media_fkey" FOREIGN KEY ("unit_id") REFERENCES "units"("id") ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS "unit_amenity" (
  "unit_id" UUID NOT NULL,
  "amenity" TEXT NOT NULL,
  "description" TEXT,

  PRIMARY KEY("unit_id", "amenity")
);
ALTER TABLE "unit_amenity" ADD CONSTRAINT "unit_id_amenity_fkey" FOREIGN KEY ("unit_id") REFERENCES "units"("id") ON DELETE CASCADE;
COMMENT ON TABLE "unit_amenity" IS 'Air conditioner, Fridge, Washing machine, ...';

END;