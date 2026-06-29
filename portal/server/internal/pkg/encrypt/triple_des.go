package encrypt

import (
	"crypto/cipher"
	"crypto/des"
	"encoding/hex"
	"errors"
	"fmt"
)

// PKCS7Padding 对数据进行填充，补足到8字节的整数倍(3DES分组固定8字节)
func PKCS7Padding(data []byte, blockSize int) []byte {
	// 计算需要填充的字节数
	padding := blockSize - len(data)%blockSize
	// 生成填充的字节：每个字节的值都是填充的长度
	padText := make([]byte, padding)
	for i := range padText {
		padText[i] = byte(padding)
	}
	// 拼接填充后的数据
	return append(data, padText...)
}

// PKCS7UnPadding 解密后去除填充的字节
func PKCS7UnPadding(data []byte) ([]byte, error) {
	length := len(data)
	if length == 0 {
		return nil, errors.New("密文长度不能为空")
	}
	// 取出最后一个字节，它的值就是填充的长度
	unPadding := int(data[length-1])
	// 校验填充的合法性
	if unPadding < 1 || unPadding > des.BlockSize {
		return nil, errors.New("填充格式错误")
	}
	// 截取真实的原文
	return data[:(length - unPadding)], nil
}

// TripleDesCbcEncrypt 3DES-CBC加密：返回 十六进制密文字符串
func TripleDesCbcEncrypt(plainText []byte, key []byte, iv []byte) (string, error) {
	// 校验密钥：3DES密钥必须是24字节
	if len(key) != 24 {
		return "", errors.New("3DES密钥长度必须是24字节")
	}
	// 校验IV向量：CBC模式IV必须等于分组长度8字节
	if len(iv) != des.BlockSize {
		return "", errors.New("CBC模式IV向量长度必须是8字节")
	}
	// 创建3DES密码块
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return "", err
	}
	// 对原文进行填充
	plainText = PKCS7Padding(plainText, des.BlockSize)
	// 创建CBC加密模式
	mode := cipher.NewCBCEncrypter(block, iv)
	// 执行加密：加密后的数据会覆盖原切片
	cipherText := make([]byte, len(plainText))
	mode.CryptBlocks(cipherText, plainText)
	// 二进制密文转十六进制字符串返回（无乱码，易存储）
	return hex.EncodeToString(cipherText), nil
}

// TripleDesCbcDecrypt 3DES-CBC解密：入参是 十六进制密文字符串，返回 原始明文
func TripleDesCbcDecrypt(cipherHex string, key []byte, iv []byte) ([]byte, error) {
	// 校验密钥和IV
	if len(key) != 24 {
		return nil, errors.New("3DES密钥长度必须是24字节")
	}
	if len(iv) != des.BlockSize {
		return nil, errors.New("CBC模式IV向量长度必须是8字节")
	}
	// 十六进制密文 转 二进制密文
	cipherText, err := hex.DecodeString(cipherHex)
	if err != nil {
		return nil, fmt.Errorf("密文格式错误，不是合法的十六进制：%w", err)
	}
	// 创建3DES密码块
	block, err := des.NewTripleDESCipher(key)
	if err != nil {
		return nil, err
	}
	// 创建CBC解密模式
	mode := cipher.NewCBCDecrypter(block, iv)
	// 执行解密
	plainText := make([]byte, len(cipherText))
	mode.CryptBlocks(plainText, cipherText)
	// 去除填充的字节，得到原始明文
	plainText, err = PKCS7UnPadding(plainText)
	if err != nil {
		return nil, fmt.Errorf("解密失败，填充校验错误：%w", err)
	}
	return plainText, nil
}
