package domain

import "errors"

var (
	ErrFieldsRequired  = errors.New("all required fields must have values")
	ErrExists          = errors.New("already exists")
	ErrBadCredentials  = errors.New("email or password is incorrect")
	ErrNotFound        = errors.New("not found")
	ErrNotExists       = errors.New("doesn't exist")
	ErrTokenNotCreated = errors.New("token didn't created")
)
