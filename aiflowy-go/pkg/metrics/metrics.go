// Package metrics provides Prometheus metrics for AIFlowy Go backend
package metrics

import (
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// HTTP request metrics
	httpRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "aiflowy_http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "aiflowy_http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10},
		},
		[]string{"method", "path"},
	)

	httpRequestsInFlight = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "aiflowy_http_requests_in_flight",
			Help: "Number of HTTP requests currently being served",
		},
	)

	// LLM/AI metrics
	llmRequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "aiflowy_llm_requests_total",
			Help: "Total number of LLM API requests",
		},
		[]string{"model", "provider", "status"},
	)

	llmRequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "aiflowy_llm_request_duration_seconds",
			Help:    "LLM API request duration in seconds",
			Buckets: []float64{0.1, 0.5, 1, 2, 5, 10, 30, 60, 120},
		},
		[]string{"model", "provider"},
	)

	llmTokensTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "aiflowy_llm_tokens_total",
			Help: "Total number of LLM tokens processed",
		},
		[]string{"model", "provider", "type"}, // type: prompt, completion
	)

	// Bot metrics
	botChatSessionsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "aiflowy_bot_chat_sessions_total",
			Help: "Total number of bot chat sessions",
		},
		[]string{"bot_id"},
	)

	// Workflow metrics
	workflowExecutionsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "aiflowy_workflow_executions_total",
			Help: "Total number of workflow executions",
		},
		[]string{"workflow_id", "status"},
	)

	workflowExecutionDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "aiflowy_workflow_execution_duration_seconds",
			Help:    "Workflow execution duration in seconds",
			Buckets: []float64{0.1, 0.5, 1, 5, 10, 30, 60, 120, 300},
		},
		[]string{"workflow_id"},
	)

	// Database metrics
	dbQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "aiflowy_db_query_duration_seconds",
			Help:    "Database query duration in seconds",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1},
		},
		[]string{"operation"},
	)

	// Application info
	appInfo = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "aiflowy_app_info",
			Help: "Application information",
		},
		[]string{"version", "go_version", "build_time"},
	)
)

// Init initializes the metrics with application info
func Init(version, goVersion, buildTime string) {
	appInfo.WithLabelValues(version, goVersion, buildTime).Set(1)
}

// MetricsHandler returns the Prometheus metrics handler for Echo
func MetricsHandler() echo.HandlerFunc {
	h := promhttp.Handler()
	return func(c echo.Context) error {
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

// Middleware returns Echo middleware for HTTP metrics
func Middleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Skip metrics endpoint
			if c.Path() == "/metrics" {
				return next(c)
			}

			start := time.Now()
			httpRequestsInFlight.Inc()

			err := next(c)

			httpRequestsInFlight.Dec()
			duration := time.Since(start).Seconds()

			status := c.Response().Status
			method := c.Request().Method
			path := c.Path()
			if path == "" {
				path = c.Request().URL.Path
			}

			// Normalize path for metrics (avoid high cardinality)
			path = normalizePath(path)

			httpRequestsTotal.WithLabelValues(method, path, strconv.Itoa(status)).Inc()
			httpRequestDuration.WithLabelValues(method, path).Observe(duration)

			return err
		}
	}
}

// normalizePath normalizes paths for metrics to avoid high cardinality
func normalizePath(path string) string {
	// Keep only the first two path segments
	// e.g., /api/v1/bot/123 -> /api/v1/bot/:id
	if len(path) > 100 {
		path = path[:100]
	}
	return path
}

// RecordLLMRequest records an LLM API request
func RecordLLMRequest(model, provider, status string, duration time.Duration) {
	llmRequestsTotal.WithLabelValues(model, provider, status).Inc()
	llmRequestDuration.WithLabelValues(model, provider).Observe(duration.Seconds())
}

// RecordLLMTokens records LLM token usage
func RecordLLMTokens(model, provider string, promptTokens, completionTokens int) {
	llmTokensTotal.WithLabelValues(model, provider, "prompt").Add(float64(promptTokens))
	llmTokensTotal.WithLabelValues(model, provider, "completion").Add(float64(completionTokens))
}

// RecordBotChatSession records a bot chat session
func RecordBotChatSession(botID string) {
	botChatSessionsTotal.WithLabelValues(botID).Inc()
}

// RecordWorkflowExecution records a workflow execution
func RecordWorkflowExecution(workflowID, status string, duration time.Duration) {
	workflowExecutionsTotal.WithLabelValues(workflowID, status).Inc()
	workflowExecutionDuration.WithLabelValues(workflowID).Observe(duration.Seconds())
}

// RecordDBQuery records a database query
func RecordDBQuery(operation string, duration time.Duration) {
	dbQueryDuration.WithLabelValues(operation).Observe(duration.Seconds())
}
