package marketchecker

import (
	"testing"
	"time"
)

func TestHKEX_MorningSession(t *testing.T) {
	hkex := NewHKEX()

	// Test morning session - Monday 10:00 AM HKT
	loc, _ := time.LoadLocation("Asia/Hong_Kong")
	morningTime := time.Date(2026, 1, 19, 10, 0, 0, 0, loc) // Monday

	if !hkex.IsOpen(morningTime) {
		t.Errorf("HKEX should be open at %v", morningTime)
	}

	status := hkex.GetStatus(morningTime)
	if status != StatusOpen {
		t.Errorf("Expected status %s, got %s", StatusOpen, status)
	}
}

func TestHKEX_AfternoonSession(t *testing.T) {
	hkex := NewHKEX()

	// Test afternoon session - Monday 2:00 PM HKT
	loc, _ := time.LoadLocation("Asia/Hong_Kong")
	afternoonTime := time.Date(2026, 1, 19, 14, 0, 0, 0, loc) // Monday

	if !hkex.IsOpen(afternoonTime) {
		t.Errorf("HKEX should be open at %v", afternoonTime)
	}

	status := hkex.GetStatus(afternoonTime)
	if status != StatusOpen {
		t.Errorf("Expected status %s, got %s", StatusOpen, status)
	}
}

func TestHKEX_LunchBreak(t *testing.T) {
	hkex := NewHKEX()

	// Test lunch break - Monday 12:30 PM HKT
	loc, _ := time.LoadLocation("Asia/Hong_Kong")
	lunchTime := time.Date(2026, 1, 19, 12, 30, 0, 0, loc) // Monday

	if hkex.IsOpen(lunchTime) {
		t.Errorf("HKEX should be closed during lunch at %v", lunchTime)
	}

	status := hkex.GetStatus(lunchTime)
	if status != StatusClosed {
		t.Errorf("Expected status %s, got %s", StatusClosed, status)
	}
}

func TestHKEX_WeekendClosed(t *testing.T) {
	hkex := NewHKEX()

	// Test weekend - Saturday 10:00 AM HKT
	loc, _ := time.LoadLocation("Asia/Hong_Kong")
	weekendTime := time.Date(2026, 1, 17, 10, 0, 0, 0, loc) // Saturday

	if hkex.IsOpen(weekendTime) {
		t.Errorf("HKEX should be closed on weekend at %v", weekendTime)
	}

	status := hkex.GetStatus(weekendTime)
	if status != StatusClosed {
		t.Errorf("Expected status %s, got %s", StatusClosed, status)
	}
}

func TestHKEX_Name(t *testing.T) {
	hkex := NewHKEX()
	if hkex.Name() != "HKEX" {
		t.Errorf("Expected name 'HKEX', got '%s'", hkex.Name())
	}
}
