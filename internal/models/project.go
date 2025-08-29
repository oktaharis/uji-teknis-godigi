package models

import "time"

type Project struct {
	ID 	  uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string     `gorm:"size:120;not null" json:"name"`
	Description *string    `json:"description,omitempty"`
	Status      string     `gorm:"size:20;not null;default:'planned'" json:"status"` // planned|in_progress|on_hold|done|canceled
	StartDate   *time.Time `json:"start_date,omitempty"`
	EndDate     *time.Time `json:"end_date,omitempty"`
	OwnerUserID *uint      `json:"owner_user_id,omitempty"`

	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at,omitempty"`

	Owner User `gorm:"foreignKey:OwnerUserID;references:ID" json:"-"`
}

func (Project) TableName() string { return "projects" }