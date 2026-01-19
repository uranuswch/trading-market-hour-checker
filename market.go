package marketchecker

import (
	"time"
)

// MarketStatus represents the trading status of a market
type MarketStatus string

const (
	// StatusClosed indicates the market is closed
	StatusClosed MarketStatus = "closed"
	// StatusOpen indicates the market is open for regular trading
	StatusOpen MarketStatus = "open"
	// StatusPremarket indicates the market is in premarket session
	StatusPremarket MarketStatus = "premarket"
	// StatusPostmarket indicates the market is in postmarket session
	StatusPostmarket MarketStatus = "postmarket"
	// StatusOvernight indicates the market is in overnight trading session
	StatusOvernight MarketStatus = "overnight"
)

// Market represents a financial market
type Market interface {
	// IsOpen checks if the market is open at the given time
	IsOpen(t time.Time) bool
	// GetStatus returns the current market status at the given time
	GetStatus(t time.Time) MarketStatus
	// Name returns the market name
	Name() string
}

// TimeRange represents a trading session time range
type TimeRange struct {
	Start time.Duration // Duration from midnight
	End   time.Duration // Duration from midnight
}

// IsWithin checks if the given time is within the range
func (tr TimeRange) IsWithin(t time.Time) bool {
	loc := t.Location()
	midnight := time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, loc)
	sinceStart := t.Sub(midnight)
	
	// Handle ranges that cross midnight
	if tr.End < tr.Start {
		return sinceStart >= tr.Start || sinceStart < tr.End
	}
	return sinceStart >= tr.Start && sinceStart < tr.End
}

// IsWeekend checks if the given time is on a weekend
func IsWeekend(t time.Time) bool {
	weekday := t.Weekday()
	return weekday == time.Saturday || weekday == time.Sunday
}
