package models

import "time"

type User struct {
	ID           uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name         string     `gorm:"size:100;not null" json:"name"`
	Email        string     `gorm:"size:255;unique;not null" json:"email"`
	PasswordHash string     `gorm:"size:255;not null" json:"-"`
	Role         string     `gorm:"size:20;not null;default:user" json:"role"`
	TokenVersion int        `gorm:"not null;default:0" json:"-"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at,omitempty"`
}

func (User) TableName() string { return "users" }
