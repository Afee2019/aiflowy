package errors

import (
	"fmt"
	"net/http"
)

// BusinessError represents a business logic error
type BusinessError struct {
	Code       int    `json:"code"`
	Message    string `json:"message"`
	HTTPStatus int    `json:"-"`
}

func (e *BusinessError) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

// New creates a new BusinessError
func New(code int, message string) *BusinessError {
	return &BusinessError{
		Code:       code,
		Message:    message,
		HTTPStatus: http.StatusOK,
	}
}

// NewWithHTTPStatus creates a new BusinessError with custom HTTP status
func NewWithHTTPStatus(code int, message string, httpStatus int) *BusinessError {
	return &BusinessError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// Common error codes (compatible with Java version)
const (
	CodeSuccess          = 200
	CodeBadRequest       = 400
	CodeUnauthorized     = 401
	CodeForbidden        = 403
	CodeNotFound         = 404
	CodeInternalError    = 500
	CodeServiceUnavailable = 503

	// Business error codes (1000+)
	CodeParamError       = 1001
	CodeDataNotFound     = 1002
	CodeDataExists       = 1003
	CodeOperationFailed  = 1004
	CodePermissionDenied = 1005

	// Auth error codes (2000+)
	CodeTokenExpired     = 2001
	CodeTokenInvalid     = 2002
	CodeLoginFailed      = 2003
	CodeCaptchaError     = 2004
	CodeAccountDisabled  = 2005
)

// Predefined errors
var (
	ErrBadRequest       = New(CodeBadRequest, "请求参数错误")
	ErrUnauthorized     = NewWithHTTPStatus(CodeUnauthorized, "未授权访问", http.StatusUnauthorized)
	ErrForbidden        = NewWithHTTPStatus(CodeForbidden, "禁止访问", http.StatusForbidden)
	ErrNotFound         = New(CodeNotFound, "资源不存在")
	ErrInternalError    = New(CodeInternalError, "服务器内部错误")

	ErrParamError       = New(CodeParamError, "参数错误")
	ErrDataNotFound     = New(CodeDataNotFound, "数据不存在")
	ErrDataExists       = New(CodeDataExists, "数据已存在")
	ErrOperationFailed  = New(CodeOperationFailed, "操作失败")
	ErrPermissionDenied = New(CodePermissionDenied, "权限不足")

	ErrTokenExpired     = NewWithHTTPStatus(CodeTokenExpired, "登录已过期，请重新登录", http.StatusUnauthorized)
	ErrTokenInvalid     = NewWithHTTPStatus(CodeTokenInvalid, "无效的令牌", http.StatusUnauthorized)
	ErrLoginFailed      = New(CodeLoginFailed, "用户名或密码错误")
	ErrCaptchaError     = New(CodeCaptchaError, "验证码错误")
	ErrAccountDisabled  = New(CodeAccountDisabled, "账号已被禁用")
)

// Wrap wraps an error with additional message
func Wrap(err error, message string) *BusinessError {
	if be, ok := err.(*BusinessError); ok {
		return &BusinessError{
			Code:       be.Code,
			Message:    message + ": " + be.Message,
			HTTPStatus: be.HTTPStatus,
		}
	}
	return New(CodeInternalError, message+": "+err.Error())
}

// Is checks if the error is a specific BusinessError
func Is(err error, target *BusinessError) bool {
	if be, ok := err.(*BusinessError); ok {
		return be.Code == target.Code
	}
	return false
}

// Helper functions for creating common errors

// BadRequest creates a 400 Bad Request error
func BadRequest(message string) *BusinessError {
	return NewWithHTTPStatus(CodeBadRequest, message, http.StatusBadRequest)
}

// Unauthorized creates a 401 Unauthorized error
func Unauthorized(message string) *BusinessError {
	return NewWithHTTPStatus(CodeUnauthorized, message, http.StatusUnauthorized)
}

// Forbidden creates a 403 Forbidden error
func Forbidden(message string) *BusinessError {
	return NewWithHTTPStatus(CodeForbidden, message, http.StatusForbidden)
}

// NotFound creates a 404 Not Found error
func NotFound(message string) *BusinessError {
	return NewWithHTTPStatus(CodeNotFound, message, http.StatusNotFound)
}

// InternalError creates a 500 Internal Server Error
func InternalError(message string) *BusinessError {
	return NewWithHTTPStatus(CodeInternalError, message, http.StatusInternalServerError)
}
