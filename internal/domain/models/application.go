package models

import "time"

type Application struct {
	ID             uint      `json:"id" db:"id"`
	CompanyName    string    `json:"company_name" db:"company_name"`
	Position       string    `json:"position" db:"position"`
	Url            string    `json:"url" db:"url"`
	JobDescription string    `json:"job_description" db:"job_description"`
	Contacts       string    `json:"contacts" db:"contacts"`
	Cv             string    `json:"cv" db:"cv"`
	CoverLetter    string    `json:"cover_letter" db:"cover_letter"`
	OfferedSalary  uint      `json:"offered_salary" db:"offered_salary"`
	Created        time.Time `json:"created" db:"created"`
	LastModified   time.Time `json:"last_modified" db:"last_modified"`
	OwnerID        uint      `json:"owner_id" db:"owner_id"`
}
