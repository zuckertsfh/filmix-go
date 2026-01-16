package services

import "errors"

var (
	ErrEmailAlreadyRegistered = errors.New("email already registered")
	ErrRoleNotFound           = errors.New("default role not found")
	ErrInternal               = errors.New("internal server error")
)
