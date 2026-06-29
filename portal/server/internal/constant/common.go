package constant

// 用户认证状态
const (
	UserVerifyStatusNotVerify = 0 // 未认证
	UserVerifyStatusVerified  = 1 // 已认证
)

// 用户账号状态
const (
	UserAccountStatusNormal  = 1 // 正常
	UserAccountStatusDisable = 2 // 禁用
)

// 用户拉黑状态
const (
	UserBlockStatusNo  = 1 // 未拉黑
	UserBlockStatusYes = 2 // 已拉黑
)

// 用户注销状态
const (
	DeactivateUserStatusNot      int32 = 0 // 未注销
	DeactivateUserStatusPending  int32 = 1 // 注销中
	DeactivateUserStatusFinished int32 = 2 // 已注销
)

// 用户角色
const (
	TouristRole = "tourist" // 访客
	VisitorRole = "visitor" // 游客
)
