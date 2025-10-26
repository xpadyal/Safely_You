// Package utils provides common utility functions used across the application.
// It includes mathematical operations, time parsing, and validation helpers.
package utils

import (
	"errors"
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

// ValidateTimestamp checks if a timestamp is within reasonable bounds
func ValidateTimestamp(t time.Time) error {
	now := time.Now()

	// Reject timestamps from more than 24 hours ago (likely data corruption)
	if t.Before(now.Add(-24 * time.Hour)) {
		return errors.New("timestamp too old (>24h)")
	}

	// Reject timestamps more than 5 minutes in the future (clock skew)
	if t.After(now.Add(5 * time.Minute)) {
		return errors.New("timestamp too far in future (>5min)")
	}

	return nil
}
