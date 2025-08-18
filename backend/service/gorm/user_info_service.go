package gorm

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"mychat-backend/internal/dao"
	"mychat-backend/internal/dto/request"
	"mychat-backend/internal/dto/respond"
	"mychat-backend/internal/model"
	"time"

	"gorm.io/gorm"
)

type UserInfoService struct {
}

// NewUserInfoService 创建用户信息服务实例
func NewUserInfoService() *UserInfoService {
	return &UserInfoService{}
}

// Register 用户注册
// 参数：注册请求信息
// 返回：成功消息、用户信息、错误信息
func (s *UserInfoService) Register(req request.RegisterRequest) (message string, userInfo respond.UserInfoResponse, err error) {
	// 1. 检查手机号是否已存在
	var existingUser model.UserInfo
	err = dao.GormDB.Where("telephone = ?", req.Telephone).First(&existingUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return "", respond.UserInfoResponse{}, fmt.Errorf("查询用户失败: %v", err)
	}
	if err == nil {
		return "", respond.UserInfoResponse{}, fmt.Errorf("手机号已存在")
	}

	// 2. 创建新用户
	newUser := &model.UserInfo{
		Uuid:         generateUUID(),
		Nickname:     req.Nickname,
		Telephone:    req.Telephone,
		Password:     hashPassword(req.Password),
		Avatar:       "https://cube.elemecdn.com/0/88/03b0d39583f48206768a7534e55bcpng.png", // 默认头像
		Gender:       0,                                                                     // 默认性别
		Status:       0,                                                                     // 正常状态
		CreatedAt:    time.Now(),
		LastOnlineAt: sql.NullTime{Time: time.Now(), Valid: true},
	}

	// 3. 保存到数据库
	if err := dao.GormDB.Create(newUser).Error; err != nil {
		return "", respond.UserInfoResponse{}, fmt.Errorf("创建用户失败: %v", err)
	}

	// 4. 转换为响应格式
	userInfo = respond.UserInfoResponse{
		ID:           uint(newUser.Id),
		UUID:         newUser.Uuid,
		Nickname:     newUser.Nickname,
		Telephone:    newUser.Telephone,
		Email:        newUser.Email,
		Avatar:       newUser.Avatar,
		Gender:       newUser.Gender,
		Signature:    newUser.Signature,
		Birthday:     newUser.Birthday,
		CreatedAt:    newUser.CreatedAt,
		LastOnlineAt: newUser.LastOnlineAt.Time,
		Status:       newUser.Status,
	}

	return "注册成功", userInfo, nil
}

// Login 用户登录（密码登录）
// 参数：登录请求信息
// 返回：成功消息、用户信息、错误信息
func (s *UserInfoService) Login(req request.LoginRequest) (message string, userInfo respond.UserInfoResponse, err error) {
	// 1. 根据手机号查找用户
	var user model.UserInfo
	if err := dao.GormDB.Where("telephone = ?", req.Telephone).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", respond.UserInfoResponse{}, fmt.Errorf("用户不存在")
		}
		return "", respond.UserInfoResponse{}, fmt.Errorf("查询用户失败: %v", err)
	}

	// 2. 检查用户状态
	if user.Status != 0 {
		return "", respond.UserInfoResponse{}, fmt.Errorf("用户已被禁用")
	}

	// 3. 验证密码
	if !verifyPassword(req.Password, user.Password) {
		return "", respond.UserInfoResponse{}, fmt.Errorf("密码错误")
	}

	// 4. 更新最后在线时间
	if err := dao.GormDB.Model(&user).Update("last_online_at", time.Now()).Error; err != nil {
		// 记录错误但不影响登录
		fmt.Printf("更新最后在线时间失败: %v\n", err)
	}

	// 5. 转换为响应格式
	userInfo = respond.UserInfoResponse{
		ID:           uint(user.Id),
		UUID:         user.Uuid,
		Nickname:     user.Nickname,
		Telephone:    user.Telephone,
		Email:        user.Email,
		Avatar:       user.Avatar,
		Gender:       user.Gender,
		Signature:    user.Signature,
		Birthday:     user.Birthday,
		CreatedAt:    user.CreatedAt,
		LastOnlineAt: user.LastOnlineAt.Time,
		Status:       user.Status,
	}

	return "登录成功", userInfo, nil
}

// PasswordLogin 密码登录（新增方法）
// 参数：手机号和密码
// 返回：成功消息、用户信息、错误信息
func (s *UserInfoService) PasswordLogin(telephone, password string) (message string, userInfo respond.UserInfoResponse, err error) {
	// 1. 根据手机号查找用户
	var user model.UserInfo
	if err := dao.GormDB.Where("telephone = ?", telephone).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "", respond.UserInfoResponse{}, fmt.Errorf("用户不存在")
		}
		return "", respond.UserInfoResponse{}, fmt.Errorf("查询用户失败: %v", err)
	}

	// 2. 检查用户状态
	if user.Status != 0 {
		return "", respond.UserInfoResponse{}, fmt.Errorf("用户已被禁用")
	}

	// 3. 验证密码
	if !verifyPassword(password, user.Password) {
		return "", respond.UserInfoResponse{}, fmt.Errorf("密码错误")
	}

	// 4. 更新最后在线时间
	if err := dao.GormDB.Model(&user).Update("last_online_at", time.Now()).Error; err != nil {
		fmt.Printf("更新最后在线时间失败: %v\n", err)
	}

	// 5. 转换为响应格式
	userInfo = respond.UserInfoResponse{
		ID:           uint(user.Id),
		UUID:         user.Uuid,
		Nickname:     user.Nickname,
		Telephone:    user.Telephone,
		Email:        user.Email,
		Avatar:       user.Avatar,
		Gender:       user.Gender,
		Signature:    user.Signature,
		Birthday:     user.Birthday,
		CreatedAt:    user.CreatedAt,
		LastOnlineAt: user.LastOnlineAt.Time,
		Status:       user.Status,
	}

	return "登录成功", userInfo, nil
}

// 工具函数

// generateUUID 生成UUID
func generateUUID() string {
	// 使用时间戳生成简单的唯一标识，确保不超过20个字符
	// 使用纳秒时间戳的后8位 + 秒时间戳的后8位 + 随机数4位
	nanos := time.Now().UnixNano()
	seconds := time.Now().Unix()

	// 取后8位数字，确保长度可控
	nanosStr := fmt.Sprintf("%08d", nanos%100000000)
	secondsStr := fmt.Sprintf("%08d", seconds%100000000)

	// 生成4位随机数
	randNum := fmt.Sprintf("%04d", (nanos+seconds)%10000)

	return nanosStr + secondsStr + randNum
}

// hashPassword 密码加密
func hashPassword(password string) string {
	hash := md5.Sum([]byte(password))
	return fmt.Sprintf("%x", hash)
}

// verifyPassword 验证密码
func verifyPassword(inputPassword, hashedPassword string) bool {
	return hashPassword(inputPassword) == hashedPassword
}
