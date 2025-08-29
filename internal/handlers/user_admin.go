package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/oktaharis/uji-teknis-godigi/internal/auth"
	"github.com/oktaharis/uji-teknis-godigi/internal/models"
	"gorm.io/gorm"
)

type UserAdminHandler struct{ DB *gorm.DB }

func NewUserAdminHandler(db *gorm.DB) *UserAdminHandler { return &UserAdminHandler{DB: db} }

type adminCreateUserReq struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=admin user"`
}

func (h *UserAdminHandler) Create(c *gin.Context) {
	var req adminCreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "VALIDATION", "message": "invalid payload"}})
		return
	}
	hash, _ := auth.HashPassword(req.Password)
	u := models.User{Name: req.Name, Email: req.Email, PasswordHash: hash, Role: req.Role}
	if err := h.DB.Create(&u).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL", "message": "failed to create user"}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": u.ID, "name": u.Name, "email": u.Email, "role": u.Role, "created_at": u.CreatedAt})
}

func (h *UserAdminHandler) List(c *gin.Context) {
	var users []models.User
	q := h.DB.Model(&models.User{})
	if v := c.Query("role"); v != "" {
		q = q.Where("role = ?", v)
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if page < 1 { page = 1 }
	if limit < 1 || limit > 100 { limit = 10 }
	var total int64
	q.Count(&total)
	if err := q.Order("created_at DESC").Offset((page - 1) * limit).Limit(limit).Find(&users).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL", "message": "failed to fetch users"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"data": users,
		"pagination": gin.H{"page": page, "limit": limit, "total": total},
	})
}

func (h *UserAdminHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var u models.User
	if err := h.DB.First(&u, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code":"NOT_FOUND","message":"user not found"}})
		return
	}
	c.JSON(http.StatusOK, u)
}

type adminUpdateUserReq struct {
	Name  *string `json:"name"`
	Email *string `json:"email" binding:"omitempty,email"`
	Role  *string `json:"role" binding:"omitempty,oneof=user admin"`
}

func (h *UserAdminHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var u models.User
	if err := h.DB.First(&u, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code":"NOT_FOUND","message":"user not found"}})
		return
	}
	var req adminUpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code":"VALIDATION","message":"invalid payload"}})
		return
	}
	if req.Name != nil { u.Name = *req.Name }
	if req.Email != nil { u.Email = *req.Email }
	if req.Role != nil { u.Role = *req.Role }
	if err := h.DB.Save(&u).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code":"INTERNAL","message":"failed to update user"}})
		return
	}
	c.JSON(http.StatusOK, u)
}

func (h *UserAdminHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.DB.Delete(&models.User{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code":"INTERNAL","message":"failed to delete user"}})
		return
	}
	c.Status(http.StatusNoContent)
}