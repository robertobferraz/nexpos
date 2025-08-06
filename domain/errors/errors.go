package errors

type ErrorCode string

// Predefined error codes
const (
	ErrInvalidInput      ErrorCode = "INVALID_INPUT"
	ErrUserAlreadyExists ErrorCode = "USER_ALREADY_EXISTS"
	ErrInternalServer    ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrUnauthorized      ErrorCode = "UNAUTHORIZED"
	ErrForbidden         ErrorCode = "FORBIDDEN"
)
