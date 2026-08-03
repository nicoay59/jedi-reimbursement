package services

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"time"

	"jedi-reimbursement-system/backend/internal/models"
)

var ErrInvalidReportFilter = errors.New("filter laporan tidak valid")

type ReportRepository interface {
	Summary(
		ctx context.Context,
		period models.ReportPeriod,
	) (models.DashboardSummary, error)
	Trend(
		ctx context.Context,
		period models.ReportPeriod,
	) ([]models.DashboardTrend, error)
	Recent(
		ctx context.Context,
		period models.ReportPeriod,
		limit int,
	) ([]models.ClaimReview, error)
	List(
		ctx context.Context,
		filter models.ReportFilter,
		limit int,
		offset int,
	) ([]models.ClaimReview, error)
	Count(
		ctx context.Context,
		filter models.ReportFilter,
	) (int64, error)
	Export(
		ctx context.Context,
		filter models.ReportFilter,
	) ([]models.ClaimReview, error)
}

type ReportValidationError struct {
	Fields map[string]string
}

func (e *ReportValidationError) Error() string {
	return ErrInvalidReportFilter.Error()
}

type ReportPage struct {
	Filter     models.ReportFilter
	Items      []models.ClaimReview
	Page       int
	Limit      int
	Total      int64
	TotalPages int
}

type ReportService struct {
	repository ReportRepository
	now        func() time.Time
}

func NewReportService(
	repository ReportRepository,
) *ReportService {
	return &ReportService{
		repository: repository,
		now:        time.Now,
	}
}

func (s *ReportService) Dashboard(
	ctx context.Context,
	startDate string,
	endDate string,
) (models.DashboardData, error) {
	period, err := s.period(startDate, endDate)
	if err != nil {
		return models.DashboardData{}, err
	}

	summary, err := s.repository.Summary(ctx, period)
	if err != nil {
		return models.DashboardData{}, err
	}

	trend, err := s.repository.Trend(ctx, period)
	if err != nil {
		return models.DashboardData{}, err
	}

	recent, err := s.repository.Recent(ctx, period, 5)
	if err != nil {
		return models.DashboardData{}, err
	}

	return models.DashboardData{
		Period:  period,
		Summary: summary,
		Trend:   trend,
		Recent:  recent,
	}, nil
}

func (s *ReportService) List(
	ctx context.Context,
	startDate string,
	endDate string,
	claimType string,
	status string,
	page int,
	limit int,
) (ReportPage, error) {
	filter, err := s.filter(
		startDate,
		endDate,
		claimType,
		status,
	)
	if err != nil {
		return ReportPage{}, err
	}

	if page < 1 {
		page = 1
	}

	if limit < 1 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := (page - 1) * limit

	items, err := s.repository.List(
		ctx,
		filter,
		limit,
		offset,
	)
	if err != nil {
		return ReportPage{}, err
	}

	total, err := s.repository.Count(ctx, filter)
	if err != nil {
		return ReportPage{}, err
	}

	totalPages := 0
	if total > 0 {
		totalPages = int(math.Ceil(float64(total) / float64(limit)))
	}

	return ReportPage{
		Filter:     filter,
		Items:      items,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: totalPages,
	}, nil
}

func (s *ReportService) ExportCSV(
	ctx context.Context,
	startDate string,
	endDate string,
	claimType string,
	status string,
) ([]byte, models.ReportFilter, error) {
	filter, err := s.filter(
		startDate,
		endDate,
		claimType,
		status,
	)
	if err != nil {
		return nil, models.ReportFilter{}, err
	}

	items, err := s.repository.Export(ctx, filter)
	if err != nil {
		return nil, models.ReportFilter{}, err
	}

	data, err := BuildClaimsCSV(items)
	if err != nil {
		return nil, models.ReportFilter{},
			fmt.Errorf("membuat CSV laporan: %w", err)
	}

	return data, filter, nil
}

func (s *ReportService) filter(
	startDate string,
	endDate string,
	claimType string,
	status string,
) (models.ReportFilter, error) {
	period, err := s.period(startDate, endDate)
	if err != nil {
		return models.ReportFilter{}, err
	}

	normalizedType, err := normalizeClaimType(claimType, true)
	if err != nil {
		return models.ReportFilter{}, &ReportValidationError{
			Fields: map[string]string{
				"type": "Jenis klaim tidak valid",
			},
		}
	}

	normalizedStatus, err := normalizeClaimStatus(status, true)
	if err != nil {
		return models.ReportFilter{}, &ReportValidationError{
			Fields: map[string]string{
				"status": "Status klaim tidak valid",
			},
		}
	}

	return models.ReportFilter{
		Period:    period,
		ClaimType: normalizedType,
		Status:    normalizedStatus,
	}, nil
}

