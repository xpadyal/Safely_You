package utils

import (
	"fmt"
	"time"
)

// Round2 rounds a float64 to 2 decimal places
func Round2(f float64) float64 {
	return float64(int(f*100+0.5)) / 100
}

// ParseRFC3339 parses an RFC3339 timestamp string and returns a UTC time
func ParseRFC3339(ts string) (time.Time, error) {
	t, err := time.Parse(time.RFC3339, ts)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid sent_at format: %s", ts)
	}
	return t.UTC(), nil
}
