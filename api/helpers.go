package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/sirjager/goth/repository"
	"github.com/sirjager/goth/repository/users"
	"github.com/sirjager/goth/vo"
)

type Validation bool

const (
	ValidationDisable Validation = false
	ValidationEnable  Validation = true
)

func (a *API) SetCookies(w http.ResponseWriter, cookies ...*http.Cookie) {
	for _, cookie := range cookies {
		http.SetCookie(w, cookie)
	}
}

func (a *API) Failure(w http.ResponseWriter, err error, statusCode ...int) {
	status := 500
	if len(statusCode) == 1 {
		status = statusCode[0]
	}
	http.Error(w, err.Error(), status)
}

func (a *API) Success(w http.ResponseWriter, response any, statusCode ...int) {
	status := 200
	if len(statusCode) == 1 {
		status = statusCode[0]
	}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		a.Failure(w, err)
	}
}

func (a *API) GetPageAndLimitFromRequest(r *http.Request, defaultPage, defaultLimit *int) {
	pageParam := r.URL.Query().Get("page")
	limitParam := r.URL.Query().Get("limit")
	if (pageParam) != "" {
		if p, err := strconv.Atoi(pageParam); err == nil && p > 0 {
			*defaultPage = p
		}
	}
	if (limitParam) != "" {
		if l, err := strconv.Atoi(limitParam); err == nil && l > 0 {
			*defaultLimit = l
		}
	}
}

func (a *API) ParseJSON(r *http.Request, v interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return err
	}
	return nil
}

func (a *API) ParseAndValidate(r *http.Request, v interface{}, validation ...Validation) error {
	validate := ValidationEnable
	if len(validation) == 1 {
		validate = validation[0]
	}
	if err := a.ParseJSON(r, v); err != nil {
		return err
	}
	if validate {
		if err := a.validate.Struct(v); err != nil {
			return err
		}
	}
	return nil
}

// fetchUserFromRepository fetches user by email or id
func fetchUserFromRepository(
	c context.Context,
	identity string,
	repo *repository.Repo,
) users.UserReadResult {
	if email, emailErr := vo.NewEmail(identity); emailErr == nil {
		return repo.UserReadByEmail(c, email)
	} else {
		return repo.UserReadByID(c, vo.MustParseID(identity))
	}
}
