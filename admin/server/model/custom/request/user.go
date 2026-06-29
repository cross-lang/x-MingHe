package request

import "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"

type UpdateUserReq struct {
	// 基本信息
	ID                  uint32            `json:"id"`
	Avatar              string            `json:"avatar"`               // 头像
	PhoneNumber         string            `json:"phone_number"`         // 手机号
	Residence           string            `json:"residence"`            // 居住地
	Profession          string            `json:"profession"`           // 职业
	GraduatedSchool     string            `json:"graduated_school"`     // 毕业学院
	WorkYears           int32             `json:"work_years"`           // 工作年限
	Contact             string            `json:"contact"`              // 联系方式
	Introduction        string            `json:"introduction"`         // 个人简介
	Skills              []string          `json:"skills"`               // 核心技能
	EducationExperience []*UserExperience `json:"education_experience"` // 教育经历
	WorkExperience      []*UserExperience `json:"work_experience"`      // 工作经历
	ProjectExperience   []*UserExperience `json:"project_experience"`   // 项目经历
	BlockStatus         int32             `json:"block_status"`         // 拉黑状态（1：未拉黑；2：已拉黑）
	AccountStatus       int32             `json:"account_status"`       // 账号状态（1：正常；2：禁用）

	// 企业和企业角色信息
	EnterpriseRoles []*struct {
		EID   uint32 `json:"e_id"` // 企业id
		PID   uint32 `json:"p_id"` // 园区id
		Roles []*struct {
			ID uint32 `json:"id"`
		} `json:"roles"` // 角色信息
	} `json:"enterprise_roles"`
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

type DetailUserReq struct {
	UserID uint32 `uri:"user_id"`
}

type DisableUserReq struct {
	UserID uint32 `uri:"user_id"`
}

type EnableUserReq struct {
	UserID uint32 `uri:"user_id"`
}

type ListUserReq struct {
	request.PageInfo
	AccountStatuses  []int32 `form:"account_statuses"`  // 账号状态列表（1：正常；2：禁用）
	BlockStatuses    []int32 `form:"block_statuses"`    // 拉黑状态列表（1：未拉黑；2：已拉黑）
	VerifyStatuses   []int32 `form:"verify_statuses"`   // 实名状态列表（0：未认证；2：已认证）
	DeactivateStatus []int32 `form:"deactivate_status"` // 用户注销状态（0：未注销，1：注销中，2：已注销）
	RegisterTimeFrom string  `form:"start_time_from"`   // 注册时间起
	RegisterTimeTo   string  `form:"start_time_to"`     // 注册时间止
	OrderKey         string  `form:"orderKey"`          // 排序
	Desc             bool    `form:"desc"`              // 排序方式:升序false(默认)|降序true
}
