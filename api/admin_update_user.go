package api

import (
	"errors"
	"net/http"

	"github.com/sirjager/gopkg/httpx"

	mw "github.com/sirjager/goth/middlewares"
)

// Admin Route for partially updating user any user
//
//	@Summary		Update User
//	@Description	Partially Update User
//	@Tags			Admin
//	@Produce		json
//	@Param			body		body	UpdateUserParams	true	"Update User Params"
//	@Param			identity	path	string				true	"Identity can either be email or id"
//	@Router			/api/admin/user/{identity} [patch]
//	@Success		200	{object}	UserResponse	"UserResponse"
func (s *Server) adminUpdateUser(w http.ResponseWriter, r *http.Request) {
	user := mw.AdminOrPanic(r)

	var params UpdateUserParams
	if err := s.ParseAndValidate(r, &params); err != nil {
		httpx.Error(w, err, http.StatusBadRequest)
		return
	}

	shouldUpdate, err := patchUser(user, &params)
	if err != nil {
		httpx.Error(w, err, http.StatusBadRequest)
		return
	}

	if !shouldUpdate {
		httpx.Error(w, errors.New("no valid updates or no changes"), http.StatusBadRequest)
		return
	}

	res := s.Repo().UserUpdate(r.Context(), user)
	if res.Error != nil {
		httpx.Error(w, res.Error, res.StatusCode)
		return
	}

	response := UserResponse{res.User.Profile()}
	httpx.Success(w, response)
}
