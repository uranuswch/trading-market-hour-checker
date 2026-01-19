package marketchecker

import (
	"fmt"
	"time"
)

// MarketType represents the type of market
type MarketType string

const (
	// MarketNASDAQ represents NASDAQ
	MarketNASDAQ MarketType = "NASDAQ"
	// MarketHKEX represents Hong Kong Exchange
	MarketHKEX MarketType = "HKEX"
	// MarketChinaAShare represents China A-Share market
	MarketChinaAShare MarketType = "ChinaAShare"
)

// Checker provides a convenient interface to check market hours
type Checker struct {
	markets map[MarketType]Market
}

// NewChecker creates a new Checker instance
func NewChecker() *Checker {
	return &Checker{
		markets: map[MarketType]Market{
			MarketNASDAQ:      NewNASDAQ(),
			MarketHKEX:        NewHKEX(),
			MarketChinaAShare: NewChinaAShare(),
		},
	}
}

// IsOpen checks if the specified market is open at the given time
func (c *Checker) IsOpen(marketType MarketType, t time.Time) (bool, error) {
	market, ok := c.markets[marketType]
	if !ok {
		return false, fmt.Errorf("unknown market type: %s", marketType)
	}
	return market.IsOpen(t), nil
}

// GetStatus returns the status of the specified market at the given time
func (c *Checker) GetStatus(marketType MarketType, t time.Time) (MarketStatus, error) {
	market, ok := c.markets[marketType]
	if !ok {
		return StatusClosed, fmt.Errorf("unknown market type: %s", marketType)
	}
	return market.GetStatus(t), nil
}

// GetMarket returns the Market interface for the specified market type
func (c *Checker) GetMarket(marketType MarketType) (Market, error) {
	market, ok := c.markets[marketType]
	if !ok {
		return nil, fmt.Errorf("unknown market type: %s", marketType)
	}
	return market, nil
}

// AddMarket allows adding a custom market to the checker
func (c *Checker) AddMarket(marketType MarketType, market Market) {
	c.markets[marketType] = market
}
