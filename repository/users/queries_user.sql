-- name: UserRead :one
select * from "users" where id = @id limit 1;

-- name: UserReadByEmail :one
select * from "users" where email = $1 limit 1;

-- name: UserCreate :one
INSERT INTO "users" (
  id, email, verified, blocked,
  provider,google_id,
  name,first_name,last_name,nick_name,
  avatar_url,location,
  created_at,updated_at
) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING *;

-- name: UserDelete :one
delete from "users" where id = $1 RETURNING id;
