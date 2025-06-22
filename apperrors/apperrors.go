package apperrors

import "errors"

// Error variables
var (
	ErrNotFound             = errors.New("not found")
	ErrUnauthorized         = errors.New("unauthorized")
	ErrWrongParam           = errors.New("wrong parameter")
	ErrServerError          = errors.New("internal server error")
	ErrParsingFailed        = errors.New("failed to parse input")
	ErrInvalidEmail         = errors.New("invalid email")
	ErrUnauthorizedToEdit   = errors.New("not authorized to edit this object")
	ErrAuthCheckFailed      = errors.New("failed to check if the user is authorized to perform this action")
	ErrTypeConversionFailed = errors.New("couldn't convert the user data into the proper type - wrong input?")
)

// Custom error type (optional)
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return e.Field + ": " + e.Message
}
