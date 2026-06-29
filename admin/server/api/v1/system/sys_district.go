package system

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	systemReq "github.com/flipped-aurora/gin-vue-admin/server/model/system/request"

	"github.com/gin-gonic/gin"
)

type DistrictApi struct {
}

// GetDistrictList 获取区域列表
// @Tags 系统管理
// @Summary 获取区域列表
// @Description 获取区域列表，支持按父级ID、级别和关键词筛选
// @Accept json
// @Produce json
// @Param pid query int false "父级ID"
// @Param level query int false "级别"
// @Param keyword query string false "关键词搜索"
// @Success 200 {object} response.Response{data=response.DistrictListResponse} "获取成功"
// @Router /system/district/list [get]
func (a *DistrictApi) GetDistrictList(c *gin.Context) {
	var req systemReq.DistrictList
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	districtList, err := districtService.GetDistrictList(c, &req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(districtList, c)
}

// GetDistrictTree 获取区域树
// @Tags 系统管理
// @Summary 获取区域树
// @Description 获取区域树结构，从指定父级ID开始
// @Accept json
// @Produce json
// @Param pid query int false "父级ID（默认0表示根节点）"
// @Success 200 {object} response.Response{data=[]response.DistrictTreeResponse} "获取成功"
// @Router /system/district/tree [get]
func (a *DistrictApi) GetDistrictTree(c *gin.Context) {
	var req systemReq.DistrictTree
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 默认从根节点开始
	if req.PID < 0 {
		req.PID = 0
	}

	districtTree, err := districtService.GetDistrictTree(c, &req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	response.OkWithData(districtTree, c)
}
