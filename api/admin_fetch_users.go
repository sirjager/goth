package api

import (
	"net/http"

	"github.com/sirjager/gopkg/httpx"

	"github.com/sirjager/goth/entity"
	mw "github.com/sirjager/goth/middlewares"
)

type UsersResponse struct {
	Users []*entity.Profile `json:"users,omitempty"`
} // @name UsersResponse

// Admin api for fetching all users
//
// @Summary		Fetch Users
// @Description	Fetch multiple users
// @Tags			Admin
// @Produce		json
// @Param			page	query		int				false	"Page number: Default 1"
// @Param			limit	query		int				false	"Per Page: Default 100"
// @Success		200		{object}	UsersResponse	"UsersResponse"
// @Router			/api/admin/users [get]
func (s *Server) adminFetchUsers(w http.ResponseWriter, r *http.Request) {
	user := mw.UserOrPanic(r)
	if !user.Master {
		panic("admin operation user must have been master")
	}

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
