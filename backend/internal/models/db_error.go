package models

import "errors"

var (
	ErrDuplicateLogin = errors.New("login already exists")
	ErrDuplicateEmail = errors.New("email already exists")
)
