package main

import (
	"fmt"
	"time"

	checker "github.com/uranuswch/trading-market-hour-checker"
)

func main() {
	// Create a new market checker
	c := checker.NewChecker()

	// Example 1: Check if NASDAQ is currently open
	fmt.Println("=== Example 1: Check Current Market Status ===")
	now := time.Now()
	isOpen, err := c.IsOpen(checker.MarketNASDAQ, now)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("NASDAQ is currently open: %v\n", isOpen)
	}

	status, err := c.GetStatus(checker.MarketNASDAQ, now)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("NASDAQ current status: %s\n", status)
	}

	// Example 2: Check specific timestamp for multiple markets
	fmt.Println("\n=== Example 2: Check Specific Timestamp ===")

	// Create a specific time: Monday, Jan 19, 2026 at 10:00 AM ET
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		fmt.Printf("Error loading timezone: %v\n", err)
		return
	}
	specificTime := time.Date(2026, 1, 19, 10, 0, 0, 0, loc)

	fmt.Printf("Checking markets at: %s\n", specificTime.Format(time.RFC3339))

	markets := []checker.MarketType{
		checker.MarketNASDAQ,
		checker.MarketHKEX,
		checker.MarketChinaAShare,
	}

	for _, market := range markets {
		status, err := c.GetStatus(market, specificTime)
		if err != nil {
			fmt.Printf("Error checking %s: %v\n", market, err)
			continue
		}
		isOpen, _ := c.IsOpen(market, specificTime)
		fmt.Printf("  %s: Status=%s, IsOpen=%v\n", market, status, isOpen)
	}

	// Example 3: Check NASDAQ extended hours
	fmt.Println("\n=== Example 3: NASDAQ Extended Hours ===")

	testTimes := []struct {
		desc string
		time time.Time
	}{
		{
			desc: "Overnight (Sunday 9:00 PM ET)",
			time: time.Date(2026, 1, 18, 21, 0, 0, 0, loc),
		},
		{
			desc: "Premarket (Monday 7:00 AM ET)",
			time: time.Date(2026, 1, 19, 7, 0, 0, 0, loc),
		},
		{
			desc: "Regular (Monday 10:00 AM ET)",
			time: time.Date(2026, 1, 19, 10, 0, 0, 0, loc),
		},
		{
			desc: "Postmarket (Monday 5:00 PM ET)",
			time: time.Date(2026, 1, 19, 17, 0, 0, 0, loc),
		},
	}

	for _, tt := range testTimes {
		status, _ := c.GetStatus(checker.MarketNASDAQ, tt.time)
		fmt.Printf("  %s: %s\n", tt.desc, status)
	}

	// Example 4: Using the Market interface directly
	fmt.Println("\n=== Example 4: Using Market Interface Directly ===")

	nasdaqMarket, err := c.GetMarket(checker.MarketNASDAQ)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Market name: %s\n", nasdaqMarket.Name())

		// Check Hong Kong market
		hkLoc, err := time.LoadLocation("Asia/Hong_Kong")
		if err != nil {
			fmt.Printf("Error loading timezone: %v\n", err)
			return
		}
		hkTime := time.Date(2026, 1, 19, 10, 30, 0, 0, hkLoc)

		hkexMarket, _ := c.GetMarket(checker.MarketHKEX)
		fmt.Printf("HKEX at %s: %s\n",
			hkTime.Format("2006-01-02 15:04 MST"),
			hkexMarket.GetStatus(hkTime))
	}

	// Example 5: Check timezone conversion
	fmt.Println("\n=== Example 5: Timezone Conversion ===")

	// Create a time in UTC and check all markets
	utcTime := time.Date(2026, 1, 19, 15, 0, 0, 0, time.UTC)
	fmt.Printf("Checking markets at: %s UTC\n", utcTime.Format("2006-01-02 15:04:05"))

	for _, market := range markets {
		status, _ := c.GetStatus(market, utcTime)
		isOpen, _ := c.IsOpen(market, utcTime)
		fmt.Printf("  %s: Status=%s, IsOpen=%v\n", market, status, isOpen)
	}
}
