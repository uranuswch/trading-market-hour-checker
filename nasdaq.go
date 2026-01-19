package marketchecker

import (
	"time"
)

// NASDAQ represents the NASDAQ stock exchange
type NASDAQ struct{}

var (
	// NASDAQ timezone (Eastern Time)
	nasdaqLocation *time.Location
)

func init() {
	var err error
	nasdaqLocation, err = time.LoadLocation("America/New_York")
	if err != nil {
		// Fallback to UTC if location loading fails
		nasdaqLocation = time.UTC
	}
}

// NewNASDAQ creates a new NASDAQ market instance
func NewNASDAQ() *NASDAQ {
	return &NASDAQ{}
}

// Name returns the market name
func (n *NASDAQ) Name() string {
	return "NASDAQ"
}

// IsOpen checks if NASDAQ is open for regular trading at the given time
func (n *NASDAQ) IsOpen(t time.Time) bool {
	status := n.GetStatus(t)
	return status == StatusOpen
}

// GetStatus returns the current market status at the given time
func (n *NASDAQ) GetStatus(t time.Time) MarketStatus {
	// Convert to Eastern Time
	localTime := t.In(nasdaqLocation)

	// Define trading session times (in Eastern Time)
	// Overnight: 8:00 PM - 4:00 AM (previous day 20:00 to current day 04:00)
	overnight := TimeRange{
		Start: 20 * time.Hour,
		End:   4 * time.Hour,
	}

	// Premarket: 4:00 AM - 9:30 AM
	premarket := TimeRange{
		Start: 4 * time.Hour,
		End:   9*time.Hour + 30*time.Minute,
	}

	// Regular trading: 9:30 AM - 4:00 PM
	regular := TimeRange{
		Start: 9*time.Hour + 30*time.Minute,
		End:   16 * time.Hour,
	}

	// Postmarket: 4:00 PM - 8:00 PM
	postmarket := TimeRange{
		Start: 16 * time.Hour,
		End:   20 * time.Hour,
	}

	// Check if it's in overnight session first (before weekend check)
	// because Sunday 8PM-midnight transitions to Monday's overnight session
	if overnight.IsWithin(localTime) {
		hour := localTime.Hour()
		if hour >= 20 { // 8 PM or later
			// Check if next day is a weekday
			nextDay := localTime.AddDate(0, 0, 1)
			if !IsWeekend(nextDay) {
				return StatusOvernight
			}
		} else if hour < 4 { // Before 4 AM
			// Current day should be a weekday
			if !IsWeekend(localTime) {
				return StatusOvernight
			}
		}
	}

	// Check if it's weekend
	if IsWeekend(localTime) {
		return StatusClosed
	}

	// Check each session
	if regular.IsWithin(localTime) {
		return StatusOpen
	}

	if premarket.IsWithin(localTime) {
		return StatusPremarket
	}

	if postmarket.IsWithin(localTime) {
		return StatusPostmarket
	}

	return StatusClosed
}
