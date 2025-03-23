package model

import "time"

type JobInformation struct {
	ApplicationID string    `json:"application_id"`
	CompanyName   string    `json:"company_name"`
	PositionName  string    `json:"position_name"`
	ApplyDate     time.Time `json:"apply_date"`
	CurrentStatus string    `json:"current_status"`
	Link          string    `json:"link"`
	ApplyChannel  string    `json:"apply_channel"`
	ApplyMethod   string    `json:"apply_method"`
	ResumeVersion string    `json:"resume_version"`
	Priority      string    `json:"priority"`
	Notes         string    `json:"notes"`
}