func (s *ReportService) period(
	startDate string,
	endDate string,
) (models.ReportPeriod, error) {
	now := s.now()
	location := now.Location()

	if strings.TrimSpace(startDate) == "" {
		startDate = time.Date(
			now.Year(),
			now.Month(),
			1,
			0,
			0,
			0,
			0,
			location,
		).Format("2006-01-02")
	}

	if strings.TrimSpace(endDate) == "" {
		endDate = dateOnly(now).Format("2006-01-02")
	}

	start, startErr := time.ParseInLocation(
		"2006-01-02",
		startDate,
		location,
	)
	end, endErr := time.ParseInLocation(
		"2006-01-02",
		endDate,
		location,
	)

	fields := make(map[string]string)

	if startErr != nil {
		fields["start_date"] = "Tanggal mulai tidak valid"
	}

	if endErr != nil {
		fields["end_date"] = "Tanggal selesai tidak valid"
	}

	if len(fields) == 0 {
		if end.Before(start) {
			fields["end_date"] =
				"Tanggal selesai tidak boleh sebelum tanggal mulai"
		}

		if end.After(dateOnly(now)) {
			fields["end_date"] =
				"Tanggal selesai tidak boleh di masa depan"
		}

		if end.Sub(start) > 365*24*time.Hour {
			fields["end_date"] =
				"Rentang laporan maksimal 366 hari"
		}
	}

	if len(fields) > 0 {
		return models.ReportPeriod{}, &ReportValidationError{
			Fields: fields,
		}
	}

	return models.ReportPeriod{
		StartDate: start,
		EndDate:   end,
	}, nil
}

func BuildClaimsCSV(
	items []models.ClaimReview,
) ([]byte, error) {
	buffer := &bytes.Buffer{}

	// BOM membantu Microsoft Excel membaca UTF-8 dengan benar.
	buffer.Write([]byte{0xEF, 0xBB, 0xBF})

	writer := csv.NewWriter(buffer)
	writer.UseCRLF = true

	header := []string{
		"Jenis",
		"ID Klaim",
		"NIK",
		"Nama Karyawan",
		"Tanggal Mulai",
		"Tanggal Selesai",
		"Ringkasan",
		"Deskripsi",
		"Nominal Parkir",
		"Waktu Mulai",
		"Waktu Selesai",
		"Durasi Jam",
		"Status",
		"Catatan Admin",
		"Pemeriksa",
		"Waktu Pemeriksaan",
		"Waktu Pengajuan",
	}

	if err := writer.Write(header); err != nil {
		return nil, err
	}

	for index := range items {
		item := items[index]

		reviewedAt := ""
		if item.ReviewedAt != nil {
			reviewedAt = item.ReviewedAt.Format(time.RFC3339)
		}

		record := []string{
			sanitizeCSVCell(item.ClaimType),
			strconv.FormatInt(item.ClaimID, 10),
			sanitizeCSVCell(item.EmployeeNumber),
			sanitizeCSVCell(item.EmployeeName),
			item.ClaimDate.Format("2006-01-02"),
			item.ClaimEndDate.Format("2006-01-02"),
			sanitizeCSVCell(item.Title),
			sanitizeCSVCell(item.Description),
			strconv.FormatFloat(item.Amount, 'f', 2, 64),
			item.StartTime,
			item.EndTime,
			strconv.FormatFloat(
				item.DurationHours,
				'f',
				2,
				64,
			),
			item.Status,
			sanitizeCSVCell(item.AdminNote),
			sanitizeCSVCell(item.ReviewerName),
			reviewedAt,
			item.CreatedAt.Format(time.RFC3339),
		}

		if err := writer.Write(record); err != nil {
			return nil, err
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func sanitizeCSVCell(value string) string {
	trimmed := strings.TrimLeft(value, " \t\r\n")
	if trimmed == "" {
		return value
	}

	switch trimmed[0] {
	case '=', '+', '-', '@':
		return "'" + value
	default:
		return value
	}
}
