package errors

// ServiceError: A custom error type that carries both HTTP status and API error message.
type ServiceError struct {
	StatusCode int
	ErrorCode  string
	Message    string
	Raw        error
}

func (e *ServiceError) Error() string {
	return e.Message
}
