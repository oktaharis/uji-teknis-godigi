package handlers

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/oktaharis/uji-teknis-godigi/internal/models"
	"github.com/oktaharis/uji-teknis-godigi/internal/response"
)

type LeadHandler struct {
	DB *gorm.DB
}

func NewLeadHandler(db *gorm.DB) *LeadHandler { return &LeadHandler{DB: db} }

type leadPayload struct {
	CompanyName string  `json:"company_name" binding:"required"`
	ContactName string  `json:"contact_name" binding:"required"`
	Email       string  `json:"email" binding:"required,email"`
	Phone       *string `json:"phone"`
	Source      *string `json:"source"`
	Industry    *string `json:"industry"`
	Region      *string `json:"region"`
	SalesRep    *string `json:"sales_rep"`
	Status      *string `json:"status"`
	Notes       *string `json:"notes"`
}

func (h *LeadHandler) Create(c *gin.Context) {
	var p leadPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		response.UnprocessableEntity(c, "Validation Error", response.ExtractValidationErrors(err))
		return
	}
	lead := models.Lead{
		CompanyName: p.CompanyName,
		ContactName: p.ContactName,
		Email:       p.Email,
		Phone:       p.Phone,
		Source:      p.Source,
		Industry:    p.Industry,
		Region:      p.Region,
		SalesRep:    p.SalesRep,
		Status:      p.Status,
		Notes:       p.Notes,
	}
	if err := h.DB.Create(&lead).Error; err != nil {
		response.InternalError(c, "Failed to create lead")
		return
	}
	response.Created(c, gin.H{
		"id": lead.LeadID, "company_name": lead.CompanyName, "status": lead.Status, "created_at": lead.CreatedAt,
	}, "Lead created")
}

func (h *LeadHandler) List(c *gin.Context) {
	var leads []models.Lead

	q := h.DB.Model(&models.Lead{})
	if v := c.Query("status"); v != "" {
		q = q.Where("status = ?", v)
	}
	if v := c.Query("source"); v != "" {
		q = q.Where("source = ?", v)
	}
	if v := c.Query("q"); v != "" {
		q = q.Where("company_name LIKE ? OR contact_name LIKE ? OR email LIKE ?", "%"+v+"%", "%"+v+"%", "%"+v+"%")
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	per, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if page < 1 {
		page = 1
	}
	if per < 1 {
		per = 10
	}
	offset := (page - 1) * per

	var total int64
	q.Count(&total)

	if err := q.Order("created_at DESC").Limit(per).Offset(offset).Find(&leads).Error; err != nil {
		response.InternalError(c, "Failed to list leads")
		return
	}
	response.OK(c, response.List(leads, page, per, total), "Lead list")
}

func (h *LeadHandler) Get(c *gin.Context) {
	id := c.Param("id")
	var lead models.Lead
	if err := h.DB.First(&lead, id).Error; err != nil {
		response.NotFound(c, "Lead not found")
		return
	}
	response.OK(c, lead, "Lead detail")
}

func (h *LeadHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var lead models.Lead
	if err := h.DB.First(&lead, id).Error; err != nil {
		response.NotFound(c, "Lead not found")
		return
	}
	var p leadPayload
	if err := c.ShouldBindJSON(&p); err != nil {
		response.UnprocessableEntity(c, "Validation Error", response.ExtractValidationErrors(err))
		return
	}
	lead.CompanyName = p.CompanyName
	lead.ContactName = p.ContactName
	lead.Email = p.Email
	lead.Phone = p.Phone
	lead.Source = p.Source
	lead.Industry = p.Industry
	lead.Region = p.Region
	lead.SalesRep = p.SalesRep
	lead.Status = p.Status
	lead.Notes = p.Notes

	if err := h.DB.Save(&lead).Error; err != nil {
		response.InternalError(c, "Failed to update lead")
		return
	}
	response.OK(c, lead, "Lead updated")
}

func (h *LeadHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.DB.Delete(&models.Lead{}, id).Error; err != nil {
		response.InternalError(c, "Failed to delete lead")
		return
	}
	response.NoContent(c, "Lead deleted")
}

// GET /leads/summary?from=YYYY-MM-DD&to=YYYY-MM-DD
func (h *LeadHandler) Summary(c *gin.Context) {
	from := c.Query("from")
	to := c.Query("to")
	var fromT, toT time.Time
	var err error

	if from != "" {
		fromT, err = time.Parse("2006-01-02", from)
		if err != nil {
			from = ""
		}
	}
	if to != "" {
		toT, err = time.Parse("2006-01-02", to)
		if err != nil {
			to = ""
		}
	}

	q := h.DB.Model(&models.Lead{})
	if from != "" {
		q = q.Where("created_at >= ?", fromT)
	}
	if to != "" {
		q = q.Where("created_at < ?", toT.Add(24*time.Hour))
	}

	// total leads
	var total int64
	q.Count(&total)

	type rowAgg struct {
		Name  *string `gorm:"column:name"`
		Count int64   `gorm:"column:count"`
	}
	by := func(field string) map[string]int64 {
		var rows []rowAgg
		res := map[string]int64{}
		sub := h.DB.Model(&models.Lead{})
		if from != "" {
			sub = sub.Where("created_at >= ?", fromT)
		}
		if to != "" {
			sub = sub.Where("created_at < ?", toT.Add(24*time.Hour))
		}
		sub.Select(field+" AS name, COUNT(*) AS count").Group(field).Scan(&rows)
		for _, r := range rows {
			k := ""
			if r.Name != nil {
				k = *r.Name
			}
			res[k] = r.Count
		}
		return res
	}

	byStatus := by("status")
	bySource := by("source")
	byRegion := by("region")

	// deals aggregate (by closed_at)
	type DealAgg struct {
		Count int64
		Total int64   `gorm:"column:total"`
		Avg   float64 `gorm:"column:avg"`
	}
	var agg DealAgg
	dq := h.DB.Model(&models.Deal{})
	if from != "" {
		dq = dq.Where("closed_at >= ?", fromT)
	}
	if to != "" {
		dq = dq.Where("closed_at < ?", toT.Add(24*time.Hour))
	}
	dq.Select("COUNT(*) as count, COALESCE(SUM(amount_idr),0) as total, COALESCE(AVG(term_months),0) as avg").Scan(&agg)

	var byStageRows []rowAgg
	dqStages := h.DB.Model(&models.Deal{})
	if from != "" {
		dqStages = dqStages.Where("closed_at >= ?", fromT)
	}
	if to != "" {
		dqStages = dqStages.Where("closed_at < ?", toT.Add(24*time.Hour))
	}
	dqStages.Select("stage AS name, COUNT(*) AS count").Group("stage").Scan(&byStageRows)
	stageMap := map[string]int64{}
	for _, r := range byStageRows {
		k := ""
		if r.Name != nil {
			k = *r.Name
		}
		stageMap[k] = r.Count
	}

	response.OK(c, gin.H{
		"total_leads": total,
		"by_status":   byStatus,
		"by_source":   bySource,
		"by_region":   byRegion,
		"deals": gin.H{
			"count":            agg.Count,
			"total_amount_idr": agg.Total,
			"avg_term_months":  agg.Avg,
			"by_stage":         stageMap,
		},
	}, "Lead summary")
}
