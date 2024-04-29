BEGIN;

DROP FUNCTION IF EXISTS "plan_rental_payments";
DROP FUNCTION IF EXISTS "get_nearest_payment_cycle";
DROP FUNCTION IF EXISTS "calculate_rental_fee";
DROP TABLE IF EXISTS "rental_payments";
DROP TYPE IF EXISTS "RENTALPAYMENTSTATUS";

END;
