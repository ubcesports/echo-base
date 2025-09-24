package models

type UpdateGamerActivityRequest struct {
	PCNumber int    `db:"pc_number" json:"pc_number"`
	ExecName string `db:"exec_name" json:"exec_name"`
}
