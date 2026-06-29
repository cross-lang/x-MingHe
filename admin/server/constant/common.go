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


// 可见性
const (
	VisibleStatusYes = 1 // 可见
	VisibleStatusNo  = 0 // 不可见
)
