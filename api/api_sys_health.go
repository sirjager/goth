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
func (a *Server) Health(w http.ResponseWriter, r *http.Request) {
	response := healthResponse{
		Timestamp: time.Now(),
		Service:   a.config.ServiceName,
		Server:    a.config.ServerName,
		Started:   a.config.StartTime,
		Status:    healthpb.HealthCheckResponse_SERVING.String(),
		Uptime:    time.Since(a.config.StartTime).String(),
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(200)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		a.logr.Error().Err(err).Msg("failed to encode health response")
		http.Error(w, err.Error(), 500)
	}
}
