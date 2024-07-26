package mw

import (
	"context"
	"net/http"

	"github.com/sirjager/goth/entity"
	"github.com/sirjager/goth/repository"
)

const (
	ContextKeyRoles       contextType = "ctx_authenticated_user_roles"
	ContextKeyPermissions contextType = "ctx_authenticated_user_permissions"
)

func AquirePermissions(repo *repository.Repo) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			user := ctx.Value(ContextKeyUser).(*entity.User)
			// fetchig permissions
			perms := repo.UserPermissions(ctx, user.ID)
			if perms.Error != nil {
				http.Error(w, perms.Error.Error(), http.StatusInternalServerError)
				return
			}
			ctx = context.WithValue(ctx, ContextKeyPermissions, perms.Permissions)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func AquireRoles(repo *repository.Repo) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			user := ctx.Value(ContextKeyUser).(*entity.User)
			// fetching roles
			roles := repo.UserRoles(ctx, user.ID)
			if roles.Error != nil {
				http.Error(w, roles.Error.Error(), http.StatusInternalServerError)
				return
			}
			ctx = context.WithValue(ctx, ContextKeyRoles, roles.Roles)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
