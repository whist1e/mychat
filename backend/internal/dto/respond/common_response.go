package respond

// CommonResponse 通用API响应结构
// 用于统一所有API接口的响应格式
type CommonResponse struct {
	Code    int         `json:"code"`            // 响应状态码：200-成功，400-参数错误，500-服务器错误
	Message string      `json:"message"`         // 响应消息
	Data    interface{} `json:"data,omitempty"`  // 响应数据，可选字段
	Error   string      `json:"error,omitempty"` // 错误信息，可选字段
}

// SuccessResponse 成功响应
// 用于返回成功的API响应
func SuccessResponse(data interface{}, message string) CommonResponse {
	return CommonResponse{
		Code:    200,
		Message: message,
		Data:    data,
	}
}

// ErrorResponse 错误响应
// 用于返回错误的API响应
func ErrorResponse(code int, message string, error string) CommonResponse {
	return CommonResponse{
		Code:    code,
		Message: message,
		Error:   error,
	}
}

// BadRequestResponse 参数错误响应
// 用于返回400状态码的响应
func BadRequestResponse(message string) CommonResponse {
	return ErrorResponse(400, message, "Bad Request")
}

// ServerErrorResponse 服务器错误响应
// 用于返回500状态码的响应
func ServerErrorResponse(message string) CommonResponse {
	return ErrorResponse(500, message, "Internal Server Error")
}
