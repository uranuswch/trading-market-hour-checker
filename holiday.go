package marketchecker

import (
	"time"
)

// HolidayProvider defines an interface for providing market holidays
type HolidayProvider interface {
	// IsHoliday checks if the given date is a market holiday
	IsHoliday(t time.Time) bool
}

// StaticHolidayProvider provides a simple static list of holidays
type StaticHolidayProvider struct {
	holidays map[string]bool // key format: "YYYY-MM-DD"
}

// NewStaticHolidayProvider creates a new static holiday provider
func NewStaticHolidayProvider(holidays []time.Time) *StaticHolidayProvider {
	holidayMap := make(map[string]bool)
	for _, h := range holidays {
		// Normalize to midnight in the holiday's location to ensure date-only comparison
		normalized := time.Date(h.Year(), h.Month(), h.Day(), 0, 0, 0, 0, h.Location())
		key := normalized.Format("2006-01-02")
		holidayMap[key] = true
	}
	return &StaticHolidayProvider{
		holidays: holidayMap,
	}
}

// IsHoliday checks if the given date is a holiday
func (p *StaticHolidayProvider) IsHoliday(t time.Time) bool {
	// Normalize to date only in the time's location
	normalized := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	key := normalized.Format("2006-01-02")
	return p.holidays[key]
}

// NASDAQ holidays for 2025
var nasdaqHolidays2025 = []time.Time{
	time.Date(2025, 1, 1, 0, 0, 0, 0, mustLoadLocation("America/New_York")),   // New Year's Day
	time.Date(2025, 1, 20, 0, 0, 0, 0, mustLoadLocation("America/New_York")),  // Martin Luther King Jr. Day
	time.Date(2025, 2, 17, 0, 0, 0, 0, mustLoadLocation("America/New_York")),  // Presidents Day
	time.Date(2025, 4, 18, 0, 0, 0, 0, mustLoadLocation("America/New_York")),  // Good Friday
	time.Date(2025, 5, 26, 0, 0, 0, 0, mustLoadLocation("America/New_York")),  // Memorial Day
	time.Date(2025, 6, 19, 0, 0, 0, 0, mustLoadLocation("America/New_York")),  // Juneteenth
	time.Date(2025, 7, 4, 0, 0, 0, 0, mustLoadLocation("America/New_York")),   // Independence Day
	time.Date(2025, 9, 1, 0, 0, 0, 0, mustLoadLocation("America/New_York")),   // Labor Day
	time.Date(2025, 11, 27, 0, 0, 0, 0, mustLoadLocation("America/New_York")), // Thanksgiving
	time.Date(2025, 12, 25, 0, 0, 0, 0, mustLoadLocation("America/New_York")), // Christmas
}

// HKEX holidays for 2025
var hkexHolidays2025 = []time.Time{
	time.Date(2025, 1, 1, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),  // New Year's Day
	time.Date(2025, 1, 29, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")), // Lunar New Year
	time.Date(2025, 1, 30, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")), // Lunar New Year
	time.Date(2025, 1, 31, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")), // Lunar New Year
	time.Date(2025, 4, 4, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),  // Ching Ming Festival
	time.Date(2025, 4, 18, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")), // Good Friday
	time.Date(2025, 4, 19, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")), // Saturday following Good Friday
	time.Date(2025, 4, 21, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")), // Easter Monday
	time.Date(2025, 5, 1, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),  // Labour Day
	time.Date(2025, 5, 5, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),  // Buddha's Birthday
	time.Date(2025, 5, 31, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")), // Tuen Ng Festival (Dragon Boat)
	time.Date(2025, 7, 1, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),  // HKSAR Establishment Day
	time.Date(2025, 10, 1, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")), // National Day
	time.Date(2025, 10, 7, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")), // Day after Mid-Autumn Festival
	time.Date(2025, 10, 11, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")), // Chung Yeung Festival
	time.Date(2025, 12, 25, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")), // Christmas
	time.Date(2025, 12, 26, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")), // Boxing Day
}

// China A-Share holidays for 2025
var chinaAShareHolidays2025 = []time.Time{
	time.Date(2025, 1, 1, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")),  // New Year's Day
	time.Date(2025, 1, 28, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // Spring Festival (Chinese New Year)
	time.Date(2025, 1, 29, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // Spring Festival
	time.Date(2025, 1, 30, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // Spring Festival
	time.Date(2025, 1, 31, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // Spring Festival
	time.Date(2025, 2, 1, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")),  // Spring Festival
	time.Date(2025, 2, 2, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")),  // Spring Festival
	time.Date(2025, 2, 3, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")),  // Spring Festival
	time.Date(2025, 4, 4, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")),  // Qingming Festival
	time.Date(2025, 4, 5, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")),  // Qingming Festival
	time.Date(2025, 4, 6, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")),  // Qingming Festival
	time.Date(2025, 5, 1, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")),  // Labour Day
	time.Date(2025, 5, 2, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")),  // Labour Day
	time.Date(2025, 5, 3, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")),  // Labour Day
	time.Date(2025, 5, 31, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // Dragon Boat Festival
	time.Date(2025, 6, 1, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")),  // Dragon Boat Festival
	time.Date(2025, 6, 2, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")),  // Dragon Boat Festival
	time.Date(2025, 10, 1, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // National Day
	time.Date(2025, 10, 2, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // National Day
	time.Date(2025, 10, 3, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // National Day
	time.Date(2025, 10, 4, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // National Day
	time.Date(2025, 10, 5, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // National Day
	time.Date(2025, 10, 6, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // National Day
	time.Date(2025, 10, 7, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // National Day
}

// mustLoadLocation loads a timezone location or panics if it fails
func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		// For testing environments, fallback to UTC with offset
		switch name {
		case "America/New_York":
			return time.FixedZone("EST", -5*3600)
		case "Asia/Hong_Kong":
			return time.FixedZone("HKT", 8*3600)
		case "Asia/Shanghai":
			return time.FixedZone("CST", 8*3600)
		default:
			panic(err)
		}
	}
	return loc
}
