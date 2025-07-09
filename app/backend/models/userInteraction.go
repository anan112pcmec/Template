package models

import (
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Username  string    `gorm:"type:varchar(100);not null;unique" json:"username"`
	Email     string    `gorm:"type:varchar(100);not null;unique" json:"email"`
	Password  string    `gorm:"type:varchar(255);not null" json:"-"`
	Role      string    `gorm:"type:varchar(50);default:'user'" json:"role"`
	IsActive  bool      `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
