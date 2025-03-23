package model

import "time"

type JobInformationResult struct {
	ApplicationID string    `json:"application_id"`
	OfferDate     time.Time `json:"offer_date"`
	RejectDate    time.Time `json:"reject_date"`
	OfferDetails  string    `json:"offer_details"`
	RejectReason  string    `json:"reject_reason"`
	Feedback      string    `json:"feedback"`
}
