package storage

import "errors"

var (
	ErrTokenAlreadyBlacklisted = errors.New("token already blacklisted")
	ErrUserAlreadyExists       = errors.New("user already exists")
	ErrUserNotFound            = errors.New("user not found")
)
