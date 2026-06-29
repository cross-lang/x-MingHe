package encrypt_test

import (
	"bytes"
	"testing"

	encrypt "portal/internal/pkg/encrypt"
)

func TestAESEncrypt(t *testing.T) {
	tests := []struct {
		name       string
		plaintext  []byte
		key        []byte
		wantErr    bool
		errMessage string
	}{
		{
			name:      "valid 32-byte key",
			plaintext: []byte("Hello, World!"),
			key:       bytes.Repeat([]byte("a"), 32),
			wantErr:   false,
		},
		{
			name:      "empty plaintext",
			plaintext: []byte(""),
			key:       bytes.Repeat([]byte("b"), 32),
			wantErr:   false,
		},
		{
			name:      "long plaintext",
			plaintext: []byte("This is a longer message that needs encryption"),
			key:       bytes.Repeat([]byte("c"), 32),
			wantErr:   false,
		},
		{
			name:       "16-byte key (invalid)",
			plaintext:  []byte("test"),
			key:        bytes.Repeat([]byte("d"), 16),
			wantErr:    true,
			errMessage: "AES key must be 32 bytes for AES-256",
		},
		{
			name:       "empty key",
			plaintext:  []byte("test"),
			key:        []byte{},
			wantErr:    true,
			errMessage: "AES key must be 32 bytes for AES-256",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encrypt.AESEncrypt(tt.plaintext, tt.key)

			if tt.wantErr {
				if err == nil {
					t.Errorf("AESEncrypt(%v, %v) expected error, got nil", tt.plaintext, tt.key)
				}
				if tt.errMessage != "" && err.Error() != tt.errMessage {
					t.Errorf("AESEncrypt(%v, %v) error = %v, want %v",
						tt.plaintext, tt.key, err.Error(), tt.errMessage)
				}
				return
			}

			if err != nil {
				t.Errorf("AESEncrypt(%v, %v) unexpected error: %v", tt.plaintext, tt.key, err)
				return
			}

			if result == "" {
				t.Error("AESEncrypt returned empty string")
			}
		})
	}
}

func TestAESDecrypt(t *testing.T) {
	validKey := bytes.Repeat([]byte("a"), 32)
	plaintext := []byte("Hello, World!")
	encrypted, err := encrypt.AESEncrypt(plaintext, validKey)
	if err != nil {
		t.Fatalf("Failed to encrypt for test: %v", err)
	}

	tests := []struct {
		name             string
		encryptedBase64  string
		key              []byte
		wantErr          bool
		errMessage       string
		shouldMatchPlain bool
		expectedPlain    []byte
	}{
		{
			name:             "valid encrypted data",
			encryptedBase64:  encrypted,
			key:              validKey,
			wantErr:          false,
			shouldMatchPlain: true,
			expectedPlain:    plaintext,
		},
		{
			name:            "wrong key",
			encryptedBase64: encrypted,
			key:             bytes.Repeat([]byte("b"), 32),
			wantErr:         true,
			errMessage:      "decryption failed",
		},
		{
			name:            "invalid base64",
			encryptedBase64: "not valid base64!!",
			key:             validKey,
			wantErr:         true,
		},
		{
			name:            "16-byte key",
			encryptedBase64: encrypted,
			key:             bytes.Repeat([]byte("c"), 16),
			wantErr:         true,
			errMessage:      "AES key must be 32 bytes",
		},
		{
			name:            "empty key",
			encryptedBase64: encrypted,
			key:             []byte{},
			wantErr:         true,
			errMessage:      "AES key must be 32 bytes",
		},
		{
			name:            "empty encrypted",
			encryptedBase64: "",
			key:             validKey,
			wantErr:         true,
			errMessage:      "ciphertext too short",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := encrypt.AESDecrypt(tt.encryptedBase64, tt.key)

			if tt.wantErr {
				if err == nil {
					t.Errorf("AESDecrypt(%q, %v) expected error, got nil", tt.encryptedBase64, tt.key)
				}
				if tt.errMessage != "" && err != nil && err.Error() != tt.errMessage &&
					!contains(err.Error(), tt.errMessage) {
					t.Errorf("AESDecrypt(%q, %v) error = %v, want to contain %v",
						tt.encryptedBase64, tt.key, err.Error(), tt.errMessage)
				}
				return
			}

			if err != nil {
				t.Errorf("AESDecrypt(%q, %v) unexpected error: %v", tt.encryptedBase64, tt.key, err)
				return
			}

			if tt.shouldMatchPlain && !bytes.Equal(result, tt.expectedPlain) {
				t.Errorf("AESDecrypt(%q, %v) = %v, want %v",
					tt.encryptedBase64, tt.key, result, tt.expectedPlain)
			}
		})
	}
}

func TestAESCiphertextUniqueness(t *testing.T) {
	t.Run("same plaintext produces different ciphertext", func(t *testing.T) {
		key := bytes.Repeat([]byte("test"), 8)
		plaintext := []byte("Hello, World!")

		encrypted1, err := encrypt.AESEncrypt(plaintext, key)
		if err != nil {
			t.Fatalf("AESEncrypt failed: %v", err)
		}

		encrypted2, err := encrypt.AESEncrypt(plaintext, key)
		if err != nil {
			t.Fatalf("AESEncrypt failed: %v", err)
		}

		// AES-GCM uses random nonce, so same plaintext should produce different ciphertext
		if encrypted1 == encrypted2 {
			t.Error("Same plaintext should produce different ciphertext due to random nonce")
		}
	})
}

func TestAESRoundTrip(t *testing.T) {
	t.Run("encrypt then decrypt", func(t *testing.T) {
		key := bytes.Repeat([]byte("roundtripkey12345678"), 2) // 32 bytes

		testCases := [][]byte{
			[]byte(""),
			[]byte("a"),
			[]byte("Hello, World!"),
			bytes.Repeat([]byte("a"), 100),
			[]byte("Special chars: !@#$%^&*()"),
		}

		for _, plaintext := range testCases {
			t.Run(string(plaintext), func(t *testing.T) {
				encrypted, err := encrypt.AESEncrypt(plaintext, key)
				if err != nil {
					t.Fatalf("AESEncrypt failed: %v", err)
				}

				decrypted, err := encrypt.AESDecrypt(encrypted, key)
				if err != nil {
					t.Fatalf("AESDecrypt failed: %v", err)
				}

				if !bytes.Equal(decrypted, plaintext) {
					t.Errorf("Round trip failed: plaintext=%q, decrypted=%q",
						string(plaintext), string(decrypted))
				}
			})
		}
	})
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) &&
		(s[:len(substr)] == substr || s[len(s)-len(substr):] == substr ||
			findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
