package middleware

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/aiflowy/aiflowy-go/internal/config"
	apierrors "github.com/aiflowy/aiflowy-go/internal/errors"
	"github.com/aiflowy/aiflowy-go/pkg/logger"
	"github.com/aiflowy/aiflowy-go/pkg/response"
	"go.uber.org/zap"
)

// SetupMiddleware configures all middleware for the Echo instance
func SetupMiddleware(e *echo.Echo) {
	// Custom error handler
	e.HTTPErrorHandler = CustomHTTPErrorHandler

	// Recovery middleware with custom handler
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         4 << 10, // 4 KB
		DisableStackAll:   false,
		DisablePrintStack: false,
		LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
			logger.Error("Panic recovered",
				zap.Error(err),
				zap.String("stack", string(stack)),
				zap.String("uri", c.Request().RequestURI),
			)
			return nil
		},
	}))

	// Request ID middleware
	e.Use(middleware.RequestID())

	// Custom logger middleware
	e.Use(LoggerMiddleware())

	// CORS middleware
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE, echo.OPTIONS},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
			"X-Requested-With",
			"X-Token",
			"X-Request-Id",
		},
		AllowCredentials: true,
		MaxAge:           86400,
	}))

	// Secure middleware (skip in development)
	if config.IsProduction() {
		e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
			XSSProtection:         "1; mode=block",
			ContentTypeNosniff:    "nosniff",
			XFrameOptions:         "SAMEORIGIN",
			HSTSMaxAge:            3600,
			ContentSecurityPolicy: "",
		}))
	}

	// Gzip compression
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
		Skipper: func(c echo.Context) bool {
			// Skip SSE endpoints
			return c.Request().Header.Get("Accept") == "text/event-stream"
		},
	}))
}

// CustomHTTPErrorHandler handles HTTP errors with unified response format
// 遵循 API 响应格式规范 v1.1
func CustomHTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	// Handle BusinessError
	if be, ok := err.(*apierrors.BusinessError); ok {
		httpStatus := be.HTTPStatus
		if httpStatus == 0 {
			httpStatus = http.StatusOK
		}
		c.JSON(httpStatus, response.Response{
			Code:    be.Code,
			Message: be.Message,
		})
		return
	}

	// Handle Echo HTTPError
	if he, ok := err.(*echo.HTTPError); ok {
		code := he.Code
		message := "未知错误"
		if m, ok := he.Message.(string); ok {
			message = m
		}
		c.JSON(code, response.Response{
			Code:    code,
			Message: message,
		})
		return
	}

	// Handle other errors - log details, return generic message
	logger.Error("Unhandled error",
		zap.Error(err),
		zap.String("uri", c.Request().RequestURI),
		zap.String("method", c.Request().Method),
		zap.String("client_ip", c.RealIP()),
		zap.String("request_id", c.Response().Header().Get(echo.HeaderXRequestID)),
	)
	c.JSON(http.StatusInternalServerError, response.Response{
		Code:    http.StatusInternalServerError,
		Message: "服务器内部错误",
	})
}

// LoggerMiddleware returns a custom logger middleware
func LoggerMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			req := c.Request()
			res := c.Response()

			// Process request
			err := next(c)
			if err != nil {
				c.Error(err)
			}

			// Log request
			latency := time.Since(start)

			// Choose log level based on status
			status := res.Status
			fields := []zap.Field{
				zap.String("method", req.Method),
				zap.String("uri", req.RequestURI),
				zap.Int("status", status),
				zap.Duration("latency", latency),
				zap.String("ip", c.RealIP()),
				zap.String("request_id", res.Header().Get(echo.HeaderXRequestID)),
			}

			if status >= 500 {
				logger.Error("HTTP Request", fields...)
			} else if status >= 400 {
				logger.Warn("HTTP Request", fields...)
			} else {
				logger.Info("HTTP Request", fields...)
			}

			return nil
		}
	}
}

// SkipperFunc returns a skipper function for middleware
type SkipperFunc func(c echo.Context) bool

// PathSkipper returns a skipper that skips specified paths
func PathSkipper(paths ...string) SkipperFunc {
	pathMap := make(map[string]bool)
	for _, p := range paths {
		pathMap[p] = true
	}
	return func(c echo.Context) bool {
		return pathMap[c.Path()]
	}
}
