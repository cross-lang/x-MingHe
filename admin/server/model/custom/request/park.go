package request

import (
	"github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
)

// CreateParkParams 创建园区请求参数
type CreateParkParams struct {
	FullName     string `json:"full_name" binding:"required,max=20"` // 园区全称
	ShortName    string `json:"short_name" binding:"max=10"`         // 园区简称
	Image        string `json:"image" binding:"required"`            // 园区图片
	Introduction string `json:"introduction" binding:"max=10000"`    // 园区简介
}

// GetParkInfoParams 获取园区信息
type GetParkInfoParams struct {
	ID uint32 `json:"id" binding:"required"` // 园区id
}

// GetParkListParams 获取园区列表请求参数
type GetParkListParams struct {
	request.PageInfo
	Status   []int32 `json:"status"`   // 园区状态（-1：禁用，1：启用）
	OrderKey string  `json:"orderKey"` // 排序
	Desc     bool    `json:"desc"`     // 排序方式:升序false(默认)|降序true
}

// UpdateParkParams 更新园区请求参数
type UpdateParkParams struct {
	ID           uint32 `json:"id" binding:"required"`
	FullName     string `json:"full_name" binding:"required,max=20"` // 园区全称
	ShortName    string `json:"short_name" binding:"max=10"`         // 园区简称
	Image        string `json:"image" binding:"required"`            // 园区图片
	Introduction string `json:"introduction" binding:"max=10000"`    // 园区简介
}

// UpdateParkStatusParams 更新园区状态请求参数
type UpdateParkStatusParams struct {
	ID     uint32 `json:"id" binding:"required"`                // 园区id
	Status int32  `json:"status" binding:"required,oneof=1 -1"` // 状态：1：启用、-1：禁用
}

// CreateParkBuildingLayoutParams 楼栋户型参数
type CreateParkBuildingLayoutParams struct {
	ID   uint32  `json:"id"`                                // 户型id
	Name string  `json:"name" binding:"required,max=10"`    // 户型名称
	Area float64 `json:"area" binding:"required,max=99999"` // 户型面积
}

// CreateParkBuildingParams 创建园区楼栋请求参数
type CreateParkBuildingParams struct {
	PID        uint32                            `json:"p_id" binding:"required"`                      // 园区id
	FullName   string                            `json:"full_name" binding:"required,max=20"`          // 楼栋全称
	ShortName  string                            `json:"short_name" binding:"max=10"`                  // 楼栋简称
	Images     []string                          `json:"images" binding:"required,max=20"`             // 楼栋图片
	LayerTotal int32                             `json:"layer_total" binding:"required,min=1,max=999"` // 楼栋总层数
	Layouts    []*CreateParkBuildingLayoutParams `json:"layouts" binding:"required,max=100,dive"`      // 户型列表
}

// GetParkBuildingListParams 获取园区楼栋列表请求参数
type GetParkBuildingListParams struct {
	request.PageInfo
	PID      uint32 `json:"p_id" binding:"required"` // 所属园区id
	OrderKey string `json:"orderKey"`                // 排序
	Desc     bool   `json:"desc"`                    // 排序方式:升序false(默认)|降序true
}

// GetParkBuildingInfoParams 获取园区楼栋信息请求参数
type GetParkBuildingInfoParams struct {
	ID uint32 `json:"id" binding:"required"`
}

// UpdateParkBuildingParams 更新园区楼栋信息
type UpdateParkBuildingParams struct {
	ID         uint32                            `json:"id" binding:"required"`                        // 园区id
	FullName   string                            `json:"full_name" binding:"required,max=20"`          // 楼栋全称
	ShortName  string                            `json:"short_name" binding:"max=10"`                  // 楼栋简称
	Images     []string                          `json:"images" binding:"required,max=20"`             // 楼栋图片
	LayerTotal int32                             `json:"layer_total" binding:"required,min=1,max=999"` // 楼栋总层数
	Layouts    []*CreateParkBuildingLayoutParams `json:"layouts" binding:"required,max=100,dive"`      // 户型列表
}

// DeleteParkBuildingParams 删除楼栋信息
type DeleteParkBuildingParams struct {
	ID uint32 `json:"id" binding:"required"`
}

