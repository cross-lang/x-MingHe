package response

type EnterpriseItem struct {
	EID   uint32 `json:"e_id"`   // 企业id
	EName string `json:"e_name"` // 企业名称
	PID   uint32 `json:"p_id"`   // 园区id
}

type RoleItem struct {
	ID           uint32 `json:"id"`
	Name         string `json:"name"`         // 角色名称
	Key          string `json:"key"`          // 角色key
	Type         int32  `json:"type"`         // 所属类型（1：企业，2：学院，3：访客，4：游客）
	Introduction string `json:"introduction"` // 角色简介
}

type EnterpriseRoleItem struct {
	EID   uint32 `json:"e_id"`   // 企业id
	EName string `json:"e_name"` // 企业名称
	EType int32  `json:"e_type"` // 企业类型，（1：企业，2：学院）
	PID   uint32 `json:"p_id"`   // 园区id

	Roles []*RoleItem `json:"roles"` // 角色信息
}

type UserItem struct {
	// 基本信息
	ID              uint32  `json:"id"`
	Name            string  `json:"name"`            // 用户名
	Avatar          string  `json:"avatar"`          // 头像
	Account         string  `json:"account"`         // 账号
	PhoneNumber     string  `json:"phone_number"`    // 手机号
	Gender          int32   `json:"gender"`          // 性别（0：保密，1：男，2：女）
	VerificationStatus int32 `json:"verification_status"` // 实名认证状态（1：已认证，0：未认证）
	BlockStatus     int32   `json:"block_status"`     // 拉黑状态（1：未拉黑；2：已拉黑）
	AccountStatus   int32   `json:"account_status"`   // 账号状态（1：正常；2：禁用）
	DeactivateStatus int32 `json:"deactivate_status"` // 用户注销状态（0：未注销，1：注销中，2：已注销）
	DeactivateAt    string  `json:"deactivate_at"`   // 用户注销时间
	CreatedAt       string  `json:"created_at"`      // 创建时间
	UpdatedAt       string  `json:"updated_at"`      // 修改时间
	DeletedAt       string  `json:"deleted_at"`      // 删除时间

	// 实名认证信息
	VChannel string `json:"v_channel"` // 认证渠道
	VName    string `json:"v_name"`    // 姓名
	VIDCard  string `json:"v_id_card"` // 身份证号码

	// 用户基本信息
	Profession       string   `json:"profession"`       // 职业
	GraduatedSchool  string   `json:"graduated_school"`  // 毕业学院
	WorkYears        float64  `json:"work_years"`        // 工作年限
	Residence        string   `json:"residence"`         // 居住地
	Skills           []string `json:"skills"`            // 核心技能
	Introduction     string   `json:"introduction"`      // 个人简介
	Contact          string   `json:"contact"`           // 联系方式
}

// UserExperience 用户经历
type UserExperience struct {
	Type           int32    `json:"type"`            // 经历过类型（1：教育经历，2：工作经历，3：项目经历）
	Subject        string   `json:"subject"`         // 经历主体
	Identity       string   `json:"identity"`        // 经历中的身份：用于教育经历中的学位、工作经历中的职位、项目经历中的身份
	StartDate      string   `json:"start_date"`      // 经历开始日期
	EndDate        string   `json:"end_date"`        // 经历结束日期
	EducationMode  string   `json:"education_mode"`  // 教育模式:全日制、非全日制、自学考试、成人教育/函授、网络/远程教育
	EducationMajor string   `json:"education_major"` // 教育专业
	Attachment     []string `json:"attachment"`      // 学历证明附件
	WorkMode       string   `json:"work_mode"`       // 工作模式：全职、兼职、实习、外包/劳务派遣
}

type UserList struct {
	Count    int64       `json:"count"`
	Page     int         `json:"page"`     // 页码
	PageSize int         `json:"pageSize"` // 每页大小
	Items    []*UserItem `json:"items"`
}
