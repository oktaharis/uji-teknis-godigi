package models

import "time"

type Deal struct {
	DealID     uint      `gorm:"column:deal_id;primaryKey;autoIncrement" json:"id"`
	LeadID     uint      `gorm:"column:lead_id;not null;index" json:"lead_id"`
	DealName   *string   `gorm:"column:deal_name;size:120" json:"deal_name,omitempty"`
	AmountIDR  int64     `gorm:"column:amount_idr;not null" json:"amount_idr"`
	Currency   string    `gorm:"column:currency;size:3;not null;default:IDR" json:"currency"`
	TermMonths int       `gorm:"column:term_months;not null;default:12" json:"term_months"`
	Stage      string    `gorm:"column:stage;size:20;not null" json:"stage"`
	ClosedAt   time.Time `gorm:"column:closed_at;not null" json:"closed_at"`

	Lead Lead `gorm:"foreignKey:LeadID;references:LeadID" json:"-"`
}

func (Deal) TableName() string { return "deals" }
