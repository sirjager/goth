package mw

import (
	"context"
	"net/http"

	"github.com/sirjager/goth/modules"
	"github.com/sirjager/goth/vo"
)

type RequestIDKey int

const (
	ContextRequestID RequestIDKey = iota
)

// RequestID attaches unique request id to each request
func RequestID(modules *modules.Modules) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID, err := vo.NewID()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			ctx := context.WithValue(r.Context(), ContextRequestID, requestID.Value().String())
			w.Header().Set("X-Server-Name", modules.Config().ServerName)
			w.Header().Set("X-Request-ID", requestID.Value().String())
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
