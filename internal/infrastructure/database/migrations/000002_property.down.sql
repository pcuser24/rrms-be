BEGIN;
DROP TABLE IF EXISTS "property_tags";
DROP TABLE IF EXISTS "property_managers";
DROP TABLE IF EXISTS "property_features";
DROP TABLE IF EXISTS "p_features";
ALTER TABLE "properties" DROP CONSTRAINT IF EXISTS "property_primary_image_fkey";
DROP TABLE IF EXISTS "property_media";
DROP TABLE IF EXISTS "new_property_manager_requests";
DROP TABLE IF EXISTS "properties";

DROP TYPE IF EXISTS "PROPERTYTYPE";
DROP TYPE IF EXISTS "MEDIATYPE";
END;
