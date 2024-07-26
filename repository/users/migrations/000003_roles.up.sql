CREATE TABLE "roles" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "name" VARCHAR(255) UNIQUE NOT NULL,
  "description" TEXT NOT NULL DEFAULT '',
  "created_by" UUID REFERENCES "users" ("id") ON DELETE CASCADE,
  "updated_by" UUID REFERENCES "users" ("id") ON DELETE CASCADE,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);


-- NOTE: Initial Roles
INSERT INTO "roles"
  (name, description)
VALUES 
  ('ADMIN', 'Admin role with all permissions, cannot manage other admins and master'),
  ('USER', 'User role with basic permissions, can only access its own resources');

