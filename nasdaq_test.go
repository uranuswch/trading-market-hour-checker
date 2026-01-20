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

func TestNASDAQ_OvernightHolidayHandling(t *testing.T) {
	nasdaq := NewNASDAQ()
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		t.Fatalf("Failed to load timezone: %v", err)
	}

	// Test Case 1: MLK Day 2026 (Monday Jan 19) at 8:00 PM
	// The overnight session from 8 PM Jan 19 to midnight should be OPEN
	// because the next day (Jan 20, Tuesday) is a normal trading day
	mlkDayEvening := time.Date(2026, 1, 19, 20, 0, 0, 0, loc) // 8:00 PM on holiday
	status := nasdaq.GetStatus(mlkDayEvening)
	if status != StatusOvernight {
		t.Errorf("Expected status %s at MLK Day 8:00 PM (next day is normal), got %s", StatusOvernight, status)
	}

	// Test Case 2: MLK Day 2026 (Monday Jan 19) at 11:00 PM
	// Should still be overnight (next day is normal trading day)
	mlkDayLateNight := time.Date(2026, 1, 19, 23, 0, 0, 0, loc) // 11:00 PM on holiday
	status = nasdaq.GetStatus(mlkDayLateNight)
	if status != StatusOvernight {
		t.Errorf("Expected status %s at MLK Day 11:00 PM (next day is normal), got %s", StatusOvernight, status)
	}

	// Test Case 3: Tuesday Jan 20, 2026 at 2:00 AM (after MLK Day)
	// The overnight session from midnight to 4 AM should be OPEN
	// because the current day (Jan 20, Tuesday) is a normal trading day
	dayAfterMLKEarlyMorning := time.Date(2026, 1, 20, 2, 0, 0, 0, loc) // 2:00 AM on normal day
	status = nasdaq.GetStatus(dayAfterMLKEarlyMorning)
	if status != StatusOvernight {
		t.Errorf("Expected status %s at 2:00 AM on Jan 20 (normal day), got %s", StatusOvernight, status)
	}

	// Test Case 4: Friday before MLK weekend at 8:00 PM (Jan 16, 2026)
	// Next day is Saturday, should be CLOSED
	fridayBeforeMLKEvening := time.Date(2026, 1, 16, 20, 0, 0, 0, loc) // 8:00 PM Friday
	status = nasdaq.GetStatus(fridayBeforeMLKEvening)
	if status != StatusClosed {
		t.Errorf("Expected status %s at Friday 8:00 PM (next day is weekend), got %s", StatusClosed, status)
	}

	// Test Case 5: Christmas Eve 2026 (Thursday Dec 24) at 8:00 PM
	// Next day is Christmas (Friday), should be CLOSED
	christmasEveEvening := time.Date(2026, 12, 24, 20, 0, 0, 0, loc) // 8:00 PM before holiday
	status = nasdaq.GetStatus(christmasEveEvening)
	if status != StatusClosed {
		t.Errorf("Expected status %s at Christmas Eve 8:00 PM (next day is holiday), got %s", StatusClosed, status)
	}

	// Test Case 6: Christmas 2026 (Friday Dec 25) at 2:00 AM
	// Current day is a holiday, should be CLOSED
	christmasMorning := time.Date(2026, 12, 25, 2, 0, 0, 0, loc) // 2:00 AM on holiday
	status = nasdaq.GetStatus(christmasMorning)
	if status != StatusClosed {
		t.Errorf("Expected status %s at Christmas 2:00 AM (current day is holiday), got %s", StatusClosed, status)
	}

	// Test Case 7: MLK Day (Monday Jan 19, 2026) at 10:00 AM
	// Regular hours on a holiday should remain closed
	mlkDay := time.Date(2026, 1, 19, 10, 0, 0, 0, loc)
	status = nasdaq.GetStatus(mlkDay)
	if status != StatusClosed {
		t.Errorf("Expected status %s at MLK Day 10:00 AM (holiday), got %s", StatusClosed, status)
	}
}
