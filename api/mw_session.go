package api

import (
	"net/http"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
)

const SessionName = "session"

// StoreUserSession stores the user in the cookies
func (a *API) StoreUserSession(w http.ResponseWriter, r *http.Request, user goth.User) error {
	session, _ := gothic.Store.Get(r, SessionName)
	session.Values["user"] = user

	if err := session.Save(r, w); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return err
	}
	return nil
}
