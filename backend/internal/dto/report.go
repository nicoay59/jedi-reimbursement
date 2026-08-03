package dto

type ReportPeriodResponse struct {
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type DashboardSummaryResponse struct {
	TotalClaims           int64   `json:"total_claims"`
	PendingClaims         int64   `json:"pending_claims"`
	ApprovedClaims        int64   `json:"approved_claims"`
	RejectedClaims        int64   `json:"rejected_claims"`
	ParkingClaims         int64   `json:"parking_claims"`
	OvertimeClaims        int64   `json:"overtime_claims"`
	TotalParkingAmount    float64 `json:"total_parking_amount"`
	ApprovedParkingAmount float64 `json:"approved_parking_amount"`
	TotalOvertimeHours    float64 `json:"total_overtime_hours"`
	ApprovedOvertimeHours float64 `json:"approved_overtime_hours"`
}

type DashboardTrendResponse struct {
	Date           string `json:"date"`
	TotalClaims    int64  `json:"total_claims"`
	ParkingClaims  int64  `json:"parking_claims"`
	OvertimeClaims int64  `json:"overtime_claims"`
	ApprovedClaims int64  `json:"approved_claims"`
}

type DashboardResponse struct {
	Period  ReportPeriodResponse     `json:"period"`
	Summary DashboardSummaryResponse `json:"summary"`
	Trend   []DashboardTrendResponse `json:"trend"`
	Recent  []ClaimReviewResponse    `json:"recent"`
}

type ReportListResponse struct {
	Period     ReportPeriodResponse  `json:"period"`
	ClaimType  string                `json:"claim_type"`
	Status     string                `json:"status"`
	Items      []ClaimReviewResponse `json:"items"`
	Pagination PaginationResponse    `json:"pagination"`
}
