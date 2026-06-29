package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/middleware"
	"github.com/gin-gonic/gin"
)

type DistrictRouter struct{}

// InitDistrictRouter 初始化省市区信息
func (s *DistrictRouter) InitDistrictRouter(Router *gin.RouterGroup) {
	districtRouter := Router.Group("district").Use(middleware.OperationRecord())
	{
		districtRouter.GET("list", districtApi.GetDistrictList)
	}
}
