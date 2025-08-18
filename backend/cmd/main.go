package main

import (
	"mychat-backend/internal/config"
	"mychat-backend/internal/https_server"
	"mychat-backend/pkg/zaplog"
	"net/http"
	"strconv"
)

func main() {
	zaplog.Info("=== MyChat Backend 服务器启动 ===")

	// 1. 加载配置文件
	zaplog.Info("1. 加载配置文件...")
	if err := config.LoadConfig(); err != nil {
		zaplog.Fatalf("配置文件加载失败: %v", err)
	}
	zaplog.Info("✅ 配置文件加载成功")

	// 2. 初始化日志系统
	zaplog.Info("2. 初始化日志系统...")
	// 日志系统会在导入zaplog包时自动初始化
	zaplog.Info("✅ 日志系统初始化成功")

	// 3. 启动HTTP服务器
	zaplog.Info("3. 启动HTTP服务器...")

	conf := config.GetConfig()
	serverAddr := conf.MainConfig.Host + ":" + strconv.Itoa(conf.MainConfig.Port)

	zaplog.Infof("🚀 服务器启动成功！监听地址: %s", serverAddr)
	zaplog.Infof("📱 注册接口: POST http://%s/register", serverAddr)
	zaplog.Infof("🔐 登录接口: POST http://%s/login", serverAddr)
	zaplog.Info("💡 提示: 使用Postman测试接口功能")
	zaplog.Info("按 Ctrl+C 停止服务器")

	// 启动服务器
	if err := http.ListenAndServe(serverAddr, https_server.GE); err != nil {
		zaplog.Fatalf("服务器启动失败: %v", err)
	}
}
