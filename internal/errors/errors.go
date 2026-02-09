package errors

import "fmt"

type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("%s: %s", e.Field, e.Message)
	}
	return e.Message
}

func NewValidationError(field, message string) error {
	return &ValidationError{Field: field, Message: message}
}

type NotFoundError struct {
	Resource string
	ID       string
}

func (e *NotFoundError) Error() string {
	if e.ID != "" {
		return fmt.Sprintf("%s %s not found", e.Resource, e.ID)
	}
	return fmt.Sprintf("%s not found", e.Resource)
}

func NewNotFoundError(resource, id string) error {
	return &NotFoundError{Resource: resource, ID: id}
}

type ForbiddenError struct {
	Message string
}

func (e *ForbiddenError) Error() string {
	return e.Message
}

func NewForbiddenError(message string) error {
	return &ForbiddenError{Message: message}
}
