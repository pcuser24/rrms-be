BEGIN;
-- CREATE TYPE "REMINDERRECURRENCEMODE" AS ENUM (
--   'NONE', 'WEEKLY', 'MONTHLY'
-- );
CREATE TABLE IF NOT EXISTS "reminders" (
  "id" BIGSERIAL PRIMARY KEY,
  "creator_id" UUID NOT NULL,
  "title" TEXT NOT NULL,
  "start_at" TIMESTAMPTZ NOT NULL,
  "end_at" TIMESTAMPTZ NOT NULL,
  "note" TEXT,
  "location" TEXT NOT NULL,
  -- "recurrence_day" INT,
  -- "recurrence_month" INT,
  -- "recurrence_mode" "REMINDERRECURRENCEMODE" NOT NULL DEFAULT 'NONE',
  -- "priority" INT NOT NULL DEFAULT 0,
  -- "resource_tag" TEXT NOT NULL,
  "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
ALTER TABLE "reminders" ADD CONSTRAINT "fk_reminders_creator_id" FOREIGN KEY ("creator_id") REFERENCES "User" ("id") ON DELETE CASCADE;
CREATE INDEX IF NOT EXISTS "idx_reminders_resource_tag" ON "reminders" ("resource_tag");
-- COMMENT ON COLUMN "reminders"."recurrence_day" IS '7-bit integer representing days in a week (0-6) when the reminder should be triggered. 0 is Sunday, 1 is Monday, and so on.';
-- COMMENT ON COLUMN "reminders"."recurrence_month" IS '32-bit integer representing days in a month (0-30) when the reminder should be triggered. 0 is the last day of the month, 1 is the first day of the month, and so on.';

END;
