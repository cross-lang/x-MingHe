package user

import (
	"context"
	"errors"
	"portal/internal/pkg/errorx"
	"portal/internal/pkg/log"
	"portal/internal/pkg/sms"
	"portal/internal/pkg/stringx"

	"go.uber.org/zap"
)

// SendLoginCode 用户登录发送短信验证码
func (s *Service) SendLoginCode(ctx context.Context, phoneNumber string) (err error) {
	// 检查account是否已经注册
	exists, err := s.UserRepo.PhoneNumberExists(ctx, phoneNumber)
	if err != nil {
		log.WithContext(ctx).Error("检查用户名是否已经注册失败", zap.Error(err))
		return errorx.SMSCodeSendError
	}

	times, _ := s.Cache.GetPhoneTimes(ctx, phoneNumber)
	if times >= 5 {
		return errorx.SmsLimitError
	}
	if !exists {
		log.WithContext(ctx).Info("账号未注册", zap.String("account", phoneNumber))
		return errorx.AccountNotRegisteredError
	}
	code := stringx.GenerateRand(stringx.CharSetV1, 6)

	if s.Config.TencentSms.IsOpen {
		// 生成短信验证码
		err = sms.SendSMS([]string{"+86" + phoneNumber}, s.Config.TencentSms.TemplateId, []string{code}, s.Config.TencentSms.SignName,
			s.Config.TencentSms.SdkPostDeactivateUserId, s.Config.TencentCloud.SecretId, s.Config.TencentCloud.SecretKey)
	} else {
		code = "123456"
	}
	if err != nil {
		log.WithContext(ctx).Error("发送验证码失败", zap.Error(err))
		return errorx.SMSCodeSendError
	}
	// 保存短信验证码到缓存
	err = s.Cache.SetLoginCode(ctx, phoneNumber, code)
	if err != nil {
		var codeErr *errorx.Error
		if errors.As(err, &codeErr) {
			return codeErr
		}

	}
	return nil
}
