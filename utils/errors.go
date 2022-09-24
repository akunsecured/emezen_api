package utils

import "errors"

var (
	ErrEmailIsAlreadyInUse      = errors.New("email is already in use")
	ErrInvalidCredentialsFormat = errors.New("bad credentials format")
	ErrNoAccountWithThisEmail   = errors.New("no account with this email address")
	ErrInvalidPassword          = errors.New("invalid password")
	ErrNotExists                = errors.New("account does not exist")
)
