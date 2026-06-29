package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

// CreateEnterpriseParams 创建企业请求参数
type CreateEnterpriseParams struct {
	PID         uint32 `json:"p_id" binding:"required"`                      // 园区id
	FullName    string `json:"full_name" binding:"required,max=40"`          // 企业全称
	ShortName   string `json:"short_name" binding:"max=20"`                  // 企业简称
	Category    []uint32 `json:"category" binding:"required"`                  // 类目
	Icon        string `json:"icon" binding:"-"`                             // logo
	Images      []string `json:"images" binding:"required,max=20"`             // 背景图
	Description string `json:"description" binding:"required,max=1000"`      // 介绍
	Labels      []string `json:"labels" binding:"max=10,dive,required,max=15"` // 标签
	Websites    []string `json:"websites" binding:"dive,required,max=40"`      // 官方网站
	Contact     []string `json:"contact" binding:"dive,required,max=40"`       // 联系方式
	Lease       []*CreateEnterpriseLease `json:"lease" binding:"required,dive"`                // 企业入驻信息
}

// CreateEnterpriseLease 创建企业入驻信息
type CreateEnterpriseLease struct {
	PID       uint32 `json:"p_id" binding:"required"`                           // 入驻的园区
	BID       uint32 `json:"b_id" binding:"required"`                           // 入驻的楼栋
	Layer     uint32 `json:"layer" binding:"required,min=1"`                    // 入驻的楼层
	RID       uint32 `json:"r_id" binding:"required"`                           // 入驻的房间
	LeaseTime string `json:"lease_time" binding:"required,datetime=2006-01-02"` // 入驻时间
}

// ListEnterpriseParams 企业列表查询参数
type ListEnterpriseParams struct {
	request.PageInfo
	OrderKey       string   `json:"orderKey"`         // 排序
	Desc           bool     `json:"desc"`             // 排序方式:升序false(默认)|降序true
	PID            uint32   `json:"p_id" binding:"-"` // 所属园区id
	Category       []int32  `json:"category"`         // 类目
	BID            []uint32 `json:"b_id"`             // 所属楼栋
	LeaseStartTime string   `json:"lease_start_time"` // 入驻时间
	LeaseEndTime   string   `json:"lease_end_time"`   // 入驻时间
}

// UpdateEnterpriseStatusParams 更新企业状态请求参数
type UpdateEnterpriseStatusParams struct {
	ID     uint32 `json:"id" binding:"required"`
	Status int32  `json:"status" binding:"required,oneof=1 -1"` // 状态：1：启用、-1：禁用
}

// TopEnterpriseParams 企业置顶请求参数
type TopEnterpriseParams struct {
	ID     uint32 `json:"id" binding:"required"`
	Status bool   `json:"status"` // 是否置顶
}
