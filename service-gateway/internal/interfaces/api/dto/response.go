package dto

// Response 统一API响应格式
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginatedResponse 分页响应格式
type PaginatedResponse struct {
	List     interface{} `json:"list"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"pageSize"`
}

// Success 成功响应
func Success(data interface{}) Response {
	return Response{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}

// SuccessWithMessage 带自定义消息的成功响应
func SuccessWithMessage(message string, data interface{}) Response {
	return Response{
		Code:    200,
		Message: message,
		Data:    data,
	}
}

// Fail 错误响应
func Fail(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
	}
}
