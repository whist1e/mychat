package respond

import "time"

// UserInfoResponse 用户信息响应结构
// 用于返回给客户端的用户信息
type UserInfoResponse struct {
	ID           uint      `json:"id"`             // 用户ID
	UUID         string    `json:"uuid"`           // 用户唯一标识
	Nickname     string    `json:"nickname"`       // 昵称
	Telephone    string    `json:"telephone"`      // 手机号
	Email        string    `json:"email"`          // 邮箱
	Avatar       string    `json:"avatar"`         // 头像URL
	Gender       int8      `json:"gender"`         // 性别：0-男，1-女
	Signature    string    `json:"signature"`      // 个性签名
	Birthday     string    `json:"birthday"`       // 生日
	CreatedAt    time.Time `json:"created_at"`     // 创建时间
	LastOnlineAt time.Time `json:"last_online_at"` // 最后在线时间
	Status       int8      `json:"status"`         // 状态：0-正常，1-禁用
}

// LoginResponse 登录响应结构
// 用于返回给客户端的登录成功信息
type LoginResponse struct {
	UserInfo  UserInfoResponse `json:"user_info"`  // 用户信息
	Token     string           `json:"token"`      // 登录令牌
	ExpiresAt time.Time        `json:"expires_at"` // 令牌过期时间
}

// RegisterResponse 注册响应结构
// 用于返回给客户端的注册成功信息
type RegisterResponse struct {
	UserInfo UserInfoResponse `json:"user_info"` // 用户信息
	Message  string           `json:"message"`   // 成功消息
}
