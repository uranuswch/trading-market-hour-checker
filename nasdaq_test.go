package marketchecker

import (
	"testing"
	"time"
)

func TestNASDAQ_RegularHours(t *testing.T) {
	nasdaq := NewNASDAQ()

	// Test regular trading hours - Tuesday 10:00 AM ET
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("Failed to load timezone: %v", err)
	}
	regularTime := time.Date(2026, 1, 20, 10, 0, 0, 0, loc) // Tuesday (not a holiday)

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

	// Test premarket - Tuesday 7:00 AM ET
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("Failed to load timezone: %v", err)
	}
	premarketTime := time.Date(2026, 1, 20, 7, 0, 0, 0, loc) // Tuesday

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

	// Test postmarket - Tuesday 5:00 PM ET
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("Failed to load timezone: %v", err)
	}
	postmarketTime := time.Date(2026, 1, 20, 17, 0, 0, 0, loc) // Tuesday

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

	// Test overnight - Sunday 11:00 PM ET (Sunday night to Monday)
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("Failed to load timezone: %v", err)
	}
	// Use Sunday Jan 25 night to Monday Jan 26 (regular trading day)
	sundayNight := time.Date(2026, 1, 25, 23, 0, 0, 0, loc) // Sunday night to Monday

	if nasdaq.IsOpen(sundayNight) {
		t.Errorf("NASDAQ regular hours should not be open at %v", sundayNight)
	}

	status := nasdaq.GetStatus(sundayNight)
	if status != StatusOvernight {
		t.Errorf("Expected status %s, got %s", StatusOvernight, status)
	}
}

func TestNASDAQ_WeekendClosed(t *testing.T) {
	nasdaq := NewNASDAQ()

	// Test weekend - Saturday 10:00 AM ET
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("Failed to load timezone: %v", err)
	}
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

func TestNASDAQ_HolidayClosed(t *testing.T) {
	nasdaq := NewNASDAQ()

	// Test Christmas 2026 - should be closed even during regular hours
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("Failed to load timezone: %v", err)
	}
	christmasTime := time.Date(2026, 12, 25, 10, 0, 0, 0, loc) // Christmas 10:00 AM ET

	if nasdaq.IsOpen(christmasTime) {
		t.Errorf("NASDAQ should be closed on Christmas at %v", christmasTime)
	}

	status := nasdaq.GetStatus(christmasTime)
	if status != StatusClosed {
		t.Errorf("Expected status %s on Christmas, got %s", StatusClosed, status)
	}

	// Test Independence Day 2026 - July 4 is a Saturday, so observed on Friday July 3
	independenceDayObserved := time.Date(2026, 7, 3, 10, 0, 0, 0, loc)
	if nasdaq.IsOpen(independenceDayObserved) {
		t.Errorf("NASDAQ should be closed on Independence Day (observed) at %v", independenceDayObserved)
	}

	status = nasdaq.GetStatus(independenceDayObserved)
	if status != StatusClosed {
		t.Errorf("Expected status %s on Independence Day (observed), got %s", StatusClosed, status)
	}
	
	// Test MLK Day 2026 - January 19
	mlkDay := time.Date(2026, 1, 19, 10, 0, 0, 0, loc)
	if nasdaq.IsOpen(mlkDay) {
		t.Errorf("NASDAQ should be closed on MLK Day at %v", mlkDay)
	}

	status = nasdaq.GetStatus(mlkDay)
	if status != StatusClosed {
		t.Errorf("Expected status %s on MLK Day, got %s", StatusClosed, status)
	}
}
