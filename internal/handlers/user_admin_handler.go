package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/oktaharis/uji-teknis-godigi/internal/auth"
	"github.com/oktaharis/uji-teknis-godigi/internal/models"
	"github.com/oktaharis/uji-teknis-godigi/internal/response"
)

type UserAdminHandler struct{ DB *gorm.DB }

func NewUserAdminHandler(db *gorm.DB) *UserAdminHandler { return &UserAdminHandler{DB: db} }

type adminCreateUserReq struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	Role     string `json:"role" binding:"required,oneof=user admin"`
}

func (h *UserAdminHandler) Create(c *gin.Context) {
	var req adminCreateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.UnprocessableEntity(c, "Validation Error", response.ExtractValidationErrors(err))
		return
	}
	hash, _ := auth.HashPassword(req.Password)
	u := models.User{Name: req.Name, Email: req.Email, PasswordHash: hash, Role: req.Role}
	if err := h.DB.Create(&u).Error; err != nil {
		response.Conflict(c, "Email already registered")
		return
	}
	response.Created(c, u, "User created")
}

func (h *UserAdminHandler) List(c *gin.Context) {
	var users []models.User
	q := h.DB.Model(&models.User{})
	if v := c.Query("q"); v != "" {
		q = q.Where("name LIKE ? OR email LIKE ?", "%"+v+"%", "%"+v+"%")
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	per, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if page < 1 {
		page = 1
	}
	if per < 1 {
		per = 10
	}
	var total int64
	q.Count(&total)
	if err := q.Order("created_at DESC").Limit(per).Offset((page-1)*per).Find(&users).Error; err != nil {
		response.InternalError(c, "Failed to list users")
		return
	}
	response.OK(c, response.List(users, page, per, total), "User list")
}

func (h *UserAdminHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var u models.User
	if err := h.DB.First(&u, id).Error; err != nil {
		response.NotFound(c, "User not found")
		return
	}
	response.OK(c, u, "User detail")
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
		response.NotFound(c, "User not found")
		return
	}
	var req adminUpdateUserReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.UnprocessableEntity(c, "Validation Error", response.ExtractValidationErrors(err))
		return
	}
	if req.Name != nil {
		u.Name = *req.Name
	}
	if req.Email != nil {
		u.Email = *req.Email
	}
	if req.Role != nil {
		u.Role = *req.Role
	}
	if err := h.DB.Save(&u).Error; err != nil {
		response.InternalError(c, "Failed to update user")
		return
	}
	response.OK(c, u, "User updated")
}

func (h *UserAdminHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.DB.Delete(&models.User{}, id).Error; err != nil {
		response.InternalError(c, "Failed to delete user")
		return
	}
	response.NoContent(c, "User deleted")
}
