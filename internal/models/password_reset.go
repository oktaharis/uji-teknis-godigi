package models

import "time"

type PasswordReset struct {
	ID        uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID    uint       `gorm:"column:user_id;not null" json:"user_id"`
	Token     string     `gorm:"size:128;not null" json:"token"`
	ExpiresAt time.Time  `gorm:"not null" json:"expires_at"`
	UsedAt    *time.Time `json:"used_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`

	User User `gorm:"foreignKey:UserID;references:ID;constraint:OnDelete:CASCADE"`
}

func (PasswordReset) TableName() string { return "password_resets" }
