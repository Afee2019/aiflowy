package auth

import (
	"strconv"

	"github.com/labstack/echo/v4"

	"github.com/aiflowy/aiflowy-go/internal/dto"
	apierrors "github.com/aiflowy/aiflowy-go/internal/errors"
	"github.com/aiflowy/aiflowy-go/internal/service"
	"github.com/aiflowy/aiflowy-go/pkg/jwt"
	"github.com/aiflowy/aiflowy-go/pkg/response"
)

// Handler handles authentication endpoints
type Handler struct {
	authService *service.AuthService
}

// NewHandler creates a new auth handler
func NewHandler() *Handler {
	return &Handler{
		authService: service.NewAuthService(),
	}
}

// Login handles user login
// POST /api/v1/auth/login
func (h *Handler) Login(c echo.Context) error {
	var req dto.LoginRequest
	if err := c.Bind(&req); err != nil {
		return apierrors.BadRequest("无效的请求参数")
	}

	// Validate required fields
	if req.Account == "" {
		return apierrors.BadRequest("账号不能为空")
	}
	if req.Password == "" {
		return apierrors.BadRequest("密码不能为空")
	}

	resp, err := h.authService.Login(c.Request().Context(), &req)
	if err != nil {
		return err
	}

	return response.Success(c, resp)
}

// Logout handles user logout
// POST /api/v1/auth/logout
func (h *Handler) Logout(c echo.Context) error {
	// For JWT, client just needs to delete the token
	// Server-side we could implement a token blacklist if needed
	return response.Success(c, nil)
}

// GetPermissions returns the permission list for current user
// GET /api/v1/auth/getPermissions
func (h *Handler) GetPermissions(c echo.Context) error {
	userID, err := getCurrentUserID(c)
	if err != nil {
		return apierrors.Unauthorized("请先登录")
	}

	permissions, err := h.authService.GetPermissions(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	return response.Success(c, permissions)
}

// GetRoles returns the role list for current user
// GET /api/v1/auth/getRoles
func (h *Handler) GetRoles(c echo.Context) error {
	userID, err := getCurrentUserID(c)
	if err != nil {
		return apierrors.Unauthorized("请先登录")
	}

	roles, err := h.authService.GetRoles(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	return response.Success(c, roles)
}

// GetMenus returns the menu list for current user
// GET /api/v1/auth/getMenus
func (h *Handler) GetMenus(c echo.Context) error {
	userID, err := getCurrentUserID(c)
	if err != nil {
		return apierrors.Unauthorized("请先登录")
	}

	menus, err := h.authService.GetMenus(c.Request().Context(), userID)
	if err != nil {
		return err
	}

	return response.Success(c, menus)
}

// GetUserInfo returns the current user info
// GET /api/v1/auth/getUserInfo
func (h *Handler) GetUserInfo(c echo.Context) error {
	claims := GetClaims(c)
	if claims == nil {
		return apierrors.Unauthorized("请先登录")
	}

	userID, _ := strconv.ParseInt(claims.UserID, 10, 64)
	account, err := h.authService.GetAccountByID(c.Request().Context(), userID)
	if err != nil || account == nil {
		return apierrors.Unauthorized("用户不存在")
	}

	// Return user info without sensitive data
	return response.Success(c, map[string]interface{}{
		"id":          account.ID,
		"loginName":   account.LoginName,
		"nickname":    account.Nickname,
		"avatar":      account.Avatar,
		"email":       account.Email,
		"mobile":      account.Mobile,
		"deptId":      account.DeptID,
		"accountType": account.AccountType,
	})
}

// getCurrentUserID extracts user ID from context
func getCurrentUserID(c echo.Context) (int64, error) {
	claims := GetClaims(c)
	if claims == nil {
		return 0, apierrors.Unauthorized("请先登录")
	}

	userID, err := strconv.ParseInt(claims.UserID, 10, 64)
	if err != nil {
		return 0, apierrors.Unauthorized("无效的用户ID")
	}

	return userID, nil
}

// GetClaims returns JWT claims from context
func GetClaims(c echo.Context) *jwt.Claims {
	claims, ok := c.Get("claims").(*jwt.Claims)
	if !ok {
		return nil
	}
	return claims
}

// GetCurrentUserID returns current user ID from context
func GetCurrentUserID(c echo.Context) int64 {
	claims := GetClaims(c)
	if claims == nil {
		return 0
	}
	userID, _ := strconv.ParseInt(claims.UserID, 10, 64)
	return userID
}

// GetCurrentTenantID returns current tenant ID from context
func GetCurrentTenantID(c echo.Context) int64 {
	claims := GetClaims(c)
	if claims == nil {
		return 0
	}
	tenantID, _ := strconv.ParseInt(claims.TenantID, 10, 64)
	return tenantID
}

// GetCurrentDeptID returns current department ID from context
func GetCurrentDeptID(c echo.Context) int64 {
	claims := GetClaims(c)
	if claims == nil {
		return 0
	}
	deptID, _ := strconv.ParseInt(claims.DeptID, 10, 64)
	return deptID
}

// GetUserContext returns userID, tenantID, deptID from context
func GetUserContext(c echo.Context) (userID, tenantID, deptID int64) {
	claims := GetClaims(c)
	if claims == nil {
		return 0, 0, 0
	}
	userID, _ = strconv.ParseInt(claims.UserID, 10, 64)
	tenantID, _ = strconv.ParseInt(claims.TenantID, 10, 64)
	deptID, _ = strconv.ParseInt(claims.DeptID, 10, 64)
	return
}
