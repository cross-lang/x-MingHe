package custom

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/custom/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type EnterpriseApi struct{}

// CreateEnterprise 创建企业
// @Tags 业务定制
// @Summary 创建企业
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.CreateEnterpriseParams true "创建企业"
// @Success 200 {object} response.Response{msg=string} "创建成功"
// @Router /enterprises/create [post]
func (e *EnterpriseApi) CreateEnterprise(c *gin.Context) {
	var params request.CreateEnterpriseParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 获取当前用户id
	userId := utils.GetUserID(c)
	// 创建楼栋
	roomId, err := enterpriseService.CreateEnterprise(c.Request.Context(), userId, params)
	if err != nil {
		global.GVA_LOG.Error("创建企业失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(roomId, c)
}

// ListEnterprise 获取园区企业列表
// @Tags 业务定制
// @Summary 获取园区企业列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param page query integer false "页码" default(1)
// @Param pageSize query integer false "每页数据量" default(10)
// @Success 200 {object} response.Response{data=response.PageResult,msg=string} "分页获取API列表,返回包括列表,总数,页码,每页数量"
// @Router /enterprises/list [get]
func (e *EnterpriseApi) ListEnterprise(c *gin.Context) {
	var params request.ListEnterpriseParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	err := utils.Verify(params.PageInfo, utils.PageInfoVerify)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	// 获取企业列表
	list, total, err := enterpriseService.ListEnterprise(c.Request.Context(), params)
	if err != nil {
		global.GVA_LOG.Error("获取企业列表失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithDetailed(response.PageResult{
		List:     list,
		Total:    total,
		Page:     params.Page,
		PageSize: params.PageSize,
	}, "获取成功", c)
}

// UpdateEnterpriseStatus 更新企业状态
// @Tags 业务定制
// @Summary 更新企业状态
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.UpdateEnterpriseStatusParams true "更新企业状态信息"
// @Success 200 {object} response.Response{msg=string} "更新成功"
// @Router /enterprises/updateStatus [post]
func (e *EnterpriseApi) UpdateEnterpriseStatus(c *gin.Context) {
	var params request.UpdateEnterpriseStatusParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 获取当前用户id
	userId := utils.GetUserID(c)

	// 更新园区信息
	err := enterpriseService.UpdateEnterpriseStatus(c.Request.Context(), userId, params)
	if err != nil {
		global.GVA_LOG.Error("更新状态失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(params.ID, c)
}

// TopEnterprise 企业置顶
// @Tags 业务定制
// @Summary 企业置顶
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.TopEnterpriseParams true "企业置顶"
// @Success 200 {object} response.Response{msg=string} "置顶成功"
// @Router /enterprises/top [post]
func (e *EnterpriseApi) TopEnterprise(c *gin.Context) {
	var params request.TopEnterpriseParams
	if err := c.ShouldBindJSON(&params); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 获取当前用户id
	userId := utils.GetUserID(c)

	err := enterpriseService.TopEnterprise(c.Request.Context(), userId, params)
	if err != nil {
		global.GVA_LOG.Error("置顶失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("success", c)
}
