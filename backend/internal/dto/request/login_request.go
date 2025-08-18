package request

type LoginRequest struct {
	Telephone string `json:"telephone" binding:"required,len=11" validate:"required,len=11"` // 手机号码，必填，11位
	Password  string `json:"password" binding:"required,min=6" validate:"required,min=6"`    // 密码，必填，最少6位
}