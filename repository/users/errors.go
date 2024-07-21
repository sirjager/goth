package users

import (
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func errorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}

// func isForeignKeyViolation(err error) bool {
// 	return errorCode(err) == "23503"
// }

func isUniqueViolation(err error) bool {
	return errorCode(err) == "23505"
}

func isRecordNotFound(err error) bool {
	return errors.Is(err, pgx.ErrNoRows)
}
