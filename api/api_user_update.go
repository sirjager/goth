package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sirjager/goth/entity"
	mw "github.com/sirjager/goth/middlewares"
)

type UpdateUserParams struct {
	Email      string `json:"email,omitempty"`
	Username   string `json:"username,omitempty"`
	Password   string `json:"password,omitempty"`
	FirstName  string `json:"firstName,omitempty"`
	LastName   string `json:"lastName,omitempty"`
	FullName   string `json:"fullName,omitempty"`
	PictureURL string `json:"pictureURL,omitempty"`
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
		s.Failure(w, err, http.StatusBadRequest)
		return
	}

	var target *entity.User   // user to update
	shouldApplyPatch := false // should apply updates or not
	if mw.IsCurrentUserIdentity(r) {
		target = user // if target is authorized user, self update
	} else {
		result := fetchUserFromRepository(r.Context(), identity, s.repo)
		if result.Error != nil {
			s.Failure(w, result.Error, result.StatusCode)
			return
		}
		target = result.User
	}

	shouldApplyPatch = patchUser(target, params)
	if shouldApplyPatch {
		result := s.repo.UserUpdate(r.Context(), target)
		if result.Error != nil {
			s.Failure(w, result.Error, result.StatusCode)
			return
		}
		target = result.User
	}

	response := UserResponse{target.Profile()}
	s.Success(w, response)
}

func patchUser(user *entity.User, params UpdateUserParams) (shouldPatch bool) {
	// only update when params are not empty and when different
	if params.FullName != "" && params.FullName != user.FullName {
		shouldPatch = true
		user.FullName = params.FullName
	}
	if params.FirstName != "" && params.FirstName != user.FirstName {
		shouldPatch = true
		user.FirstName = params.FirstName
	}
	if params.LastName != "" && params.LastName != user.LastName {
		shouldPatch = true
		user.LastName = params.LastName
	}
	if params.PictureURL != "" && params.PictureURL != user.PictureURL {
		shouldPatch = true
		user.PictureURL = params.PictureURL
	}
	return shouldPatch
}