// UpdateParkBuildingStatusParams 修改园区楼栋状态
type UpdateParkBuildingStatusParams struct {
	ID     uint32 `json:"id" binding:"required"`                // 楼栋id
	Status int32  `json:"status" binding:"required,oneof=1 -1"` // 状态：1：启用、-1：禁用
}

// CreateParkBuildingRoomParams 创建园区楼栋房间请求参数
type CreateParkBuildingRoomParams struct {
	PID         uint32   `json:"p_id" binding:"required"`                   // 所属园区
	BID         uint32   `json:"b_id" binding:"required"`                   // 所属楼栋
	LID         uint32   `json:"l_id" binding:"required"`                   // 关联户型
	FullName    string   `json:"full_name" binding:"required,max=20"`       // 房间全称
	ShortName   string   `json:"short_name" binding:"max=10"`               // 房间简称
	Layer       int32    `json:"layer" binding:"required,min=1,max=999"`    // 所属楼层
	AllowSettle int32    `json:"allow_settle" binding:"oneof=0 1"`          // 是否允许入驻
	Images      []string `json:"images" binding:"required,max=20"`          // 房间图片
	PublicSpace int32    `json:"public_space" binding:"oneof=0 1 2 3 4"`    // 是否是公共空间（0：非公共空间，1：会议室，2：VIP接待室，3：影棚，4：茶室）
	Rent        float64  `json:"rent" binding:"max=99999"`                  // 租金，非公共空间时有值
	Supporting  []string `json:"supporting" binding:"dive,required,max=20"` // 配套设置，非公共空间时有值
}

// GetParkBuildingRoomInfoParams 获取房间信息请求参数
type GetParkBuildingRoomInfoParams struct {
	ID uint32 `json:"id" binding:"required"` // 房间id
}

// GetParkBuildingRoomListParams 获取园区楼栋列表请求参数
type GetParkBuildingRoomListParams struct {
	request.PageInfo
	PID         uint32   `json:"p_id" binding:"required"` // 所属园区id
	BID         []uint32 `json:"b_id"`                    // 所属楼栋id
	AllowSettle []int32  `json:"allow_settle"`            // 是否运行入驻，1：允许、0：不允许
	PublicSpace []int32  `json:"public_space"`            // 公共空间类型
	OrderKey    string   `json:"orderKey"`                // 排序
	Desc        bool     `json:"desc"`                    // 排序方式:升序false(默认)|降序true
}

// UpdateParkBuildingRoomParams 更新楼栋房间信息请求参数
type UpdateParkBuildingRoomParams struct {
	ID          uint32   `json:"id" binding:"required"`                     // 房间id
	BID         uint32   `json:"b_id" binding:"required"`                   // 所属楼栋
	LID         uint32   `json:"l_id" binding:"required"`                   // 关联户型
	FullName    string   `json:"full_name" binding:"required,max=20"`       // 房间全称
	ShortName   string   `json:"short_name" binding:"max=10"`               // 房间简称
	Layer       int32    `json:"layer" binding:"required,min=1,max=999"`    // 所属楼层
	AllowSettle int32    `json:"allow_settle" binding:"oneof=0 1"`          // 是否允许入驻
	Images      []string `json:"images" binding:"required,max=20"`          // 房间图片
	PublicSpace int32    `json:"public_space" binding:"oneof=0 1 2 3 4"`    // 是否是公共空间（0：非公共空间，1：会议室，2：VIP接待室，3：影棚，4：茶室）
	Rent        float64  `json:"rent" binding:"max=99999"`                  // 租金，非公共空间时有值
	Supporting  []string `json:"supporting" binding:"dive,required,max=20"` // 配套设置，非公共空间时有值
}

// DeleteParkBuildingRoomParams 删除楼栋房间请求参数
type DeleteParkBuildingRoomParams struct {
	ID uint32 `json:"id" binding:"required"`
}

// UpdateParkBuildingRoomStatusParams 修改园区楼栋状态
type UpdateParkBuildingRoomStatusParams struct {
	ID     uint32 `json:"id" binding:"required"`                // 楼栋id
	Status int32  `json:"status" binding:"required,oneof=1 -1"` // 状态：1：启用、-1：禁用
}
