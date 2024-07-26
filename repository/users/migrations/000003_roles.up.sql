CREATE TABLE "roles" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "name" VARCHAR(255) UNIQUE NOT NULL,
  "description" TEXT NOT NULL DEFAULT '',
  "created_by" UUID REFERENCES "users" ("id") ON DELETE CASCADE,
  "updated_by" UUID REFERENCES "users" ("id") ON DELETE CASCADE,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);
