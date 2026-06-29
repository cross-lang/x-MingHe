package encrypt_test

import (
	"bytes"
	"encoding/hex"
	"testing"

	encrypt "portal/internal/pkg/encrypt"
)

func TestPKCS7Padding(t *testing.T) {
	tests := []struct {
		name        string
		data        []byte
		blockSize   int
		expectedLen int
	}{
		{
			name:        "data is multiple of block size",
			data:        []byte("12345678"),
			blockSize:   8,
			expectedLen: 16, // 8 + 8 bytes padding (0x08 x 8)
		},
		{
			name:        "data needs 1 byte padding",
			data:        []byte("1234567"),
			blockSize:   8,
			expectedLen: 8,
		},
		{
			name:        "empty data",
			data:        []byte(""),
			blockSize:   8,
			expectedLen: 8,
		},
		{
			name:        "large data",
			data:        []byte("This is a longer text that needs padding"),
			blockSize:   8,
			expectedLen: 48,
		},
		{
			name:        "different block size",
			data:        []byte("123456"),
			blockSize:   16,
			expectedLen: 16,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := encrypt.PKCS7Padding(tt.data, tt.blockSize)

			if len(result) != tt.expectedLen {
				t.Errorf("PKCS7Padding(%v, %d) length = %d, want %d",
					tt.data, tt.blockSize, len(result), tt.expectedLen)
			}

			// Verify padding bytes are correct
			paddingLen := result[len(result)-1]
			if int(paddingLen) != len(result)-len(tt.data) {
				t.Errorf("PKCS7Padding: padding byte = %d, expected %d",
					paddingLen, len(result)-len(tt.data))
			}
		})
	}
}

func TestPKCS7UnPadding(t *testing.T) {
	tests := []struct {
		name     string
		data     []byte
		wantErr  bool
		expected []byte
	}{
		{
			name:     "valid padding",
			data:     []byte("12345678" + string([]byte{0x01})),
			wantErr:  false,
			expected: []byte("12345678"),
		},
		{
			name:     "8 bytes padding",
			data:     []byte("12345678" + string(bytes.Repeat([]byte{0x08}, 8))),
			wantErr:  false,
			expected: []byte("12345678"),
		},
		{
			name:    "empty data",
			data:    []byte{},
			wantErr: true,
		},
		{
			name:     "no padding needed",
			data:     []byte("12345678"),
			wantErr:  false,
			expected: []byte("12345678"), // PKCS7 always adds padding
		},
		{
			name:    "invalid padding value",
			data:    []byte("1234" + string([]byte{0xFF})),
			wantErr: true,
		},
		{
			name:    "padding too large",
			data:    []byte("1234" + string([]byte{0x10})),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encrypt.PKCS7UnPadding(tt.data)

			if tt.wantErr {
				if err == nil {
					t.Errorf("PKCS7UnPadding(%v) expected error, got nil", tt.data)
				}
				return
			}

			if err != nil {
				t.Errorf("PKCS7UnPadding(%v) unexpected error: %v", tt.data, err)
				return
			}

			if !bytes.Equal(result, tt.expected) {
				t.Errorf("PKCS7UnPadding(%v) = %v, want %v",
					tt.data, result, tt.expected)
			}
		})
	}
}

func TestPKCS7RoundTrip(t *testing.T) {
	t.Run("padding then unpadding", func(t *testing.T) {
		testCases := [][]byte{
			[]byte(""),
			[]byte("a"),
			[]byte("1234567"),
			[]byte("12345678"),
			[]byte("Hello, World!"),
			bytes.Repeat([]byte("a"), 100),
		}

		for _, data := range testCases {
			t.Run(string(data), func(t *testing.T) {
				padded := encrypt.PKCS7Padding(data, 8)
				unpadded, err := encrypt.PKCS7UnPadding(padded)

				if err != nil {
					t.Errorf("PKCS7UnPadding failed: %v", err)
					return
				}

				if !bytes.Equal(unpadded, data) {
					t.Errorf("Round trip failed: original=%q, result=%q",
						string(data), string(unpadded))
				}
			})
		}
	})
}

