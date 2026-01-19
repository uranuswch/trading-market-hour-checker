package marketchecker

import (
	"testing"
	"time"
)

func TestChinaAShare_MorningSession(t *testing.T) {
	china := NewChinaAShare()
	
	// Test morning session - Monday 10:00 AM CST
	loc, _ := time.LoadLocation("Asia/Shanghai")
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
	loc, _ := time.LoadLocation("Asia/Shanghai")
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
	loc, _ := time.LoadLocation("Asia/Shanghai")
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
	loc, _ := time.LoadLocation("Asia/Shanghai")
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
