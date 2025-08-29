package middleware

import (
	"github.com/gin-gonic/gin"

	"github.com/oktaharis/uji-teknis-godigi/internal/response"
)

// NotFoundHandler returns a middleware that handles 404 errors with JSON response
func NotFoundHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		response.NotFound(c, "Route not found")
	}
}