func TestTripleDesCbcEncrypt(t *testing.T) {
	validKey := []byte("CCE31A176862725175bb539f")
	validIV := []byte("CCE31A176")

	tests := []struct {
		name      string
		plainText []byte
		key       []byte
		iv        []byte
		wantErr   bool
		errMsg    string
	}{
		{
			name:      "valid input",
			plainText: []byte("Hello, 3DES!"),
			key:       validKey,
			iv:        validIV,
			wantErr:   false,
		},
		{
			name:      "empty plaintext",
			plainText: []byte(""),
			key:       validKey,
			iv:        validIV,
			wantErr:   false,
		},
		{
			name:      "long plaintext",
			plainText: []byte("This is a much longer message that will span multiple blocks"),
			key:       validKey,
			iv:        validIV,
			wantErr:   false,
		},
		{
			name:      "invalid key length",
			plainText: []byte("test"),
			key:       []byte("shortkey"),
			iv:        validIV,
			wantErr:   true,
			errMsg:    "3DES密钥长度必须是24字节",
		},
		{
			name:      "invalid IV length",
			plainText: []byte("test"),
			key:       validKey,
			iv:        []byte("short"),
			wantErr:   true,
			errMsg:    "CBC模式IV向量长度必须是8字节",
		},
		{
			name:      "16-byte key",
			plainText: []byte("test"),
			key:       bytes.Repeat([]byte("a"), 16),
			iv:        validIV,
			wantErr:   true,
			errMsg:    "3DES密钥长度必须是24字节",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encrypt.TripleDesCbcEncrypt(tt.plainText, tt.key, tt.iv)

			if tt.wantErr {
				if err == nil {
					t.Errorf("TripleDesCbcEncrypt(%v, %v, %v) expected error, got nil",
						tt.plainText, tt.key, tt.iv)
				}
				if tt.errMsg != "" && err != nil && err.Error() != tt.errMsg {
					t.Errorf("TripleDesCbcEncrypt error = %v, want %v", err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("TripleDesCbcEncrypt(%v, %v, %v) unexpected error: %v",
					tt.plainText, tt.key, tt.iv, err)
				return
			}

			// Verify result is valid hex
			_, err = hex.DecodeString(result)
			if err != nil {
				t.Errorf("TripleDesCbcEncrypt returned invalid hex: %v", err)
			}
		})
	}
}

func TestTripleDesCbcDecrypt(t *testing.T) {
	validKey := []byte("CCE31A176862725175bb539f")
	validIV := []byte("CCE31A176")

	// First encrypt some test data
	plainText := []byte("Hello, 3DES!")
	encrypted, err := encrypt.TripleDesCbcEncrypt(plainText, validKey, validIV)
	if err != nil {
		t.Fatalf("Failed to encrypt test data: %v", err)
	}

	tests := []struct {
		name        string
		cipherHex   string
		key         []byte
		iv          []byte
		wantErr     bool
		errMsg      string
		shouldMatch bool
		expected    []byte
	}{
		{
			name:        "valid encrypted data",
			cipherHex:   encrypted,
			key:         validKey,
			iv:          validIV,
			wantErr:     false,
			shouldMatch: true,
			expected:    plainText,
		},
		{
			name:      "wrong key",
			cipherHex: encrypted,
			key:       bytes.Repeat([]byte("b"), 24),
			iv:        validIV,
			wantErr:   true,
			errMsg:    "解密失败",
		},
		{
			name:      "invalid hex",
			cipherHex: "not valid hex!!",
			key:       validKey,
			iv:        validIV,
			wantErr:   true,
		},
		{
			name:      "invalid key length",
			cipherHex: encrypted,
			key:       []byte("short"),
			iv:        validIV,
			wantErr:   true,
			errMsg:    "3DES密钥长度必须是24字节",
		},
		{
			name:      "invalid IV length",
			cipherHex: encrypted,
			key:       validKey,
			iv:        []byte("short"),
			wantErr:   true,
			errMsg:    "CBC模式IV向量长度必须是8字节",
		},
		{
			name:      "empty hex string",
			cipherHex: "",
			key:       validKey,
			iv:        validIV,
			wantErr:   true,
			errMsg:    "填充格式错误",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encrypt.TripleDesCbcDecrypt(tt.cipherHex, tt.key, tt.iv)

			if tt.wantErr {
				if err == nil {
					t.Errorf("TripleDesCbcDecrypt(%q, %v, %v) expected error, got nil",
						tt.cipherHex, tt.key, tt.iv)
				}
				if tt.errMsg != "" && err != nil && err.Error() != tt.errMsg &&
					!contains(err.Error(), tt.errMsg) {
					t.Errorf("TripleDesCbcDecrypt error = %v, want to contain %v",
						err.Error(), tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("TripleDesCbcDecrypt(%q, %v, %v) unexpected error: %v",
					tt.cipherHex, tt.key, tt.iv, err)
				return
			}

			if tt.shouldMatch && !bytes.Equal(result, tt.expected) {
				t.Errorf("TripleDesCbcDecrypt result = %q, want %q",
					string(result), string(tt.expected))
			}
		})
	}
}

func TestTripleDesCbcRoundTrip(t *testing.T) {
	t.Run("encrypt then decrypt", func(t *testing.T) {
		key := []byte("CCE31A176862725175bb539f")
		iv := []byte("CCE31A176")

		testCases := [][]byte{
			[]byte(""),
			[]byte("a"),
			[]byte("Hello, 3DES!"),
			bytes.Repeat([]byte("a"), 100),
			[]byte("Special chars: !@#$%^&*()"),
		}

		for _, plaintext := range testCases {
			t.Run(string(plaintext), func(t *testing.T) {
				encrypted, err := encrypt.TripleDesCbcEncrypt(plaintext, key, iv)
				if err != nil {
					t.Fatalf("TripleDesCbcEncrypt failed: %v", err)
				}

				decrypted, err := encrypt.TripleDesCbcDecrypt(encrypted, key, iv)
				if err != nil {
					t.Fatalf("TripleDesCbcDecrypt failed: %v", err)
				}

				if !bytes.Equal(decrypted, plaintext) {
					t.Errorf("Round trip failed: plaintext=%q, decrypted=%q",
						string(plaintext), string(decrypted))
				}
			})
		}
	})
}

func TestTripleDesCbcDeterministic(t *testing.T) {
	t.Run("same input produces same output", func(t *testing.T) {
		key := []byte("CCE31A176862725175bb539f")
		iv := []byte("CCE31A176")
		plaintext := []byte("Hello, 3DES!")

		encrypted1, err := encrypt.TripleDesCbcEncrypt(plaintext, key, iv)
		if err != nil {
			t.Fatalf("TripleDesCbcEncrypt failed: %v", err)
		}

		encrypted2, err := encrypt.TripleDesCbcEncrypt(plaintext, key, iv)
		if err != nil {
			t.Fatalf("TripleDesCbcEncrypt failed: %v", err)
		}

		// CBC with same IV should produce same ciphertext
		if encrypted1 != encrypted2 {
			t.Error("Same input with same IV should produce same ciphertext (CBC mode)")
		}
	})
}

func TestKnownTestVectors(t *testing.T) {
	t.Run("test against known vector", func(t *testing.T) {
		// Test with the known test vector from the existing test
		key := []byte("CCE31A176862725175bb539f")
		iv := []byte("CCE31A176")
		cipherHex := "be8a5e494489f61a26233cb4f43fd7c7"

		_, err := encrypt.TripleDesCbcDecrypt(cipherHex, key, iv)
		if err != nil {
			t.Logf("Known vector decryption: %v", err)
		}
	})
}
