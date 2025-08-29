package auth

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"

	"github.com/oktaharis/uji-teknis-godigi/internal/config"
	"github.com/oktaharis/uji-teknis-godigi/internal/models"
	"github.com/oktaharis/uji-teknis-godigi/internal/response"
)

func AuthRequired(cfg *config.Config, db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Unauthorized(c, "missing bearer token")
			c.Abort()
			return
		}

		tokenStr := parts[1]
		token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !token.Valid {
			response.Unauthorized(c, "invalid or expired token")
			c.Abort()
			return
		}

		claims, ok := token.Claims.(*Claims)
		if !ok {
			response.Unauthorized(c, "invalid claims")
			c.Abort()
			return
		}

		var user models.User
		if err := db.First(&user, claims.UserID).Error; err != nil {
			response.Unauthorized(c, "user not found")
			c.Abort()
			return
		}

		// Token-version check (logout invalidation)
		if user.TokenVersion != claims.TokenVersion {
			response.Unauthorized(c, "token revoked")
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Next()
	}
}
