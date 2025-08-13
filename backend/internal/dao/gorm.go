package dao

import (
	"fmt"
	"mychat-backend/internal/config"
	"mychat-backend/internal/model"
	"mychat-backend/pkg/zaplog"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var GormDB *gorm.DB

func init() {
	conf := config.GetConfig()
	user := conf.MysqlConfig.User
	password := conf.MysqlConfig.Password
	host := conf.MysqlConfig.Host
	port := conf.MysqlConfig.Port
	databaseName := conf.MysqlConfig.DatabaseName

	// 构建数据库连接字符串
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		user, password, host, port, databaseName)

	var err error
	GormDB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		zaplog.Fatalf("Failed to connect to database: %v", err)
	}

	// 自动迁移数据库表结构
	err = GormDB.AutoMigrate(&model.UserInfo{}, &model.GroupInfo{}, &model.UserContact{}, &model.Session{}, &model.ContactApply{}, &model.Message{})
	if err != nil {
		zaplog.Fatalf("Failed to auto migrate database: %v", err)
	}

	zaplog.Info("Database connected and migrated successfully")
}
