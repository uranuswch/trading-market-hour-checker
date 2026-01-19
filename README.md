# Trading Market Hour Checker

A Go library for checking whether financial markets are open at a given timestamp. Supports multiple exchanges including NASDAQ (with extended hours), HKEX, and China A-Share markets.

## Features

- ✅ **NASDAQ**: Regular, premarket, postmarket, and overnight trading sessions
- ✅ **HKEX** (Hong Kong Exchange): Morning and afternoon trading sessions
- ✅ **China A-Share**: SSE (Shanghai Stock Exchange) and SZSE (Shenzhen Stock Exchange)
- ✅ Automatic timezone conversion for each market
- ✅ Weekend awareness
- ✅ Holiday awareness for exchange-specific holidays
- ✅ Simple and intuitive API

## Installation

```bash
go get github.com/uranuswch/trading-market-hour-checker
```

## Quick Start

```go
package main

import (
    "fmt"
    "time"
    
    checker "github.com/uranuswch/trading-market-hour-checker"
)

func main() {
    // Create a new market checker
    c := checker.NewChecker()
    
    // Check if NASDAQ is currently open
    now := time.Now()
    isOpen, err := c.IsOpen(checker.MarketNASDAQ, now)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    
    fmt.Printf("NASDAQ is open: %v\n", isOpen)
    
    // Get detailed status
    status, _ := c.GetStatus(checker.MarketNASDAQ, now)
    fmt.Printf("NASDAQ status: %s\n", status)
}
```

## Market Trading Hours

### NASDAQ
- **Overnight**: 8:00 PM - 4:00 AM ET (Sunday-Thursday nights)
- **Premarket**: 4:00 AM - 9:30 AM ET
- **Regular**: 9:30 AM - 4:00 PM ET
- **Postmarket**: 4:00 PM - 8:00 PM ET

### HKEX (Hong Kong Exchange)
- **Morning Session**: 9:30 AM - 12:00 PM HKT
- **Afternoon Session**: 1:00 PM - 4:00 PM HKT

### China A-Share (SSE/SZSE)
- **Morning Session**: 9:30 AM - 11:30 AM CST
- **Afternoon Session**: 1:00 PM - 3:00 PM CST

## Usage Examples

### Check Multiple Markets

```go
c := checker.NewChecker()

markets := []checker.MarketType{
    checker.MarketNASDAQ,
    checker.MarketHKEX,
    checker.MarketChinaAShare,
}

timestamp := time.Now()

for _, market := range markets {
    isOpen, _ := c.IsOpen(market, timestamp)
    status, _ := c.GetStatus(market, timestamp)
    fmt.Printf("%s: Open=%v, Status=%s\n", market, isOpen, status)
}
```

### Check NASDAQ Extended Hours

```go
c := checker.NewChecker()
loc, _ := time.LoadLocation("America/New_York")

// Check premarket
premarketTime := time.Date(2026, 1, 19, 7, 0, 0, 0, loc)
status, _ := c.GetStatus(checker.MarketNASDAQ, premarketTime)
fmt.Printf("Status: %s\n", status) // Output: premarket

// Check overnight
overnightTime := time.Date(2026, 1, 18, 21, 0, 0, 0, loc)
status, _ = c.GetStatus(checker.MarketNASDAQ, overnightTime)
fmt.Printf("Status: %s\n", status) // Output: overnight
```

### Using Market Interface Directly

```go
c := checker.NewChecker()

// Get the market instance
nasdaqMarket, err := c.GetMarket(checker.MarketNASDAQ)
if err != nil {
    fmt.Printf("Error: %v\n", err)
    return
}

// Use the market interface methods
fmt.Printf("Market: %s\n", nasdaqMarket.Name())
isOpen := nasdaqMarket.IsOpen(time.Now())
status := nasdaqMarket.GetStatus(time.Now())
```

### Timezone Conversion

The library automatically handles timezone conversions for each market:

```go
c := checker.NewChecker()

// Create a time in UTC
utcTime := time.Date(2026, 1, 19, 15, 0, 0, 0, time.UTC)

// Check all markets - they will be converted to their local timezones
for _, market := range []checker.MarketType{checker.MarketNASDAQ, checker.MarketHKEX, checker.MarketChinaAShare} {
    status, _ := c.GetStatus(market, utcTime)
    fmt.Printf("%s: %s\n", market, status)
}
```

## API Reference

### Types

#### MarketType
```go
const (
    MarketNASDAQ      MarketType = "NASDAQ"
    MarketHKEX        MarketType = "HKEX"
    MarketChinaAShare MarketType = "ChinaAShare"
)
```

#### MarketStatus
```go
const (
    StatusClosed     MarketStatus = "closed"
    StatusOpen       MarketStatus = "open"
    StatusPremarket  MarketStatus = "premarket"
    StatusPostmarket MarketStatus = "postmarket"
    StatusOvernight  MarketStatus = "overnight"
)
```

### Checker

#### NewChecker() *Checker
Creates a new Checker instance with all supported markets.

#### IsOpen(marketType MarketType, t time.Time) (bool, error)
Checks if the specified market is open for regular trading at the given time.

#### GetStatus(marketType MarketType, t time.Time) (MarketStatus, error)
Returns the detailed status of the specified market at the given time.

#### GetMarket(marketType MarketType) (Market, error)
Returns the Market interface for the specified market type.

#### AddMarket(marketType MarketType, market Market)
Allows adding a custom market implementation to the checker.

### Market Interface

```go
type Market interface {
    IsOpen(t time.Time) bool
    GetStatus(t time.Time) MarketStatus
    Name() string
}
```

## Running Tests

```bash
go test -v
```

## Running the Example

```bash
cd examples/basic
go run main.go
```

## Notes

- All markets are closed on weekends
- The library includes dynamic holiday calculation and calendars:
  - **NASDAQ**: US federal holidays are calculated dynamically for any year (New Year's Day, MLK Day, Presidents Day, Good Friday, Memorial Day, Juneteenth, Independence Day, Labor Day, Thanksgiving, Christmas). Observed holidays on weekends are automatically handled.
  - **HKEX**: Hong Kong market holidays for 2025-2026 (including Lunar New Year, Ching Ming Festival, Easter, Buddha's Birthday, Dragon Boat Festival, National Day, Mid-Autumn Festival, Chung Yeung Festival, Christmas). Lunar calendar holidays require manual specification.
  - **China A-Share**: Mainland China market holidays for 2025-2026 (including Spring Festival/Chinese New Year, Qingming Festival, Labour Day, Dragon Boat Festival, National Day Golden Week). Lunar calendar holidays require manual specification.
- **Holiday Limitations**: NASDAQ holidays are calculated dynamically for any year. HKEX and China A-Share holidays (which depend on lunar calendar) are pre-specified for 2025-2026. To extend support beyond 2026, add additional years to the holiday lists in `holiday.go`
- NASDAQ overnight trading requires the next trading day to be a weekday and not a holiday
- Timezone data is loaded from the system's timezone database

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.