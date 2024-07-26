package users

import (
	"context"
)

type UserUpdateParams struct {
	FirstName  string
	LastName   string
	NickName   string
	Name       string
	PictureURL string
}

func (r *repo) UserUpdate(ctx context.Context, id string, p *UserUpdateParams) {
 
}
