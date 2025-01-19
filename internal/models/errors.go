package models

import "github.com/pkg/errors"

// Ошибки для событий
var (
	ErrEventNotFound       = errors.New("event not found")
	ErrEventAlreadyDeleted = errors.New("event is already deleted")
)

// Ошибки для публикаций
var (
	ErrPublicationNotFound      = errors.New("publication not found")
	ErrPublicationAlreadyDone   = errors.New("publication is already processed")
	ErrPublicationDraftNotFound = errors.New("draft publication not found")
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)
