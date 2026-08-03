package dto

import "time"

type ParkingClaimResponse struct {
	ID                  int64   `json:"id"`
	ParkingStartDate    string  `json:"parking_start_date"`
	ParkingEndDate      string  `json:"parking_end_date"`
	ParkingLocation     string  `json:"parking_location"`
	Amount              float64 `json:"amount"`
	Description         string  `json:"description"`
	ReceiptOriginalName string  `json:"receipt_original_name"`
	ReceiptMIMEType     string  `json:"receipt_mime_type"`
	ReceiptSize         int64   `json:"receipt_size"`
	ReceiptAvailable    bool    `json:"receipt_available"`
	Status              string  `json:"status"`
	AdminNote           string  `json:"admin_note"`
	CreatedAt           string  `json:"created_at"`
	UpdatedAt           string  `json:"updated_at"`
}

type PaginationResponse struct {
	Page       int   `json:"page"`
	Limit      int   `json:"limit"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
}

type ParkingClaimListResponse struct {
	Items      []ParkingClaimResponse `json:"items"`
	Pagination PaginationResponse     `json:"pagination"`
}

func NewParkingClaimResponse(
	claimID int64,
	parkingStartDate time.Time,
	parkingEndDate time.Time,
	parkingLocation string,
	amount float64,
	description string,
	receiptOriginalName string,
	receiptMIMEType string,
	receiptSize int64,
	status string,
	adminNote string,
	createdAt time.Time,
	updatedAt time.Time,
) ParkingClaimResponse {
	return ParkingClaimResponse{
		ID:                  claimID,
		ParkingStartDate:    parkingStartDate.Format("2006-01-02"),
		ParkingEndDate:      parkingEndDate.Format("2006-01-02"),
		ParkingLocation:     parkingLocation,
		Amount:              amount,
		Description:         description,
		ReceiptOriginalName: receiptOriginalName,
		ReceiptMIMEType:     receiptMIMEType,
		ReceiptSize:         receiptSize,
		ReceiptAvailable:    receiptOriginalName != "",
		Status:              status,
		AdminNote:           adminNote,
		CreatedAt:           createdAt.Format(time.RFC3339),
		UpdatedAt:           updatedAt.Format(time.RFC3339),
	}
}
