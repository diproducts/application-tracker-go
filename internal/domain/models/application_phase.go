package models

import "time"

type ApplicationPhase struct {
	ID            int64     `json:"id" db:"id"`
	Name          string    `json:"name" db:"name"`
	Date          time.Time `json:"date" db:"date"`
	Created       time.Time `json:"created" db:"created"`
	Notes         string    `json:"notes" db:"notes"`
	ApplicationID int64     `json:"application_id" db:"application_id"`
}
