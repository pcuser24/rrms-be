BEGIN;

CREATE TYPE "PROPERTYTYPE" AS ENUM (
  'APARTMENT', 'PRIVATE',
  'ROOM', 
  'STORE', 'OFFICE',
  'VILLA',
  'MINIAPARTMENT'
);

CREATE TABLE IF NOT EXISTS "properties" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "creator_id" UUID NOT NULL,
  "name" TEXT NOT NULL DEFAULT '',
  "building" TEXT,
  "project" TEXT,
  "area" REAL NOT NULL DEFAULT 0,
  "number_of_floors" INTEGER,
  "year_built" INTEGER,
  "orientation" VARCHAR(4),
  "entrance_width" REAL,
  "facade" REAL,
  "full_address" TEXT NOT NULL,
  "city" TEXT NOT NULL,
  "district" TEXT NOT NULL,
  "ward" TEXT,
  "lat" DOUBLE PRECISION,
  "lng" DOUBLE PRECISION,
  "primary_image" BIGINT,
  "description" TEXT,
  "type" "PROPERTYTYPE" NOT NULL,
  "is_public" BOOLEAN NOT NULL DEFAULT FALSE,
  "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
ALTER TABLE "properties" ADD CONSTRAINT "property_user_fkey" FOREIGN KEY ("creator_id") REFERENCES "User"("id") ON DELETE CASCADE;
COMMENT ON COLUMN "properties"."orientation" IS 'n,s,w,e,nw,ne,sw,se';

CREATE TABLE IF NOT EXISTS "property_managers" (
  "property_id" UUID NOT NULL,
  "manager_id" UUID NOT NULL,
  "role" TEXT NOT NULL,

  PRIMARY KEY("property_id", "manager_id")
);
ALTER TABLE "property_managers" ADD CONSTRAINT "property_managers_property_fkey" FOREIGN KEY ("property_id") REFERENCES "properties"("id") ON DELETE CASCADE;
ALTER TABLE "property_managers" ADD CONSTRAINT "property_managers_manager_fkey" FOREIGN KEY ("manager_id") REFERENCES "User"("id") ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS "p_features" (
  "id" BIGSERIAL PRIMARY KEY,
  "feature" TEXT NOT NULL
);
COMMENT ON TABLE "p_features" IS 'Security guard, Parking, Gym, ...';
ALTER TABLE "p_features" ADD CONSTRAINT "p_features_feature_unique" UNIQUE ("feature");
INSERT INTO p_features (feature) VALUES 
('p-feature_security'),
('p-feature_fire-alarm'),
('p-feature_gym'),
('p-feature_fitness-center'),
('p-feature_swimming-pool'),
('p-feature_community-rooms'),
('p-feature_public-library'),
('p-feature_parking'),
('p-feature_outdoor-common-area'),
('p-feature_services'),
('p-feature_facilities');

CREATE TABLE IF NOT EXISTS "property_features" (
  "property_id" UUID NOT NULL,
  "feature_id" BIGINT NOT NULL,
  "description" TEXT,

  PRIMARY KEY("property_id", "feature_id")
);
ALTER TABLE "property_features" ADD CONSTRAINT "property_features_property_id_fkey" FOREIGN KEY ("property_id") REFERENCES "properties"("id") ON DELETE CASCADE;
ALTER TABLE "property_features" ADD CONSTRAINT "property_features_feature_id_fkey" FOREIGN KEY ("feature_id") REFERENCES "p_features"("id") ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS "property_tags" (
  "id" BIGSERIAL PRIMARY KEY,
  "property_id" UUID NOT NULL,
  "tag" TEXT NOT NULL
);
ALTER TABLE "property_tags" ADD CONSTRAINT "property_id_tag_fkey" FOREIGN KEY ("property_id") REFERENCES "properties"("id") ON DELETE CASCADE;

CREATE TYPE "MEDIATYPE" AS ENUM (
  'IMAGE',
  'VIDEO'
);
CREATE TABLE IF NOT EXISTS "property_media" (
  "id" BIGSERIAL PRIMARY KEY,
  "property_id" UUID NOT NULL,
  "url" TEXT NOT NULL,
  "type" "MEDIATYPE" NOT NULL,
  "description" TEXT
);
ALTER TABLE "property_media" ADD CONSTRAINT "property_id_media_fkey" FOREIGN KEY ("property_id") REFERENCES "properties"("id") ON DELETE CASCADE;
ALTER TABLE "properties" ADD CONSTRAINT "property_primary_image_fkey" FOREIGN KEY ("primary_image") REFERENCES "property_media"("id") ON DELETE SET NULL;

CREATE TABLE IF NOT EXISTS "new_property_manager_requests" (
  "id" BIGSERIAL PRIMARY KEY,
  "creator_id" UUID NOT NULL,
  "property_id" UUID NOT NULL,
  "user_id" UUID,
  "email" TEXT NOT NULL,
  "approved" BOOLEAN NOT NULL DEFAULT FALSE,
  "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
ALTER TABLE "new_property_manager_requests" ADD CONSTRAINT "new_property_manager_requests_creator_fkey" FOREIGN KEY ("creator_id") REFERENCES "User"("id") ON DELETE CASCADE;
ALTER TABLE "new_property_manager_requests" ADD CONSTRAINT "new_property_manager_requests_property_fkey" FOREIGN KEY ("property_id") REFERENCES "properties"("id") ON DELETE CASCADE;

CREATE TYPE "PROPERTYVERIFICATIONSTATUS" AS ENUM (
  'PENDING',
  'APPROVED',
  'REJECTED'
);
CREATE TABLE IF NOT EXISTS "property_verification_requests" (
  "id" BIGSERIAL PRIMARY KEY,
  "creator_id" UUID NOT NULL,
  "property_id" UUID NOT NULL,
  "video_url" TEXT NOT NULL,
  "house_ownership_certificate" TEXT,
  "certificate_of_landuse_right" TEXT,
  "front_idcard" TEXT NOT NULL,
  "back_idcard" TEXT NOT NULL,
  "note" TEXT,
  "feedback" TEXT,
  "status" "PROPERTYVERIFICATIONSTATUS" NOT NULL DEFAULT 'PENDING',
  "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
ALTER TABLE "property_verification_requests" ADD CONSTRAINT "property_verification_requests_property_fkey" FOREIGN KEY ("property_id") REFERENCES "properties"("id") ON DELETE CASCADE;
ALTER TABLE "property_verification_requests" ADD CONSTRAINT "property_verification_requests_creator_fkey" FOREIGN KEY ("creator_id") REFERENCES "User"("id") ON DELETE CASCADE;

END;
