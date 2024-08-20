-- name: UsersRead :many
select * from "users" limit sqlc.narg('limit') offset sqlc.narg('offset');

-- name: UserRead :one
select * from "users" where id = @id limit 1;

-- name: UserReadByEmail :one
select * from "users" where email = $1 limit 1;

-- name: UserReadByUsername :one
select * from "users" where username = $1 limit 1;

-- name: UserCreate :one
INSERT INTO "users" (
  id,email,username,password,verified,blocked,
  provider,google_id,
  full_name,first_name,last_name,
  avatar_url,picture_url,master,
  created_at,updated_at
) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) RETURNING *;

-- name: UserDelete :one
DELETE from "users" WHERE id = $1 RETURNING id;

-- name: UserUpdate :one
UPDATE "users" SET
  full_name = $1,
  first_name = $2, 
  last_name = $3,
  picture_url = $4,
  avatar_url = $5
WHERE id = @id RETURNING *;


-- name: UserReadMaster :one 
SELECT * FROM "users" where master = true LIMIT 1;

-- name: UserUpdateVerified :one
UPDATE "users" SET verified = $1 WHERE id = @id RETURNING *;
