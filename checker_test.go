package marketchecker

import (
	"testing"
	"time"
)

func TestChecker_IsOpen(t *testing.T) {
	checker := NewChecker()
	
	// Test NASDAQ during regular hours
	loc, _ := time.LoadLocation("America/New_York")
	nasdaqTime := time.Date(2026, 1, 19, 10, 0, 0, 0, loc) // Monday 10 AM ET
	
	isOpen, err := checker.IsOpen(MarketNASDAQ, nasdaqTime)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !isOpen {
		t.Errorf("NASDAQ should be open at %v", nasdaqTime)
	}
}

func TestChecker_GetStatus(t *testing.T) {
	checker := NewChecker()
	
	// Test HKEX during morning session
	loc, _ := time.LoadLocation("Asia/Hong_Kong")
	hkexTime := time.Date(2026, 1, 19, 10, 0, 0, 0, loc) // Monday 10 AM HKT
	
	status, err := checker.GetStatus(MarketHKEX, hkexTime)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if status != StatusOpen {
		t.Errorf("Expected status %s, got %s", StatusOpen, status)
	}
}

func TestChecker_UnknownMarket(t *testing.T) {
	checker := NewChecker()
	
	_, err := checker.IsOpen("UnknownMarket", time.Now())
	if err == nil {
		t.Error("Expected error for unknown market type")
	}
	
	_, err = checker.GetStatus("UnknownMarket", time.Now())
	if err == nil {
		t.Error("Expected error for unknown market type")
	}
}

func TestChecker_GetMarket(t *testing.T) {
	checker := NewChecker()
	
	market, err := checker.GetMarket(MarketNASDAQ)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if market == nil {
		t.Error("Expected market instance, got nil")
	}
	if market.Name() != "NASDAQ" {
		t.Errorf("Expected market name 'NASDAQ', got '%s'", market.Name())
	}
}

func TestChecker_AddMarket(t *testing.T) {
	checker := NewChecker()
	
	// Add a custom market
	customMarket := NewNASDAQ()
	checker.AddMarket("CustomNASDAQ", customMarket)
	
	market, err := checker.GetMarket("CustomNASDAQ")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if market == nil {
		t.Error("Expected market instance, got nil")
	}
}

func TestChecker_MultipleMarkets(t *testing.T) {
	checker := NewChecker()
	
	// Create a time that's weekday in all timezones
	loc, _ := time.LoadLocation("America/New_York")
	testTime := time.Date(2026, 1, 19, 10, 0, 0, 0, loc) // Monday 10 AM ET
	
	// Test all markets
	markets := []MarketType{MarketNASDAQ, MarketHKEX, MarketChinaAShare}
	
	for _, marketType := range markets {
		_, err := checker.GetStatus(marketType, testTime)
		if err != nil {
			t.Errorf("Failed to get status for %s: %v", marketType, err)
		}
	}
}
