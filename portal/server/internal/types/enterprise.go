package types

// EnterpriseItem 企业信息
type EnterpriseItem struct {
	ID          uint32   `json:"id"`          // 企业ID
	Type        int32    `json:"type"`        // 企业类型，（1：企业，2：学院）
	FullName    string   `json:"full_name"`   // 企业全称
	ShortName   string   `json:"short_name"`  // 企业简称
	Description string   `json:"description"` // 企业简介
	Icon        string   `json:"icon"`        // 企业icon
	Images      []string `json:"images"`      // 背景图片
	Contact     []string `json:"contact"`     // 类型方式
	Website     []string `json:"website"`     // 企业官网
	Labels      []string `json:"labels"`      // 标签
	Category    []uint32 `json:"category"`    // 类目
	LeaseTime   string   `json:"lease_time"`  // 时间
	Top         int32    `json:"top"`         // 是否置顶
	CreatedAt   string   `json:"created_at"`  // 创建时间
	UpdatedAt   string   `json:"updated_at"`  // 修改时间
}

// EnterpriseList 企业列表结果
type EnterpriseList struct {
	Count int64                 `json:"count"`
	Items []*EnterpriseItem `json:"items"`
}

// PostJoinEnterpriseReq 加入企业请求
type PostJoinEnterpriseReq struct {
	EnterpriseId uint32 `json:"enterprise_id" binding:"required"` // 企业ID
	Code         string `json:"code" binding:"required"`          // 企业入驻码
}

// ListEnterpriseReq 获取企业列表请求参数
type ListEnterpriseReq struct {
	ParkId   uint32 `uri:"park_id" binding:"required"`
	Size     int    `form:"size"`
	Page     int    `form:"page"`
	Keyword  string `form:"keyword"`
	Category int    `form:"category"`
	Type     int    `form:"type"`
}
