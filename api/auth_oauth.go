package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"

	"github.com/sirjager/gopkg/httpx"
	mw "github.com/sirjager/goth/middlewares"
	"github.com/sirjager/goth/repository"
	"github.com/sirjager/goth/vo"

	"net/url"
)

// OAuth Provider
//
//	@Summary		OAuth Provider
//	@Description	Authenticates a user with a specified oauth provider
//	@Tags			Auth
//	@Produce		json
//	@Router			/api/auth/{provider} [get]
//	@Param			provider	path		string			true	"OAuth provider name [google,github]"	Enums(google, github)
//	@Success		200			{object}	UserResponse	"User object"
func (a *Server) oauthProvider(w http.ResponseWriter, r *http.Request) {
	refererURL := r.Header.Get("Referer")
	parsedURL, err := url.Parse(refererURL)
	if err != nil {
		http.Error(w, "invalid refer url", http.StatusBadRequest)
		return
	}
	// Reconstruct the base URL
	refererURL = parsedURL.Scheme + "://" + parsedURL.Host
	provider := chi.URLParam(r, "provider")
	if user, authenticated := mw.IsAuthenticated(r, a.App); authenticated {
		response := UserResponse{User: user.Profile()}
		httpx.Success(w, response)
		return
	}
	req := r.WithContext(context.WithValue(r.Context(), "provider", provider))
	gothic.SetState = func(req *http.Request) string {
		return refererURL
	}
	gothic.BeginAuthHandler(w, req)
}

func (s *Server) oauthCallback(w http.ResponseWriter, r *http.Request) {
	redirectURL := gothic.GetState(r) // set by AuthProvider, original url of calling client

	provider := chi.URLParam(r, "provider")
	req := r.WithContext(context.WithValue(r.Context(), "provider", provider))
	gothUser, err := gothic.CompleteUserAuth(w, req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// NOTE: saving user in database
	// IF EXISTS      : fetch from database using email, and return it
	// IF NOT EXISTS  : create and save user object, and return it
	newUser := GothUserToEntityUser(gothUser)
	// If master user does not exists, we make newUser a Master User.
	exists, existsErr := masterUserExists(r.Context(), s.Repo())
	if existsErr != nil {
		http.Error(w, existsErr.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		newUser.Master = true
	}

	result := s.Repo().UserCreate(r.Context(), newUser)
	if result.Error != nil {
		// if its not, user already exits error, return it
		if result.StatusCode != http.StatusConflict {
			httpx.Error(w, result.Error, result.StatusCode)
			return
		}
		// if its user already exits error, get user and return it
		result = s.Repo().UserGetByEmail(r.Context(), vo.MustParseEmail(gothUser.Email))
	}

	if result.Error != nil {
		httpx.Error(w, result.Error, result.StatusCode)
		return
	}

	// store user in cookie
	if err := storeUserSession(w, r, gothUser); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, redirectURL, http.StatusFound)
}

// StoreUserSession stores the user in the cookies
func storeUserSession(w http.ResponseWriter, r *http.Request, user goth.User) error {
	session, _ := gothic.Store.Get(r, mw.CookieGothicSession)
	session.Values["email"] = user.Email
	if err := session.Save(r, w); err != nil {
		return err
	}
	return nil
}

func masterUserExists(c context.Context, repo repository.Repo) (bool, error) {
	master := repo.UserGetMaster(c)
	if master.Error != nil {
		if master.StatusCode != http.StatusNotFound {
			return false, master.Error
		}
		// error is UserNotFoundError
		return false, nil
	}
	return master.User.Master, nil
}
