package sms

import (
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"
)

func SendSMS(
	phoneNumbers []string,
	templateID string,
	templateParamSet []string,
	signName string,
	sdkPostDeactivateUserID string,
	secretId string,
	secretKey string,
) error {
	// 1. 配置客户端凭证
	credential := common.NewCredential(secretId, secretKey)

	// 2. 配置客户端选项（指定地域，短信服务默认ap-wuhan）
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "sms.tencentcloudapi.com" // 短信服务接入点
	client, err := sms.NewClient(credential, "ap-wuhan", cpf)
	if err != nil {
		return fmt.Errorf("创建短信客户端失败: %v", err)
	}

	// 3. 构造发送请求参数
	request := sms.NewSendSmsRequest()
	request.SignName = common.StringPtr(signName)
	request.TemplateId = common.StringPtr(templateID)
	request.PhoneNumberSet = common.StringPtrs(phoneNumbers)
	request.TemplateParamSet = common.StringPtrs(templateParamSet) // 模板参数（无参数则传空切片）

	// 4. 发送请求并处理响应
	response, err := client.SendSms(request)
	if err != nil {
		// 捕获腾讯云SDK错误
		if sdkErr, ok := err.(*errors.TencentCloudSDKError); ok {
			return fmt.Errorf("SDK错误: 代码[%s] 消息[%s] 请求ID[%s]", sdkErr.Code, sdkErr.Message, sdkErr.RequestId)
		}
		return fmt.Errorf("发送短信失败: %v", err)
	}

	// 5. 检查单个手机号的发送结果
	for _, status := range response.Response.SendStatusSet {
		if *status.Code != "Ok" {
			return fmt.Errorf("手机号[%s]发送失败: 代码[%s] 消息[%s]", *status.PhoneNumber, *status.Code, *status.Message)
		}
	}

	fmt.Println("短信发送成功，响应:", response.ToJsonString())
	return nil
}
