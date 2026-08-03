package models

import "time"

type ParkingClaim struct {
	ID                  int64      `json:"id"`
	EmployeeID          int64      `json:"employee_id"`
	ParkingStartDate    time.Time  `json:"parking_start_date"`
	ParkingEndDate      time.Time  `json:"parking_end_date"`
	ParkingLocation     string     `json:"parking_location"`
	Amount              float64    `json:"amount"`
	Description         string     `json:"description"`
	ReceiptPath         string     `json:"-"`
	ReceiptOriginalName string     `json:"receipt_original_name"`
	ReceiptMIMEType     string     `json:"receipt_mime_type"`
	ReceiptSize         int64      `json:"receipt_size"`
	Status              string     `json:"status"`
	AdminNote           string     `json:"admin_note"`
	ReviewedBy          *int64     `json:"reviewed_by"`
	ReviewedAt          *time.Time `json:"reviewed_at"`
	CreatedAt           time.Time  `json:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at"`
}
