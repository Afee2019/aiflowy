package middleware

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"

	apierrors "github.com/aiflowy/aiflowy-go/internal/errors"
	"github.com/aiflowy/aiflowy-go/pkg/logger"
)

func init() {
	// Initialize logger for tests
	logger.Init(&logger.Config{
		Level:  "error",
		Format: "console",
		Output: "console",
	})
}

func TestCustomHTTPErrorHandler_BusinessError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := apierrors.New(1001, "test business error")
	CustomHTTPErrorHandler(err, c)

	if rec.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rec.Code)
	}

	// Check response body contains error code
	body := rec.Body.String()
	if body == "" {
		t.Error("expected non-empty response body")
	}
}

func TestCustomHTTPErrorHandler_BusinessErrorWithHTTPStatus(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := apierrors.NewWithHTTPStatus(401, "unauthorized", http.StatusUnauthorized)
	CustomHTTPErrorHandler(err, c)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("expected status 401, got %d", rec.Code)
	}
}

func TestCustomHTTPErrorHandler_EchoHTTPError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := echo.NewHTTPError(http.StatusNotFound, "resource not found")
	CustomHTTPErrorHandler(err, c)

	if rec.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rec.Code)
	}
}

func TestCustomHTTPErrorHandler_GenericError(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := errors.New("generic error")
	CustomHTTPErrorHandler(err, c)

	if rec.Code != http.StatusInternalServerError {
		t.Errorf("expected status 500, got %d", rec.Code)
	}
}

func TestCustomHTTPErrorHandler_AlreadyCommitted(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Write something to commit the response
	c.Response().Write([]byte("already written"))

	err := errors.New("error after commit")
	CustomHTTPErrorHandler(err, c)

	// Should return early without changing anything
	// Just verify it doesn't panic
}

func TestLoggerMiddleware(t *testing.T) {
	e := echo.New()
	e.Use(LoggerMiddleware())

	tests := []struct {
		name           string
		handler        echo.HandlerFunc
		expectedStatus int
	}{
		{
			name: "successful request",
			handler: func(c echo.Context) error {
				return c.String(http.StatusOK, "ok")
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "client error",
			handler: func(c echo.Context) error {
				return echo.NewHTTPError(http.StatusBadRequest, "bad request")
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "server error",
			handler: func(c echo.Context) error {
				return echo.NewHTTPError(http.StatusInternalServerError, "server error")
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e.GET("/test", tt.handler)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			if rec.Code != tt.expectedStatus {
				t.Errorf("expected status %d, got %d", tt.expectedStatus, rec.Code)
			}
		})
	}
}

func TestPathSkipper(t *testing.T) {
	tests := []struct {
		name           string
		paths          []string
		requestPath    string
		expectedResult bool
	}{
		{
			name:           "skip matching path",
			paths:          []string{"/health", "/metrics"},
			requestPath:    "/health",
			expectedResult: true,
		},
		{
			name:           "don't skip non-matching path",
			paths:          []string{"/health", "/metrics"},
			requestPath:    "/api/users",
			expectedResult: false,
		},
		{
			name:           "empty paths",
			paths:          []string{},
			requestPath:    "/any",
			expectedResult: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			skipper := PathSkipper(tt.paths...)

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, tt.requestPath, nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetPath(tt.requestPath)

			result := skipper(c)
			if result != tt.expectedResult {
				t.Errorf("expected %v, got %v", tt.expectedResult, result)
			}
		})
	}
}

func TestSkipperFunc(t *testing.T) {
	// Test that SkipperFunc is correctly typed
	var skipper SkipperFunc = func(c echo.Context) bool {
		return c.Request().Method == http.MethodOptions
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodOptions, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	if !skipper(c) {
		t.Error("expected skipper to return true for OPTIONS request")
	}

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	c = e.NewContext(req, rec)
	if skipper(c) {
		t.Error("expected skipper to return false for GET request")
	}
}

func BenchmarkCustomHTTPErrorHandler_BusinessError(b *testing.B) {
	e := echo.New()
	err := apierrors.New(1001, "test error")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)
		CustomHTTPErrorHandler(err, c)
	}
}

func BenchmarkLoggerMiddleware(b *testing.B) {
	e := echo.New()
	e.Use(LoggerMiddleware())
	e.GET("/test", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/test", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
	}
}

func BenchmarkPathSkipper(b *testing.B) {
	skipper := PathSkipper("/health", "/metrics", "/ready", "/live")
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/health")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		skipper(c)
	}
}
