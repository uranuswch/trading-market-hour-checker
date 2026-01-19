package marketchecker

import (
	"time"
)

// ChinaAShare represents the China A-Share market (SSE and SZSE)
// Both Shanghai Stock Exchange and Shenzhen Stock Exchange have the same trading hours
type ChinaAShare struct{
	holidayProvider HolidayProvider
}

var (
	// China A-Share timezone (China Standard Time)
	chinaLocation *time.Location
)

func init() {
	var err error
	chinaLocation, err = time.LoadLocation("Asia/Shanghai")
	if err != nil {
		// Fallback to UTC+8 if location loading fails
		chinaLocation = time.FixedZone("CST", 8*3600)
	}
}

// NewChinaAShare creates a new China A-Share market instance
func NewChinaAShare() *ChinaAShare {
	return &ChinaAShare{
		holidayProvider: NewStaticHolidayProvider(chinaAShareHolidays2025),
	}
}

// Name returns the market name
func (c *ChinaAShare) Name() string {
	return "China A-Share"
}

// IsOpen checks if China A-Share market is open for trading at the given time
func (c *ChinaAShare) IsOpen(t time.Time) bool {
	status := c.GetStatus(t)
	return status == StatusOpen
}

// GetStatus returns the current market status at the given time
func (c *ChinaAShare) GetStatus(t time.Time) MarketStatus {
	// Convert to China Standard Time
	localTime := t.In(chinaLocation)

	// Check if it's a holiday first
	if c.holidayProvider != nil && c.holidayProvider.IsHoliday(localTime) {
		return StatusClosed
	}

	// Check if it's weekend
	if IsWeekend(localTime) {
		return StatusClosed
	}

	// China A-Share trading hours (CST):
	// Morning session: 9:30 AM - 11:30 AM
	// Afternoon session: 1:00 PM - 3:00 PM

	morningSession := TimeRange{
		Start: 9*time.Hour + 30*time.Minute,
		End:   11*time.Hour + 30*time.Minute,
	}

	afternoonSession := TimeRange{
		Start: 13 * time.Hour,
		End:   15 * time.Hour,
	}

	if morningSession.IsWithin(localTime) || afternoonSession.IsWithin(localTime) {
		return StatusOpen
	}

	return StatusClosed
}
