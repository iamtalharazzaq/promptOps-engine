package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/promptops/backend/pkg/metrics"
)

// Metrics returns middleware that records HTTP metrics for Prometheus.
func Metrics() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Capture status code by wrapping response writer
			rw := &responseWriterMetrics{ResponseWriter: w, statusCode: http.StatusOK}

			next.ServeHTTP(rw, r)

			duration := time.Since(start).Seconds()
			statusCode := strconv.Itoa(rw.statusCode)

			// Update Prometheus metrics
			metrics.HTTPRequestsTotal.WithLabelValues(statusCode, r.Method, r.URL.Path).Inc()
			metrics.HTTPRequestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
		})
	}
}

type responseWriterMetrics struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriterMetrics) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}
