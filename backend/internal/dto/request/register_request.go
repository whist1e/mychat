package request

// RegisterRequest 用户注册请求结构
// 用于接收客户端发送的用户注册信息
type RegisterRequest struct {
	Telephone string `json:"telephone" binding:"required,len=11"`            // 手机号码，必填，11位
	Password  string `json:"password" binding:"required,min=6,max=20"` // 密码，必填，6-20位
	Nickname  string `json:"nickname" binding:"required,min=2,max=20"` // 昵称，必填，2-20位
	SmsCode   string `json:"sms_code" binding:"required,len=6"`         // 短信验证码，必填，6位
}
