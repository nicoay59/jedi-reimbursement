package models

import "time"

type OvertimeClaim struct {
	ID              int64      `json:"id"`
	EmployeeID      int64      `json:"employee_id"`
	OvertimeDate    time.Time  `json:"overtime_date"`
	StartTime       string     `json:"start_time"`
	EndTime         string     `json:"end_time"`
	DurationHours   float64    `json:"duration_hours"`
	WorkDescription string     `json:"work_description"`
	Status          string     `json:"status"`
	AdminNote       string     `json:"admin_note"`
	ReviewedBy      *int64     `json:"reviewed_by"`
	ReviewedAt      *time.Time `json:"reviewed_at"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}
