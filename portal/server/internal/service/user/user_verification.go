package user

import (
	"context"
	"crypto/sha256"
	"fmt"
	"portal/internal/model"
	"portal/internal/pkg/encrypt"
	"portal/internal/pkg/errorx"
	"portal/internal/pkg/gormx"
	"portal/internal/pkg/log"
	"portal/internal/types"
	"strings"

	v20180301 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/faceid/v20180301"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// 实名认证响应结果和错误信息
var idCardVerificationResultMap = map[string]string{
	"0":  "姓名和身份证号一致",
	"-1": "姓名和身份证号不一致",
	"-2": "非法身份证号（长度、校验位等不正确）",
	"-3": "非法姓名（长度、格式等不正确）",
	"-4": "证件库服务异常",
	"-5": "证件库中无此身份证记录",
	"-6": "权威比对系统升级中，请稍后再试",
	"-7": "认证次数超过当日限制",
}

// UserVerification 用户实名认证
func (s *Service) UserVerification(ctx context.Context, user *model.XUser, form *types.UserVerificationReq) error {
	// 检查用户是否已完成认证
	if user.VerificationStatus == 1 {
		log.WithContext(ctx).Info("用户已完成实名认证", zap.Any("user_id", user.ID))
		return nil
	}
	// 检查身份证号码是否已被认证
	idCardHash := s.idCardHash(form.IdCard)
	idCardHasVerify, err := s.UserRepo.IDCardHasVerify(ctx, idCardHash)
	if err != nil {
		log.WithContext(ctx).Error("检查身份证号码是否已被认证失败", zap.Error(err))
		return errorx.UserVerifyError
	}
	if idCardHasVerify {
		return errorx.IDCardHasVerifyError
	}

	// 创建腾讯云身份信息认证client失
	secretId := s.Config.TencentCloud.SecretId
	secretKey := s.Config.TencentCloud.SecretKey
	client, err := v20180301.NewClientWithSecretId(secretId, secretKey, "")
	if err != nil {
		log.WithContext(ctx).Error("创建腾讯云身份信息认证client失败", zap.Error(err))
		return errorx.UserVerifyError
	}
	req := v20180301.NewIdCardVerificationRequest()
	req.IdCard = &form.IdCard
	req.Name = &form.Name
	resp, err := client.IdCardVerification(req)
	if err != nil {
		log.WithContext(ctx).Error("请求腾讯云身份信息认证接口失败", zap.Error(err))
		return errorx.UserVerifyError
	}
	// 认证结果处理
	if *resp.Response.Result != "0" {
		errMsg, ok := idCardVerificationResultMap[*resp.Response.Result]
		if ok {
			return errorx.NewError(errorx.UserVerifyError.Code, errMsg)
		}
		return errorx.UserVerifyError
	}

	// 对身份证号码进行加密存储
	idCard, err := encrypt.AESEncrypt([]byte(form.IdCard), []byte(s.Config.DataEncryptKey))
	if err != nil {
		log.WithContext(ctx).Error("用户实名认证信息脱敏失败", zap.Error(err))
		return errorx.UserVerifyError
	}

	// 保存用户实名认证信息
	err = gormx.Transaction(ctx, func(tx *gorm.DB) error {
		ctx = gormx.NewContext(ctx, tx)
		// 更新用户实名认证状态
		err = s.UserRepo.UpdateVerifyStatus(ctx, user.ID, form.Name)
		if err != nil {
			return err
		}
		// 保存用户实名认证信息
		return s.UserRepo.SaveVerify(ctx, &model.XUserIdentityVerification{
			UID:     user.ID,
			Channel: "TencentCloud",
			Name:    form.Name,
			IDCard:  idCard,
			Hash:    idCardHash,
		})
	})
	if err != nil {
		log.WithContext(ctx).Error("保存用户实名认证信息失败", zap.Error(err))
		return errorx.UserVerifyError
	}
	return nil
}

// 身份证号码hash
func (s *Service) idCardHash(idCard string) string {
	// 清理输入：去除空格等
	idCard = strings.TrimSpace(idCard)
	// 计算 SHA-256
	hash := sha256.Sum256([]byte(idCard))
	// 返回十六进制字符串
	return fmt.Sprintf("%x", hash)
}
