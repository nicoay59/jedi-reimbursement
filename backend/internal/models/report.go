package models

import "time"

type ReportPeriod struct {
	StartDate time.Time
	EndDate   time.Time
}

type DashboardSummary struct {
	TotalClaims           int64
	PendingClaims         int64
	ApprovedClaims        int64
	RejectedClaims        int64
	ParkingClaims         int64
	OvertimeClaims        int64
	TotalParkingAmount    float64
	ApprovedParkingAmount float64
	TotalOvertimeHours    float64
	ApprovedOvertimeHours float64
}

type DashboardTrend struct {
	Date           time.Time
	TotalClaims    int64
	ParkingClaims  int64
	OvertimeClaims int64
	ApprovedClaims int64
}

type DashboardData struct {
	Period  ReportPeriod
	Summary DashboardSummary
	Trend   []DashboardTrend
	Recent  []ClaimReview
}

type ReportFilter struct {
	Period    ReportPeriod
	ClaimType string
	Status    string
}
