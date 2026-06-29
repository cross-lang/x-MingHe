package errorx

import (
	"fmt"
)

type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("(%d)%s", e.Code, e.Message)
}

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// UnknownErrorCode 未知错误
const UnknownErrorCode = 1000100

var (
	ParameterError = NewError(1000101, "参数错误")

	// 通用错误
	RecordNotFoundError = NewError(1000102, "记录不存在")
	DBError             = NewError(1000103, "数据库操作失败")
	InvalidActionError  = NewError(1000104, "非法的动作")

	RateLimitError           = NewError(1000202, "请求过于频繁，请稍后再试")
	UserLoginStatusException = NewError(1000203, "用户登录状态异常")
	UserNotFoundError        = NewError(1000205, "用户不存在")
	SmsLimitError            = NewError(1000206, "短信发送过于频繁，请10分钟后再试")
	CacheSystemError         = NewError(1000207, "缓存系统错误")

	FailedEnterpriseListError      = NewError(1000301, "查询企业列表失败")
	UserNotJoinEd                  = NewError(1000302, "用户未入职任何企业")
	FailedParkListError            = NewError(1000302, "获取园区列表失败")
	FailedParkInfoError            = NewError(1000304, "获取园区信息失败")
	ParkDisableError               = NewError(1000305, "园区已禁用")
	EnterpriseDisableError         = NewError(1000306, "企业已禁用")
	UserRegisterError              = NewError(1000310, "用户注册失败")
	VerificationCodeError          = NewError(1000311, "验证码错误")
	AccountRegisteredError         = NewError(1000312, "账号已注册")
	UserLoginError                 = NewError(1000313, "用户登录失败")
	UserLoginPasswordError         = NewError(1000314, "用户登录密码错误")
	FailedUserInfoError            = NewError(1000315, "用户信息获取失败")
	SMSCodeSendError               = NewError(1000316, "短信验证码获取失败")
	FailedEnterpriseError          = NewError(1000317, "获取企业信息失败")
	FailedEnterpriseMemberError    = NewError(1000318, "获取企业成员列表失败")
	FailedEnterpriseCodeError      = NewError(1000319, "获取企业入驻码失败")
	FailedEnterpriseLeaseError     = NewError(1000320, "获取企业租聘信息失败")
	FailedEnterpriseRoleError      = NewError(1000321, "获取企业角色失败")
	BatchEnterpriseRemoveUserError = NewError(1000322, "批量移除企业成员失败")
	BatchAuthorizationError        = NewError(1000323, "批量分配企业角色失败")
	FailedUserEnterpriseInfoError  = NewError(1000324, "获取用户的企业信息失败")
	FailedUserListError            = NewError(1000349, "获取用户列表失败")
	FailedBulletinError            = NewError(1000325, "获取公告详情失败")
	UserLogoutError                = NewError(1000326, "用户退出登录失败")
	AccountNotRegisteredError      = NewError(1000327, "账号未注册")
	EnterpriseCodeError            = NewError(1000328, "无效的入驻码")
	UserVerifyError                = NewError(1000328, "用户实名认证失败")
	IDCardHasVerifyError           = NewError(1000329, "身份证号码已被认证")
	IDCardNotVerifyError           = NewError(1000330, "用户未认证")
	JoinEnterpriseError            = NewError(1000330, "加入企业或者学院失败")
	UserHasEnterpriseError         = NewError(1000331, "已经入职了企业或者学院")
	SaveEnterpriseError            = NewError(1000332, "保存企业信息失败")
	FailedUserRoleError            = NewError(1000334, "获取用户角色失败")
	FailedSwitchError              = NewError(1000335, "切换失败")
	GetCosCredentialError          = NewError(1000336, "获取临时密钥失败")
	NoPermissionError              = NewError(1000337, "没有权限")
	GetEnterpriseCategoryError     = NewError(1000338, "获取企业类目失败")
	CosUploadFileError             = NewError(1000339, "文件上传失败")
	UserNoVerifyError              = NewError(1000340, "用户未完成实名认证")
	UserBlockedError               = NewError(1000341, "您暂时不能使用该功能，如有疑问请联系园区管理员")
	AccountDisableError            = NewError(1000342, "该账号已被禁用，如有疑问请联系园区管理员")
	QrcodeTypeError                = NewError(1000343, "二维码类型错误")
	FailedRoomNotFoundError        = NewError(1000344, "房间不存在")
	FailedRoomInfoError            = NewError(1000345, "获取房间信息失败")
	FailedBuildingNotFoundError    = NewError(1000346, "楼栋不存在")
	FailedBuildingInfoError        = NewError(1000347, "获取楼栋信息失败")
	ChangePhoneNumberError         = NewError(1000348, "换绑手机号失败")

	// 用户反馈相关
	UserFeedbackSubmitError = NewError(1000358, "提交用户反馈失败")
	UserFeedbackListError   = NewError(1000359, "获取用户反馈列表失败")
	UserFeedbackDetailError = NewError(1000360, "获取用户反馈详情失败")

	// 用户信息修改申请相关
	UserInfoUpdatePostDeactivateUserlyError        = NewError(1000361, "提交用户信息修改申请失败")
	UserInfoUpdatePostDeactivateUserlyPendingError = NewError(1000362, "存在未审批的用户信息修改申请")

	// 人才相关
	FailedTalentListError                            = NewError(1000350, "获取人才列表失败")
	FailedTalentDetailError                          = NewError(1000351, "获取人才详情失败")
	FailedTalentContactRecordDetailError             = NewError(1000352, "获取人才联系记录详情失败")
	FailedTalentContactError                         = NewError(1000353, "联系人才失败")
	FailedTalentContactRecordFoundError              = NewError(1000354, "用户的人才联系记录已存在")
	FailedTalentExperienceDetailError                = NewError(1000355, "获取人才经历详情失败")
	FailedTalentContactPostDeactivateUserroveInvaild = NewError(1000356, "非法的人才建联审批动作")
	FailedTalentContactRecordUpdateError             = NewError(1000357, "更新人才联系记录失败")
	FailedUserSpaceError                             = NewError(1000363, "获取用户的企业失败")

	// 课程/活动相关
	FailedActivityListError                 = NewError(1000400, "获取园区的活动/课程列表失败")
	FailedActivityDetailError               = NewError(1000401, "获取活动/课程详情失败")
	FailedActivityNotFoundError             = NewError(1000402, "用户的活动/课程不存在")
	FailedActivityFoundError                = NewError(1000403, "用户的活动/课程已存在")
	FailedActivitySignUpRecordDetailError   = NewError(1000404, "获取活动/课程报名记录详情失败")
	FailedActivityExecuteSignUpError        = NewError(1000405, "报名活动/课程失败")
	FailedActivityCancelSignUpError         = NewError(1000406, "取消报名活动/课程失败")
	FailedActivitySignUpRecordNotFoundError = NewError(1000407, "用户的活动/课程报名记录不存在")
	FailedActivitySignUpRecordFoundError    = NewError(1000408, "用户的活动/课程报名记录已存在")
	FailedActivityParkListError             = NewError(1000409, "查询园区活动/课程列表失败")

	// 消息相关
	FailedMessageListError         = NewError(1000450, "查询消息列表失败")
	FailedMessageDetailError       = NewError(1000451, "获取消息详情失败")
	FailedMessageUpdateStatusError = NewError(1000452, "更新消息状态失败")
	UserMessageSaveError           = NewError(1000453, "保存消息记录失败")

	// 服务相关
	FailedServiceDetailError                                        = NewError(1000500, "获取服务详情失败")
	FailedServiceOrderDetailError                                   = NewError(1000501, "获取服务单详情失败")
	FailedServiceProcessNodeDetailError                             = NewError(1000502, "获取服务流程节点详情失败")
	FailedServiceProcessNodePostDeactivateUserroverDetailError      = NewError(1000503, "获取服务流程节点审批人详情失败")
	FailedServiceOrderPlaceError                                    = NewError(1000504, "提交（创建）服务单失败")
	FailedServiceOrderProcessNodePostDeactivateUserroverDeleteError = NewError(1000505, "删除服务单流程节点审批人失败")
	FailedServiceOrderProcessNodeDeleteError                        = NewError(1000506, "删除服务单流程节点失败")
	FailedServiceOrderRecallError                                   = NewError(1000507, "删除服务单失败")
	FailedServiceOrderUpdateError                                   = NewError(1000508, "更新服务单失败")
	FailedServiceOrderPostDeactivateUserroverActionExistsError      = NewError(1000509, "服务单审批人存在审批行为")
	FailedServiceCategoryError                                      = NewError(1000510, "获取服务大类失败")
	FailedServiceSubcategoryError                                   = NewError(1000511, "获取服务小类失败")
	FailedServiceListError                                          = NewError(1000512, "获取服务列表失败")
	FailedServiceOrderListError                                     = NewError(1000513, "获取服务单列表失败")
	FailedServiceOrderProcessNodePostDeactivateUserroverUpdateError = NewError(1000514, "更新服务单流程节点审批人失败")
	FailedServiceOrderProcessNodeUpdateError                        = NewError(1000515, "更新服务单流程节点失败")
	FailedServiceNotFoundError                                      = NewError(1000516, "服务不存在")
	FailedServiceOrderNotFoundError                                 = NewError(1000517, "服务单不存在")

	// 商户相关
	FailedMerchantLoginError = NewError(1000600, "商户登录失败")
	LockedMerchantLoginError = NewError(1000601, "商户账号已被冻结，5分钟后重试")
	UserBalanceNotEnough     = NewError(1000602, "用户余额不足")

	// 开票申请相关
	FailedInvoiceServiceOrderError = NewError(1000700, "提交开票服务单失败")

	// 公共空间预约相关
	FailedPublicSpaceServiceOrderError = NewError(1000702, "提交公共空间服务单失败")
	FailedPublicSpaceOccupationError   = NewError(1000701, "提交公共空间占用失败")
	PublicSpaceAlreadyOccupiedError    = NewError(1000703, "该公共空间在所选时间段已被占用")

	// ai聊天相关
	GetSessionFailedError         = NewError(1000800, "获取用户会话失败")
	GetSessionHistoryFailedError  = NewError(1000801, "获取用户会话历史失败")
	RecordChatMessagesFailedError = NewError(1000802, "保存聊天记录失败")
	CallLLMFailedError            = NewError(1000803, "调用知识库问答服务失败")

	// 合作需求上传相关
	FailedCORequirementServiceOrderError = NewError(1000850, "提交合作需求服务单失败")
	FailedCORequirementRepeatContact     = NewError(1000851, "重复对接合作需求")

	// 合作实训上传相关
	FailedCOPracticalTrainingServiceOrderError = NewError(1000900, "提交合作实训服务单失败")
	FailedCOPracticalTrainingRepeatContact     = NewError(1000901, "重复对接合作实训")
	COPracticalTrainingCapacityFullError       = NewError(1000907, "报名人数已满")

	// 用户注销相关
	UserAlreadyDeactivatedError = NewError(1000902, "用户已注销或正在注销，不能重复操作")
	DeactivateUserdError        = NewError(1000903, "用户已注销失败")

	// 企业信息修改申请相关
	EnterpriseInfoUpdatePostDeactivateUserlyError        = NewError(1000904, "提交企业信息修改申请失败")
	EnterpriseInfoUpdatePostDeactivateUserlyPendingError = NewError(1000905, "存在未审批的企业信息修改申请")
	FailedDistrictListError                              = NewError(1000906, "查询地区列表失败")
)
