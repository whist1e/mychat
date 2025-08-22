// Package dao 数据访问对象层
// 该文件负责数据库连接管理和初始化
// 使用GORM作为ORM框架，提供全局数据库连接实例
package dao

import (
	"fmt"
	"mychat-backend/internal/config"
	"mychat-backend/internal/model"
	"mychat-backend/pkg/zaplog"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// GormDB 全局GORM数据库连接实例
// 该变量在整个应用中被其他包使用，用于执行数据库操作
var GormDB *gorm.DB

// init 包初始化函数
// 当dao包被导入时自动执行，负责：
// 1. 建立数据库连接
// 2. 自动迁移表结构
// 3. 初始化全局数据库实例
func init() {
	// 获取应用配置信息
	conf := config.GetConfig()

	// 从配置中提取MySQL连接参数
	user := conf.MysqlConfig.User                 // 数据库用户名
	password := conf.MysqlConfig.Password         // 数据库密码
	host := conf.MysqlConfig.Host                 // 数据库主机地址
	port := conf.MysqlConfig.Port                 // 数据库端口
	databaseName := conf.MysqlConfig.DatabaseName // 数据库名称

	// 构建MySQL数据库连接字符串（DSN）
	// 其中：
	//   - %s:%s 表示用户名:密码
	//   - @tcp(%s:%d) 表示连接到主机:端口
	//   - /%s 表示数据库名
	//   - charset=utf8mb4 支持完整UTF-8字符集（含emoji）
	//   - parseTime=True 自动解析MySQL的时间类型为Go的time.Time
	//   - loc=Local 使用本地时区
	// dsn最终格式示例: "root:password@tcp(127.0.0.1:3306)/mychat?charset=utf8mb4&parseTime=True&loc=Local"
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, databaseName)

	var err error
	// 使用GORM连接MySQL数据库
	// mysql.Open(dsn): 返回带有dsn配置的驱动
	// &gorm.Config{}: 使用默认GORM配置
	GormDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		// 连接失败时记录错误日志并终止程序
		zaplog.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移数据库表结构
	// AutoMigrate会自动创建不存在的表，并更新表结构以匹配模型定义
	// 这是GORM的一个强大功能，可以自动管理数据库schema
	err = GormDB.AutoMigrate(
		&model.UserInfo{},     // 用户信息表
		&model.GroupInfo{},    // 群组信息表
		&model.UserContact{},  // 用户联系人表
		&model.Session{},      // 会话表
		&model.ContactApply{}, // 联系人申请表
		&model.Message{},      // 消息表
	)
	if err != nil {
		// 迁移失败时记录错误日志并终止程序
		zaplog.Fatalf("Failed to auto migrate database: %v", err)
	}

	// 记录成功日志
	zaplog.Info("Database connected and migrated successfully")
}
