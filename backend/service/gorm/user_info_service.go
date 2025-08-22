package gorm

import (
	"fmt"
	"crypto/md5"
	"database/sql"
	"mychat-backend/internal/dao"
	"mychat-backend/internal/dto/request"
	"mychat-backend/internal/dto/respond"
	"mychat-backend/internal/model"
	"mychat-backend/pkg/constants"
	"mychat-backend/pkg/zaplog"
	"time"

	"gorm.io/gorm"
)

type userInfoService struct{}

var UserInfoService = new(userInfoService)

// Register 用户注册
// 参数：注册请求信息
// 返回：成功消息、用户信息、错误信息
func (s *userInfoService) Register(req request.RegisterRequest) (string, *respond.RegisterRespond, int) {
	// 1. 检查手机号是否已存在
	var existingUser model.UserInfo
	err := dao.GormDB.Where("telephone = ?", req.Telephone).First(&existingUser).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return constants.SYSTEM_ERROR, nil, -2
	}
	if err == nil {
		return "手机号已存在", nil, -2
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
		return constants.SYSTEM_ERROR, nil, -1
	}

	// 4. 转换为响应格式
	userInfo := respond.RegisterRespond{
		Uuid:         newUser.Uuid,
		Nickname:     newUser.Nickname,
		Telephone:    newUser.Telephone,
		Email:        newUser.Email,
		Avatar:       newUser.Avatar,
		Gender:       newUser.Gender,
		Signature:    newUser.Signature,
		Birthday:     newUser.Birthday,
		CreatedAt: newUser.CreatedAt.Format("2006-01-02 15:04:05"),
		IsAdmin:   newUser.IsAdmin,
		Status:    newUser.Status,
	}

	return "注册成功", &userInfo, 0
}

// Login 用户登录（密码登录）
// 参数：登录请求信息
// 返回：成功消息、用户信息、错误信息
func (s *userInfoService) Login(req request.LoginRequest) (string, *respond.LoginRespond, int) {
	// 1. 根据手机号查找用户
	var user model.UserInfo
	if err := dao.GormDB.Where("telephone = ?", req.Telephone).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return constants.SYSTEM_ERROR, nil, -2
		}
		return constants.SYSTEM_ERROR, nil, -2
	}

	// 2. 检查用户状态
	if user.Status != 0 {
		return "用户已被禁用", nil, -2
	}

	// 3. 验证密码
	if !verifyPassword(req.Password, user.Password) {
		return "密码错误", nil, -2
	}

	// 4. 更新最后在线时间
	if err := dao.GormDB.Model(&user).Update("last_online_at", time.Now()).Error; err != nil {
		// 记录错误但不影响登录
		zaplog.Error(err.Error())
	}

	// 5. 转换为响应格式
	userInfo := respond.LoginRespond{
		Uuid:         user.Uuid,
		Nickname:     user.Nickname,
		Telephone:    user.Telephone,
		Email:        user.Email,
		Avatar:       user.Avatar,
		Gender:       user.Gender,
		Signature:    user.Signature,
		Birthday:     user.Birthday,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		IsAdmin:   user.IsAdmin,
		Status:    user.Status,
	}

	return "登录成功", &userInfo, 0
}

// PasswordLogin 密码登录（新增方法）
// 参数：手机号和密码
// 返回：成功消息、用户信息、错误信息
func (s *userInfoService) PasswordLogin(telephone, password string) (string, *respond.LoginRespond, int) {
	// 1. 根据手机号查找用户
	var user model.UserInfo
	if err := dao.GormDB.Where("telephone = ?", telephone).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return "用户不存在", nil, -2
		}
		return constants.SYSTEM_ERROR, nil, -2
	}

	// 2. 检查用户状态
	if user.Status != 0 {
		return "用户已被禁用", nil, -2
	}

	// 3. 验证密码
	if !verifyPassword(password, user.Password) {
		return "密码错误", nil, -2
	}

	// 4. 更新最后在线时间
	if err := dao.GormDB.Model(&user).Update("last_online_at", time.Now()).Error; err != nil {
		zaplog.Error(err.Error())
	}

	// 5. 转换为响应格式
	userInfo := respond.LoginRespond{
		Uuid:         user.Uuid,
		Nickname:     user.Nickname,
		Telephone:    user.Telephone,
		Email:        user.Email,
		Avatar:       user.Avatar,
		Gender:       user.Gender,
		Signature:    user.Signature,
		Birthday:     user.Birthday,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		IsAdmin:   user.IsAdmin,
		Status:    user.Status,
	}

	return "登录成功", &userInfo, 0
}

func (s *userInfoService) UpdateUserInfo(req request.GetUserInfo) (string, *respond.GetUserInfoRespond, int) {
	return "更新用户信息成功", nil, 0
}

func (s *userInfoService) GetUserInfo(uuid string) (string, *respond.GetUserInfoRespond, int) {
	var userInfo model.UserInfo
	if res := dao.GormDB.Where("uuid = ?", uuid).Find(&userInfo); res.Error != nil {
		zaplog.Error(res.Error.Error())
		return constants.SYSTEM_ERROR, nil, -1
	}
	rsp := respond.GetUserInfoRespond{
		Uuid:      userInfo.Uuid,
		Nickname:  userInfo.Nickname,
		Telephone: userInfo.Telephone,
		Avatar:    userInfo.Avatar,
		Email:     userInfo.Email,
		Gender:    userInfo.Gender,
		Birthday:  userInfo.Birthday,
		Signature: userInfo.Signature,
		CreatedAt: userInfo.CreatedAt.Format("2006-01-02 15:04:05"),
		IsAdmin:   userInfo.IsAdmin,
		Status:    userInfo.Status,
	}
	return "获取用户信息成功", &rsp, 0
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
