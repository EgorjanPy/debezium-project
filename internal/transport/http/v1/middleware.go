package v1

import (
	"context"
	"debez/pkg/logger"
	"net/http"

	"github.com/google/uuid"
)

func LoggingMiddleware(ctx context.Context) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get("x-request-id")
			if requestID != "" {
				logger.WithRequestID(r.Context(), requestID)
			} else {
				logger.WithRequestID(r.Context(), uuid.New().String())
			}

			next.ServeHTTP(w, r)
		})
	}
}
