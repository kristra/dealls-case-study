package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCountWeekdays(t *testing.T) {
	tests := []struct {
		start    string
		end      string
		expected int
	}{
		{"2025-06-01", "2025-06-30", 21}, // June 2025 has 21 weekdays
		{"2025-06-01", "2025-06-01", 0},  // Sunday
		{"2025-06-03", "2025-06-03", 1},  // Tuesday
		{"2025-06-08", "2025-06-14", 5},  // full week (Sun to Sat)
		{"2025-06-14", "2025-06-08", 5},  // reverse order
		{"2025-06-15", "2025-06-15", 0},  // single Sunday
		{"2025-06-16", "2025-06-16", 1},  // single Monday
	}

	for _, tt := range tests {
		start, _ := time.Parse("2006-01-02", tt.start)
		end, _ := time.Parse("2006-01-02", tt.end)

		t.Run(tt.start+" to "+tt.end, func(t *testing.T) {
			actual := CountWeekdays(start, end)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
