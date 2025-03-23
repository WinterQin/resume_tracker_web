package model

type ContactInfo struct {
	ApplicationID string `json:"application_id"`
	ContactName   string `json:"contact_name"`
	ContactInfo   string `json:"contact_info"`
	Notes         string `json:"notes"`
}
