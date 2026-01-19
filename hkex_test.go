package marketchecker

import (
	"testing"
	"time"
)

func TestHKEX_MorningSession(t *testing.T) {
	hkex := NewHKEX()

	// Test morning session - Monday 10:00 AM HKT
	loc, err := time.LoadLocation("Asia/Hong_Kong")
	if err != nil {
		t.Fatalf("Failed to load timezone: %v", err)
	}
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
	loc, err := time.LoadLocation("Asia/Hong_Kong")
	if err != nil {
		t.Fatalf("Failed to load timezone: %v", err)
	}
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
	loc, err := time.LoadLocation("Asia/Hong_Kong")
	if err != nil {
		t.Fatalf("Failed to load timezone: %v", err)
	}
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
	loc, err := time.LoadLocation("Asia/Hong_Kong")
	if err != nil {
		t.Fatalf("Failed to load timezone: %v", err)
	}
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

func TestHKEX_HolidayClosed(t *testing.T) {
	hkex := NewHKEX()

	// Test Lunar New Year 2025 - should be closed even during regular hours
	loc, err := time.LoadLocation("Asia/Hong_Kong")
	if err != nil {
		t.Fatalf("Failed to load timezone: %v", err)
	}
	lunarNewYear := time.Date(2025, 1, 29, 10, 0, 0, 0, loc) // Jan 29, 2025 10:00 AM HKT

	if hkex.IsOpen(lunarNewYear) {
		t.Errorf("HKEX should be closed on Lunar New Year at %v", lunarNewYear)
	}

	status := hkex.GetStatus(lunarNewYear)
	if status != StatusClosed {
		t.Errorf("Expected status %s on Lunar New Year, got %s", StatusClosed, status)
	}

	// Test Christmas 2025
	christmas := time.Date(2025, 12, 25, 10, 0, 0, 0, loc)
	if hkex.IsOpen(christmas) {
		t.Errorf("HKEX should be closed on Christmas at %v", christmas)
	}

	status = hkex.GetStatus(christmas)
	if status != StatusClosed {
		t.Errorf("Expected status %s on Christmas, got %s", StatusClosed, status)
	}
}
