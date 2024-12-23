package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sirjager/gopkg/httpx"
	"github.com/sirjager/goth/entity"
	mw "github.com/sirjager/goth/middlewares"
	"github.com/sirjager/goth/repository/users"
)

type UserResponse struct {
	User *entity.Profile `json:"user,omitempty"`
} //	@name	UserResponse

// @Summary		Single User
// @Description	Fetch specific user
// @Tags			Resources
// @Produce		json
// @Param			identity	path		string			true	"Identity can either be email or id"
// @Success		200			{object}	UserResponse	"UserResponse"
// @Router			/users/{identity} [get]
func (s *Server) UserGet(w http.ResponseWriter, r *http.Request) {
	user := mw.UserOrPanic(r)
	identity := chi.URLParam(r, "identity")
	var result users.UserReadResult

	// if asking for own user document, then return authenticated user
	if mw.IsCurrentUserIdentity(r) {
		result = users.UserReadResult{Error: nil, User: user, StatusCode: 200}
	} else {
		result = fetchUserFromRepository(r.Context(), identity, s.Repo())
	}

	if result.Error != nil {
		http.Error(w, result.Error.Error(), result.StatusCode)
		return
	}

	// resolved master role request
	response := UserResponse{result.User.Profile()}
	httpx.Success(w, response, result.StatusCode)
}

type UsersResponse struct {
	Users []*entity.Profile `json:"users"`
} //	@name	UsersResponse

// @Summary		Multiple Users
// @Description	Fetch multiple users
// @Tags			Resources
// @Produce		json
// @Param			page	query		int				false	"Page number: Default 1"
// @Param			limit	query		int				false	"Per Page: Default 100"
// @Success		200		{object}	UsersResponse	"UsersResponse"
// @Router			/users [get]
func (s *Server) UsersGet(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 100
	s.GetPageAndLimitFromRequest(r, &page, &limit)
	result := s.Repo().UserGetAll(r.Context(), limit, page)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), result.StatusCode)
		return
	}

	response := UsersResponse{EntitiesToProfiles(result.Users)}
	httpx.Success(w, response, result.StatusCode)
}
