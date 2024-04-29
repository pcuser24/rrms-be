BEGIN;

CREATE TYPE "RENTALPAYMENTSTATUS" AS ENUM ('PLAN', 'ISSUED', 'PENDING', 'REQUEST2PAY', 'PAID', 'CANCELLED');

CREATE TABLE IF NOT EXISTS "rental_payments" (
  "id" BIGSERIAL PRIMARY KEY,
  "code" VARCHAR(50) NOT NULL,
  "rental_id" BIGINT NOT NULL,
  "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  "start_date" DATE NOT NULL,
  "end_date" DATE NOT NULL,
  "expiry_date" DATE,
  "payment_date" DATE,
  "updated_by" UUID,
  "status" "RENTALPAYMENTSTATUS" NOT NULL DEFAULT 'PLAN',
  "amount" REAL NOT NULL CHECK (amount >= 0),
  "discount" REAL CHECK (discount >= 0),
  CHECK (amount >= discount),
  "note" TEXT,

  UNIQUE ("code")
);
ALTER TABLE "rental_payments" ADD CONSTRAINT "fk_rental_payments_rental_id" FOREIGN KEY ("rental_id") REFERENCES "rentals" ("id") ON DELETE CASCADE;
ALTER TABLE "rental_payments" ADD CONSTRAINT "fk_rental_payments_updated_by" FOREIGN KEY ("updated_by") REFERENCES "User" ("id") ON DELETE SET NULL;
COMMENT ON COLUMN "rental_payments"."code" IS '{payment.id}_{ELECTRICITY | WATER | RENTAL | DEPOSIT | SERVICES{id}}_{payment.created_at}';
COMMENT ON COLUMN "rental_payments"."payment_date" IS 'the date the payment gets paid';

-- helper function to calculate rental fee of a billing cycle
CREATE OR REPLACE FUNCTION calculate_rental_fee(start_date DATE, end_date DATE, basis INT, price REAL)
RETURNS REAL AS $$
DECLARE
  rental_duration INT;
  basis_in_days INT;
BEGIN
  rental_duration := (end_date - start_date);

  IF start_date + basis * INTERVAL '1 month' > end_date THEN
    RETURN (price * (rental_duration::NUMERIC / (basis * 30)))::REAL;
  ELSE
    RETURN price;
  END IF;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION get_nearest_payment_cycle(start_date DATE, end_date DATE, basis INT, is_prepaid BOOLEAN)
RETURNS DATE AS $$
DECLARE
  cycle DATE;
  next_cycle DATE;
BEGIN
  IF is_prepaid THEN
    cycle := start_date;
  ELSE
    cycle := start_date + basis * INTERVAL '1 month';
    IF cycle > end_date THEN
      cycle = start_date
    END IF;
  END IF;
  LOOP
    next_cycle := cycle + basis * INTERVAL '1 month';
    EXIT WHEN next_cycle >= end_date;
    cycle := next_cycle;
  END LOOP;
  RETURN cycle;
END;
$$ LANGUAGE plpgsql;

-- function to generate rental payments, called on daily basis by cronjob
CREATE OR REPLACE FUNCTION plan_rental_payments() 
RETURNS SETOF BIGINT AS
$BODY$
DECLARE
  rental_id BIGINT;
BEGIN
  FOR rental_id IN
    SELECT "id" FROM "rentals" WHERE (rentals.start_date + INTERVAL '1 month' * rentals.rental_period) >= CURRENT_DATE
  LOOP 
    RETURN QUERY SELECT * FROM plan_rental_payment(rental_id);
  END LOOP;
END;
$BODY$ LANGUAGE plpgsql;

CREATE OR REPLACE FUNCTION plan_rental_payment(rental_id BIGINT) 
RETURNS SETOF BIGINT AS 
$BODY$
DECLARE
  rental_record RECORD;
  payment_id BIGINT;
  start_date DATE;
  end_date DATE;
  nearest_cycle DATE;
  payment_code VARCHAR(50);
  amount NUMERIC;
  rental_service RECORD;
