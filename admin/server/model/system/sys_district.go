package system

// SysDistrict mapped from table <sys_district>
type SysDistrict struct {
	ID       int    `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	PID      int    `gorm:"column:pid;not null" json:"pid"`             // 父级ID
	CityName string `gorm:"column:city_name;not null" json:"city_name"` // 区域名称
	Level    int    `gorm:"column:level;not null" json:"level"`         // 级别
}

// TableName SysDistrict's table name
func (*SysDistrict) TableName() string {
	return "sys_district"
}
