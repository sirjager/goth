package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type welcomeResponse struct {
	Message string `json:"message,omitempty"`
	Docs    string `json:"docs,omitempty"`
} // @name WelcomeResponse

// Welcome Docs
//
//	@Summary		Welcome
//	@Description	Welcome
//	@Tags			System
//	@Accept			json
//	@Produce		json
//	@Router			/ [get]
//	@Success		200	{object}	welcomeResponse
func (a *API) SysWelcome(w http.ResponseWriter, r *http.Request) {
	response := welcomeResponse{
		Message: fmt.Sprintf("Welcome to %s", a.config.ServiceName),
		Docs:    "/swagger",
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		a.logr.Error().Err(err).Msg("failed to encode welcome response")
		http.Error(w, err.Error(), 500)
	}
}
