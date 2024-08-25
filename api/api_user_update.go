package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sirjager/gopkg/httpx"
	"github.com/sirjager/goth/entity"
	mw "github.com/sirjager/goth/middlewares"
	"github.com/sirjager/goth/vo"
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

	isUserUpdated := false
	// only update when params are not empty and when different

	if params.Username != "" && params.Username != user.Username.Value() {
		newUsername, err := vo.NewUsername(params.Username)
		if err != nil {
			httpx.Error(w, err, http.StatusBadRequest)
			return
		}
		isUserUpdated = true
		user.Username = newUsername
	}

	if params.FullName != "" && params.FullName != user.FullName {
		isUserUpdated = true
		user.FullName = params.FullName
	}

	if params.FirstName != "" && params.FirstName != user.FirstName {
		isUserUpdated = true
		user.FirstName = params.FirstName
	}
	if params.LastName != "" && params.LastName != user.LastName {
		isUserUpdated = true
		user.LastName = params.LastName
	}
	if params.PictureURL != "" && params.PictureURL != user.PictureURL {
		isUserUpdated = true
		user.PictureURL = params.PictureURL
	}

	if isUserUpdated {
		result := s.Repo().UserUpdate(r.Context(), target)
		if result.Error != nil {
			httpx.Error(w, result.Error, result.StatusCode)
			return
		}
		target = result.User
	}

	response := UserResponse{target.Profile()}
	httpx.Success(w, response)
}
