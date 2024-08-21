package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"

	"github.com/sirjager/goth/repository/users"
)

type Repo interface {
	users.UserRepository
}

type repo struct {
	users.UserRepository
}

func NewRepository(conn *pgxpool.Pool, pgURL string, logger zerolog.Logger) (Repo, error) {
	users, err := users.NewUsersRepo(conn, pgURL, logger)
	if err != nil {
		return nil, err
	}

	return &repo{users}, nil
}
