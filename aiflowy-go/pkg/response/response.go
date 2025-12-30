package response

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Response 统一响应结构 (遵循 API 响应格式规范 v1.1)
// code = 0 表示成功，非零表示错误（与 HTTP 状态码一致）
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PageData 分页响应数据结构
type PageData struct {
	Items interface{} `json:"items"`
	Total int64       `json:"total"`
	Page  int         `json:"page"`
	Size  int         `json:"size"`
}

// ValidationError 字段验证错误
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code,omitempty"`
}

// ValidationErrorData 验证错误响应数据
type ValidationErrorData struct {
	Errors []ValidationError `json:"errors"`
}

// ----- 成功响应 -----

// Success 成功响应 (code=0)
func Success(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// SuccessWithMessage 带自定义消息的成功响应
func SuccessWithMessage(c echo.Context, message string, data interface{}) error {
	return c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: message,
		Data:    data,
	})
}

// Created 201 - 资源创建成功
func Created(c echo.Context, data interface{}) error {
	return c.JSON(http.StatusCreated, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// NoContent 204 - 删除成功（无响应体）
func NoContent(c echo.Context) error {
	return c.NoContent(http.StatusNoContent)
}

// PageSuccess 分页列表响应
func PageSuccess(c echo.Context, items interface{}, total int64, page, size int) error {
	return c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data: PageData{
			Items: items,
			Total: total,
			Page:  page,
			Size:  size,
		},
	})
}

// ----- 错误响应 -----

// Error 通用错误响应
func Error(c echo.Context, httpStatus int, message string) error {
	return c.JSON(httpStatus, Response{
		Code:    httpStatus,
		Message: message,
	})
}

// BadRequest 400 - 请求参数错误
func BadRequest(c echo.Context, message string) error {
	return c.JSON(http.StatusBadRequest, Response{
		Code:    http.StatusBadRequest,
		Message: message,
	})
}

// Unauthorized 401 - 认证失败
func Unauthorized(c echo.Context, message string) error {
	return c.JSON(http.StatusUnauthorized, Response{
		Code:    http.StatusUnauthorized,
		Message: message,
	})
}

// Forbidden 403 - 授权失败
func Forbidden(c echo.Context, message string) error {
	return c.JSON(http.StatusForbidden, Response{
		Code:    http.StatusForbidden,
		Message: message,
	})
}

// NotFound 404 - 资源不存在
func NotFound(c echo.Context, message string) error {
	return c.JSON(http.StatusNotFound, Response{
		Code:    http.StatusNotFound,
		Message: message,
	})
}

// Conflict 409 - 资源冲突
func Conflict(c echo.Context, message string) error {
	return c.JSON(http.StatusConflict, Response{
		Code:    http.StatusConflict,
		Message: message,
	})
}

// UnprocessableEntity 422 - 验证错误（带详情）
func UnprocessableEntity(c echo.Context, errors []ValidationError) error {
	return c.JSON(http.StatusUnprocessableEntity, Response{
		Code:    http.StatusUnprocessableEntity,
		Message: "请求参数验证失败",
		Data:    ValidationErrorData{Errors: errors},
	})
}

// TooManyRequests 429 - 请求频率超限
func TooManyRequests(c echo.Context, message string) error {
	return c.JSON(http.StatusTooManyRequests, Response{
		Code:    http.StatusTooManyRequests,
		Message: message,
	})
}

// InternalError 500 - 服务器内部错误
func InternalError(c echo.Context, message string) error {
	return c.JSON(http.StatusInternalServerError, Response{
		Code:    http.StatusInternalServerError,
		Message: message,
	})
}

// ServiceUnavailable 503 - 服务暂时不可用
func ServiceUnavailable(c echo.Context, message string) error {
	return c.JSON(http.StatusServiceUnavailable, Response{
		Code:    http.StatusServiceUnavailable,
		Message: message,
	})
}

// ----- 兼容旧版本的别名 -----

// Result 是 Response 的别名，保持兼容性
type Result = Response

// OK 是 Success 的别名
func OK(c echo.Context, data interface{}) error {
	return Success(c, data)
}

// OKWithMessage 是 SuccessWithMessage 的别名
func OKWithMessage(c echo.Context, message string, data interface{}) error {
	return SuccessWithMessage(c, message, data)
}

// Fail 是 Error 的别名
func Fail(c echo.Context, code int, message string) error {
	return c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
	})
}
