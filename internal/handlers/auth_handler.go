package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	goMysql "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/oktaharis/uji-teknis-godigi/internal/auth"
	"github.com/oktaharis/uji-teknis-godigi/internal/config"
	"github.com/oktaharis/uji-teknis-godigi/internal/models"
)

type AuthHandler struct {
	Cfg *config.Config
	DB  *gorm.DB
}

func NewAuthHandler(cfg *config.Config, db *gorm.DB) *AuthHandler {
	return &AuthHandler{Cfg: cfg, DB: db}
}

type registerReq struct {
	Name     string `json:"name" binding:"required,min=2"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "VALIDATION", "message": "invalid payload"}})
		return
	}
	var exists int64
	h.DB.Model(&models.User{}).Where("email = ?", req.Email).Count(&exists)
	if exists > 0 {
		c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "EMAIL_EXISTS", "message": "email already registered"}})
		return
	}
	hash, _ := auth.HashPassword(req.Password)
	u := models.User{Name: req.Name, Email: req.Email, PasswordHash: hash, Role: "user"}
	if err := h.DB.Create(&u).Error; err != nil {
		var me *goMysql.MySQLError
		if errors.As(err, &me) && me.Number == 1062 { // duplicate key
			c.JSON(http.StatusConflict, gin.H{"error": gin.H{"code": "EMAIL_EXISTS", "message": "email already registered"}})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL", "message": "failed to create user"}})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"id": u.ID, "name": u.Name, "email": u.Email, "created_at": u.CreatedAt})
}

type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "VALIDATION", "message": "invalid payload"}})
		return
	}
	var u models.User
	if err := h.DB.Where("email = ?", req.Email).First(&u).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "INVALID_CREDENTIALS", "message": "email or password is incorrect"}})
		return
	}
	if !auth.CheckPassword(u.PasswordHash, req.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": gin.H{"code": "INVALID_CREDENTIALS", "message": "email or password is incorrect"}})
		return
	}
	tok, exp, err := auth.SignJWT(h.Cfg.JWTSecret, u.ID, u.TokenVersion, h.Cfg.JWTExpires)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL", "message": "failed to sign token"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tok, "expires_in": h.Cfg.JWTExpires, "expires_at": exp})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	u := c.MustGet("user").(models.User)
	// Increment token_version to revoke existing tokens
	h.DB.Model(&models.User{}).Where("id = ?", u.ID).Update("token_version", gorm.Expr("token_version + 1"))
	c.Status(http.StatusNoContent)
}

type forgotReq struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req forgotReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "VALIDATION", "message": "invalid payload"}})
		return
	}
	var u models.User
	if err := h.DB.Where("email = ?", req.Email).First(&u).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code": "USER_NOT_FOUND", "message": "email not found"}})
		return
	}
	token := uuid.NewString()
	pr := models.PasswordReset{
		UserID:    u.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}
	if err := h.DB.Create(&pr).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL", "message": "failed to create reset token"}})
		return
	}
	// NOTE: test mode — kirim token di response.
	c.JSON(http.StatusOK, gin.H{"message": "reset token generated (test mode)", "reset_token": token})
}

type resetReq struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req resetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code": "VALIDATION", "message": "invalid payload"}})
		return
	}
	var pr models.PasswordReset
	err := h.DB.Where("token = ?", req.Token).First(&pr).Error
	if err != nil || pr.UsedAt != nil || time.Now().After(pr.ExpiresAt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{"code": "BAD_TOKEN", "message": "reset token invalid or expired"}})
		return
	}
	hash, _ := auth.HashPassword(req.NewPassword)
	if err := h.DB.Model(&models.User{}).Where("id = ?", pr.UserID).Updates(map[string]any{
		"password_hash": hash,
		"token_version": gorm.Expr("token_version + 1"),
	}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code": "INTERNAL", "message": "failed to reset password"}})
		return
	}
	now := time.Now()
	h.DB.Model(&pr).Update("used_at", &now)
	c.JSON(http.StatusOK, gin.H{"message": "password updated"})
}
