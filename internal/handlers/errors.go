package handlers

import "fmt"

type HTTPError struct {
	Status  int
	Message string
	Err     error
}

func (e *HTTPError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

// Helper constructor
func NewHTTPError(status int, message string, err ...error) *HTTPError {
	var internal error
	if len(err) > 0 {
		internal = err[0]
	}

	return &HTTPError{
		Status:  status,
		Message: message,
		Err:     internal,
	}
}
