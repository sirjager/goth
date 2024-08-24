// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: queries_user.sql

package sqlc

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

const userCreate = `-- name: UserCreate :one
INSERT INTO "users" (
  id,email,username,password,verified,blocked,
  provider,google_id,
  full_name,first_name,last_name,
  avatar_url,picture_url,master,
  created_at,updated_at
) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16) RETURNING id, email, username, password, verified, blocked, provider, google_id, full_name, first_name, last_name, avatar_url, picture_url, master, created_at, updated_at
`

type UserCreateParams struct {
	ID         uuid.UUID `json:"id"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	Password   string    `json:"password"`
	Verified   bool      `json:"verified"`
	Blocked    bool      `json:"blocked"`
	Provider   string    `json:"provider"`
	GoogleID   string    `json:"google_id"`
	FullName   string    `json:"full_name"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	AvatarUrl  string    `json:"avatar_url"`
	PictureUrl string    `json:"picture_url"`
	Master     bool      `json:"master"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (q *Queries) UserCreate(ctx context.Context, arg UserCreateParams) (User, error) {
	row := q.db.QueryRow(ctx, userCreate,
		arg.ID,
		arg.Email,
		arg.Username,
		arg.Password,
		arg.Verified,
		arg.Blocked,
		arg.Provider,
		arg.GoogleID,
		arg.FullName,
		arg.FirstName,
		arg.LastName,
		arg.AvatarUrl,
		arg.PictureUrl,
		arg.Master,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Password,
		&i.Verified,
		&i.Blocked,
		&i.Provider,
		&i.GoogleID,
		&i.FullName,
		&i.FirstName,
		&i.LastName,
		&i.AvatarUrl,
		&i.PictureUrl,
		&i.Master,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const userDelete = `-- name: UserDelete :one
DELETE from "users" WHERE id = $1 RETURNING id
`

func (q *Queries) UserDelete(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, userDelete, id)
	err := row.Scan(&id)
	return id, err
}

const userRead = `-- name: UserRead :one
select id, email, username, password, verified, blocked, provider, google_id, full_name, first_name, last_name, avatar_url, picture_url, master, created_at, updated_at from "users" where id = $1 limit 1
`

func (q *Queries) UserRead(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, userRead, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Password,
		&i.Verified,
		&i.Blocked,
		&i.Provider,
		&i.GoogleID,
		&i.FullName,
		&i.FirstName,
		&i.LastName,
		&i.AvatarUrl,
		&i.PictureUrl,
		&i.Master,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const userReadByEmail = `-- name: UserReadByEmail :one
select id, email, username, password, verified, blocked, provider, google_id, full_name, first_name, last_name, avatar_url, picture_url, master, created_at, updated_at from "users" where email = $1 limit 1
`

func (q *Queries) UserReadByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, userReadByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Password,
		&i.Verified,
		&i.Blocked,
		&i.Provider,
		&i.GoogleID,
		&i.FullName,
		&i.FirstName,
		&i.LastName,
		&i.AvatarUrl,
		&i.PictureUrl,
		&i.Master,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const userReadByUsername = `-- name: UserReadByUsername :one
select id, email, username, password, verified, blocked, provider, google_id, full_name, first_name, last_name, avatar_url, picture_url, master, created_at, updated_at from "users" where username = $1 limit 1
`

func (q *Queries) UserReadByUsername(ctx context.Context, username string) (User, error) {
	row := q.db.QueryRow(ctx, userReadByUsername, username)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Password,
		&i.Verified,
		&i.Blocked,
		&i.Provider,
		&i.GoogleID,
		&i.FullName,
		&i.FirstName,
		&i.LastName,
		&i.AvatarUrl,
		&i.PictureUrl,
		&i.Master,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const userReadMaster = `-- name: UserReadMaster :one
SELECT id, email, username, password, verified, blocked, provider, google_id, full_name, first_name, last_name, avatar_url, picture_url, master, created_at, updated_at FROM "users" where master = true LIMIT 1
`

func (q *Queries) UserReadMaster(ctx context.Context) (User, error) {
	row := q.db.QueryRow(ctx, userReadMaster)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Password,
		&i.Verified,
		&i.Blocked,
		&i.Provider,
		&i.GoogleID,
		&i.FullName,
		&i.FirstName,
		&i.LastName,
		&i.AvatarUrl,
		&i.PictureUrl,
		&i.Master,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const userUpdate = `-- name: UserUpdate :one
UPDATE "users" SET
  full_name = $1,
  first_name = $2, 
  last_name = $3,
  picture_url = $4,
  avatar_url = $5,
  username = $6
WHERE id = $7 RETURNING id, email, username, password, verified, blocked, provider, google_id, full_name, first_name, last_name, avatar_url, picture_url, master, created_at, updated_at
`

type UserUpdateParams struct {
	FullName   string    `json:"full_name"`
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	PictureUrl string    `json:"picture_url"`
	AvatarUrl  string    `json:"avatar_url"`
	Username   string    `json:"username"`
	ID         uuid.UUID `json:"id"`
}

func (q *Queries) UserUpdate(ctx context.Context, arg UserUpdateParams) (User, error) {
	row := q.db.QueryRow(ctx, userUpdate,
		arg.FullName,
		arg.FirstName,
		arg.LastName,
		arg.PictureUrl,
		arg.AvatarUrl,
		arg.Username,
		arg.ID,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Password,
		&i.Verified,
		&i.Blocked,
		&i.Provider,
		&i.GoogleID,
		&i.FullName,
		&i.FirstName,
		&i.LastName,
		&i.AvatarUrl,
		&i.PictureUrl,
		&i.Master,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const userUpdatePassword = `-- name: UserUpdatePassword :one
UPDATE "users" SET password = $1 WHERE id = $2 RETURNING id, email, username, password, verified, blocked, provider, google_id, full_name, first_name, last_name, avatar_url, picture_url, master, created_at, updated_at
`

type UserUpdatePasswordParams struct {
	Password string    `json:"password"`
	ID       uuid.UUID `json:"id"`
}

func (q *Queries) UserUpdatePassword(ctx context.Context, arg UserUpdatePasswordParams) (User, error) {
	row := q.db.QueryRow(ctx, userUpdatePassword, arg.Password, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Password,
		&i.Verified,
		&i.Blocked,
		&i.Provider,
		&i.GoogleID,
		&i.FullName,
		&i.FirstName,
		&i.LastName,
		&i.AvatarUrl,
		&i.PictureUrl,
		&i.Master,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const userUpdateVerified = `-- name: UserUpdateVerified :one
UPDATE "users" SET verified = $1 WHERE id = $2 RETURNING id, email, username, password, verified, blocked, provider, google_id, full_name, first_name, last_name, avatar_url, picture_url, master, created_at, updated_at
`

type UserUpdateVerifiedParams struct {
	Verified bool      `json:"verified"`
	ID       uuid.UUID `json:"id"`
}

func (q *Queries) UserUpdateVerified(ctx context.Context, arg UserUpdateVerifiedParams) (User, error) {
	row := q.db.QueryRow(ctx, userUpdateVerified, arg.Verified, arg.ID)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Username,
		&i.Password,
		&i.Verified,
		&i.Blocked,
		&i.Provider,
		&i.GoogleID,
		&i.FullName,
		&i.FirstName,
		&i.LastName,
		&i.AvatarUrl,
		&i.PictureUrl,
		&i.Master,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const usersRead = `-- name: UsersRead :many
select id, email, username, password, verified, blocked, provider, google_id, full_name, first_name, last_name, avatar_url, picture_url, master, created_at, updated_at from "users" limit $2 offset $1
`

type UsersReadParams struct {
	Offset pgtype.Int4 `json:"offset"`
	Limit  pgtype.Int4 `json:"limit"`
}

func (q *Queries) UsersRead(ctx context.Context, arg UsersReadParams) ([]User, error) {
	rows, err := q.db.Query(ctx, usersRead, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []User{}
	for rows.Next() {
		var i User
		if err := rows.Scan(
			&i.ID,
			&i.Email,
			&i.Username,
			&i.Password,
			&i.Verified,
			&i.Blocked,
			&i.Provider,
			&i.GoogleID,
			&i.FullName,
			&i.FirstName,
			&i.LastName,
			&i.AvatarUrl,
			&i.PictureUrl,
			&i.Master,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
