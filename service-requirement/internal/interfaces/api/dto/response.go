package dto

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一API响应结构
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// Success 成功响应
func Success(data interface{}) Response {
	return Response{
		Code:    http.StatusOK,
		Message: "success",
		Data:    data,
	}
}

// PageSuccess 分页成功响应
func PageSuccess(list interface{}, total int64, page, size int) Response {
	return Response{
		Code:    http.StatusOK,
		Message: "success",
		Data: gin.H{
			"list":  list,
			"total": total,
			"page":  page,
			"size":  size,
		},
	}
}

// Error 错误响应
func Error(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
	}
}

// BadRequest 参数错误响应
func BadRequest(message string) Response {
	return Error(http.StatusBadRequest, message)
}

// NotFound 未找到响�?func NotFound(message string) Response {
	return Error(http.StatusNotFound, message)
}

// InternalError 内部错误响应
func InternalError(message string) Response {
	return Error(http.StatusInternalServerError, message)
}
