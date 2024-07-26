package users

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/sirjager/goth/repository/users/sqlc"
)

type UserRolesResult struct {
	Error      error
	Roles      []sqlc.Role
	StatusCode int
}

func (r *UserRepo) UserRoles(ctx context.Context, userID string) (res UserRolesResult) {
	roles, err := r.store.UserRoles(ctx, uuid.MustParse(userID))
	if err != nil {
		res.Error = err
		res.StatusCode = http.StatusInternalServerError
		return
	}

	res.StatusCode = http.StatusOK
	res.Roles = roles
	return
}
