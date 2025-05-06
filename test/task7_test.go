package test

import (
	"testing"
	"time"
)

func isMorning(now time.Time) bool {
	hour := now.Hour()
	return hour >= 6 && hour < 12
}

func TestIsMorning(t *testing.T) {
	tests := []struct {
		name string
		time time.Time
		want bool
	}{
		{"morning", time.Date(2023, 1, 1, 8, 0, 0, 0, time.UTC), true},
		{"afternoon", time.Date(2023, 1, 1, 14, 0, 0, 0, time.UTC), false},
		{"night", time.Date(2023, 1, 1, 2, 0, 0, 0, time.UTC), false},
		{"edge 6am", time.Date(2023, 1, 1, 6, 0, 0, 0, time.UTC), true},
		{"edge 11:59am", time.Date(2023, 1, 1, 11, 59, 59, 0, time.UTC), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isMorning(tt.time); got != tt.want {
				t.Errorf("isMorning() = %v, want %v", got, tt.want)
			}
		})
	}
}
