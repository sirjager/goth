package api

import (
	"encoding/json"
	"net/http"

	"github.com/sirjager/goth/entity"
)

// @Summary		Get User
// @Description	Get User
// @Tags			Resources
// @Produce		json
// @Success		200	{object}	UserResponse	"User Response"
// @Router			/users [get]
func (a *API) getUser(w http.ResponseWriter, r *http.Request) {
	user := r.Context().Value(userContext).(*entity.User)
	res := UserResponse{User: EntityToUser(user)}
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), 500)
	}
}
