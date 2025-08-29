package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/oktaharis/uji-teknis-godigi/internal/models"
)

type UserHandler struct{}

func NewUserHandler() *UserHandler { return &UserHandler{} }

func (h *UserHandler) Me(c *gin.Context) {
	u := c.MustGet("user").(models.User)
	c.JSON(http.StatusOK, gin.H{
		"id":         u.ID,
		"name":       u.Name,
		"email":      u.Email,
		"role":       u.Role,
		"created_at": u.CreatedAt,
	})
}
