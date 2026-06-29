package custom

import (
	"github.com/flipped-aurora/gin-vue-admin/server/global"
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/response"
	"github.com/flipped-aurora/gin-vue-admin/server/model/custom/request"
	"github.com/flipped-aurora/gin-vue-admin/server/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type UserApi struct{}

// UpdateUser 编辑用户
// @Tags 业务定制
// @Summary 编辑用户
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param data body request.UpdateUserReq true "编辑用户"
// @Success 200 {object} response.Response{data=string, msg=string} "编辑成功"
// @Router /custom/users/update [post]
func (p *UserApi) UpdateUser(c *gin.Context) {
	var req request.UpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	global.GVA_LOG.Info(">>>>>>>> UpdateUser", zap.Any("req", req))
	err := userService.UpdateUser(c.Request.Context(), &req)
	if err != nil {
		global.GVA_LOG.Error("编辑用户失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("", c)
}

// DisableUser 禁用用户
// @Tags 业务定制
// @Summary 禁用用户
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param user_id path integer true "用户ID"
// @Success 200 {object} response.Response{data=string, msg=string} "禁用成功"
// @Router /custom/users/{user_id}/disable [post]
func (p *UserApi) DisableUser(c *gin.Context) {
	var req request.DisableUserReq
	if err := c.ShouldBindUri(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	global.GVA_LOG.Info(">>>>>>>> DisableUserReq", zap.Any("req", req))

	err := userService.DisableUser(c.Request.Context(), &req)
	if err != nil {
		global.GVA_LOG.Error("禁用用户失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("", c)
}

// EnableUser 启用用户
// @Tags 业务定制
// @Summary 启用用户
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param user_id path integer true "用户ID"
// @Success 200 {object} response.Response{data=string, msg=string} "启用成功"
// @Router /custom/users/{user_id}/enable [post]
func (p *UserApi) EnableUser(c *gin.Context) {
	var req request.EnableUserReq
	if err := c.ShouldBindUri(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	global.GVA_LOG.Info(">>>>>>>> EnableUser", zap.Any("req", req))
	err := userService.EnableUser(c.Request.Context(), &req)
	if err != nil {
		global.GVA_LOG.Error("启用用户失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData("", c)
}

// DetailUser 查询用户详情
// @Tags 业务定制
// @Summary 查询用户详情
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param user_id path integer true "用户ID"
// @Success 200 {object} response.Response{data=response.UserItem, msg=string} "查询成功"
// @Router /custom/users/{user_id}/detail [get]
func (p *UserApi) DetailUser(c *gin.Context) {
	var req request.DetailUserReq
	if err := c.ShouldBindUri(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userId := utils.GetUserID(c)
	global.GVA_LOG.Info(">>>>>>>> DetailUser", zap.Any("req", req))
	UserId, err := userService.DetailUser(c.Request.Context(), userId, &req)
	if err != nil {
		global.GVA_LOG.Error("查询某用户失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(UserId, c)
}

// ListUser 查询用户列表
// @Tags 业务定制
// @Summary 查询用户列表
// @Security ApiKeyAuth
// @Accept application/json
// @Produce application/json
// @Param keyword query string false "关键字（匹配用户名称等）"
// @Param account_statuses query []integer false "账号状态列表（1：正常；2：禁用）" collectionFormat(multi)
// @Param block_statuses query []integer false "拉黑状态列表（1：未拉黑；2：已拉黑）" collectionFormat(multi)
// @Param verify_statuses query []integer false "实名状态列表（0：未认证；2：已认证）" collectionFormat(multi)
// @Param start_time_from query string false "注册日期起"
// @Param start_time_to query string false "注册日期止"
// @Param page query integer false "页码" default(1)
// @Param pageSize query integer false "每页数据量" default(10)
// @Param orderKey query string false "排序字段"
// @Param desc query boolean false "是否降序" default(false)
// @Success 200 {object} response.Response{data=response.UserList,msg=string} "分页获取用户列表"
// @Router /custom/users/list [get]
func (p *UserApi) ListUser(c *gin.Context) {
	var req request.ListUserReq
	if err := c.ShouldBindQuery(&req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	userId := utils.GetUserID(c)
	global.GVA_LOG.Info(">>>>>>>> ListUser", zap.Any("req", req))
	UserList, err := userService.ListUser(c.Request.Context(), userId, &req)
	if err != nil {
		global.GVA_LOG.Error("查询用户列表失败", zap.Error(err))
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(UserList, c)
}
