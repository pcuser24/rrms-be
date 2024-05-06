BEGIN;

CREATE TYPE "RENTALCOMPLAINTSTATUS" AS ENUM ('PENDING', 'RESOLVED', 'CLOSED');
CREATE TYPE "RENTALCOMPLAINTTYPE" AS ENUM ('REPORT', 'SUGGESTION');
CREATE TABLE IF NOT EXISTS "rental_complaints" (
  "id" BIGSERIAL PRIMARY KEY,
  "rental_id" BIGINT NOT NULL,
  "creator_id" UUID NOT NULL,
  "title" TEXT NOT NULL,
  "content" TEXT NOT NULL,
  "suggestion" TEXT,
  "media" TEXT[],
  "occurred_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  "updated_by" UUID NOT NULL,
  "type" "RENTALCOMPLAINTTYPE" NOT NULL,
  "status" "RENTALCOMPLAINTSTATUS" NOT NULL DEFAULT 'PENDING'
);
ALTER TABLE "rental_complaints" ADD CONSTRAINT "fk_rental_complaints_rental_id" FOREIGN KEY ("rental_id") REFERENCES "rentals" ("id") ON DELETE CASCADE;
ALTER TABLE "rental_complaints" ADD CONSTRAINT "fk_rental_complaints_creator_id" FOREIGN KEY ("creator_id") REFERENCES "User" ("id") ON DELETE CASCADE;

CREATE TABLE IF NOT EXISTS "rental_complaint_replies" (
  "complaint_id" BIGINT NOT NULL,
  "replier_id" UUID NOT NULL,
  "reply" TEXT NOT NULL,
  "media" TEXT[],
  "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
ALTER TABLE "rental_complaint_replies" ADD CONSTRAINT "fk_rental_complaint_replies_complaint_id" FOREIGN KEY ("complaint_id") REFERENCES "rental_complaints" ("id") ON DELETE CASCADE;
ALTER TABLE "rental_complaint_replies" ADD CONSTRAINT "fk_rental_complaint_replies_replier_id" FOREIGN KEY ("replier_id") REFERENCES "User" ("id") ON DELETE CASCADE;

END;
