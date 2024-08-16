CREATE TABLE "users" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),

  "email" VARCHAR(255) UNIQUE NOT NULL,
  "username" VARCHAR(255) UNIQUE NOT NULL,
  "password" TEXT NOT NULL DEFAULT '',
  "verified" BOOL NOT NULL DEFAULT false,
  "blocked" BOOL NOT NULL DEFAULT false,

  "provider" VARCHAR(255) NOT NULL DEFAULT '',
  "google_id" VARCHAR(255) NOT NULL DEFAULT '',

  "name" VARCHAR(255) NOT NULL DEFAULT '',
  "first_name" VARCHAR(255) NOT NULL DEFAULT '',
  "last_name" VARCHAR(255) NOT NULL DEFAULT '',
  "nick_name" VARCHAR(255) NOT NULL DEFAULT '',

  "avatar_url" TEXT NOT NULL DEFAULT '',
  "picture_url" TEXT NOT NULL DEFAULT '',
  "location" TEXT NOT NULL DEFAULT '',

  "master" BOOL NOT NULL DEFAULT FALSE,
  CHECK (master IN (TRUE, FALSE)),

  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);


CREATE INDEX ON "users" ("created_at");

CREATE UNIQUE INDEX "only_one_master" ON "users" (master) WHERE master = TRUE;

CREATE TRIGGER trg_update_updated_at BEFORE UPDATE ON "users"
FOR EACH ROW WHEN (OLD.id = NEW.id) EXECUTE FUNCTION fn_update_timestamp();
