package timex_test

import (
	"testing"
	"time"

	timex "portal/internal/pkg/timex"
)

func TestParseChineseDateTime(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "AM time",
			input:   "2024年1月15日 上午09:30",
			wantErr: false,
		},
		{
			name:    "PM time",
			input:   "2024年12月25日 下午15:45",
			wantErr: false,
		},
		{
			name:    "midnight",
			input:   "2024年6月30日 上午00:00",
			wantErr: false,
		},
		{
			name:    "midnight PM",
			input:   "2024年6月30日 下午00:00",
			wantErr: false,
		},
		{
			name:    "invalid format",
			input:   "2024/01/15 09:30",
			wantErr: true,
		},
		{
			name:    "empty string",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := timex.ParseChineseDateTime(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseChineseDateTime(%q) expected error, got nil", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseChineseDateTime(%q) unexpected error: %v", tt.input, err)
				return
			}

			if result.IsZero() {
				t.Errorf("ParseChineseDateTime(%q) returned zero time", tt.input)
			}
		})
	}
}

func TestTimestampToTime(t *testing.T) {
	tests := []struct {
		name      string
		timestamp int64
		wantYear  int
		wantMonth time.Month
		wantDay   int
	}{
		{
			name:      "unix epoch",
			timestamp: 0,
			wantYear:  1970,
			wantMonth: time.January,
			wantDay:   1,
		},
		{
			name:      "current time around 2024",
			timestamp: 1704067200,
			wantYear:  2024,
			wantMonth: time.January,
			wantDay:   1,
		},
		{
			name:      "future timestamp",
			timestamp: 1735689600,
			wantYear:  2025,
			wantMonth: time.January,
			wantDay:   1,
		},
		{
			name:      "negative timestamp",
			timestamp: -86400,
			wantYear:  1969,
			wantMonth: time.December,
			wantDay:   31,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := timex.TimestampToTime(tt.timestamp)

			year, month, day := result.Date()
			if year != tt.wantYear || month != tt.wantMonth || day != tt.wantDay {
				t.Errorf("TimestampToTime(%d) = %v, want year=%d month=%d day=%d",
					tt.timestamp, result, tt.wantYear, tt.wantMonth, tt.wantDay)
			}
		})
	}
}

func TestTimestampMsToTime(t *testing.T) {
	tests := []struct {
		name      string
		timestamp int64
		wantYear  int
		wantMonth time.Month
		wantDay   int
	}{
		{
			name:      "milli timestamp 2024-01-01 12:30:45",
			timestamp: 1704113445000,
			wantYear:  2024,
			wantMonth: time.January,
			wantDay:   1,
		},
		{
			name:      "zero millisecond timestamp",
			timestamp: 0,
			wantYear:  1970,
			wantMonth: time.January,
			wantDay:   1,
		},
		{
			name:      "large millisecond timestamp",
			timestamp: 17356896000000,
			wantYear:  2025,
			wantMonth: time.January,
			wantDay:   1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := timex.TimestampMsToTime(tt.timestamp)

			year, month, day := result.Date()
			if year != tt.wantYear || month != tt.wantMonth || day != tt.wantDay {
				t.Errorf("TimestampMsToTime(%d) = %v, want %d-%02d-%02d",
					tt.timestamp, result, tt.wantYear, tt.wantMonth, tt.wantDay)
			}
		})
	}
}

func TestTimestampConsistency(t *testing.T) {
	t.Run("seconds and milliseconds produce same result", func(t *testing.T) {
		now := time.Now()
		secTimestamp := now.Unix()
		msTimestamp := now.UnixMilli()

		secTime := timex.TimestampToTime(secTimestamp)
		msTime := timex.TimestampMsToTime(msTimestamp * 1000)

		diff := secTime.Sub(msTime)
		if diff < 0 {
			diff = -diff
		}

		// Difference should be less than 1 second
		if diff >= time.Second {
			t.Errorf("TimestampToTime and TimestampMsToTime differ by %v, expected < 1s", diff)
		}
	})
}
