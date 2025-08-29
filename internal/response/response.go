package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func JSON(c *gin.Context, status int, success bool, message string, data interface{}) {
	c.JSON(status, APIResponse{
		Success: success,
		Message: message,
		Data:    data,
	})
}

func OK(c *gin.Context, data interface{}, message string) {
	JSON(c, http.StatusOK, true, orDefault(message, "Success"), data)
}

func Created(c *gin.Context, data interface{}, message string) {
	JSON(c, http.StatusCreated, true, orDefault(message, "Success"), data)
}

func NoContent(c *gin.Context, message string) {
	JSON(c, http.StatusNoContent, true, orDefault(message, "Success"), nil)
}

func BadRequest(c *gin.Context, message string, details interface{}) {
	JSON(c, http.StatusBadRequest, false, orDefault(message, "Bad Request"), details)
}

func Unauthorized(c *gin.Context, message string) {
	JSON(c, http.StatusUnauthorized, false, orDefault(message, "Unauthorized"), nil)
}

func Forbidden(c *gin.Context, message string) {
	JSON(c, http.StatusForbidden, false, orDefault(message, "Forbidden"), nil)
}

func NotFound(c *gin.Context, message string) {
	JSON(c, http.StatusNotFound, false, orDefault(message, "Not Found"), nil)
}

func Conflict(c *gin.Context, message string) {
	JSON(c, http.StatusConflict, false, orDefault(message, "Conflict"), nil)
}

func UnprocessableEntity(c *gin.Context, message string, details interface{}) {
	JSON(c, http.StatusUnprocessableEntity, false, orDefault(message, "Validation Error"), details)
}

func InternalError(c *gin.Context, message string) {
	JSON(c, http.StatusInternalServerError, false, orDefault(message, "Internal Server Error"), nil)
}

func orDefault(got, def string) string {
	if got == "" {
		return def
	}
	return got
}
