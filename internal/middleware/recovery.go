package middleware

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/oktaharis/uji-teknis-godigi/internal/response"
)

// JSONRecovery returns a middleware that recovers from panics and returns JSON error response
func JSONRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		if err, ok := recovered.(string); ok {
			log.Printf("Panic recovered: %s", err)
		} else {
			log.Printf("Panic recovered: %v", recovered)
		}
		response.InternalError(c, "Internal server error")
	})
}