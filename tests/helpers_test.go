package tests

import (
	"math"
	"testing"
	"time"

	"github.com/xpadyal/Safely_You/internal/models"
	"github.com/xpadyal/Safely_You/internal/store"
)

func mustTime(s string) time.Time {
	t, err := store.ParseRFC3339(s)
	if err != nil {
		panic(err)
	}
	return t
}

func TestUniqueMinuteCount(t *testing.T) {
	ts := []time.Time{
		mustTime("2025-10-25T10:00:01Z"),
		mustTime("2025-10-25T10:00:50Z"),
		mustTime("2025-10-25T10:01:00Z"),
	}
	if got := store.UniqueMinuteCount(ts); got != 2 {
		t.Fatalf("want 2 unique minutes, got %d", got)
	}
}

func TestMinutesBetweenFirstAndLast(t *testing.T) {
	ts := []time.Time{
		mustTime("2025-10-25T10:00:00Z"),
		mustTime("2025-10-25T10:02:00Z"),
		mustTime("2025-10-25T10:01:00Z"),
	}
	span, err := store.MinutesBetweenFirstAndLast(ts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if span != 3 {
		t.Fatalf("want span 3 (inclusive), got %d", span)
	}
}

func TestComputeUptimeBasic(t *testing.T) {
	// Heartbeats at 10:00 and 10:02, missing 10:01 → uptime=66.67%
	d := &models.Device{
		Heartbeats: []time.Time{
			mustTime("2025-10-25T10:00:00Z"),
			mustTime("2025-10-25T10:02:00Z"),
		},
	}
	u := store.ComputeUptime(d)
	if math.Abs(u-66.6667) > 0.2 {
		t.Fatalf("want ~66.67, got %v", u)
	}
}

func TestComputeUptimeDrop(t *testing.T) {
	// Heartbeats at 10:00 and 10:03, missing 10:01 and 10:02 → uptime=50%
	d := &models.Device{
		Heartbeats: []time.Time{
			mustTime("2025-10-25T10:00:00Z"),
			mustTime("2025-10-25T10:03:00Z"),
		},
	}
	u := store.ComputeUptime(d)
	if math.Abs(u-50.0) > 0.2 {
		t.Fatalf("want ~50.0, got %v", u)
	}
}

func TestComputeUptimeSingleMinute(t *testing.T) {
	// Multiple heartbeats in same minute → 100% uptime
	d := &models.Device{
		Heartbeats: []time.Time{
			mustTime("2025-10-25T10:00:05Z"),
			mustTime("2025-10-25T10:00:50Z"),
		},
	}
	u := store.ComputeUptime(d)
	if u != 100 {
		t.Fatalf("want 100 for single-minute window, got %v", u)
	}
}

func TestComputeAvgUpload(t *testing.T) {
	d := &models.Device{UploadTimes: []int64{4, 6}}
	avg, ok := store.ComputeAvgUpload(d)
	if !ok {
		t.Fatalf("expected ok=true")
	}
	// Average of 4 and 6 nanoseconds should be "5ns"
	if avg != "5ns" {
		t.Fatalf("want 5ns, got %q", avg)
	}

	d2 := &models.Device{UploadTimes: []int64{5, 6, 5}}
	avg2, ok2 := store.ComputeAvgUpload(d2)
	if !ok2 {
		t.Fatalf("expected ok2=true")
	}
	// Average of 5, 6, 5 nanoseconds should be "5ns"
	if avg2 != "5ns" {
		t.Fatalf("want 5ns, got %q", avg2)
	}
}
