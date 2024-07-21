CREATE TABLE "users" (
  "id" UUID PRIMARY KEY,

  "email" VARCHAR(255) UNIQUE NOT NULL,
  "verified" BOOL NOT NULL DEFAULT false,
  "blocked" BOOL NOT NULL DEFAULT false,

  "provider" VARCHAR(255) NOT NULL DEFAULT '',
  "google_id" VARCHAR(255) NOT NULL DEFAULT '',

  "name" VARCHAR(255) NOT NULL DEFAULT '',
  "first_name" VARCHAR(255) NOT NULL DEFAULT '',
  "last_name" VARCHAR(255) NOT NULL DEFAULT '',
  "nick_name" VARCHAR(255) NOT NULL DEFAULT '',

  "avatar_url" TEXT NOT NULL DEFAULT '',
  "location" TEXT NOT NULL DEFAULT '',

  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);

CREATE INDEX ON "users" ("created_at");

CREATE TRIGGER trg_update_updated_at BEFORE UPDATE ON "users"
FOR EACH ROW WHEN (OLD.id = NEW.id) EXECUTE FUNCTION fn_update_timestamp();
