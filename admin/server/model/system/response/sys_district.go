package response

// DistrictDetail 区域详情响应
type DistrictDetail struct {
	ID       int    `json:"id"`        // 区域ID
	PID      int    `json:"pid"`       // 父级ID
	CityName string `json:"city_name"` // 区域名称
	Level    int    `json:"level"`     // 级别
}

// DistrictListResponse 区域列表响应
type DistrictListResponse struct {
	Items []DistrictDetail `json:"items"` // 区域列表
}

// DistrictTreeResponse 区域树响应
type DistrictTreeResponse struct {
	ID       int                    `json:"id"`                 // 区域ID
	PID      int                    `json:"pid"`                // 父级ID
	CityName string                 `json:"city_name"`          // 区域名称
	Level    int                    `json:"level"`              // 级别
	Children []DistrictTreeResponse `json:"children,omitempty"` // 子节点
}
