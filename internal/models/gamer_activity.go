package models

import (
	"time"

	"github.com/google/uuid"
)

type GamerActivity struct {
	ID            uuid.UUID  `db:"id" json:"id"`
	StudentNumber string     `db:"student_number" json:"student_number"`
	PCNumber      int        `db:"pc_number" json:"pc_number"`
	Game          string     `db:"game" json:"game"`
	StartedAt     time.Time  `db:"started_at" json:"started_at"`
	EndedAt       *time.Time `db:"ended_at" json:"ended_at,omitempty"`
	ExecName      *string    `db:"exec_name" json:"exec_name,omitempty"`
}
