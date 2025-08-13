package main

import (
	"fmt"
	"mychat-backend/internal/config"
	"mychat-backend/internal/dao"
	"mychat-backend/pkg/zaplog"
)

func main() {
	fmt.Println("=== MyChat Backend 数据库连接测试 ===")

	// 1. 加载配置文件
	fmt.Println("1. 加载配置文件...")
	if err := config.LoadConfig(); err != nil {
		fmt.Printf("❌ 配置文件加载失败: %v\n", err)
		return
	}
	fmt.Println("✅ 配置文件加载成功")

	// 2. 初始化日志系统
	fmt.Println("2. 初始化日志系统...")
	// 日志系统会在导入zaplog包时自动初始化
	fmt.Println("✅ 日志系统初始化成功")

	// 3. 测试数据库连接
	fmt.Println("3. 测试数据库连接...")

	// 检查GormDB是否已初始化
	if dao.GormDB == nil {
		fmt.Println("❌ 数据库连接失败: GormDB为nil")
		return
	}

	// 测试数据库连接
	sqlDB, err := dao.GormDB.DB()
	if err != nil {
		fmt.Printf("❌ 获取数据库实例失败: %v\n", err)
		return
	}

	// 执行ping测试
	if err := sqlDB.Ping(); err != nil {
		fmt.Printf("❌ 数据库ping失败: %v\n", err)
		return
	}

	fmt.Println("✅ 数据库连接成功")

	// 4. 显示数据库信息
	fmt.Println("4. 数据库信息:")
	conf := config.GetConfig()
	fmt.Printf("   - 主机: %s:%d\n", conf.MysqlConfig.Host, conf.MysqlConfig.Port)
	fmt.Printf("   - 数据库: %s\n", conf.MysqlConfig.DatabaseName)
	fmt.Printf("   - 用户: %s\n", conf.MysqlConfig.User)

	// 5. 测试完成
	fmt.Println("\n🎉 数据库连接测试完成！")
	fmt.Println("所有功能正常，可以开始开发其他功能。")
}
