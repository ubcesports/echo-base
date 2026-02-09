package models

import "time"

type GamerProfile struct {
	ID                   string     `json:"id"`
	StudentNumber        string     `json:"student_number"`
	FirstName            string     `json:"first_name"`
	LastName             string     `json:"last_name"`
	MembershipTier       int        `json:"membership_tier"`
	Banned               *bool      `json:"banned,omitempty"`
	Notes                *string    `json:"notes,omitempty"`
	CreatedAt            time.Time  `json:"created_at"`
	MembershipExpiryDate *time.Time `json:"membership_expiry_date,omitempty"`
}

type GamerActivity struct {
	ID            string     `json:"id"`
	StudentNumber string     `json:"student_number"`
	PCNumber      int        `json:"pc_number"`
	Game          string     `json:"game"`
	StartedAt     time.Time  `json:"started_at"`
	EndedAt       *time.Time `json:"ended_at,omitempty"`
	ExecName      *string    `json:"exec_name,omitempty"`
	FirstName     *string    `json:"first_name,omitempty"`
	LastName      *string    `json:"last_name,omitempty"`
}

type CreateGamerProfileRequest struct {
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	StudentNumber  string  `json:"student_number"`
	MembershipTier int     `json:"membership_tier"`
	Banned         *bool   `json:"banned,omitempty"`
	Notes          *string `json:"notes,omitempty"`
}

type CreateActivityRequest struct {
	StudentNumber string `json:"student_number"`
	PCNumber      int    `json:"pc_number"`
	Game          string `json:"game"`
}

type UpdateActivityRequest struct {
	PCNumber int    `json:"pc_number"`
	ExecName string `json:"exec_name"`
}
