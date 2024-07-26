CREATE TABLE "permissions" (
  "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
  "name" VARCHAR(255) UNIQUE NOT NULL,
  "description" TEXT NOT NULL DEFAULT '',
  "created_by" UUID REFERENCES "users" ("id") ON DELETE CASCADE,
  "updated_by" UUID REFERENCES "users" ("id") ON DELETE CASCADE,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (now()),
  "updated_at" TIMESTAMPTZ NOT NULL DEFAULT (now())
);


-- NOTE: Initial Permissions
INSERT INTO "permissions" 
  (name, description)
VALUES 
  ('create:permissions', 'creation of new permissions'),
  ('read:permissions', 'reading permissions'),
  ('update:permissions', 'updating permissions'),
  ('delete:permissions', 'deletion of permissions'),

  ('create:roles', 'creation of new roles'),
  ('read:roles', 'reading roles'),
  ('update:roles', 'updating roles'),
  ('delete:roles', 'deletion of roles'),

  ('create:users', 'creation of new users'),
  ('read:users', 'users'),
  ('update:users', 'updating users'),
  ('delete:users', 'deletion of users');


