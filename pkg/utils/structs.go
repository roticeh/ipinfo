package utils


type ErrorDetails struct {
	// Status    int         `json:"status"`
	ErrorType string      `json:"error_type"`
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details,omitempty"`
}

type APIResponse struct {
	Success   bool          `json:"success"`
	Status    int           `json:"status,omitempty"`
	Message   string        `json:"message,omitempty"`
	Data      interface{}   `json:"data,omitempty"`
	Error     *ErrorDetails `json:"error,omitempty"`
	Timestamp int64         `json:"timestamp"` // Unix Milliseconds
}
