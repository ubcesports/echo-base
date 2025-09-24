package models

import "time"

type MembershipResult struct {
	MembershipExpiryDate time.Time `db:"membership_expiry_date" json:"membership_expiry_date"`
	MembershipTier       int       `db:"membership_tier" json:"membership_tier"`
}
