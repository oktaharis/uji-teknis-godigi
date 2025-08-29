package handlers

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	goMysql "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/oktaharis/uji-teknis-godigi/internal/auth"
	"github.com/oktaharis/uji-teknis-godigi/internal/config"
	"github.com/oktaharis/uji-teknis-godigi/internal/models"
	"github.com/oktaharis/uji-teknis-godigi/internal/response"
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
		response.UnprocessableEntity(c, "Validation Error", response.ExtractValidationErrors(err))
		return
	}

	hash, _ := auth.HashPassword(req.Password)
	u := models.User{Name: req.Name, Email: req.Email, PasswordHash: hash, Role: "user"}

	if err := h.DB.Create(&u).Error; err != nil {
		var me *goMysql.MySQLError
		if errors.As(err, &me) && me.Number == 1062 {
			response.Conflict(c, "Email already registered")
			return
		}
		response.InternalError(c, "Failed to create user")
		return
	}

	response.Created(c, gin.H{
		"id": u.ID, "name": u.Name, "email": u.Email, "created_at": u.CreatedAt,
	}, "User registered")
}

type loginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.UnprocessableEntity(c, "Validation Error", response.ExtractValidationErrors(err))
		return
	}
	var u models.User
	if err := h.DB.Where("email = ?", req.Email).First(&u).Error; err != nil {
		response.Unauthorized(c, "Email or password is incorrect")
		return
	}
	if !auth.CheckPassword(u.PasswordHash, req.Password) {
		response.Unauthorized(c, "Email or password is incorrect")
		return
	}
	tok, exp, err := auth.SignJWT(h.Cfg.JWTSecret, u.ID, u.TokenVersion, h.Cfg.JWTExpires)
	if err != nil {
		response.InternalError(c, "Failed to sign token")
		return
	}
	response.OK(c, gin.H{
		"token": tok, "expires_in": h.Cfg.JWTExpires, "expires_at": exp,
	}, "Login success")
}

func (h *AuthHandler) Logout(c *gin.Context) {
	u := c.MustGet("user").(models.User)
	h.DB.Model(&models.User{}).Where("id = ?", u.ID).Update("token_version", gorm.Expr("token_version + 1"))
	response.NoContent(c, "Logged out")
}

type forgotReq struct {
	Email string `json:"email" binding:"required,email"`
}

func (h *AuthHandler) ForgotPassword(c *gin.Context) {
	var req forgotReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.UnprocessableEntity(c, "Validation Error", response.ExtractValidationErrors(err))
		return
	}
	var u models.User
	if err := h.DB.Where("email = ?", req.Email).First(&u).Error; err != nil {
		response.NotFound(c, "Email not found")
		return
	}
	token := uuid.NewString()
	pr := models.PasswordReset{
		UserID:    u.ID,
		Token:     token,
		ExpiresAt: time.Now().Add(30 * time.Minute),
	}
	if err := h.DB.Create(&pr).Error; err != nil {
		response.InternalError(c, "Failed to create reset token")
		return
	}
	response.OK(c, gin.H{
		"reset_token": token,
	}, "Reset token generated (test mode)")
}

type resetReq struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=6"`
}

func (h *AuthHandler) ResetPassword(c *gin.Context) {
	var req resetReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.UnprocessableEntity(c, "Validation Error", response.ExtractValidationErrors(err))
		return
	}
	var pr models.PasswordReset
	err := h.DB.Where("token = ?", req.Token).First(&pr).Error
	if err != nil || pr.UsedAt != nil || time.Now().After(pr.ExpiresAt) {
		response.BadRequest(c, "Reset token invalid or expired", nil)
		return
	}
	hash, _ := auth.HashPassword(req.NewPassword)
	if err := h.DB.Model(&models.User{}).Where("id = ?", pr.UserID).Updates(map[string]any{
		"password_hash": hash,
		"token_version": gorm.Expr("token_version + 1"),
	}).Error; err != nil {
		response.InternalError(c, "Failed to reset password")
		return
	}
	now := time.Now()
	h.DB.Model(&pr).Update("used_at", &now)
	response.OK(c, nil, "Password updated")
}
