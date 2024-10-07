package models

type User struct {
	ID             int64  `json:"id" db:"id"`
	HashedPassword string `json:"password" db:"password"`
	Email          string `json:"email" db:"email"`
	Name           string `json:"name" db:"name"`
}
