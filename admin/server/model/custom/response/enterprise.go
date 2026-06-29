package response

// ListEnterpriseItem 获取企业列表item
type ListEnterpriseItem struct {
	ID        uint32   `json:"id"`
	FullName  string   `json:"full_name"` // 企业全称
	ShortName string   `json:"short_name"` // 企业简称
	Category  []string `json:"category"`   // 类目
	BNAME     string   `json:"b_name"`     // 入驻楼栋名称
	Layer     uint32   `json:"layer"`      // 入驻楼层
	RName     string   `json:"r_name"`    // 入驻房间名称
	LeaseTime string   `json:"lease_time"` // 入驻时间
	Websites  []string `json:"websites"`   // 官方网站
	Contact   []string `json:"contact"`    // 联系方式
	Status    int32    `json:"status"`     // 状态：1：启用、-1：禁用
	Weight    int32    `json:"weight"`     // 权重
	Top       int32    `json:"top"`        // 置顶
}

// DetailEnterprise 获取企业信息返回参数
type DetailEnterprise struct {
	ID          uint32   `json:"id"`
	FullName    string   `json:"full_name"`   // 企业全称
	ShortName   string   `json:"short_name"`  // 企业简称
	Icon        string   `json:"icon"`        // logo
	Description string   `json:"description"` // 介绍
	Labels      []string `json:"labels"`      // 标签
	Websites    []string `json:"websites"`    // 官方网站
	Contact     []string `json:"contact"`     // 联系方式
	Weight      int32    `json:"weight"`      // 权重
}
