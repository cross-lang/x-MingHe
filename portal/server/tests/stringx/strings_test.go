package stringx_test

import (
	"reflect"
	"testing"

	stringx "portal/internal/pkg/stringx"
)

func TestSplit(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		sep      string
		expected []string
	}{
		{
			name:     "normal split",
			s:        "a,b,c",
			sep:      ",",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "empty string",
			s:        "",
			sep:      ",",
			expected: []string{},
		},
		{
			name:     "no separator",
			s:        "abc",
			sep:      ",",
			expected: []string{"abc"},
		},
		{
			name:     "multiple separators",
			s:        "a,,b,,c",
			sep:      ",",
			expected: []string{"a", "", "b", "", "c"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stringx.Split(tt.s, tt.sep)
			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("Split(%q, %q) = %v, want %v", tt.s, tt.sep, result, tt.expected)
			}
		})
	}
}

func TestStrLike(t *testing.T) {
	tests := []struct {
		name     string
		str      string
		expected string
	}{
		{
			name:     "normal string",
			str:      "test",
			expected: "%test%",
		},
		{
			name:     "string with percent",
			str:      "100%",
			expected: "%100\\%%",
		},
		{
			name:     "string with underscore",
			str:      "test_data",
			expected: "%test\\_data%",
		},
		{
			name:     "string with both percent and underscore",
			str:      "100%_data",
			expected: "%100\\%%\\_data%",
		},
		{
			name:     "empty string",
			str:      "",
			expected: "%%",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := stringx.StrLike(tt.str)
			if result != tt.expected {
				t.Errorf("StrLike(%q) = %q, want %q", tt.str, result, tt.expected)
			}
		})
	}
}

func TestGetTradeNo(t *testing.T) {
	t.Run("generates unique trade no", func(t *testing.T) {
		uniq := "TEST"
		tradeNo1 := stringx.GetTradeNo(uniq)
		tradeNo2 := stringx.GetTradeNo(uniq)

		// Check if the trade no has correct prefix
		if len(tradeNo1) < 12 {
			t.Errorf("GetTradeNo(%q) returned too short: %q", uniq, tradeNo1)
		}

		// Check if it contains the uniq prefix
		if tradeNo1[0:8] != "TRTEST"+uniq {
			t.Errorf("GetTradeNo(%q) should start with TRTEST%%s, got %s", uniq, tradeNo1[:8])
		}

		// Check if two calls generate different results
		if tradeNo1 == tradeNo2 {
			t.Errorf("GetTradeNo(%q) should generate different results, got same: %s", uniq, tradeNo1)
		}
	})
}
