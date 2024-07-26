CREATE TABLE "role_permissions" (
    "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
    "role_id" UUID UNIQUE NOT NULL REFERENCES "roles" ("id") ON DELETE CASCADE,
    "permission_id" UUID UNIQUE NOT NULL REFERENCES "permissions" ("id") ON DELETE CASCADE
);

-- NOTE: Initial Roles
INSERT INTO "roles"
  (name, description)
VALUES 
  ('ADMIN', 'Admin role with all permissions, cannot manage other admins and master'),
  ('USER', 'User role with basic permissions, can only access its own resources');

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


