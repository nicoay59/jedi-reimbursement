package dto

import "time"

type ReviewClaimRequest struct {
	Status string `json:"status"`
	Note   string `json:"note"`
}

type ClaimReviewResponse struct {
	ClaimType           string  `json:"claim_type"`
	ClaimID             int64   `json:"claim_id"`
	EmployeeID          int64   `json:"employee_id"`
	EmployeeNumber      string  `json:"employee_number"`
	EmployeeName        string  `json:"employee_name"`
	ClaimDate           string  `json:"claim_date"`
	ClaimEndDate        string  `json:"claim_end_date"`
	Title               string  `json:"title"`
	Description         string  `json:"description"`
	Amount              float64 `json:"amount"`
	StartTime           string  `json:"start_time"`
	EndTime             string  `json:"end_time"`
	DurationHours       float64 `json:"duration_hours"`
	Status              string  `json:"status"`
	AdminNote           string  `json:"admin_note"`
	ReviewerName        string  `json:"reviewer_name"`
	ReviewedAt          string  `json:"reviewed_at"`
	CreatedAt           string  `json:"created_at"`
	UpdatedAt           string  `json:"updated_at"`
	ReceiptAvailable    bool    `json:"receipt_available"`
	ReceiptOriginalName string  `json:"receipt_original_name"`
	ReceiptMIMEType     string  `json:"receipt_mime_type"`
	ReceiptSize         int64   `json:"receipt_size"`
}

type ClaimReviewListResponse struct {
	Items      []ClaimReviewResponse `json:"items"`
	Pagination PaginationResponse    `json:"pagination"`
}

type ClaimHistoryResponse struct {
	ID             int64  `json:"id"`
	PreviousStatus string `json:"previous_status"`
	NewStatus      string `json:"new_status"`
	Note           string `json:"note"`
	UpdatedBy      int64  `json:"updated_by"`
	UpdatedByName  string `json:"updated_by_name"`
	CreatedAt      string `json:"created_at"`
}

type ClaimHistoryListResponse struct {
	Items []ClaimHistoryResponse `json:"items"`
}

func FormatOptionalTime(value *time.Time) string {
	if value == nil {
		return ""
	}
	return value.Format(time.RFC3339)
}
