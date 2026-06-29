package app

import (
	_ "portal/docs" // ⭐ 非常重要

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func (a *App) RegisterSwagger() {
	a.Engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
