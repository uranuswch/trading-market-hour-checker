package marketchecker

import (
	"time"
)

// HKEX represents the Hong Kong Stock Exchange
type HKEX struct{}

var (
	// HKEX timezone (Hong Kong Time)
	hkexLocation *time.Location
)

func init() {
	var err error
	hkexLocation, err = time.LoadLocation("Asia/Hong_Kong")
	if err != nil {
		// Fallback to UTC+8 if location loading fails
		hkexLocation = time.FixedZone("HKT", 8*3600)
	}
}

// NewHKEX creates a new HKEX market instance
func NewHKEX() *HKEX {
	return &HKEX{}
}

// Name returns the market name
func (h *HKEX) Name() string {
	return "HKEX"
}

// IsOpen checks if HKEX is open for trading at the given time
func (h *HKEX) IsOpen(t time.Time) bool {
	status := h.GetStatus(t)
	return status == StatusOpen
}

// GetStatus returns the current market status at the given time
func (h *HKEX) GetStatus(t time.Time) MarketStatus {
	// Convert to Hong Kong Time
	localTime := t.In(hkexLocation)
	
	// Check if it's weekend
	if IsWeekend(localTime) {
		return StatusClosed
	}
	
	// HKEX trading hours (Hong Kong Time):
	// Morning session: 9:30 AM - 12:00 PM
	// Afternoon session: 1:00 PM - 4:00 PM
	
	morningSession := TimeRange{
		Start: 9*time.Hour + 30*time.Minute,
		End:   12 * time.Hour,
	}
	
	afternoonSession := TimeRange{
		Start: 13 * time.Hour,
		End:   16 * time.Hour,
	}
	
	if morningSession.IsWithin(localTime) || afternoonSession.IsWithin(localTime) {
		return StatusOpen
	}
	
	return StatusClosed
}
