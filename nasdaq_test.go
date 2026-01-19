package marketchecker

import (
	"testing"
	"time"
)

func TestNASDAQ_RegularHours(t *testing.T) {
	nasdaq := NewNASDAQ()
	
	// Test regular trading hours - Monday 10:00 AM ET
	loc, _ := time.LoadLocation("America/New_York")
	regularTime := time.Date(2026, 1, 19, 10, 0, 0, 0, loc) // Monday
	
	if !nasdaq.IsOpen(regularTime) {
		t.Errorf("NASDAQ should be open at %v", regularTime)
	}
	
	status := nasdaq.GetStatus(regularTime)
	if status != StatusOpen {
		t.Errorf("Expected status %s, got %s", StatusOpen, status)
	}
}

func TestNASDAQ_Premarket(t *testing.T) {
	nasdaq := NewNASDAQ()
	
	// Test premarket - Monday 7:00 AM ET
	loc, _ := time.LoadLocation("America/New_York")
	premarketTime := time.Date(2026, 1, 19, 7, 0, 0, 0, loc) // Monday
	
	if nasdaq.IsOpen(premarketTime) {
		t.Errorf("NASDAQ regular hours should not be open at %v", premarketTime)
	}
	
	status := nasdaq.GetStatus(premarketTime)
	if status != StatusPremarket {
		t.Errorf("Expected status %s, got %s", StatusPremarket, status)
	}
}

func TestNASDAQ_Postmarket(t *testing.T) {
	nasdaq := NewNASDAQ()
	
	// Test postmarket - Monday 5:00 PM ET
	loc, _ := time.LoadLocation("America/New_York")
	postmarketTime := time.Date(2026, 1, 19, 17, 0, 0, 0, loc) // Monday
	
	if nasdaq.IsOpen(postmarketTime) {
		t.Errorf("NASDAQ regular hours should not be open at %v", postmarketTime)
	}
	
	status := nasdaq.GetStatus(postmarketTime)
	if status != StatusPostmarket {
		t.Errorf("Expected status %s, got %s", StatusPostmarket, status)
	}
}

func TestNASDAQ_Overnight(t *testing.T) {
	nasdaq := NewNASDAQ()
	
	// Test overnight - Monday 11:00 PM ET (Sunday night to Monday)
	loc, _ := time.LoadLocation("America/New_York")
	overnightTime := time.Date(2026, 1, 18, 23, 0, 0, 0, loc) // Sunday night
	
	if nasdaq.IsOpen(overnightTime) {
		t.Errorf("NASDAQ regular hours should not be open at %v", overnightTime)
	}
	
	status := nasdaq.GetStatus(overnightTime)
	if status != StatusOvernight {
		t.Errorf("Expected status %s, got %s", StatusOvernight, status)
	}
}

func TestNASDAQ_WeekendClosed(t *testing.T) {
	nasdaq := NewNASDAQ()
	
	// Test weekend - Saturday 10:00 AM ET
	loc, _ := time.LoadLocation("America/New_York")
	weekendTime := time.Date(2026, 1, 17, 10, 0, 0, 0, loc) // Saturday
	
	if nasdaq.IsOpen(weekendTime) {
		t.Errorf("NASDAQ should be closed on weekend at %v", weekendTime)
	}
	
	status := nasdaq.GetStatus(weekendTime)
	if status != StatusClosed {
		t.Errorf("Expected status %s, got %s", StatusClosed, status)
	}
}

func TestNASDAQ_Name(t *testing.T) {
	nasdaq := NewNASDAQ()
	if nasdaq.Name() != "NASDAQ" {
		t.Errorf("Expected name 'NASDAQ', got '%s'", nasdaq.Name())
	}
}
