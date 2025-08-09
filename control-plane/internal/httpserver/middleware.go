package httpserver

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/hywgb/pjSipKF/control-plane/internal/metrics"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func WithMiddlewares(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			defer func() {
				if rec := recover(); rec != nil {
					logger.Error("panic recovered", zap.Any("panic", rec))
					http.Error(w, "internal error", http.StatusInternalServerError)
				}
			}()
			rw := middleware.NewWrapResponseWriter(rec, r.ProtoMajor)
			next.ServeHTTP(rw, r)
			code := rw.Status()
			metrics.RequestCounter.WithLabelValues(r.URL.Path, r.Method, strconv.Itoa(code)).Inc()
			logger.Info("http_request",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Int("status", code),
				zap.Duration("dur_ms", time.Since(start)),
			)
		})
	}
}