package request

// DistrictList 区域列表请求
type DistrictList struct {
	PID     int    `json:"pid" form:"pid"`         // 父级ID
	Level   int    `json:"level" form:"level"`     // 级别
	Keyword string `json:"keyword" form:"keyword"` // 关键词搜索
}

// DistrictTree 区域树请求
type DistrictTree struct {
	PID int `json:"pid" form:"pid"` // 父级ID（默认0表示根节点）
}
