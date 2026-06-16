package errors

import (
	"fmt"
	"net/http"
)

// AppError 应用统一错误类型
type AppError struct {
	Code    int    `json:"code"`    // 业务错误码（对应HTTP状态码）
	Message string `json:"message"` // 错误信息
	Err     error `json:"-"`       // 原始错误（不序列化到JSON）
}

// Error 实现error接口
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Unwrap 支持errors.Is/As解包
func (e *AppError) Unwrap() error {
	return e.Err
}

// HTTPStatus 返回HTTP状态码
func (e *AppError) HTTPStatus() int {
	return e.Code
}

// 预定义标准错误
var (
	// 认证相关错误
	ErrUnauthorized = &AppError{Code: http.StatusUnauthorized, Message: "未授权，请先登录"}
	ErrTokenExpired = &AppError{Code: http.StatusUnauthorized, Message: "Token已过期，请重新登录"}
	ErrTokenInvalid = &AppError{Code: http.StatusUnauthorized, Message: "Token无效"}
	ErrForbidden    = &AppError{Code: http.StatusForbidden, Message: "无权限访问该资源"}

	// 限流相关错误
	ErrTooManyRequests = &AppError{Code: http.StatusTooManyRequests, Message: "请求过于频繁，请稍后再试"}

	// 参数校验错误
	ErrBadRequest = &AppError{Code: http.StatusBadRequest, Message: "请求参数错误"}

	// 服务端错误
	ErrInternalServer      = &AppError{Code: http.StatusInternalServerError, Message: "服务器内部错误"}
	ErrServiceUnavailable  = &AppError{Code: http.StatusServiceUnavailable, Message: "服务暂时不可用，请稍后重试"}
	ErrBadGateway          = &AppError{Code: http.StatusBadGateway, Message: "上游服务不可达"}

	// 认证业务错误
	ErrInvalidCredentials = &AppError{Code: http.StatusUnauthorized, Message: "用户名或密码错误"}
	ErrUserNotFound       = &AppError{Code: http.StatusNotFound, Message: "用户不存在"}
	ErrUserDisabled       = &AppError{Code: http.StatusForbidden, Message: "账号已被禁用"}
)

// New 创建新的应用错误
func New(code int, message string, err error) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}

// Wrap 包装已有错误
func Wrap(err error, code int, message string) *AppError {
	return &AppError{Code: code, Message: message, Err: err}
}
