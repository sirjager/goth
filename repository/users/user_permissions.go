package users

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"github.com/sirjager/goth/repository/users/sqlc"
)

type UserPermissionsResult struct {
	Error       error
	Permissions []sqlc.Permission
	StatusCode  int
}

func (r *UserRepo) UserPermissions(ctx context.Context, userID string) (res UserPermissionsResult) {
	perms, err := r.store.UserPermissions(ctx, uuid.MustParse(userID))
	if err != nil {
		res.Error = err
		res.StatusCode = http.StatusInternalServerError
		return
	}

	res.StatusCode = http.StatusOK
	res.Permissions = perms
	return
}
