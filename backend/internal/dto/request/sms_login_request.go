package request

// SmsLoginRequest 短信登录请求结构
// 用于接收客户端发送的短信登录信息
type SmsLoginRequest struct {
	Telephone string `json:"telephone" binding:"required,len=11"` // 手机号码，必填，11位
	SmsCode   string `json:"sms_code" binding:"required,len=6"`   // 短信验证码，必填，6位
}