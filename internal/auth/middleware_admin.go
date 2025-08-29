package auth

import (

	"github.com/gin-gonic/gin"
	"github.com/oktaharis/uji-teknis-godigi/internal/models"
	"github.com/oktaharis/uji-teknis-godigi/internal/response"
)

func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		v, ok := c.Get("user")
		if !ok {
			response.Unauthorized(c, "Unauthorized")
			c.Abort()
			return
		}
		user := v.(models.User)
		if user.Role != "admin" {
			response.Forbidden(c, "Admin only")
			c.Abort()
			return
		}
		c.Next()
	}
}
