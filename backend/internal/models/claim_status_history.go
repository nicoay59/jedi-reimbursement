package models

import "time"

type ClaimStatusHistory struct {
	ID             int64     `json:"id"`
	ClaimType      string    `json:"claim_type"`
	ClaimID        int64     `json:"claim_id"`
	PreviousStatus *string   `json:"previous_status"`
	NewStatus      string    `json:"new_status"`
	Note           string    `json:"note"`
	UpdatedBy      int64     `json:"updated_by"`
	CreatedAt      time.Time `json:"created_at"`
}
