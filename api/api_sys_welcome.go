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

func welcomeMessaage(serviceName string) string {
	return fmt.Sprintf("Welcome to %s", serviceName)
}

const docsPath = "/swagger"

// Welcome Docs
//
//	@Summary		Welcome
//	@Description	Welcome
//	@Tags			System
//	@Accept			json
//	@Produce		json
//	@Router			/ [get]
//	@Success		200	{object}	welcomeResponse
func (a *Server) Welcome(w http.ResponseWriter, r *http.Request) {
	response := welcomeResponse{
		Message: welcomeMessaage(a.config.ServiceName),
		Docs:    docsPath,
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		a.logr.Error().Err(err).Msg("failed to encode welcome response")
		http.Error(w, err.Error(), 500)
	}
}
