package domain

import "errors"

var (
	ErrPersonNotFound      = errors.New("person not found")
	ErrPersonAlreadyExists = errors.New("person already exists")
	ErrNameSurnameRequired = errors.New("name and surname are required")

	ErrInvalidUUID        = errors.New("invalid uuid format")
	ErrInvalidRequestBody = errors.New("invalid request body")

	ErrEmptyName      = errors.New("name cannot be empty")
	ErrNon200Response = errors.New("API returned non-200 status code")

	ErrInvalidConfiguration      = errors.New("invalid database configuration: required fields missing")
	ErrMigrationPathNotSpecified = errors.New("migration path not specified")
)
