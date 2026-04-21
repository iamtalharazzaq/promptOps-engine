package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// HTTPRequestsTotal tracks the total number of HTTP requests processed by the backend.
	HTTPRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "promptops_http_requests_total",
			Help: "Total number of HTTP requests processed, partitioned by status code, method and path.",
		},
		[]string{"code", "method", "path"},
	)

	// HTTPRequestDuration tracks the latency of HTTP requests in seconds.
	HTTPRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "promptops_http_request_duration_seconds",
			Help:    "Latency of HTTP requests in seconds.",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	// OllamaTokenUsageTotal tracks the total number of tokens used (input + output).
	OllamaTokenUsageTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "promptops_ollama_token_usage_total",
			Help: "Total number of tokens consumed from Ollama, partitioned by model and type (prompt/completion).",
		},
		[]string{"model", "type"},
	)

	// OllamaRequestDuration tracks the latency of Ollama generation requests.
	OllamaRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "promptops_ollama_request_duration_seconds",
			Help:    "Latency of Ollama generation requests in seconds.",
			Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 30, 60},
		},
		[]string{"model"},
	)
)
