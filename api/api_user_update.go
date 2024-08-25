package api

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sirjager/gopkg/httpx"

	"github.com/sirjager/goth/entity"
	mw "github.com/sirjager/goth/middlewares"
)

type UpdateUserParams struct {
	Username        string `json:"username,omitempty"`
	FirstName       string `json:"firstName,omitempty"`
	LastName        string `json:"lastName,omitempty"`
	FullName        string `json:"fullName,omitempty"`
	PictureURL      string `json:"pictureURL,omitempty"`
	NewPassword     string `json:"newPassword,omitempty"`
	CurrentPassword string `json:"currentPassword,omitempty"`
} //	@name	UpdateUserParams

// @Summary		Update User
// @Description	Partially Update User
// @Tags			Resources
// @Produce		json
// @Param			identity	path	string				true	"Identity can either be email or id"
// @Param			body		body	UpdateUserParams	true	"Update User Params"
// @Router			/users/{identity} [patch]
// @Success		200	{object}	UserResponse	"UserResponse"
func (s *Server) UserUpdate(w http.ResponseWriter, r *http.Request) {
	user := mw.UserOrPanic(r)
	identity := chi.URLParam(r, "identity")
	var params UpdateUserParams
	if err := s.ParseAndValidate(r, &params); err != nil {
		httpx.Error(w, err, http.StatusBadRequest)
		return
	}

	var target *entity.User // user to update
	if mw.IsCurrentUserIdentity(r) {
		target = user // if target is authorized user, self update
	} else {
		result := fetchUserFromRepository(r.Context(), identity, s.Repo())
		if result.Error != nil {
			httpx.Error(w, result.Error, result.StatusCode)
			return
		}
		target = result.User
	}

	shouldUpdate, err := patchUser(target, &params)
	if err != nil {
		httpx.Error(w, err, http.StatusBadRequest)
		return
	}

	if !shouldUpdate {
		httpx.Error(w, errors.New("no valid updates or no changes"), http.StatusBadRequest)
		return
	}

	res := s.Repo().UserUpdate(r.Context(), target)
	if res.Error != nil {
		httpx.Error(w, res.Error, res.StatusCode)
		return
	}

	response := UserResponse{res.User.Profile()}
	httpx.Success(w, response)
}
