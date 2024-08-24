package api

import (
	"net/http"

	"github.com/sirjager/gopkg/httpx"
	mw "github.com/sirjager/goth/middlewares"
)

// @Summary		User
// @Description	Returns Authenticated User
// @Tags			Auth
// @Produce		json
// @Success		200	{object}	UserResponse	"UserResponse"
// @Router			/auth/user [get]
func (a *Server) AuthUser(w http.ResponseWriter, r *http.Request) {
	user := mw.UserOrPanic(r)
	response := UserResponse{user.Profile()}
	httpx.Success(w, response)
}
