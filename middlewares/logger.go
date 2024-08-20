package mw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/rs/zerolog"

	"github.com/sirjager/goth/config"
)

const (
	boldGreen  = "\033[1;32m"
	boldRed    = "\033[1;31m"
	boldYellow = "\033[1;33m"
	boldCyan   = "\033[1;36m"
	reset      = "\033[0m"
)

type ResponseRecorder struct {
	http.ResponseWriter
	Body       *bytes.Buffer
	StatusCode int
}

func (rec *ResponseRecorder) WriteHeader(statusCode int) {
	rec.StatusCode = statusCode
	rec.ResponseWriter.WriteHeader(statusCode)
}

func (rec *ResponseRecorder) Write(b []byte) (int, error) {
	rec.Body.Write(b)
	return rec.ResponseWriter.Write(b)
}

func Logger(logr zerolog.Logger, config *config.Config) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now() // Start timer
			path := r.URL.Path

			rec := &ResponseRecorder{ResponseWriter: w, StatusCode: 200, Body: &bytes.Buffer{}}
			next.ServeHTTP(rec, r)

			duration := time.Since(start)
			event := logr.Info()

			if rec.StatusCode != http.StatusOK {
				var data map[string]interface{}
				if err := json.Unmarshal(rec.Body.Bytes(), &data); err != nil {
					data = map[string]interface{}{}
				}
				event = logr.Error().Interface("error", data["message"])
			}

			if rec.StatusCode >= 400 && rec.StatusCode < 500 {
				event = logr.Warn()
			} else if rec.StatusCode >= 500 {
				event = logr.Error()
			}

			shortenedPath := shortenPath(path, 20)
			icon := getIcon(rec.StatusCode)
			coloredIcon := getColoredIcon(rec.StatusCode)

			event.
				Str("method", r.Method).
				Str("path", shortenedPath).
				Dur("latency", duration).
				Int("code", rec.StatusCode)

			if config.LoggerLogfile != "" {
				event.Msg(icon)
			} else {
				if strings.Contains(config.GoEnv, "test") || strings.Contains(config.GoEnv, "prod") {
					event.Msg(icon)
				} else {
					event.Msg(coloredIcon)
				}
			}
		})
	}
}

func getIcon(code int) string {
	switch {
	case code >= 200 && code < 300:
		return ""
	case code >= 300 && code < 400:
		return ""
	case code >= 400 && code < 500:
		return ""
	default:
		return ""
	}
}

func getColoredIcon(code int) string {
	switch {
	case code >= 200 && code < 300:
		return fmt.Sprintf("%s %s", boldGreen, reset)
	case code >= 300 && code < 400:
		return fmt.Sprintf("%s %s", boldCyan, reset)
	case code >= 400 && code < 500:
		return fmt.Sprintf("%s %s", boldYellow, reset)
	default:
		return fmt.Sprintf("%s %s", boldRed, reset)
	}
}

func shortenPath(path string, max int) string {
	if len(path) > max {
		return path[0:max] + "..."
	}
	return path
}
