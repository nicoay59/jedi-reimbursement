package services

import (
	"bytes"
	"context"
	"testing"
	"time"

	"jedi-reimbursement-system/backend/internal/models"
)

type fakeReportRepository struct {
	summary models.DashboardSummary
	trend   []models.DashboardTrend
	recent  []models.ClaimReview
	items   []models.ClaimReview
}

func (f *fakeReportRepository) Summary(
	context.Context,
	models.ReportPeriod,
) (models.DashboardSummary, error) {
	return f.summary, nil
}

func (f *fakeReportRepository) Trend(
	context.Context,
	models.ReportPeriod,
) ([]models.DashboardTrend, error) {
	return f.trend, nil
}

func (f *fakeReportRepository) Recent(
	context.Context,
	models.ReportPeriod,
	int,
) ([]models.ClaimReview, error) {
	return f.recent, nil
}

func (f *fakeReportRepository) List(
	context.Context,
	models.ReportFilter,
	int,
	int,
) ([]models.ClaimReview, error) {
	return f.items, nil
}

func (f *fakeReportRepository) Count(
	context.Context,
	models.ReportFilter,
) (int64, error) {
	return int64(len(f.items)), nil
}

func (f *fakeReportRepository) Export(
	context.Context,
	models.ReportFilter,
) ([]models.ClaimReview, error) {
	return f.items, nil
}

func TestDashboardUsesCurrentMonthByDefault(t *testing.T) {
	repository := &fakeReportRepository{
		summary: models.DashboardSummary{
			TotalClaims: 4,
		},
	}

	service := NewReportService(repository)
	service.now = func() time.Time {
		return time.Date(
			2026,
			7,
			18,
			12,
			0,
			0,
			0,
			time.Local,
		)
	}

	result, err := service.Dashboard(
		context.Background(),
		"",
		"",
	)
	if err != nil {
		t.Fatalf("Dashboard() error = %v", err)
	}

	if result.Period.StartDate.Format("2006-01-02") != "2026-07-01" {
		t.Fatalf("StartDate = %v", result.Period.StartDate)
	}

	if result.Period.EndDate.Format("2006-01-02") != "2026-07-18" {
		t.Fatalf("EndDate = %v", result.Period.EndDate)
	}
}

func TestReportRejectsLongPeriod(t *testing.T) {
	service := NewReportService(&fakeReportRepository{})
	service.now = func() time.Time {
		return time.Date(
			2026,
			12,
			31,
			12,
			0,
			0,
			0,
			time.UTC,
		)
	}

	_, err := service.List(
		context.Background(),
		"2025-01-01",
		"2026-12-31",
		"ALL",
		"ALL",
		1,
		10,
	)
	if err == nil {
		t.Fatal("List() seharusnya menolak rentang panjang")
	}
}

func TestBuildClaimsCSVAddsBOMAndSanitizesFormula(t *testing.T) {
	data, err := BuildClaimsCSV([]models.ClaimReview{
		{
			ClaimType:      models.ClaimTypeParking,
			ClaimID:        1,
			EmployeeNumber: "=2+2",
			EmployeeName:   "+cmd",
			ClaimDate:      time.Date(2026, 7, 18, 0, 0, 0, 0, time.UTC),
			Title:          "@lokasi",
			Description:    "-deskripsi",
			Status:         models.ClaimStatusApproved,
			CreatedAt:      time.Date(2026, 7, 18, 10, 0, 0, 0, time.UTC),
		},
	})
	if err != nil {
		t.Fatalf("BuildClaimsCSV() error = %v", err)
	}

	if !bytes.HasPrefix(data, []byte{0xEF, 0xBB, 0xBF}) {
		t.Fatal("CSV tidak memiliki UTF-8 BOM")
	}

	for _, expected := range [][]byte{
		[]byte("'=2+2"),
		[]byte("'+cmd"),
		[]byte("'@lokasi"),
		[]byte("'-deskripsi"),
	} {
		if !bytes.Contains(data, expected) {
			t.Fatalf("CSV tidak berisi %q", expected)
		}
	}
}
