package utils

import (
	"time"
)

// GetNextMayFirst calculates the next May 1st
// If current month >= May, returns May 1st of next year
// Otherwise returns May 1st of current year
// This is because when memberships are created they expire on the
// very next May first
func GetNextMayFirst() (*time.Time, error) {
	now := time.Now()
	year := now.Year()

	if now.Month() >= time.May {
		year++
	}

	expiry := time.Date(year, time.May, 1, 0, 0, 0, 0, time.UTC)
	return &expiry, nil
}

// IsDateExpired compares expiryDate with today
// Returns true if today is after the expiry date, false if not expired or expiryDate is nil
func IsDateExpired(expiryDate *time.Time) (bool, error) {
	if expiryDate == nil {
		return false, nil
	}

	// Get start of day for both dates for proper date-only comparison
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

	expiry := expiryDate
	expiryStart := time.Date(expiry.Year(), expiry.Month(), expiry.Day(), 0, 0, 0, 0, time.UTC)

	return today.After(expiryStart), nil
}

func TruncateToDate(t time.Time) time.Time {
	return t.Truncate(24 * time.Hour)
}
