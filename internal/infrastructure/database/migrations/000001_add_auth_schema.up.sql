BEGIN;
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE IF NOT EXISTS "User" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4() ,
  "email" VARCHAR(45) NOT NULL,
  "password" VARCHAR(200) DEFAULT NULL,
  "group_id" UUID DEFAULT NULL,
  "created_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  "updated_at" TIMESTAMPTZ DEFAULT NOW() NOT NULL,
  "created_by" UUID DEFAULT NULL,
  "updated_by" UUID DEFAULT NULL,
  "deleted_f" BOOL DEFAULT FALSE NOT NULL,

  UNIQUE ("email")
);
COMMENT ON TABLE "User" IS 'Bang user';
COMMENT ON COLUMN "User".deleted_f IS '1: deleted, 0: not deleted';

CREATE TABLE "Account" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "userId" UUID NOT NULL,
  "type" TEXT NOT NULL,
  "provider" TEXT NOT NULL,
  "providerAccountId" TEXT NOT NULL,
  "refresh_token" TEXT,
  "access_token" TEXT,
  "expires_at" INTEGER,
  "token_type" TEXT,
  "scope" TEXT,
  "id_token" TEXT,
  "session_state" TEXT,

  CONSTRAINT "account_user_id_fkey" FOREIGN KEY ("userId") REFERENCES "User"("id") ON DELETE CASCADE,
  UNIQUE ("provider", "providerAccountId")
);

CREATE TABLE "Session" (
  "id" UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
  "sessionToken" TEXT NOT NULL,
  "userId" UUID NOT NULL,
  "expires" TIMESTAMPTZ NOT NULL,

  "user_agent" TEXT,
  "client_ip" TEXT,
  "is_blocked" BOOLEAN NOT NULL DEFAULT FALSE,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

  CONSTRAINT "session_user_id_fkey" FOREIGN KEY ("userId") REFERENCES "User"("id") ON DELETE CASCADE
);

CREATE TABLE "VerificationToken" (
  "identifier"  TEXT NOT NULL,
  "token"       TEXT NOT NULL,
  "expires"     TIMESTAMPTZ NOT NULL,

  UNIQUE("identifier", "token")
);
END;
