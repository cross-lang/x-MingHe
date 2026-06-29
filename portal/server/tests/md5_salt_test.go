package tests

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"testing"

	pkg "portal/internal/pkg"
)

func TestMD5WithSalt(t *testing.T) {
	tests := []struct {
		name     string
		password string
		salt     string
		expected string
	}{
		{
			name:     "simple password with salt",
			password: "password123",
			salt:     "salt1",
			expected: md5Hash("password123$salt1"),
		},
		{
			name:     "empty password with salt",
			password: "",
			salt:     "mysalt",
			expected: md5Hash("$mysalt"),
		},
		{
			name:     "password with empty salt",
			password: "mypassword",
			salt:     "",
			expected: md5Hash("mypassword$"),
		},
		{
			name:     "both empty",
			password: "",
			salt:     "",
			expected: md5Hash("$"),
		},
		{
			name:     "password with special chars",
			password: "p@$$w0rd!123",
			salt:     "#salt#",
			expected: md5Hash("p@$$w0rd!123$#salt#"),
		},
		{
			name:     "unicode characters",
			password: "密码123",
			salt:     "盐值",
			expected: md5Hash("密码123$盐值"),
		},
		{
			name:     "very long password",
			password: string(bytes.Repeat([]byte("a"), 1000)),
			salt:     "longsalt",
			expected: md5Hash(string(bytes.Repeat([]byte("a"), 1000)) + "$longsalt"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := pkg.MD5WithSalt(tt.password, tt.salt)

			if result != tt.expected {
				t.Errorf("MD5WithSalt(%q, %q) = %q, want %q",
					tt.password, tt.salt, result, tt.expected)
			}
		})
	}
}

func TestMD5WithSaltDeterministic(t *testing.T) {
	t.Run("same input produces same output", func(t *testing.T) {
		password := "testpassword"
		salt := "testsalt"

		result1 := pkg.MD5WithSalt(password, salt)
		result2 := pkg.MD5WithSalt(password, salt)
		result3 := pkg.MD5WithSalt(password, salt)

		if result1 != result2 || result2 != result3 {
			t.Error("MD5WithSalt should be deterministic")
		}
	})
}

func TestMD5WithSaltDifferent(t *testing.T) {
	t.Run("different inputs produce different outputs", func(t *testing.T) {
		password := "testpassword"
		salt := "testsalt"

		result1 := pkg.MD5WithSalt(password, salt)
		result2 := pkg.MD5WithSalt(password+"x", salt)
		result3 := pkg.MD5WithSalt(password, salt+"x")

		if result1 == result2 {
			t.Error("Different passwords should produce different hashes")
		}
		if result1 == result3 {
			t.Error("Different salts should produce different hashes")
		}
	})
}

func TestMD5WithSaltOrder(t *testing.T) {
	t.Run("password$ salt is correct order", func(t *testing.T) {
		password := "password123"
		salt := "salt1"

		result := pkg.MD5WithSalt(password, salt)
		expected := md5Hash(password + "$" + salt)

		if result != expected {
			t.Errorf("MD5WithSalt should hash password$salt, got different result")
		}

		// Verify that salt$password produces different result
		reverseExpected := md5Hash(salt + "$" + password)
		if reverseExpected == expected {
			t.Error("password$salt and salt$password should produce different hashes")
		}
	})
}

func TestMD5WithSaltOutputFormat(t *testing.T) {
	t.Run("output is valid hex string", func(t *testing.T) {
		result := pkg.MD5WithSalt("test", "salt")

		// MD5 produces 32 hex characters
		if len(result) != 32 {
			t.Errorf("MD5WithSalt should return 32 characters, got %d", len(result))
		}

		// Verify all characters are valid hex
		for _, c := range result {
			if !((c >= '0' && c <= '9') || (c >= 'a' && c <= 'f')) {
				t.Errorf("MD5WithSalt returned invalid hex character: %c", c)
			}
		}
	})
}

func TestMD5WithSaltPerformance(t *testing.T) {
	t.Run("can handle many hashes quickly", func(t *testing.T) {
		password := "testpassword"
		salt := "testsalt"
		iterations := 1000

		for i := 0; i < iterations; i++ {
			_ = pkg.MD5WithSalt(password, salt)
		}
		// Test passes if it doesn't timeout
	})
}

func md5Hash(input string) string {
	hash := md5.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}
