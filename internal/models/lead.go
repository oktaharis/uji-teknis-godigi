package models

import "time"

type Lead struct {
	LeadID      uint      `gorm:"column:lead_id;primaryKey;autoIncrement" json:"id"`
	CreatedAt   time.Time `gorm:"column:created_at" json:"created_at"`
	CompanyName string    `gorm:"column:company_name;size:255;not null" json:"company_name"`
	ContactName string    `gorm:"column:contact_name;size:100;not null" json:"contact_name"`
	Email       string    `gorm:"column:email;size:255;not null" json:"email"`
	Phone       *string   `gorm:"column:phone;size:30" json:"phone,omitempty"`
	Source      *string   `gorm:"column:source;size:50" json:"source,omitempty"`
	Industry    *string   `gorm:"column:industry;size:50" json:"industry,omitempty"`
	Region      *string   `gorm:"column:region;size:50" json:"region,omitempty"`
	SalesRep    *string   `gorm:"column:sales_rep;size:50" json:"sales_rep,omitempty"`
	Status      *string   `gorm:"column:status;size:30" json:"status,omitempty"`
	Notes       *string   `gorm:"column:notes" json:"notes,omitempty"`

	Deals []Deal `gorm:"foreignKey:LeadID;references:LeadID" json:"deals,omitempty"`
}

func (Lead) TableName() string { return "leads" }
