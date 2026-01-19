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

// DynamicHolidayProvider dynamically calculates US federal holidays
type DynamicHolidayProvider struct {
	location *time.Location
}

// NewDynamicHolidayProvider creates a new dynamic holiday provider for US federal holidays
func NewDynamicHolidayProvider(location *time.Location) *DynamicHolidayProvider {
	return &DynamicHolidayProvider{
		location: location,
	}
}

// IsHoliday checks if the given date is a US federal holiday
func (p *DynamicHolidayProvider) IsHoliday(t time.Time) bool {
	// Convert to provider's timezone
	t = t.In(p.location)
	year := t.Year()
	month := t.Month()
	day := t.Day()
	
	// New Year's Day (January 1, or observed on nearby weekday if on weekend)
	if isObservedHoliday(year, time.January, 1, t, p.location) {
		return true
	}
	
	// Martin Luther King Jr. Day (3rd Monday in January)
	if month == time.January && day == nthWeekdayOfMonth(year, time.January, time.Monday, 3) {
		return true
	}
	
	// Presidents Day (3rd Monday in February)
	if month == time.February && day == nthWeekdayOfMonth(year, time.February, time.Monday, 3) {
		return true
	}
	
	// Good Friday (Friday before Easter - calculated)
	if isGoodFriday(year, month, day) {
		return true
	}
	
	// Memorial Day (last Monday in May)
	if month == time.May && day == lastWeekdayOfMonth(year, time.May, time.Monday) {
		return true
	}
	
	// Juneteenth (June 19, or observed on nearby weekday if on weekend)
	if isObservedHoliday(year, time.June, 19, t, p.location) {
		return true
	}
	
	// Independence Day (July 4, or observed on nearby weekday if on weekend)
	if isObservedHoliday(year, time.July, 4, t, p.location) {
		return true
	}
	
	// Labor Day (1st Monday in September)
	if month == time.September && day == nthWeekdayOfMonth(year, time.September, time.Monday, 1) {
		return true
	}
	
	// Thanksgiving (4th Thursday in November)
	if month == time.November && day == nthWeekdayOfMonth(year, time.November, time.Thursday, 4) {
		return true
	}
	
	// Christmas (December 25, or observed on nearby weekday if on weekend)
	if isObservedHoliday(year, time.December, 25, t, p.location) {
		return true
	}
	
	return false
}

// nthWeekdayOfMonth returns the day of month for the nth occurrence of a weekday
// Returns -1 if the nth occurrence doesn't exist in the month
func nthWeekdayOfMonth(year int, month time.Month, weekday time.Weekday, n int) int {
	// Start at the first day of the month
	t := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	
	// Find the first occurrence of the target weekday
	for t.Weekday() != weekday {
		t = t.AddDate(0, 0, 1)
	}
	
	// Add (n-1) weeks
	t = t.AddDate(0, 0, (n-1)*7)
	
	// Validate we're still in the target month
	if t.Month() != month {
		return -1
	}
	
	return t.Day()
}

// lastWeekdayOfMonth returns the day of month for the last occurrence of a weekday
func lastWeekdayOfMonth(year int, month time.Month, weekday time.Weekday) int {
	// Start at the last day of the month
	t := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)
	
	// Find the last occurrence of the target weekday
	for t.Weekday() != weekday {
		t = t.AddDate(0, 0, -1)
	}
	
	return t.Day()
}

// isObservedHoliday checks if a date or its observed substitute is a holiday
// If the holiday falls on a weekend, it's observed on the nearest weekday
func isObservedHoliday(year int, month time.Month, day int, t time.Time, loc *time.Location) bool {
	holiday := time.Date(year, month, day, 0, 0, 0, 0, loc)
	
	// Check if it's the actual holiday
	if t.Year() == holiday.Year() && t.Month() == holiday.Month() && t.Day() == holiday.Day() {
		return true
	}
	
	// If holiday is on Saturday, observed on Friday
	if holiday.Weekday() == time.Saturday {
		observed := holiday.AddDate(0, 0, -1)
		if t.Year() == observed.Year() && t.Month() == observed.Month() && t.Day() == observed.Day() {
			return true
		}
	}
	
	// If holiday is on Sunday, observed on Monday
	if holiday.Weekday() == time.Sunday {
		observed := holiday.AddDate(0, 0, 1)
		if t.Year() == observed.Year() && t.Month() == observed.Month() && t.Day() == observed.Day() {
			return true
		}
	}
	
	return false
}

// isGoodFriday checks if the date is Good Friday
func isGoodFriday(year int, month time.Month, day int) bool {
	easter := calculateEaster(year)
	goodFriday := easter.AddDate(0, 0, -2)
	return month == goodFriday.Month() && day == goodFriday.Day()
}

