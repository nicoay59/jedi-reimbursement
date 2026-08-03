package dto

import "time"

type CreateOvertimeClaimRequest struct {
	OvertimeDate    string `json:"overtime_date"`
	StartTime       string `json:"start_time"`
	EndTime         string `json:"end_time"`
	WorkDescription string `json:"work_description"`
}

type OvertimeClaimResponse struct {
	ID              int64   `json:"id"`
	OvertimeDate    string  `json:"overtime_date"`
	StartTime       string  `json:"start_time"`
	EndTime         string  `json:"end_time"`
	DurationHours   float64 `json:"duration_hours"`
	WorkDescription string  `json:"work_description"`
	Status          string  `json:"status"`
	AdminNote       string  `json:"admin_note"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
}

type OvertimeClaimListResponse struct {
	Items      []OvertimeClaimResponse `json:"items"`
	Pagination PaginationResponse      `json:"pagination"`
}

func NewOvertimeClaimResponse(
	claimID int64,
	overtimeDate time.Time,
	startTime string,
	endTime string,
	durationHours float64,
	workDescription string,
	status string,
	adminNote string,
	createdAt time.Time,
	updatedAt time.Time,
) OvertimeClaimResponse {
	return OvertimeClaimResponse{
		ID:              claimID,
		OvertimeDate:    overtimeDate.Format("2006-01-02"),
		StartTime:       normalizeClock(startTime),
		EndTime:         normalizeClock(endTime),
		DurationHours:   durationHours,
		WorkDescription: workDescription,
		Status:          status,
		AdminNote:       adminNote,
		CreatedAt:       createdAt.Format(time.RFC3339),
		UpdatedAt:       updatedAt.Format(time.RFC3339),
	}
}

func normalizeClock(value string) string {
	if len(value) >= 5 {
		return value[:5]
	}
	return value
}
