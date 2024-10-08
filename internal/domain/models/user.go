package models

type User struct {
	ID             int64
	HashedPassword string
	Email          string
	Name           string
}
