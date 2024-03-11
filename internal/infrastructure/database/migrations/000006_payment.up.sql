BEGIN;

CREATE TYPE "PAYMENTSTATUS" AS ENUM (
  'PENDING', 'SUCCESS', 'FAILED'
);

CREATE TABLE IF NOT EXISTS "payments" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" UUID NOT NULL,
  "order_id" TEXT NOT NULL,
  "order_info" TEXT NOT NULL,
  "amount" BIGINT NOT NULL,
  "status" "PAYMENTSTATUS" NOT NULL DEFAULT 'PENDING',
  "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
ALTER TABLE "payments" ADD CONSTRAINT "fk_payments_user_id" FOREIGN KEY ("user_id") REFERENCES "User" ("id") ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS "payment_items" (
  "payment_id" BIGINT NOT NULL,
  "name" TEXT NOT NULL,
  "price" BIGINT NOT NULL,
  "quantity" INTEGER NOT NULL,
  "discount" INTEGER NOT NULL
);
ALTER TABLE "payment_items" ADD CONSTRAINT "fk_payment_items_payment_id" FOREIGN KEY ("payment_id") REFERENCES "payments" ("id") ON DELETE CASCADE;

END;
