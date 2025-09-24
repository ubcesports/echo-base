package models

type GamerActivityWithName struct {
	GamerActivity
	FirstName string `db:"first_name" json:"first_name"`
	LastName  string `db:"last_name" json:"last_name"`
}
