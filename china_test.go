package marketchecker

import (
	"testing"
	"time"
)

func TestChinaAShare_MorningSession(t *testing.T) {
	china := NewChinaAShare()

	// Test morning session - Monday 10:00 AM CST
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatalf("Failed to load timezone: %v", err)
	}
	morningTime := time.Date(2026, 1, 19, 10, 0, 0, 0, loc) // Monday

	if !china.IsOpen(morningTime) {
		t.Errorf("China A-Share should be open at %v", morningTime)
	}

	status := china.GetStatus(morningTime)
	if status != StatusOpen {
		t.Errorf("Expected status %s, got %s", StatusOpen, status)
	}
}

func TestChinaAShare_AfternoonSession(t *testing.T) {
	china := NewChinaAShare()

	// Test afternoon session - Monday 2:00 PM CST
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatalf("Failed to load timezone: %v", err)
	}
	afternoonTime := time.Date(2026, 1, 19, 14, 0, 0, 0, loc) // Monday

	if !china.IsOpen(afternoonTime) {
		t.Errorf("China A-Share should be open at %v", afternoonTime)
	}

	status := china.GetStatus(afternoonTime)
	if status != StatusOpen {
		t.Errorf("Expected status %s, got %s", StatusOpen, status)
	}
}

func TestChinaAShare_LunchBreak(t *testing.T) {
	china := NewChinaAShare()

	// Test lunch break - Monday 12:00 PM CST
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatalf("Failed to load timezone: %v", err)
	}
	lunchTime := time.Date(2026, 1, 19, 12, 0, 0, 0, loc) // Monday

	if china.IsOpen(lunchTime) {
		t.Errorf("China A-Share should be closed during lunch at %v", lunchTime)
	}

	status := china.GetStatus(lunchTime)
	if status != StatusClosed {
		t.Errorf("Expected status %s, got %s", StatusClosed, status)
	}
}

func TestChinaAShare_WeekendClosed(t *testing.T) {
	china := NewChinaAShare()

	// Test weekend - Saturday 10:00 AM CST
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatalf("Failed to load timezone: %v", err)
	}
	weekendTime := time.Date(2026, 1, 17, 10, 0, 0, 0, loc) // Saturday

	if china.IsOpen(weekendTime) {
		t.Errorf("China A-Share should be closed on weekend at %v", weekendTime)
	}

	status := china.GetStatus(weekendTime)
	if status != StatusClosed {
		t.Errorf("Expected status %s, got %s", StatusClosed, status)
	}
}

func TestChinaAShare_Name(t *testing.T) {
	china := NewChinaAShare()
	if china.Name() != "China A-Share" {
		t.Errorf("Expected name 'China A-Share', got '%s'", china.Name())
	}
}

func TestChinaAShare_HolidayClosed(t *testing.T) {
	china := NewChinaAShare()

	// Test Spring Festival 2026 - should be closed even during regular hours
	loc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		t.Fatalf("Failed to load timezone: %v", err)
	}
	springFestival := time.Date(2026, 2, 17, 10, 0, 0, 0, loc) // Feb 17, 2026 10:00 AM CST

	if china.IsOpen(springFestival) {
		t.Errorf("China A-Share should be closed on Spring Festival at %v", springFestival)
	}

	status := china.GetStatus(springFestival)
	if status != StatusClosed {
		t.Errorf("Expected status %s on Spring Festival, got %s", StatusClosed, status)
	}

	// Test National Day 2026
	nationalDay := time.Date(2026, 10, 1, 10, 0, 0, 0, loc)
	if china.IsOpen(nationalDay) {
		t.Errorf("China A-Share should be closed on National Day at %v", nationalDay)
	}

	status = china.GetStatus(nationalDay)
	if status != StatusClosed {
		t.Errorf("Expected status %s on National Day, got %s", StatusClosed, status)
	}
}
