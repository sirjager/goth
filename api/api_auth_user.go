package api

import (
	"errors"
	"net/http"

	"github.com/sirjager/gopkg/httpx"
	"github.com/sirjager/gopkg/utils"

	"github.com/sirjager/goth/entity"
	mw "github.com/sirjager/goth/middlewares"
	"github.com/sirjager/goth/vo"
)

// @Summary		User
// @Description	Get Authenticated User
// @Tags			Auth
// @Produce		json
// @Success		200	{object}	UserResponse	"UserResponse"
// @Router			/auth/user [get]
func (s *Server) AuthUser(w http.ResponseWriter, r *http.Request) {
	user := mw.UserOrPanic(r)
	response := UserResponse{user.Profile()}
	httpx.Success(w, response)
}

// @Summary		User
// @Description	Update Authenticated User
// @Tags			Auth
// @Produce		json
// @Router			/auth/user [patch]
// @Param			body	body		UpdateUserParams	true	"Update User Params"
// @Success		200		{object}	UserResponse		"UserResponse"
func (s *Server) AuthUserPatch(w http.ResponseWriter, r *http.Request) {
	user := mw.UserOrPanic(r)
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

func patchUser(user *entity.User, params *UpdateUserParams) (bool, error) {
	shouldApplyPatch := false
	if !utils.IsEmpty(params.Username) && !user.Username.IsEqual(params.Username) {
		username, err := vo.NewUsername(params.Username)
		if err != nil {
			return shouldApplyPatch, err
		}
		user.Username = username
	}
	if !utils.IsEmpty(params.FullName) && user.FullName != params.FullName {
		user.FirstName = params.FullName
		shouldApplyPatch = true
	}
	if !utils.IsEmpty(params.FirstName) && user.FirstName != params.FirstName {
		user.FirstName = params.FirstName
		shouldApplyPatch = true
	}
	if !utils.IsEmpty(params.LastName) && user.LastName != params.LastName {
		user.LastName = params.LastName
		shouldApplyPatch = true
	}
	if !utils.IsEmpty(params.PictureURL) && user.PictureURL != params.PictureURL {
		user.PictureURL = params.PictureURL
		shouldApplyPatch = true
	}
	// to change to a new password, one must provide current password
	if !utils.IsEmpty(params.NewPassword) {
		// return err if weak/invalid password
		newPassword, err := vo.NewPassword(params.NewPassword)
		if err != nil {
			return shouldApplyPatch, err
		}
		// match current password with users's password
		if err = user.Password.VerifyPassword(params.CurrentPassword); err != nil {
			return shouldApplyPatch, errors.New("invalid current password")
		}
		// hashing new password and updating to users password
		hashedNewPassword, err := newPassword.HashPassword()
		if err != nil {
			return shouldApplyPatch, err
		}
		shouldApplyPatch = true
		user.Password = hashedNewPassword
	}
	return shouldApplyPatch, nil
}