BEGIN
  SELECT "id", "movein_date", "rental_period", "rental_payment_basis", "rental_price", "payment_type", "electricity_setup_by", "electricity_payment_type", "electricity_price", "water_setup_by", "water_payment_type", "water_price", (rentals.start_date + INTERVAL '1 month' * rentals.rental_period) AS expiry_date INTO rental_record FROM "rentals" WHERE id = rental_id;
  
  -- plan rental payment
  nearest_cycle := get_nearest_payment_cycle(rental_record.movein_date, CURRENT_DATE, rental_record.rental_payment_basis, rental_record.payment_type = 'PREPAID');
  IF nearest_cycle != rental_record.movein_date THEN
    IF rental_record.payment_type = 'PREPAID' THEN 
      start_date := nearest_cycle;
      end_date := start_date + INTERVAL '1 month' * rental_record.rental_payment_basis;
      IF end_date > rental_record.expiry_date THEN
        end_date := rental_record.expiry_date;
      END IF;
    ELSE
      start_date := nearest_cycle - INTERVAL '1 month' * rental_record.rental_payment_basis;
      if start_date < rental_record.movein_date THEN
        start_date := rental_record.movein_date;
      END IF;
      end_date = nearest_cycle;
    END IF;
    amount := calculate_rental_fee(start_date, end_date, rental_record.rental_payment_basis, rental_record.rental_price);
    payment_code := rental_record.id || '_RENTAL_' || LPAD(EXTRACT(MONTH FROM start_date)::TEXT, 2, '0') || EXTRACT(YEAR FROM start_date)|| LPAD(EXTRACT(MONTH FROM end_date)::TEXT, 2, '0') || EXTRACT(YEAR FROM end_date) || '_A';
    SELECT id FROM "rental_payments" INTO payment_id WHERE "code" = payment_code;
    IF not found THEN
      INSERT INTO "rental_payments" ("code", "rental_id", "status", "amount", "start_date", "end_date") VALUES (payment_code, rental_record.id, 'PLAN', amount, start_date, end_date) RETURNING id INTO payment_id;
      RETURN NEXT payment_id;
    END IF;
  END IF;
  -- plan service payments
  nearest_cycle := get_nearest_payment_cycle(rental_record.movein_date, CURRENT_DATE, 1, FALSE);
  IF nearest_cycle = rental_record.movein_date THEN
    RETURN;
  END IF;
  start_date := nearest_cycle - INTERVAL '1 month';
  end_date = nearest_cycle;
  -- plan electricity payment
  IF rental_record.electricity_setup_by = 'LANDLORD' THEN
  payment_code := rental_record.id || '_ELECTRICITY_' || LPAD(EXTRACT(MONTH FROM start_date)::TEXT, 2, '0') || EXTRACT(YEAR FROM start_date)|| LPAD(EXTRACT(MONTH FROM end_date)::TEXT, 2, '0') || EXTRACT(YEAR FROM end_date) || '_A';
  SELECT id FROM "rental_payments" INTO payment_id WHERE "code" = payment_code LIMIT 1;
  IF not found THEN
    INSERT INTO "rental_payments" ("code", "rental_id", "status", "amount", "start_date", "end_date") VALUES (payment_code, rental_record.id, 'PLAN', 0, start_date, end_date) RETURNING id INTO payment_id;
    RETURN NEXT payment_id;
  END IF; 
  END IF; 
  -- plan water payment
  IF rental_record.water_setup_by = 'LANDLORD' THEN
  payment_code := rental_record.id || '_WATER_' || LPAD(EXTRACT(MONTH FROM start_date)::TEXT, 2, '0') || EXTRACT(YEAR FROM start_date)|| LPAD(EXTRACT(MONTH FROM end_date)::TEXT, 2, '0') || EXTRACT(YEAR FROM end_date) || '_A';
  SELECT id FROM "rental_payments" INTO payment_id WHERE "code" = payment_code LIMIT 1;
  IF not found THEN
    INSERT INTO "rental_payments" ("code", "rental_id", "status", "amount", "start_date", "end_date") VALUES (payment_code, rental_record.id, 'PLAN', 0, start_date, end_date) RETURNING id INTO payment_id;
    RETURN NEXT payment_id;
  END IF; 
  END IF;
  -- plan service payments
  FOR rental_service IN
    SELECT "id", "name", "setup_by", "provider", "price" FROM "rental_services" WHERE "rental_services"."rental_id" = rental_record.id AND "rental_services"."setup_by" = 'LANDLORD'
  LOOP
    CONTINUE WHEN rental_service.setup_by = 'TENANT';
    payment_code := rental_record.id || '_SERVICE_' || rental_service.id || '_' || LPAD(EXTRACT(MONTH FROM start_date)::TEXT, 2, '0') || EXTRACT(YEAR FROM start_date)|| LPAD(EXTRACT(MONTH FROM end_date)::TEXT, 2, '0') || EXTRACT(YEAR FROM end_date) || '_A';
    SELECT id FROM "rental_payments" INTO payment_id WHERE "code" = payment_code LIMIT 1;
    IF not found THEN
      amount := calculate_rental_fee(start_date, end_date, 1, rental_service.price);
      INSERT INTO "rental_payments" ("code", "rental_id", "status", "amount", "start_date", "end_date") VALUES (payment_code, rental_record.id, 'PLAN', amount, start_date, end_date) RETURNING id INTO payment_id;
      RETURN NEXT payment_id;
    END IF;
  END LOOP;
END;
$BODY$ LANGUAGE plpgsql;

END;
