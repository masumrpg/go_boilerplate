package utils

import "github.com/gofiber/fiber/v2"

// APIResponse represents a standardized API response
type APIResponse struct {
	Code    int         `json:"code"`
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    any         `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// SuccessResponse sends a successful response
func SuccessResponse(c *fiber.Ctx, statusCode int, data any, message string) error {
	return c.Status(statusCode).JSON(APIResponse{
		Code:    statusCode,
		Success: true,
		Message: message,
		Data:    data,
	})
}

// ErrorResponse sends an error response
func ErrorResponse(c *fiber.Ctx, statusCode int, message string, err error) error {
	errorMsg := message
	if err != nil {
		errorMsg = message + ": " + err.Error()
	}

	return c.Status(statusCode).JSON(APIResponse{
		Code:    statusCode,
		Success: false,
		Error:   errorMsg,
	})
}

// PagedResponse represents a paginated response
type PagedResponse struct {
	Code    int              `json:"code"`
	Success bool             `json:"success"`
	Data    any              `json:"data"`
	Message string           `json:"message,omitempty"`
	Meta    *PaginationMeta  `json:"meta,omitempty"`
}

// PaginationMeta contains pagination metadata
type PaginationMeta struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	Total     int `json:"total"`
	TotalPages int `json:"total_pages"`
}

// SuccessPagedResponse sends a successful paginated response
func SuccessPagedResponse(c *fiber.Ctx, statusCode int, data any, message string, meta *PaginationMeta) error {
	return c.Status(statusCode).JSON(PagedResponse{
		Code:    statusCode,
		Success: true,
		Data:    data,
		Message: message,
		Meta:    meta,
	})
}
