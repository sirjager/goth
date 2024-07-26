CREATE TABLE "roles_permissions" (
    "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
    "role_id" UUID UNIQUE NOT NULL REFERENCES "roles" ("id") ON DELETE CASCADE,
    "permission_id" UUID UNIQUE NOT NULL REFERENCES "permissions" ("id") ON DELETE CASCADE
);