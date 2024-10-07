package usecase

type UserCreateInput struct {
	email    string
	password string
	name     string
}

type UserLogin struct {
	email    string
	password string
}
