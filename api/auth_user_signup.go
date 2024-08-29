package api

import (
	"net/http"
	"strings"

	"github.com/sirjager/gopkg/httpx"
	"github.com/sirjager/goth/entity"
	"github.com/sirjager/goth/vo"
)

type signupReqParams struct {
	Email    string `json:"email,omitempty"    validate:"required,gte=3"`
	Password string `json:"password,omitempty" validate:"required"`
} // @name SignUpRequestParams

// Route for signing up a new user
//
//	@Summary		SignUp User
//	@Description	Sign up a new user using email and password
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Router			/api/auth/signup [post]
//	@Param			body	body		SignUpRequestParams	true	"sign up params : email and password"
//	@Success		201		{object}	UserResponse		"User object"
func (s *Server) authUserSignUp(w http.ResponseWriter, r *http.Request) {
	var params signupReqParams
	if err := s.ParseAndValidate(r, &params); err != nil {
		httpx.Error(w, err, http.StatusBadRequest)
		return
	}

	var username *vo.Username
	email, err := vo.NewEmail(params.Email)
	if err != nil {
		httpx.Error(w, err, http.StatusBadRequest)
		return
	}
	password, err := vo.NewPassword(params.Password)
	if err != nil {
		httpx.Error(w, err, http.StatusBadRequest)
		return
	}

	username, err = vo.NewUsername(strings.Split(params.Email, "@")[0])
	if err != nil {
		httpx.Error(w, err, http.StatusBadRequest)
		return
	}

	hashedPassword, err := password.HashPassword()
	if err != nil {
		httpx.Error(w, err)
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
	exists, existsErr := masterUserExists(r.Context(), s.Repo())
	if existsErr != nil {
		httpx.Error(w, existsErr)
		return
	}
	if !exists {
		newUser.Master = true
	}

	result := s.Repo().UserCreate(r.Context(), newUser)
	if result.Error != nil {
		httpx.Error(w, result.Error, result.StatusCode)
		return
	}

	response := UserResponse{result.User.Profile()}
	httpx.Success(w, response, result.StatusCode)
}
