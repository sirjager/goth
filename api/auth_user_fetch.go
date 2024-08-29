package api

import (
	"net/http"

	"github.com/sirjager/gopkg/httpx"

	"github.com/sirjager/goth/entity"
	mw "github.com/sirjager/goth/middlewares"
)

type UserResponse struct {
	User *entity.Profile `json:"user,omitempty"`
} // @name UserResponse

// Authenticated user route for fetching authenticated user
//
// @Summary		User Fetch
// @Description	Get Authenticated User
// @Tags			Auth
// @Produce		json
// @Success		200	{object}	UserResponse	"UserResponse"
// @Router			/api/auth/user [get]
func (s *Server) authUserFetch(w http.ResponseWriter, r *http.Request) {
	user := mw.UserOrPanic(r)
	response := UserResponse{user.Profile()}
	httpx.Success(w, response)
}
