package models

import "time"

const (
	ClaimTypeAll      = "ALL"
	ClaimTypeParking  = "PARKING"
	ClaimTypeOvertime = "OVERTIME"
)

type ClaimReview struct {
	ClaimType           string
	ClaimID             int64
	EmployeeID          int64
	EmployeeNumber      string
	EmployeeName        string
	ClaimDate           time.Time
	ClaimEndDate        time.Time
	Title               string
	Description         string
	Amount              float64
	StartTime           string
	EndTime             string
	DurationHours       float64
	Status              string
	AdminNote           string
	ReviewedBy          *int64
	ReviewerName        string
	ReviewedAt          *time.Time
	CreatedAt           time.Time
	UpdatedAt           time.Time
	ReceiptPath         string
	ReceiptOriginalName string
	ReceiptMIMEType     string
	ReceiptSize         int64
}

type ClaimReviewHistory struct {
	ID             int64
	ClaimType      string
	ClaimID        int64
	PreviousStatus string
	NewStatus      string
	Note           string
	UpdatedBy      int64
	UpdatedByName  string
	CreatedAt      time.Time
}
