package api

import (
	"net/http"
)

type UpdateUserParams struct {
	FirstName  string `json:"first_name,omitempty"  validate:"min=3,max=30"`
	LastName   string `json:"last_name,omitempty"   validate:"min=1,max=30"`
	NickName   string `json:"nick_name,omitempty"   validate:"min=1,max=30"`
	Name       string `json:"name,omitempty"        validate:"min=1,max=100"`
	PictureURL string `json:"picture_url,omitempty" validate:"min=10"`
}

// @Summary		Single User
// @Description	Fetch specific user
// @Tags			Resources
// @Produce		json
// @Param			identity	path		string			true	"Identity can either be email or id"
// @Success		200			{object}	UserResponse	"User Response"
// @Router			/users/{identity} [get]
func (a *API) UserUpdate(w http.ResponseWriter, r *http.Request) {
	var params UpdateUserParams
	if err := a.ParseAndValidate(r, params); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	
}
