// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.19.1
// source: queries_user.sql

package sqlc

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const userCreate = `-- name: UserCreate :one
INSERT INTO "users" (
  id, email, verified, blocked,
  provider,google_id,
  name,first_name,last_name,nick_name,
  avatar_url,location,
  created_at,updated_at
) 
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING id, email, verified, blocked, provider, google_id, name, first_name, last_name, nick_name, avatar_url, location, created_at, updated_at
`

type UserCreateParams struct {
	ID        uuid.UUID `json:"id"`
	Email     string    `json:"email"`
	Verified  bool      `json:"verified"`
	Blocked   bool      `json:"blocked"`
	Provider  string    `json:"provider"`
	GoogleID  string    `json:"google_id"`
	Name      string    `json:"name"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	NickName  string    `json:"nick_name"`
	AvatarUrl string    `json:"avatar_url"`
	Location  string    `json:"location"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (q *Queries) UserCreate(ctx context.Context, arg UserCreateParams) (User, error) {
	row := q.db.QueryRow(ctx, userCreate,
		arg.ID,
		arg.Email,
		arg.Verified,
		arg.Blocked,
		arg.Provider,
		arg.GoogleID,
		arg.Name,
		arg.FirstName,
		arg.LastName,
		arg.NickName,
		arg.AvatarUrl,
		arg.Location,
		arg.CreatedAt,
		arg.UpdatedAt,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Verified,
		&i.Blocked,
		&i.Provider,
		&i.GoogleID,
		&i.Name,
		&i.FirstName,
		&i.LastName,
		&i.NickName,
		&i.AvatarUrl,
		&i.Location,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const userDelete = `-- name: UserDelete :one
delete from "users" where id = $1 RETURNING id
`

func (q *Queries) UserDelete(ctx context.Context, id uuid.UUID) (uuid.UUID, error) {
	row := q.db.QueryRow(ctx, userDelete, id)
	err := row.Scan(&id)
	return id, err
}

const userRead = `-- name: UserRead :one
select id, email, verified, blocked, provider, google_id, name, first_name, last_name, nick_name, avatar_url, location, created_at, updated_at from "users" where id = $1 limit 1
`

func (q *Queries) UserRead(ctx context.Context, id uuid.UUID) (User, error) {
	row := q.db.QueryRow(ctx, userRead, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Verified,
		&i.Blocked,
		&i.Provider,
		&i.GoogleID,
		&i.Name,
		&i.FirstName,
		&i.LastName,
		&i.NickName,
		&i.AvatarUrl,
		&i.Location,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const userReadByEmail = `-- name: UserReadByEmail :one
select id, email, verified, blocked, provider, google_id, name, first_name, last_name, nick_name, avatar_url, location, created_at, updated_at from "users" where email = $1 limit 1
`

func (q *Queries) UserReadByEmail(ctx context.Context, email string) (User, error) {
	row := q.db.QueryRow(ctx, userReadByEmail, email)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Email,
		&i.Verified,
		&i.Blocked,
		&i.Provider,
		&i.GoogleID,
		&i.Name,
		&i.FirstName,
		&i.LastName,
		&i.NickName,
		&i.AvatarUrl,
		&i.Location,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
