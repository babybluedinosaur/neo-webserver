package main

import (
	"testing"
)

func TestGetDateIntervalByWeek(t *testing.T) {
	// Calendar week 4 of 2026: 19.01.2026 - 25.01.2026
	start, end := getDateIntervalByWeek(4)

	expectedStart := "2026-01-19"
	expectedEnd := "2026-01-25"

	if start.Format("2006-01-02") != expectedStart {
		t.Errorf("expected start %s, got %s", expectedStart, start.Format("2006-01-02"))
	}

	if end.Format("2006-01-02") != expectedEnd {
		t.Errorf("expected end %s, got %s", expectedEnd, end.Format("2006-01-02"))
	}
}
