package encrypt

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"io"
)

// 加密算法：AES（Advanced Encryption Standard）
// 密匙长度：32 字节 → 对应 AES-256
// 工作模式：GCM
// 填充方式：无填充
// 编码方式: Base64
// Nonce（IV）:由 gcm.NonceSize() 决定，默认为 12 字节（96 bits）

func AESEncrypt(plaintext []byte, key []byte) (string, error) {
	if len(key) != 32 {
		return "", errors.New("AES key must be 32 bytes for AES-256")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)
	result := append(nonce, ciphertext...)

	return base64.StdEncoding.EncodeToString(result), nil
}

func AESDecrypt(encryptedBase64 string, key []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("AES key must be 32 bytes for AES-256")
	}

	data, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, errors.New("decryption failed: invalid ciphertext or key")
	}

	return plaintext, nil
}
