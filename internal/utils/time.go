package utils

import (
	"fmt"
	"time"
)

func GetPacificLocation() (*time.Location, error) {
	loc, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return nil, fmt.Errorf("failed to load Pacific timezone: %w", err)
	}
	return loc, nil
}

func NowInPacific() (time.Time, error) {
	loc, err := GetPacificLocation()
	if err != nil {
		return time.Time{}, err
	}
	return time.Now().In(loc), nil
}

// GetNextMayFirst calculates the next May 1st in Pacific timezone
// If current month >= May, returns May 1st of next year
// Otherwise returns May 1st of current year
// This is because when memberships are created they expire on the 
// very next May first
func GetNextMayFirst() (*time.Time, error) {
	loc, err := GetPacificLocation()
	if err != nil {
		return nil, err
	}

	now := time.Now().In(loc)
	year := now.Year()

	if now.Month() >= time.May {
		year++
	}

	expiry := time.Date(year, time.May, 1, 0, 0, 0, 0, loc)
	return &expiry, nil
}

// IsDateExpired compares expiryDate with today at day granularity in Pacific timezone
// Returns true if today is after the expiry date, false if not expired or expiryDate is nil
func IsDateExpired(expiryDate *time.Time) (bool, error) {
	if expiryDate == nil {
		return false, nil
	}

	loc, err := GetPacificLocation()
	if err != nil {
		return false, err
	}

	// Truncate to start of day for date-only comparison
	today := time.Now().In(loc).Truncate(24 * time.Hour)
	expiry := expiryDate.In(loc).Truncate(24 * time.Hour)

	return today.After(expiry), nil
}

func TruncateToDate(t time.Time) time.Time {
	return t.Truncate(24 * time.Hour)
}