// calculateEaster calculates Easter Sunday using the Gregorian Easter Algorithm (also known as Anonymous Gregorian algorithm)
func calculateEaster(year int) time.Time {
	a := year % 19
	b := year / 100
	c := year % 100
	d := b / 4
	e := b % 4
	f := (b + 8) / 25
	g := (b - f + 1) / 3
	h := (19*a + b - d - g + 15) % 30
	i := c / 4
	k := c % 4
	l := (32 + 2*e + 2*i - h - k) % 7
	m := (a + 11*h + 22*l) / 451
	month := (h + l - 7*m + 114) / 31
	day := ((h + l - 7*m + 114) % 31) + 1
	
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

// HKEX holidays - includes both 2025 and 2026 (lunar calendar dates require manual specification)
var hkexHolidays = []time.Time{
	// 2025
	time.Date(2025, 1, 1, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),   // New Year's Day
	time.Date(2025, 1, 29, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),  // Lunar New Year
	time.Date(2025, 1, 30, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),  // Lunar New Year
	time.Date(2025, 1, 31, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),  // Lunar New Year
	time.Date(2025, 4, 4, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),   // Ching Ming Festival
	time.Date(2025, 4, 18, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),  // Good Friday
	time.Date(2025, 4, 19, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),  // Saturday following Good Friday
	time.Date(2025, 4, 21, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),  // Easter Monday
	time.Date(2025, 5, 1, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),   // Labour Day
	time.Date(2025, 5, 5, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),   // Buddha's Birthday
	time.Date(2025, 5, 31, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),  // Tuen Ng Festival (Dragon Boat)
	time.Date(2025, 7, 1, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),   // HKSAR Establishment Day
	time.Date(2025, 10, 1, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),  // National Day
	time.Date(2025, 10, 7, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),  // Day after Mid-Autumn Festival
	time.Date(2025, 10, 11, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")), // Chung Yeung Festival
	time.Date(2025, 12, 25, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")), // Christmas
	time.Date(2025, 12, 26, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")), // Boxing Day
	// 2026
	time.Date(2026, 1, 1, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),   // New Year's Day
	time.Date(2026, 2, 17, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),  // Lunar New Year
	time.Date(2026, 2, 18, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),  // Lunar New Year
	time.Date(2026, 2, 19, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),  // Lunar New Year
	time.Date(2026, 4, 3, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),   // Good Friday
	time.Date(2026, 4, 4, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),   // Saturday following Good Friday
	time.Date(2026, 4, 5, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),   // Ching Ming Festival
	time.Date(2026, 4, 6, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),   // Easter Monday
	time.Date(2026, 5, 1, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),   // Labour Day
	time.Date(2026, 5, 24, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),  // Buddha's Birthday
	time.Date(2026, 6, 19, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),  // Tuen Ng Festival (Dragon Boat)
	time.Date(2026, 7, 1, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),   // HKSAR Establishment Day
	time.Date(2026, 9, 26, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),  // Day after Mid-Autumn Festival
	time.Date(2026, 10, 1, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")),  // National Day
	time.Date(2026, 10, 21, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")), // Chung Yeung Festival
	time.Date(2026, 12, 25, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")), // Christmas
	time.Date(2026, 12, 26, 0, 0, 0, 0, mustLoadLocation("Asia/Hong_Kong")), // Boxing Day
}

// China A-Share holidays - includes both 2025 and 2026 (lunar calendar dates require manual specification)
var chinaAShareHolidays = []time.Time{
	// 2025
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
	// 2026
	time.Date(2026, 1, 1, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")),  // New Year's Day
	time.Date(2026, 1, 2, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")),  // New Year's Day (observed)
	time.Date(2026, 2, 16, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // Spring Festival Eve
	time.Date(2026, 2, 17, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // Spring Festival (Chinese New Year)
	time.Date(2026, 2, 18, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // Spring Festival
	time.Date(2026, 2, 19, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // Spring Festival
	time.Date(2026, 2, 20, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // Spring Festival
	time.Date(2026, 2, 21, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // Spring Festival
	time.Date(2026, 2, 22, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // Spring Festival
	time.Date(2026, 4, 4, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")),  // Qingming Festival
	time.Date(2026, 4, 5, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")),  // Qingming Festival
	time.Date(2026, 4, 6, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")),  // Qingming Festival
	time.Date(2026, 5, 1, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")),  // Labour Day
	time.Date(2026, 5, 2, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")),  // Labour Day
	time.Date(2026, 5, 3, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")),  // Labour Day
	time.Date(2026, 6, 18, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // Dragon Boat Festival
	time.Date(2026, 6, 19, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // Dragon Boat Festival
	time.Date(2026, 6, 20, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // Dragon Boat Festival
	time.Date(2026, 9, 25, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // Mid-Autumn Festival
	time.Date(2026, 9, 26, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // Mid-Autumn Festival
	time.Date(2026, 9, 27, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // Mid-Autumn Festival
	time.Date(2026, 10, 1, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // National Day
	time.Date(2026, 10, 2, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // National Day
	time.Date(2026, 10, 3, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // National Day
	time.Date(2026, 10, 4, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // National Day
	time.Date(2026, 10, 5, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // National Day
	time.Date(2026, 10, 6, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // National Day
	time.Date(2026, 10, 7, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // National Day
	time.Date(2026, 10, 8, 0, 0, 0, 0, mustLoadLocation("Asia/Shanghai")), // National Day
}

// mustLoadLocation loads a timezone location or panics if it fails
func mustLoadLocation(name string) *time.Location {
	loc, err := time.LoadLocation(name)
	if err != nil {
		// Fallback for environments without timezone data.
		// Note: These fixed offsets don't account for DST, but since holidays
		// are compared by date only (not time), this is acceptable for the fallback case.
		switch name {
		case "America/New_York":
			return time.FixedZone("ET", -5*3600) // Approximation without DST
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
