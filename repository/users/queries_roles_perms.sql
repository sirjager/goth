-- name: UserReadMaster :one 
SELECT * FROM "users" where master = true LIMIT 1;

-- name: UserRoles :many
SELECT  r.* FROM "users_roles" ur
JOIN "roles" r ON ur.role_id = r.id 
WHERE ur.user_id = $1;

-- name: RolePermissions :many
SELECT  p.* FROM "roles_permissions" rp
JOIN "permissions" r ON rp.permission_id = r.id 
WHERE rp.role_id = $1;

-- name: UserPermissions :many
SELECT p.*
FROM "users_roles" ur
JOIN "roles_permissions" rp ON ur.role_id = rp.role_id
JOIN "permissions" p ON rp.permission_id = p.id
WHERE ur.user_id = $1;
