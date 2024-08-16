package api

import (
	"context"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"

	mw "github.com/sirjager/goth/middlewares"
	"github.com/sirjager/goth/repository"
	"github.com/sirjager/goth/vo"
)

func (a *API) AuthCallback(w http.ResponseWriter, r *http.Request) {
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
	exists, existsErr := masterUserExists(r.Context(), a.repo)
	if existsErr != nil {
		http.Error(w, existsErr.Error(), http.StatusInternalServerError)
		return
	}
	if !exists {
		newUser.Master = true
	}

	result := a.repo.UserCreate(r.Context(), newUser)
	if result.Error != nil {
		// if its not, user already exits error, return it
		if result.StatusCode != http.StatusConflict {
			a.Failure(w, result.Error, result.StatusCode)
			return
		}
		// if its user already exits error, get user and return it
		result = a.repo.UserReadByEmail(r.Context(), vo.MustParseEmail(gothUser.Email))
	}

	if result.Error != nil {
		a.Failure(w, result.Error, result.StatusCode)
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
	session, _ := gothic.Store.Get(r, mw.SessionCookieName)
	session.Values["email"] = user.Email
	if err := session.Save(r, w); err != nil {
		return err
	}
	return nil
}

func masterUserExists(c context.Context, repo *repository.Repo) (bool, error) {
	master := repo.UserReadMaster(c)
	if master.Error != nil {
		if master.StatusCode != http.StatusNotFound {
			return false, master.Error
		}
		// error is UserNotFoundError
		return false, nil
	}
	return master.User.Master, nil
}
