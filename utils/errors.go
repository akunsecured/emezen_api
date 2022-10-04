package utils

import "errors"

var (
	ErrEmailIsAlreadyInUse             = errors.New("email is already in use")
	ErrInvalidCredentialsFormat        = errors.New("bad credentials format")
	ErrNoAccountWithThisEmail          = errors.New("no account with this email address")
	ErrInvalidPassword                 = errors.New("invalid password")
	ErrNotExists                       = errors.New("account does not exist")
	ErrInvalidTokenFormat              = errors.New("invalid token format in header")
	ErrMissingAuthToken                = errors.New("missing token")
	ErrExpiredAuthToken                = errors.New("expired token")
	ErrTokenParseError                 = errors.New("parse error")
	ErrInsertedIDIsNotObjectID         = errors.New("the variable InsertedID is not in the correct type")
	ErrNoMatchedDocumentFoundForDelete = errors.New("no matched document found for delete")
	ErrUnimplementedMethod             = errors.New("unimplemented method")
)
