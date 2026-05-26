package errors

const (
	// ErrInternalServerError
	// Old v1: server/internal-server-error
	// Old v2: INTERNAL_SERVER_ERROR
	ErrInternalServerError = "internal_server_error"

	// ErrTimeoutExceeded
	// Old v1: TIMEOUT_EXCEEDED
	ErrTimeoutExceeded = "timeout_exceeded"

	// ErrBadRequest
	// Old v1: request/bad-request
	// Old v2: BAD_REQUEST
	ErrBadRequest = "bad_request"

	// ErrBadGateway
	ErrBadGateWay = "bad_gateway"

	// ErrUserNotFound
	// Old v1: auth/user-not-found
	// Old v2: USER_NOT_FOUND
	Err404NotFound = "404_not_found"

	// ErrMalformedRequest
	// Old v1: request/invalid-json
	// Old v2: INVALID_JSON_FORMAT
	ErrMalformedRequest = "invalid_json_format"

	// ErrRateLimitExceeded
	// Old v1: request/rate_limit_exceeded
	// Old v2: RATE_LIMIT_EXCEEDED
	ErrRateLimitExceeded = "rate_limit_exceeded"

	// ErrTooManyAttempts
	ErrTooManyAttempts = "too_many_attempts"

	// ErrValidationFailed
	// Old v1: request/validation-failed
	// Old v2: VALIDATION_FAILED
	ErrValidationFailed = "validation_failed"
)

const (

	// ErrForbidden
	// Old v1: auth/forbidden
	// Old v2: ACCESS_DENIED
	ErrForbidden = "access_denied"

	// ErrIpBlocked
	// Old v1: security/ip-blocked
	// Old v2: IP_BLOCKED
	ErrIpBlocked = "ip_blocked"
)
