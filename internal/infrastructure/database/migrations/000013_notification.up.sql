BEGIN;

-- CREATE TYPE "NOTIFICATIONTYPE" AS ENUM ('SYSTEM', 'ALERT', 'REMINDER', 'PROMOTION', 'NEWS');
-- CREATE TABLE IF NOT EXISTS "user_notification_preferences" (
--   "user_id" UUID NOT NULL,
--   "type" "NOTIFICATIONTYPE" NOT NULL
-- );
-- ALTER TABLE "user_notification_preferences" ADD CONSTRAINT "user_notification_fk" FOREIGN KEY ("user_id") REFERENCES "User" ("id") ON DELETE CASCADE;

CREATE TYPE "PLATFORM" AS ENUM ('WEB', 'IOS', 'ANDROID');
CREATE TABLE IF NOT EXISTS "user_notification_devices" (
  "user_id" UUID NOT NULL,
  "session_id" UUID NOT NULL,
  "token" TEXT NOT NULL,
  "platform" "PLATFORM" NOT NULL,
  "last_accessed" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,

  UNIQUE("user_id", "session_id")
);
ALTER TABLE "user_notification_devices" ADD CONSTRAINT "user_notification_tokens_fk" FOREIGN KEY ("user_id") REFERENCES "User" ("id") ON DELETE CASCADE;
COMMENT ON TABLE "user_notification_devices" IS 'This table stores the devices that the user uses to receive push notifications.';

CREATE TYPE "NOTIFICATIONCHANNEL" AS ENUM ('EMAIL', 'PUSH');
CREATE TABLE IF NOT EXISTS "notifications" (
  "id" BIGSERIAL PRIMARY KEY,
  "user_id" UUID,
  "title" TEXT NOT NULL,
  "content" TEXT NOT NULL,
  "data" JSONB NOT NULL DEFAULT '{}'::JSONB,
  "seen" BOOLEAN DEFAULT FALSE NOT NULL,
  "target" TEXT NOT NULL,
  "channel" "NOTIFICATIONCHANNEL" NOT NULL,
  "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
ALTER TABLE "notifications" ADD CONSTRAINT "notifications_user_id_fkey" FOREIGN KEY ("user_id") REFERENCES "User" ("id") ON DELETE CASCADE;
COMMENT ON COLUMN "notifications"."target" IS 'The target of the notification. It can be an email, notification token, phone number.';

END;
