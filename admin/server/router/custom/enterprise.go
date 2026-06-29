package custom

import "github.com/gin-gonic/gin"

type EnterpriseRouter struct{}

// InitEnterpriseRouter 注册企业路由
func (e *EnterpriseRouter) InitEnterpriseRouter(router *gin.RouterGroup) {
	enterpriseGroup := router.Group("enterprises")
	{
		// 创建企业
		enterpriseGroup.POST("create", enterpriseApi.CreateEnterprise)
		// 置顶企业
		enterpriseGroup.POST("top", enterpriseApi.TopEnterprise)
		// 获取企业列表
		enterpriseGroup.GET("list", enterpriseApi.ListEnterprise)
		// 更新企业状态
		enterpriseGroup.POST("updateStatus", enterpriseApi.UpdateEnterpriseStatus)
	}
}
