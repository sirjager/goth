CREATE TABLE "users_roles" (
    "id" UUID PRIMARY KEY DEFAULT (uuid_generate_v4()),
    "role_id" UUID UNIQUE NOT NULL REFERENCES "roles" ("id") ON DELETE CASCADE,
    "user_id" UUID UNIQUE NOT NULL REFERENCES "users" ("id") ON DELETE CASCADE
);
