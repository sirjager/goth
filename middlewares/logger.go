package mw

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/sirjager/goth/modules"
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

func Logger(modules *modules.Modules) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now() // Start timer
			path := r.URL.Path

			rec := &ResponseRecorder{ResponseWriter: w, StatusCode: 200, Body: &bytes.Buffer{}}
			next.ServeHTTP(rec, r)

			duration := time.Since(start)
			event := modules.Logger().Info()

			if rec.StatusCode != http.StatusOK {
				var data map[string]interface{}
				if err := json.Unmarshal(rec.Body.Bytes(), &data); err != nil {
					data = map[string]interface{}{}
				}
				event = modules.Logger().Error().Interface("error", data["message"])
			}

			if rec.StatusCode >= 400 && rec.StatusCode < 500 {
				event = modules.Logger().Warn()
			} else if rec.StatusCode >= 500 {
				event = modules.Logger().Error()
			}

			shortenedPath := shortenPath(path, 20)
			icon := getIcon(rec.StatusCode)
			coloredIcon := getColoredIcon(rec.StatusCode)

			event.
				Str("method", r.Method).
				Str("path", shortenedPath).
				Dur("latency", duration).
				Int("code", rec.StatusCode)

			if modules.Config().LoggerLogfile != "" {
				event.Msg(icon)
			} else {
				goEnv := modules.Config().GoEnv
				if strings.Contains(goEnv, "test") || strings.Contains(goEnv, "prod") {
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
