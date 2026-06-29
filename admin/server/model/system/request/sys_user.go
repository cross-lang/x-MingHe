package request

import (
	common "github.com/flipped-aurora/gin-vue-admin/server/model/common/request"
	"github.com/flipped-aurora/gin-vue-admin/server/model/system"
)

// Register User register structure
type Register struct {
	Username     string `json:"userName" example:"用户名" binding:"required,max=20"`
	Password     string `json:"passWord" example:"密码" binding:"required,max=20"`
	NickName     string `json:"nickName" example:"昵称" binding:"max=20"`
	HeaderImg    string `json:"headerImg" example:"头像链接"`
	AuthorityId  uint   `json:"authorityId" swaggertype:"string" example:"int 角色id"`
	Enable       int    `json:"enable" swaggertype:"string" example:"int 是否启用"`
	AuthorityIds []uint `json:"authorityIds" swaggertype:"string" example:"[]uint 角色id"`
	Phone        string `json:"phone" example:"电话号码" binding:"required,max=20"`
	Email        string `json:"email" example:"电子邮箱" binding:"max=255"`
}

// Login User login structure
type Login struct {
	Username  string `json:"username"`  // 用户名
	Password  string `json:"password"`  // 密码
	Captcha   string `json:"captcha"`   // 验证码
	CaptchaId string `json:"captchaId"` // 验证码ID
}

// ChangePasswordReq Modify password structure
type ChangePasswordReq struct {
	ID          uint   `json:"-"`           // 从 JWT 中提取 user id，避免越权
	Password    string `json:"password"`    // 密码
	NewPassword string `json:"newPassword"` // 新密码
}

type ResetPassword struct {
	ID       uint   `json:"ID" form:"ID"`
	Password string `json:"password" form:"password" gorm:"comment:用户登录密码"` // 用户登录密码
}

// SetUserAuth Modify user's auth structure
type SetUserAuth struct {
	AuthorityId uint `json:"authorityId"` // 角色ID
}

// SetUserAuthorities Modify user's auth structure
type SetUserAuthorities struct {
	ID           uint
	AuthorityIds []uint `json:"authorityIds"` // 角色ID
}

type ChangeUserInfo struct {
	ID           uint                  `gorm:"primarykey"`                                                                           // 主键ID
	NickName     string                `json:"nickName" gorm:"default:系统用户;comment:用户昵称" binding:"max=20"`                           // 用户昵称
	Phone        string                `json:"phone"  gorm:"comment:用户手机号" binding:"max=20"`                                         // 用户手机号
	AuthorityIds []uint                `json:"authorityIds" gorm:"-"`                                                                // 角色ID
	Email        string                `json:"email"  gorm:"comment:用户邮箱" binding:"max=255"`                                         // 用户邮箱
	HeaderImg    string                `json:"headerImg" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"` // 用户头像
	Enable       int                   `json:"enable" gorm:"comment:冻结用户"`                                                           //冻结用户
	Authorities  []system.SysAuthority `json:"-" gorm:"many2many:sys_user_authority;"`
}
type ChangeSelfInfo struct {
	ID           uint                  `gorm:"primarykey"`                                                                           // 主键ID
	NickName     string                `json:"nickName" gorm:"default:系统用户;comment:用户昵称" binding:"min=0,max=20"`                     // 用户昵称
	Phone        string                `json:"phone"  gorm:"comment:用户手机号" binding:"min=0,max=20"`                                   // 用户手机号
	AuthorityIds []uint                `json:"authorityIds" gorm:"-"`                                                                // 角色ID
	Email        string                `json:"email"  gorm:"comment:用户邮箱" binding:"min=0,max=255"`                                   // 用户邮箱
	HeaderImg    string                `json:"headerImg" gorm:"default:https://qmplusimg.henrongyi.top/gva_header.jpg;comment:用户头像"` // 用户头像
	Enable       int                   `json:"enable" gorm:"comment:冻结用户"`                                                           //冻结用户
	Authorities  []system.SysAuthority `json:"-" gorm:"many2many:sys_user_authority;"`
}

type GetUserList struct {
	common.PageInfo
	Username string `json:"username" form:"username"`
	NickName string `json:"nickName" form:"nickName"`
	Phone    string `json:"phone" form:"phone"`
	Email    string `json:"email" form:"email"`
}

type ChangeParkReq struct {
	Pid int `json:"p_id"`
}
