package api

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/sirjager/goth/repository/users"
	"github.com/sirjager/goth/vo"
)

type UserResponse struct {
	User User `json:"user,omitempty"`
} //	@name	UserResponse

// @Summary		Single User
// @Description	Fetch specific user
// @Tags			Resources
// @Produce		json
// @Param			identity	path		string			true	"Identity can either be email or id"
// @Success		200			{object}	UserResponse	"User Response"
// @Router			/users/{identity} [get]
func (a *API) UserGet(w http.ResponseWriter, r *http.Request) {
	identity := chi.URLParam(r, "identity")
	var result users.UserReadResult

	if email, emailErr := vo.NewEmail(identity); emailErr == nil {
		result = a.repo.UserReadByEmail(r.Context(), email.Value())
	} else {
		result = a.repo.UserReadByID(r.Context(), identity)
	}
	if result.Error != nil {
		if result.StatusCode == http.StatusNotFound {
			http.Error(w, result.Error.Error(), result.StatusCode)
			return
		}
		http.Error(w, result.Error.Error(), result.StatusCode)
		return
	}

	response := UserResponse{EntityToUser(result.User)}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type UsersResponse struct {
	Users []User `json:"users"`
} //	@name	UsersResponse

// @Summary		Multiple Users
// @Description	Fetch multiple users
// @Tags			Resources
// @Produce		json
// @Param			page	query		int				false	"Page number: Default 1"
// @Param			limit	query		int				false	"Per Page: Default 100"
// @Success		200		{object}	UsersResponse	"Users Response"
// @Router			/users [get]
func (a *API) UsersGet(w http.ResponseWriter, r *http.Request) {
	page := 1
	limit := 100
	getPageAndLimitFromRequest(r, &page, &limit)
	result := a.repo.UsersRead(r.Context(), limit, page)
	if result.Error != nil {
		http.Error(w, result.Error.Error(), result.StatusCode)
		return
	}

	response := UsersResponse{EntitiesToUsers(result.Users)}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
