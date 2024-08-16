package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sirjager/goth/entity"
	mw "github.com/sirjager/goth/middlewares"
)

type UpdateUserParams struct {
	FirstName  string `json:"first_name,omitempty"  validate:"max=30"`
	LastName   string `json:"last_name,omitempty"   validate:"max=30"`
	Name       string `json:"name,omitempty"        validate:"max=100"`
	PictureURL string `json:"picture_url,omitempty" validate:""`
} //	@name	UpdateUserParams

// @Summary		Update User
// @Description	Partially Update User
// @Tags			Resources
// @Produce		json
// @Param			identity	path	string				true	"Identity can either be email or id"
// @Param			body		body	UpdateUserParams	true	"Update User Params"
// @Router			/users/{identity} [patch]
// @Success		200	{object}	UserResponse	"UserResponse"
func (a *API) UserUpdate(w http.ResponseWriter, r *http.Request) {
	user := mw.UserOrPanic(r)
	identity := chi.URLParam(r, "identity")
	var params UpdateUserParams
	if err := a.ParseAndValidate(r, &params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var target *entity.User   // user to update
	shouldApplyPatch := false // should apply updates or not
	if mw.IsCurrentUserIdentity(r) {
		target = user // if target is authorized user, self update
	} else {
		result := fetchUserFromRepository(r.Context(), identity, a.repo)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), result.StatusCode)
			return
		}
		target = result.User
	}

	shouldApplyPatch = patchUser(target, params)
	if shouldApplyPatch {
		result := a.repo.UserUpdate(r.Context(), target)
		if result.Error != nil {
			http.Error(w, result.Error.Error(), result.StatusCode)
			return
		}
		target = result.User
	}

	response := UserResponse{target.Profile()}
	a.Success(w, response)
}

func patchUser(user *entity.User, params UpdateUserParams) (shouldPatch bool) {
	// only update when params are not empty and when different
	if params.Name != "" && params.Name != user.Name {
		shouldPatch = true
		user.Name = params.Name
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
