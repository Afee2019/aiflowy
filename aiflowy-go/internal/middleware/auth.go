package middleware

import (
	"strings"

	"github.com/labstack/echo/v4"

	apierrors "github.com/aiflowy/aiflowy-go/internal/errors"
	"github.com/aiflowy/aiflowy-go/pkg/jwt"
	"github.com/aiflowy/aiflowy-go/pkg/response"
)

const (
	// AuthHeader is the header name for token (compatible with Java version)
	AuthHeader = "aiflowy-token"
	// AuthBearerHeader is the standard Authorization header
	AuthBearerHeader = "Authorization"
	// BearerPrefix is the prefix for Bearer tokens
	BearerPrefix = "Bearer "
)

// JWTAuth returns a JWT authentication middleware
func JWTAuth() echo.MiddlewareFunc {
	return JWTAuthWithConfig(JWTConfig{})
}

// JWTConfig defines the config for JWT middleware
type JWTConfig struct {
	// Skipper defines a function to skip middleware
	Skipper func(c echo.Context) bool
}

// JWTAuthWithConfig returns a JWT authentication middleware with config
func JWTAuthWithConfig(config JWTConfig) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check skipper
			if config.Skipper != nil && config.Skipper(c) {
				return next(c)
			}

			// Extract token from request
			tokenString := extractToken(c)
			if tokenString == "" {
				return c.JSON(401, response.Response{
					Code:    401,
					Message: "请先登录",
				})
			}

			// Parse and validate token
			claims, err := jwt.ParseToken(tokenString)
			if err != nil {
				if err == jwt.ErrTokenExpired {
					return c.JSON(401, response.Response{
						Code:    401,
						Message: "登录已过期，请重新登录",
					})
				}
				return c.JSON(401, response.Response{
					Code:    401,
					Message: "无效的登录凭证",
				})
			}

			// Store claims in context for later use
			c.Set("claims", claims)
			c.Set("userId", claims.UserID)
			c.Set("loginName", claims.LoginName)
			c.Set("tenantId", claims.TenantID)

			return next(c)
		}
	}
}

// extractToken extracts token from request headers
// Supports both "aiflowy-token" header and "Authorization: Bearer <token>" format
func extractToken(c echo.Context) string {
	// Try aiflowy-token header first (compatible with Java version)
	token := c.Request().Header.Get(AuthHeader)
	if token != "" {
		return token
	}

	// Try Authorization header with Bearer prefix
	auth := c.Request().Header.Get(AuthBearerHeader)
	if strings.HasPrefix(auth, BearerPrefix) {
		return strings.TrimPrefix(auth, BearerPrefix)
	}

	// Try X-Token header
	token = c.Request().Header.Get("X-Token")
	if token != "" {
		return token
	}

	// Try query parameter (for WebSocket/SSE)
	token = c.QueryParam("token")
	if token != "" {
		return token
	}

	return ""
}

// RequireAuth returns a middleware that requires authentication
// This is an alias for JWTAuth() for clearer naming
func RequireAuth() echo.MiddlewareFunc {
	return JWTAuth()
}

// RequireRole returns a middleware that requires specific roles
func RequireRole(roles ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// First check if user is authenticated
			claims, ok := c.Get("claims").(*jwt.Claims)
			if !ok || claims == nil {
				return apierrors.Unauthorized("请先登录")
			}

			// TODO: Implement role checking from database or cache
			// For now, just allow all authenticated users
			// This can be enhanced with Casbin or custom role checking

			return next(c)
		}
	}
}

// RequirePermission returns a middleware that requires specific permissions
func RequirePermission(permissions ...string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// First check if user is authenticated
			claims, ok := c.Get("claims").(*jwt.Claims)
			if !ok || claims == nil {
				return apierrors.Unauthorized("请先登录")
			}

			// TODO: Implement permission checking from database or cache
			// For now, just allow all authenticated users
			// This can be enhanced with Casbin

			return next(c)
		}
	}
}
