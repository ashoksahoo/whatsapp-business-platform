package utils

import (
	"net/http"

	"github.com/ashoksahoo/whatsapp-business-platform/pkg/errors"
	"github.com/gin-gonic/gin"
)

// SuccessResponse represents a successful API response
type SuccessResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data"`
}

// ListResponse represents a list API response with pagination
type ListResponse struct {
	Success    bool               `json:"success"`
	Data       interface{}        `json:"data"`
	Pagination PaginationResponse `json:"pagination"`
}

// ErrorResponse represents an error API response
type ErrorResponse struct {
	Success bool                   `json:"success"`
	Error   *errors.AppError       `json:"error"`
}

// SuccessJSON sends a successful JSON response
func SuccessJSON(c *gin.Context, status int, data interface{}) {
	c.JSON(status, SuccessResponse{
		Success: true,
		Data:    data,
	})
}

// ErrorJSON sends an error JSON response
func ErrorJSON(c *gin.Context, err *errors.AppError) {
	c.JSON(err.StatusCode, ErrorResponse{
		Success: false,
		Error:   err,
	})
}

// ListJSON sends a list response with pagination
func ListJSON(c *gin.Context, data interface{}, pagination *Pagination) {
	c.JSON(http.StatusOK, ListResponse{
		Success:    true,
		Data:       data,
		Pagination: pagination.ToResponse(),
	})
}

// CreatedJSON sends a 201 Created response
func CreatedJSON(c *gin.Context, data interface{}) {
	SuccessJSON(c, http.StatusCreated, data)
}

// NoContentJSON sends a 204 No Content response
func NoContentJSON(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// BadRequestJSON sends a 400 Bad Request error
func BadRequestJSON(c *gin.Context, message string) {
	ErrorJSON(c, errors.NewBadRequest(message))
}

// UnauthorizedJSON sends a 401 Unauthorized error
func UnauthorizedJSON(c *gin.Context, message string) {
	ErrorJSON(c, errors.NewUnauthorized(message))
}

// NotFoundJSON sends a 404 Not Found error
func NotFoundJSON(c *gin.Context, resource, id string) {
	ErrorJSON(c, errors.NewNotFound(resource, id))
}

// InternalErrorJSON sends a 500 Internal Server Error
func InternalErrorJSON(c *gin.Context, err error) {
	ErrorJSON(c, errors.NewInternalError(err))
}
