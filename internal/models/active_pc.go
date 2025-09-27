package models

import "time"

type ActivePC struct {
	GamerActivity
	FirstName      string     `db:"first_name" json:"first_name"`
	LastName       string     `db:"last_name" json:"last_name"`
	MembershipTier int        `db:"membership_tier" json:"membership_tier"`
	Banned         *bool      `db:"banned" json:"banned,omitempty"`
	Notes          *string    `db:"notes" json:"notes,omitempty"`
	CreatedAt      *time.Time `db:"created_at" json:"created_at,omitempty"`
}
