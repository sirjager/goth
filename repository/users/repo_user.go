package users

import (
	"context"

	"github.com/sirjager/goth/entity"
)

type UsersReadResult struct {
	Error      error
	Users      []*entity.User
	StatusCode int
}

type UserReadResult struct {
	Error      error
	User       *entity.User
	StatusCode int
}

type UserDeleteResult struct {
	Error      error
	StatusCode int
}

type UsersRepo interface {
	UserCreate(c context.Context, u *entity.User) UserReadResult
	UsersRead(ctx context.Context, optionalLimit, optionalPage int) UsersReadResult
	UserReadByID(c context.Context, uid string) UserReadResult
	UserReadByEmail(c context.Context, email string) UserReadResult
	UserReadMaster(c context.Context) UserReadResult
	UserDelete(c context.Context, uid string) UserDeleteResult
}
