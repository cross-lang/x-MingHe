package stringx_test

import (
	"regexp"
	"testing"

	stringx "portal/internal/pkg/stringx"
)

func TestGenerateRand(t *testing.T) {
	tests := []struct {
		name        string
		charset     string
		length      int
		shouldFail  bool
		expectedLen int
	}{
		{
			name:        "numeric only",
			charset:     stringx.CharSetV1,
			length:      6,
			expectedLen: 6,
		},
		{
			name:        "lowercase letters",
			charset:     stringx.CharSetV2,
			length:      10,
			expectedLen: 10,
		},
		{
			name:        "uppercase letters",
			charset:     stringx.CharSetV3,
			length:      8,
			expectedLen: 8,
		},
		{
			name:        "mixed alphanumeric",
			charset:     stringx.CharSetV6,
			length:      12,
			expectedLen: 12,
		},
		{
			name:        "default charset",
			charset:     "",
			length:      10,
			expectedLen: 10,
		},
		{
			name:       "zero length",
			charset:    stringx.CharSetV1,
			length:     0,
			shouldFail: true,
		},
		{
			name:       "negative length",
			charset:    stringx.CharSetV1,
			length:     -5,
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stringx.GenerateRand(tt.charset, tt.length)

			if tt.shouldFail {
				if result != "" {
					t.Errorf("GenerateRand(%q, %d) should return empty string for invalid input, got %q",
						tt.charset, tt.length, result)
				}
				return
			}

			if len(result) != tt.expectedLen {
				t.Errorf("GenerateRand(%q, %d) returned length %d, want %d",
					tt.charset, tt.length, len(result), tt.expectedLen)
			}

			// Verify all characters are from the charset
			validCharset := tt.charset
			if validCharset == "" {
				validCharset = stringx.CharSetV6
			}
			for _, char := range result {
				if !contains(validCharset, byte(char)) {
					t.Errorf("GenerateRand(%q, %d) returned invalid character %c",
						tt.charset, tt.length, char)
				}
			}
		})
	}
}

func contains(s string, char byte) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == char {
			return true
		}
	}
	return false
}

func TestGenerateRandUniqueness(t *testing.T) {
	t.Run("generates unique values", func(t *testing.T) {
		length := 1000
		results := make(map[string]bool)

		for i := 0; i < length; i++ {
			result := stringx.GenerateRand(stringx.CharSetV6, 12)
			if results[result] {
				t.Logf("Warning: duplicate value generated: %s", result)
			}
			results[result] = true
		}

		if len(results) < length*90/100 {
			t.Errorf("GenerateRand produced too few unique values: %d out of %d", len(results), length)
		}
	})
}

func TestGenerateRandNumeric(t *testing.T) {
	t.Run("generates only numbers", func(t *testing.T) {
		result := stringx.GenerateRand(stringx.CharSetV1, 6)
		matched, _ := regexp.MatchString("^[0-9]{6}$", result)
		if !matched {
			t.Errorf("GenerateRand with CharSetV1 should return only digits, got: %s", result)
		}
	})
}
