package models

import (
	"testing"
	"time"

	"github.com/ubcesports/echo-base/internal/utils"
)

func TestNewMembershipTier(t *testing.T) {
	tests := []struct {
		name        string
		tierNumber  int
		wantName    string
		wantErr     bool
	}{
		{"Tier 0", 0, "No Membership", false},
		{"Tier 1", 1, "Tier 1", false},
		{"Tier 2", 2, "Tier 2", false},
		{"Tier 3", 3, "Premier", false},
		{"Invalid tier -1", -1, "", true},
		{"Invalid tier 4", 4, "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tier, err := NewMembershipTier(tt.tierNumber)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMembershipTier() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tier.GetName() != tt.wantName {
				t.Errorf("GetName() = %v, want %v", tier.GetName(), tt.wantName)
			}
		})
	}
}

func TestTierSessionDurations(t *testing.T) {
	tests := []struct {
		tierNumber int
		wantMs     int64
	}{
		{0, 0},
		{1, 60 * 60 * 1000},
		{2, 2 * 60 * 60 * 1000},
		{3, 5 * 60 * 60 * 1000},
	}

	for _, tt := range tests {
		tier, err := NewMembershipTier(tt.tierNumber)
		if err != nil {
			t.Fatalf("NewMembershipTier(%d) error = %v", tt.tierNumber, err)
		}
		if got := tier.GetSessionDurationMs(); got != tt.wantMs {
			t.Errorf("Tier %d GetSessionDurationMs() = %v, want %v", tt.tierNumber, got, tt.wantMs)
		}
	}
}

func TestTierDailyLimits(t *testing.T) {
	tests := []struct {
		tierNumber int
		wantLimit  bool
	}{
		{0, false},
		{1, true},
		{2, false},
		{3, false},
	}

	for _, tt := range tests {
		tier, err := NewMembershipTier(tt.tierNumber)
		if err != nil {
			t.Fatalf("NewMembershipTier(%d) error = %v", tt.tierNumber, err)
		}
		if got := tier.HasDailyLimit(); got != tt.wantLimit {
			t.Errorf("Tier %d HasDailyLimit() = %v, want %v", tt.tierNumber, got, tt.wantLimit)
		}
	}
}

func TestTierExpiryDates(t *testing.T) {
	loc, err := utils.GetPacificLocation()
	if err != nil {
		t.Fatalf("GetPacificLocation() error = %v", err)
	}

	now := time.Now().In(loc)
	currentYear := now.Year()
	expectedYear := currentYear
	if now.Month() >= time.May {
		expectedYear = currentYear + 1
	}

	tests := []struct {
		tierNumber   int
		wantExpiry   bool
		expectedYear int
	}{
		{0, false, 0},
		{1, true, expectedYear},
		{2, true, expectedYear},
		{3, true, expectedYear},
	}

	for _, tt := range tests {
		tier, err := NewMembershipTier(tt.tierNumber)
		if err != nil {
			t.Fatalf("NewMembershipTier(%d) error = %v", tt.tierNumber, err)
		}

		expiry, err := tier.GetExpiryDate()
		if err != nil {
			t.Fatalf("Tier %d GetExpiryDate() error = %v", tt.tierNumber, err)
		}

		if !tt.wantExpiry {
			if expiry != nil {
				t.Errorf("Tier %d GetExpiryDate() = %v, want nil", tt.tierNumber, expiry)
			}
		} else {
			if expiry == nil {
				t.Errorf("Tier %d GetExpiryDate() = nil, want non-nil", tt.tierNumber)
				continue
			}
			if expiry.Year() != tt.expectedYear {
				t.Errorf("Tier %d GetExpiryDate() year = %d, want %d", tt.tierNumber, expiry.Year(), tt.expectedYear)
			}
			if expiry.Month() != time.May || expiry.Day() != 1 {
				t.Errorf("Tier %d GetExpiryDate() = %v, want May 1", tt.tierNumber, expiry)
			}
		}
	}
}

func TestTierIsExpired(t *testing.T) {
	loc, err := utils.GetPacificLocation()
	if err != nil {
		t.Fatalf("GetPacificLocation() error = %v", err)
	}

	yesterday := time.Now().In(loc).AddDate(0, 0, -1)
	tomorrow := time.Now().In(loc).AddDate(0, 0, 1)
	lastYear := time.Now().In(loc).AddDate(-1, 0, 0)

	tests := []struct {
		name       string
		tierNumber int
		expiryDate *time.Time
		wantExpired bool
	}{
		{"Tier 0 with nil expiry", 0, nil, false},
		{"Tier 1 with nil expiry", 1, nil, false},
		{"Tier 1 expired yesterday", 1, &yesterday, true},
		{"Tier 1 expires tomorrow", 1, &tomorrow, false},
		{"Tier 2 expired last year", 2, &lastYear, true},
		{"Tier 3 expires tomorrow", 3, &tomorrow, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tier, err := NewMembershipTier(tt.tierNumber)
			if err != nil {
				t.Fatalf("NewMembershipTier(%d) error = %v", tt.tierNumber, err)
			}

			expired, err := tier.IsExpired(tt.expiryDate)
			if err != nil {
				t.Fatalf("IsExpired() error = %v", err)
			}

			if expired != tt.wantExpired {
				t.Errorf("IsExpired() = %v, want %v", expired, tt.wantExpired)
			}
		})
	}
}

func TestTierExpiryBoundary(t *testing.T) {
	loc, err := utils.GetPacificLocation()
	if err != nil {
		t.Fatalf("GetPacificLocation() error = %v", err)
	}

	mayFirst2024 := time.Date(2024, time.May, 1, 0, 0, 0, 0, loc)
	april30_2024 := time.Date(2024, time.April, 30, 23, 59, 59, 0, loc)
	may2_2024 := time.Date(2024, time.May, 2, 0, 0, 1, 0, loc)

	tier, err := NewMembershipTier(1)
	if err != nil {
		t.Fatalf("NewMembershipTier(1) error = %v", err)
	}

	tests := []struct {
		name        string
		expiryDate  time.Time
		checkDate   time.Time
		wantExpired bool
	}{
		{"Check on expiry date", mayFirst2024, mayFirst2024, false},
		{"Check day before expiry", mayFirst2024, april30_2024, false},
		{"Check day after expiry", mayFirst2024, may2_2024, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expired, err := tier.IsExpired(&tt.expiryDate)
			if err != nil {
				t.Fatalf("IsExpired() error = %v", err)
			}

			if expired != tt.wantExpired {
				t.Errorf("IsExpired(%v) at %v = %v, want %v",
					tt.expiryDate.Format("2006-01-02"),
					tt.checkDate.Format("2006-01-02"),
					expired, tt.wantExpired)
			}
		})
	}
}
