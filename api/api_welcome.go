package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type WelcomResponse struct {
	Message string `json:"message,omitempty"`
	Docs    string `json:"docs,omitempty"`
} // @name WelcomeResponse

func welcomeMessaage(serviceName string) string {
	return fmt.Sprintf("Welcome to %s", serviceName)
}

// Welcome message
//
//	@Summary		Welcome
//	@Description	Welcome message
//	@Tags			System
//	@Accept			json
//	@Produce		json
//	@Router			/api [get]
//	@Success		200	{object}	WelcomeResponse
func (s *Server) apiWelcome(w http.ResponseWriter, r *http.Request) {
	response := WelcomResponse{
		Message: welcomeMessaage(s.Config().ServiceName),
		Docs:    "/api/docs",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.Logger().Error().Err(err).Msg("failed to encode welcome response")
		http.Error(w, err.Error(), 500)
	}
}
