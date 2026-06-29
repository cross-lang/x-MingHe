package initialize

import (
	"github.com/flipped-aurora/gin-vue-admin/server/router"
	"github.com/gin-gonic/gin"
)

// 占位方法，保证文件可以正确加载，避免go空变量检测报错，请勿删除。
func holder(routers ...*gin.RouterGroup) {
	_ = routers
	_ = router.RouterGroupApp
}

func initBizRouter(routers ...*gin.RouterGroup) {
	privateGroup := routers[0]
	publicGroup := routers[1]

	holder(publicGroup, privateGroup)
	{
		// 定制API路由
		customGroup := privateGroup.Group("custom")
		// 注册用户相关路由
		router.RouterGroupApp.Custom.InitUserRouter(customGroup)
		// 注册企业修改路由
		router.RouterGroupApp.Custom.InitEnterpriseRouter(customGroup)
	}
}
