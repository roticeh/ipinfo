package utils

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

// responseError writes a standardized JSON error response to the client.
// It optionally returns the ApiError struct for internal error handling flows.
// ResponseError: Centralized high-performance error responder with internal logging.
func ResponseError(c *fiber.Ctx, statusCode int, errorCode, message string, details ...interface{}) error {

	errResponse := &ErrorDetails{
		// Status:    statusCode,
		ErrorType: http.StatusText(statusCode),
		Code:      errorCode,
		Message:   message,
	}

	if len(details) > 0 {
		errResponse.Details = details[0]
	}

	// return c.Status(status).JSON(apiErr)
	return c.Status(statusCode).JSON(APIResponse{
		Success:   false,
		Status:    statusCode,
		Error:     errResponse,
		Timestamp: time.Now().UnixMilli(),
	})
}

// ResponseSuccess: Standardized successful response helper.
func ResponseSuccess(c *fiber.Ctx, statusCode int, message string, data interface{}) error {

	return c.Status(statusCode).JSON(APIResponse{
		Success:   true,
		Status:    statusCode,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
	})
}

func Response(c *fiber.Ctx, statusCode int, message string, data interface{}, success bool) error {

	return c.Status(statusCode).JSON(APIResponse{
		Success:   success,
		Status:    statusCode,
		Message:   message,
		Data:      data,
		Timestamp: time.Now().UnixMilli(),
	})
}
