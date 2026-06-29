package types

import (
	"github.com/golang-jwt/jwt/v5"
)

type UserLoginJwtClaims struct {
	UserId uint32 `json:"user_id"`
	jwt.RegisteredClaims
}

// UserRegisterReq 用户注册请求参数
type UserRegisterReq struct {
	PhoneNumber string `json:"phone_number" binding:"required"` // 账号
	Code        string `json:"code" binding:"required"`         // 验证码
}

// UserLoginReq 用户登录请求参数
type UserLoginReq struct {
	PhoneNumber string `json:"phone_number" binding:"required"` // 账号
	Code        string `json:"code" binding:"required"`         // 验证码
}

// UserLoginResp 用户登录返回参数
type UserLoginResp struct {
	Token     string `json:"token"`      // token
	ExpiresAt int64  `json:"expires_at"` // 过期时间
}

// UserItem 用户详情
type UserItem struct {
	ID                  uint32            `json:"id"`                   // 用户ID
	Name                string            `json:"name"`                 // 用户名称
	Avatar              string            `json:"avatar"`               // 用户头像
	Account             string            `json:"account"`              // 账号
	PhoneNumber         string            `json:"phone_number"`         // 手机号
	Gender              int32             `json:"gender"`               // 性别，0：保密；1：男；2：女
	VerificationStatus  int32             `json:"verification_status"`  // 实名认证状态
	AccountStatus       int32             `json:"account_status"`       // 禁用状态（1：正常，2：禁用）
	BlockStatus         int32             `json:"block_status"`         // 拉黑状态（1：未拉黑，2：已拉黑）
	Roles               []*UserRole       `json:"roles"`                // 用户角色
	PID                 uint32            `json:"p_id"`                 // 用户当前切换的园区
	EID                 uint32            `json:"e_id"`                 // 用户当前切换的企业
	Profession          string            `json:"profession"`           // 职业
	GraduatedSchool     string            `json:"graduated_school"`     // 毕业学院
	WorkYears           int32             `json:"work_years"`           // 工作年限
	Residence           string            `json:"residence"`            // 居住地
	Skills              []string          `json:"skills"`               // 核心技能
	Introduction        string            `json:"introduction"`         // 个人简介
	Contact             string            `json:"contact"`              // 联系方式
	EducationExperience []*UserExperience `json:"education_experience"` // 教育经历
	WorkExperience      []*UserExperience `json:"work_experience"`      // 工作经历
	ProjectExperience   []*UserExperience `json:"project_experience"`   // 项目经历
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

type UserRole struct {
	Key  string `json:"key"`  // 角色key，(enterprise_staff：企业用户角色，enterprise_office：企业行政角色，enterprise_personnel：企业人力角色，enterprise_accounting：企业财务角色，enterprise_supervisor：企业管理员角色，campus_student：学院学生角色，campus_teacher：学院教师角色，campus_admin：学院管理员角色，tourist：访客，visitor：游客)
	Name string `json:"name"` // 角色名称
}

// UserVerificationReq 用户实名认证参数
type UserVerificationReq struct {
	Name   string `json:"name" binding:"required"`
	IdCard string `json:"id_card" binding:"required"`
}

// DetailDeactivateUserCheckResp 查询用户注销准入评估结果
type DetailDeactivateUserCheckResp struct {
	CanDeactivate bool     `json:"can_deactivate"` // 是否允许注销
	Items         []string `json:"items"`          // 未通过的指标列表，每个字符串表示一个未通过的指标key
}
