package api

import (
	"encoding/json"
	"net/http"
	"time"

	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

type healthResponse struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
	Started   time.Time `json:"started,omitempty"`
	Service   string    `json:"service,omitempty"`
	Server    string    `json:"server,omitempty"`
	Status    string    `json:"status,omitempty"`
	Uptime    string    `json:"uptime,omitempty"`
} // @name HealthResponse

// @Summary		Health
// @Description	Health Check
// @Tags			System
// @Accept			json
// @Produce		json
// @Router			/health [get]
// @Success		200	{object}	healthResponse
func (s *Server) Health(w http.ResponseWriter, r *http.Request) {
	response := healthResponse{
		Timestamp: time.Now(),
		Service:   s.Config().ServiceName,
		Server:    s.Config().ServerName,
		Started:   s.Config().StartTime,
		Status:    healthpb.HealthCheckResponse_SERVING.String(),
		Uptime:    time.Since(s.Config().StartTime).String(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		s.Logger().Error().Err(err).Msg("failed to encode health response")
		http.Error(w, err.Error(), 500)
	}
}
