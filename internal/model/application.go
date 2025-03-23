package model

import (
	"time"

	"gorm.io/gorm"
)

// ApplicationStatus 申请状态
type ApplicationStatus string

const (
	StatusSubmitted ApplicationStatus = "submitted" // 已投递
	StatusWritten   ApplicationStatus = "written"   // 笔试中
	StatusInterview ApplicationStatus = "interview" // 面试中
	StatusAccepted  ApplicationStatus = "accepted"  // 已录用
	StatusRejected  ApplicationStatus = "rejected"  // 已拒绝
)

// Application 实习申请记录
type Application struct {
	gorm.Model
	UserID      uint              `gorm:"not null" json:"user_id"`
	Company     string            `gorm:"type:varchar(128);not null" json:"company"`
	Position    string            `gorm:"type:varchar(128);not null" json:"position"`
	Status      ApplicationStatus `gorm:"type:varchar(32);not null" json:"status"`
	ApplyDate   time.Time         `json:"apply_date"`
	NextEvent   *time.Time        `json:"next_event"` // 下一个面试/笔试时间
	EventType   string            `json:"event_type"` // 事件类型（笔试/面试）
	EventLink   string            `json:"event_link"` // 面试链接
	Notes       string            `gorm:"type:text" json:"notes"`
	Salary      string            `json:"salary"`       // 薪资
	Location    string            `json:"location"`     // 工作地点
	ContactInfo string            `json:"contact_info"` // 联系人信息
}
