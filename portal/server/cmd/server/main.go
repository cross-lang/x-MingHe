package main

import (
	"flag"
	"fmt"
	"portal/internal/app"
	"portal/internal/config"
)

// @title 明河（MingHe）门户系统后端
// @version 1.0
// @description 明河（MingHe）门户系统后端的 API 文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	diContainer, err := wireDIContainer()
	if err != nil {
		panic(err)
	}

	// 验证配置
	if err := diContainer.Config.Validate(); err != nil {
		panic(fmt.Sprintf("配置验证失败: %v", err))
	}

	application := app.InitApp(diContainer)
	application.Run()
}

// 初始化配置文件地址
func init() {
	configPath := flag.String("c", "config.yaml", "配置文件")
	flag.Parse()
	if *configPath != "" {
		config.SerConfigPath(*configPath)
	}
}
