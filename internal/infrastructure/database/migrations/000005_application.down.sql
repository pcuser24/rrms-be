BEGIN;

DROP TABLE IF EXISTS "application_units";
DROP TABLE IF EXISTS "application_minors";
DROP TABLE IF EXISTS "application_coaps";
DROP TABLE IF EXISTS "application_pets";
DROP TABLE IF EXISTS "application_vehicles";
DROP TABLE IF EXISTS "applications";

DROP TYPE IF EXISTS "APPLICATION_STATUS";
DROP TYPE IF EXISTS "TENANTTYPE";
END;
