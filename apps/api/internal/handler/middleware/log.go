package middleware

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/kobayashiyabako16g/passkey-auth-example/internal/contexts"
	"github.com/kobayashiyabako16g/passkey-auth-example/pkg/logger"
)

func LogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		start := time.Now()
		id, err := uuid.NewV7()
		if err != nil {
			logger.Error(ctx, "Failed to generate request ID", err)
			id = uuid.New()
		}

		ctx = contexts.SetRequestID(ctx, id.String())

		logger.Info(ctx, "HTTP Request",
			"request_id", id.String(),
			"method", r.Method,
			"path", r.URL.Path,
			"user_agent", r.UserAgent(),
			"ip", r.RemoteAddr,
		)

		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)

		logger.Info(ctx, "HTTP Request Finished",
			"request_id", id.String(),
			"method", r.Method,
			"path", r.URL.Path,
			"status", lrw.statusCode,
			"duration", time.Since(start),
		)

	})
}

// loggingResponseWriter はレスポンスのステータスコードを取得するためのラッパー
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
