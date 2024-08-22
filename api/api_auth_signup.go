package api

import (
	"net/http"
	"strings"

	"github.com/sirjager/goth/entity"
	"github.com/sirjager/goth/vo"
)

type SignUpRequestParams struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email,omitempty"    validate:"required,gte=3"`
	Password string `json:"password,omitempty" validate:"required"`
} // @name SignUpRequestParams

// Signup Request
//
//	@Summary		Signup
//	@Description	Signup using email and password
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Router			/auth/signup [post]
//	@Param			body	body		SignUpRequestParams	true	"Signup request params"
//	@Success		201		{object}	UserResponse		"User object"
func (a *Server) Signup(w http.ResponseWriter, r *http.Request) {
	var params SignUpRequestParams
	if err := a.ParseAndValidate(r, &params); err != nil {
		a.Failure(w, err, http.StatusBadRequest)
		return
	}

	var username *vo.Username
	email, err := vo.NewEmail(params.Email)
	if err != nil {
		a.Failure(w, err, http.StatusBadRequest)
		return
	}
	password, err := vo.NewPassword(params.Password)
	if err != nil {
		a.Failure(w, err, http.StatusBadRequest)
		return
	}

	if params.Username == "" {
		params.Username = strings.Split(params.Email, "@")[0]
	}

	username, err = vo.NewUsername(params.Username)
	if err != nil {
		a.Failure(w, err, http.StatusBadRequest)
		return
	}

	hashedPassword, err := password.HashPassword()
	if err != nil {
		a.Failure(w, err)
		return
	}

	newUser := &entity.User{
		Email:    email,
		Verified: false,
		Username: username,
		Provider: "credentials",
		Password: hashedPassword,
	}

	// If master user does not exists, we make newUser a Master User.
	exists, existsErr := masterUserExists(r.Context(), a.repo)
	if existsErr != nil {
		a.Failure(w, existsErr)
		return
	}
	if !exists {
		newUser.Master = true
	}

	result := a.repo.UserCreate(r.Context(), newUser)
	if result.Error != nil {
		a.Failure(w, result.Error, result.StatusCode)
		return
	}

	response := UserResponse{result.User.Profile()}
	a.Success(w, response, result.StatusCode)
}
