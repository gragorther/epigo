package apperrors

import "errors"

// Error variables
var (
	ErrNotFound               = errors.New("not found")
	ErrUserNotFound           = errors.New("user not found")
	ErrUnauthorized           = errors.New("unauthorized")
	ErrWrongParam             = errors.New("wrong parameter")
	ErrServerError            = errors.New("internal server error")
	ErrParsingFailed          = errors.New("failed to parse input")
	ErrInvalidEmail           = errors.New("invalid email")
	ErrUnauthorizedToEdit     = errors.New("not authorized to edit this object")
	ErrAuthCheckFailed        = errors.New("failed to check if the user is authorized to perform this action")
	ErrTypeConversionFailed   = errors.New("couldn't convert the user data into the proper type - wrong input?")
	ErrNoUsers                = errors.New("no users")
	ErrCreationOfObjectFailed = errors.New("failed to put this in the database")
	ErrDatabaseFetchFailed    = errors.New("failed to fetch from the database")
	ErrNoGroups               = errors.New("no groups")
	ErrDeleteFailed           = errors.New("failed to delete")
	ErrUserAlreadyExists      = errors.New("user already exists")
	ErrHashingFailed          = errors.New("failed to hash password")
	ErrHashCheckFailed        = errors.New("failed to check password hash")
	ErrInvalidPassword        = errors.New("invalid password")
	ErrJWTCreationError       = errors.New("failed to generate JWT token")
	ErrMissingAuthHeader      = errors.New("auth header is missing")
	ErrInvalidAuthTokenFormat = errors.New("invalid auth token format")
	ErrInvalidToken           = errors.New("invalid or expired token")
	ErrExpiredToken           = errors.New("expired token")
	ErrFailedToGetUserID      = errors.New("failed to get user ID")
)
