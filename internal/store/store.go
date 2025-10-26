package store

import (
	"errors"
	"time"

	"github.com/xpadyal/Safely_You/internal/models"
)

// NewStore creates a new Store instance
func NewStore() *models.Store {
	return &models.Store{
		Devices: make(map[string]*models.Device),
	}
}

// EnsureDevice creates a device if it doesn't exist and returns it
func EnsureDevice(s *models.Store, id string) *models.Device {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	if d, ok := s.Devices[id]; ok {
		return d
	}
	d := &models.Device{}
	s.Devices[id] = d
	return d
}

// Counts unique minute buckets in UTC
// Example: 10:00:05 and 10:00:55 are in the same minute
func UniqueMinuteCount(times []time.Time) int {
	if len(times) == 0 {
		return 0
	}
	seen := make(map[string]struct{}, len(times))
	for _, t := range times {
		key := t.UTC().Format("2006-01-02T15:04")
		seen[key] = struct{}{}
	}
	return len(seen)
}

// Calculates the inclusive minute span between first and last timestamps
// This includes both the first and last minute buckets
func MinutesBetweenFirstAndLast(times []time.Time) (int, error) {
	if len(times) < 2 {
		return 0, errors.New("need at least two timestamps to form a window")
	}
	minT := times[0]
	maxT := times[0]
	for _, t := range times[1:] {
		if t.Before(minT) {
			minT = t
		}
		if t.After(maxT) {
			maxT = t
		}
	}
	if !maxT.After(minT) {
		return 0, errors.New("non-positive window")
	}
	span := int(maxT.Sub(minT)/time.Minute) + 1
	return span, nil
}

// Formula: (unique_minute_heartbeats / total_minute_span) * 100
func ComputeUptime(d *models.Device) (float64, error) {
	if d == nil || len(d.Heartbeats) == 0 {
		return 0.0, nil
	}
	uniq := UniqueMinuteCount(d.Heartbeats)
	windowMinutes, err := MinutesBetweenFirstAndLast(d.Heartbeats)
	if err != nil || windowMinutes == 0 {
		// All heartbeats in the same minute = 100% uptime
		if uniq > 0 {
			return 100.0, nil
		}
		return 0.0, nil
	}
	return (float64(uniq) / float64(windowMinutes)) * 100.0, nil
}

// Calculates average upload time and returns it as a duration string
// Returns ("0s", nil) if no data available
func ComputeAvgUpload(d *models.Device) (string, error) {
	if d == nil || len(d.UploadTimes) == 0 {
		return "0s", nil
	}
	var sum int64
	for _, v := range d.UploadTimes {
		sum += v
	}
	avg := sum / int64(len(d.UploadTimes))
	// Convert nanoseconds to duration and format as string
	duration := time.Duration(avg)
	return duration.String(), nil
}

// Device operations

func AddHeartbeat(s *models.Store, id string, t time.Time) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	d, ok := s.Devices[id]
	if !ok {
		return errors.New("device not found")
	}
	d.Heartbeats = append(d.Heartbeats, t.UTC())
	return nil
}

func AddUploadTime(s *models.Store, id string, v int64) error {
	s.Mu.Lock()
	defer s.Mu.Unlock()

	d, ok := s.Devices[id]
	if !ok {
		return errors.New("device not found")
	}
	d.UploadTimes = append(d.UploadTimes, v)
	return nil
}

// Returns the device data
func SnapshotDevice(s *models.Store, id string) (*models.Device, bool) {
	s.Mu.RLock()
	defer s.Mu.RUnlock()

	d, ok := s.Devices[id]
	return d, ok
}
