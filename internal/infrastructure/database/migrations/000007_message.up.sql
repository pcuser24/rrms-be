BEGIN;

CREATE TABLE IF NOT EXISTS "msg_groups" (
  "group_id" BIGSERIAL PRIMARY KEY,
  "name" TEXT NOT NULL,
  "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  "created_by" UUID NOT NULL
);
ALTER TABLE "msg_groups" ADD CONSTRAINT "fk_msg_groups_created_by" FOREIGN KEY ("created_by") REFERENCES "User" ("id") ON DELETE CASCADE;
-- CREATE INDEX IF NOT EXISTS "idx_msg_groups_name" ON "msg_groups" ("name");

CREATE TABLE IF NOT EXISTS "msg_group_members" (
  "group_id" BIGINT NOT NULL,
  "user_id" UUID NOT NULL
);
ALTER TABLE "msg_group_members" ADD CONSTRAINT "pk_group_members" PRIMARY KEY ("group_id", "user_id");
ALTER TABLE "msg_group_members" ADD CONSTRAINT "fk_group_members_group_id" FOREIGN KEY ("group_id") REFERENCES "msg_groups" ("group_id") ON DELETE CASCADE;
ALTER TABLE "msg_group_members" ADD CONSTRAINT "fk_group_members_user_id" FOREIGN KEY ("user_id") REFERENCES "User" ("id") ON DELETE CASCADE;

CREATE TYPE "MESSAGESTATUS" AS ENUM (
  'ACTIVE', 'DELETED'
);
CREATE TYPE "MESSAGETYPE" AS ENUM (
  'TEXT', 'IMAGE', 'FILE'
);
CREATE TABLE IF NOT EXISTS "messages" (
  "id" BIGSERIAL PRIMARY KEY,
  "group_id" BIGINT NOT NULL,
  "from_user" UUID NOT NULL,
  "content" TEXT NOT NULL,
  "status" "MESSAGESTATUS" NOT NULL DEFAULT 'ACTIVE',
  "type" "MESSAGETYPE" NOT NULL DEFAULT 'TEXT',
  "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL
);
ALTER TABLE "messages" ADD CONSTRAINT "fk_messages_group_id" FOREIGN KEY ("group_id") REFERENCES "msg_groups" ("group_id") ON DELETE CASCADE;
ALTER TABLE "messages" ADD CONSTRAINT "fk_messages_from" FOREIGN KEY ("from_user") REFERENCES "User" ("id") ON DELETE CASCADE;

END;
