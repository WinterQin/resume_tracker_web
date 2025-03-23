package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username    string     `gorm:"type:varchar(32);uniqueIndex;not null" json:"username"`
	Password    string     `gorm:"type:varchar(256);not null" json:"-"`
	Email       string     `gorm:"type:varchar(128);uniqueIndex;not null" json:"email"`
	Age         int        `json:"age"`
	Gender      string     `json:"gender"`
	Phone       string     `json:"phone"`
	LastLoginAt *time.Time `json:"last_login_at"`
}
