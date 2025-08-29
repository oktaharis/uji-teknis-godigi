package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/oktaharis/uji-teknis-godigi/internal/models"
)

type ProjectHandler struct{ DB *gorm.DB }

func NewProjectHandler(db *gorm.DB) *ProjectHandler { return &ProjectHandler{DB: db} }

type projectPayload struct {
	Name        string  `json:"name" binding:"required,min=2"`
	Description *string `json:"description"`
	Status      *string `json:"status" binding:"omitempty,oneof=planned in_progress on_hold done canceled"`
	StartDate   *string `json:"start_date"` // "YYYY-MM-DD"
	EndDate     *string `json:"end_date"`
	OwnerUserID *uint   `json:"owner_user_id"`
}

func parseDatePtr(s *string) *time.Time {
	if s == nil || *s == "" { return nil }
	t, err := time.Parse("2006-01-02", *s)
	if err != nil { return nil }
	return &t
}

func (h *ProjectHandler) Create(c *gin.Context) {
	var p projectPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code":"VALIDATION","message":"invalid payload"}})
		return
	}
	status := "planned"
	if p.Status != nil { status = *p.Status }
	proj := models.Project{
		Name:        p.Name,
		Description: p.Description,
		Status:      status,
		StartDate:   parseDatePtr(p.StartDate),
		EndDate:     parseDatePtr(p.EndDate),
		OwnerUserID: p.OwnerUserID,
	}
	if err := h.DB.Create(&proj).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code":"INTERNAL","message":"failed to create project"}})
		return
	}
	c.JSON(http.StatusCreated, proj)
}

func (h *ProjectHandler) List(c *gin.Context) {
	var items []models.Project
	q := h.DB.Model(&models.Project{})
	if v := c.Query("status"); v != "" {
		q = q.Where("status = ?", v)
	}
	if v := c.Query("q"); v != "" {
		q = q.Where("name LIKE ? OR description LIKE ?", "%"+v+"%", "%"+v+"%")
	}
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	per, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if page < 1 { page = 1 }
	if per  < 1 { per  = 10 }
	var total int64
	q.Count(&total)
	if err := q.Order("created_at DESC").Limit(per).Offset((page-1)*per).Find(&items).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code":"INTERNAL","message":"failed to list projects"}})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": items, "pagination": gin.H{"page": page, "per_page": per, "total": total}})
}

func (h *ProjectHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var item models.Project
	if err := h.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code":"NOT_FOUND","message":"project not found"}})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ProjectHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var item models.Project
	if err := h.DB.First(&item, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{"code":"NOT_FOUND","message":"project not found"}})
		return
	}
	var p projectPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": gin.H{"code":"VALIDATION","message":"invalid payload"}})
		return
	}
	if p.Name != "" { item.Name = p.Name }
	if p.Description != nil { item.Description = p.Description }
	if p.Status != nil { item.Status = *p.Status }
	if p.StartDate != nil { item.StartDate = parseDatePtr(p.StartDate) }
	if p.EndDate != nil { item.EndDate = parseDatePtr(p.EndDate) }
	if p.OwnerUserID != nil { item.OwnerUserID = p.OwnerUserID }
	if err := h.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code":"INTERNAL","message":"failed to update project"}})
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ProjectHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.DB.Delete(&models.Project{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"code":"INTERNAL","message":"failed to delete project"}})
		return
	}
	c.Status(http.StatusNoContent)
}
