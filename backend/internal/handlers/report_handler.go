package handlers

import (
	"errors"
	"fmt"
	"mime"
	"net/http"

	"jedi-reimbursement-system/backend/internal/dto"
	"jedi-reimbursement-system/backend/internal/responses"
	"jedi-reimbursement-system/backend/internal/services"
)

type ReportHandler struct {
	service *services.ReportService
}

func NewReportHandler(
	service *services.ReportService,
) *ReportHandler {
	return &ReportHandler{service: service}
}

func (h *ReportHandler) Dashboard(
	w http.ResponseWriter,
	r *http.Request,
) {
	result, err := h.service.Dashboard(
		r.Context(),
		r.URL.Query().Get("start_date"),
		r.URL.Query().Get("end_date"),
	)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}

	trend := make(
		[]dto.DashboardTrendResponse,
		0,
		len(result.Trend),
	)
	for index := range result.Trend {
		trend = append(
			trend,
			dto.DashboardTrendResponse{
				Date: result.Trend[index].Date.Format(
					"2006-01-02",
				),
				TotalClaims:    result.Trend[index].TotalClaims,
				ParkingClaims:  result.Trend[index].ParkingClaims,
				OvertimeClaims: result.Trend[index].OvertimeClaims,
				ApprovedClaims: result.Trend[index].ApprovedClaims,
			},
		)
	}

	recent := make(
		[]dto.ClaimReviewResponse,
		0,
		len(result.Recent),
	)
	for index := range result.Recent {
		recent = append(
			recent,
			claimReviewResponse(&result.Recent[index]),
		)
	}

	responses.WriteJSON(
		w,
		http.StatusOK,
		responses.APIResponse{
			Success: true,
			Message: "Dashboard berhasil dimuat",
			Data: dto.DashboardResponse{
				Period: dto.ReportPeriodResponse{
					StartDate: result.Period.StartDate.Format(
						"2006-01-02",
					),
					EndDate: result.Period.EndDate.Format(
						"2006-01-02",
					),
				},
				Summary: dto.DashboardSummaryResponse{
					TotalClaims:           result.Summary.TotalClaims,
					PendingClaims:         result.Summary.PendingClaims,
					ApprovedClaims:        result.Summary.ApprovedClaims,
					RejectedClaims:        result.Summary.RejectedClaims,
					ParkingClaims:         result.Summary.ParkingClaims,
					OvertimeClaims:        result.Summary.OvertimeClaims,
					TotalParkingAmount:    result.Summary.TotalParkingAmount,
					ApprovedParkingAmount: result.Summary.ApprovedParkingAmount,
					TotalOvertimeHours:    result.Summary.TotalOvertimeHours,
					ApprovedOvertimeHours: result.Summary.ApprovedOvertimeHours,
				},
				Trend:  trend,
				Recent: recent,
			},
		},
	)
}

func (h *ReportHandler) List(
	w http.ResponseWriter,
	r *http.Request,
) {
	result, err := h.service.List(
		r.Context(),
		r.URL.Query().Get("start_date"),
		r.URL.Query().Get("end_date"),
		r.URL.Query().Get("type"),
		r.URL.Query().Get("status"),
		positiveQueryInt(r, "page", 1),
		positiveQueryInt(r, "limit", 10),
	)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}

	items := make(
		[]dto.ClaimReviewResponse,
		0,
		len(result.Items),
	)
	for index := range result.Items {
		items = append(
			items,
			claimReviewResponse(&result.Items[index]),
		)
	}

	responses.WriteJSON(
		w,
		http.StatusOK,
		responses.APIResponse{
			Success: true,
			Message: "Laporan berhasil dimuat",
			Data: dto.ReportListResponse{
				Period: dto.ReportPeriodResponse{
					StartDate: result.Filter.Period.StartDate.Format(
						"2006-01-02",
					),
					EndDate: result.Filter.Period.EndDate.Format(
						"2006-01-02",
					),
				},
				ClaimType: result.Filter.ClaimType,
				Status:    result.Filter.Status,
				Items:     items,
				Pagination: dto.PaginationResponse{
					Page:       result.Page,
					Limit:      result.Limit,
					Total:      result.Total,
					TotalPages: result.TotalPages,
				},
			},
		},
	)
}

func (h *ReportHandler) Export(
	w http.ResponseWriter,
	r *http.Request,
) {
	data, filter, err := h.service.ExportCSV(
		r.Context(),
		r.URL.Query().Get("start_date"),
		r.URL.Query().Get("end_date"),
		r.URL.Query().Get("type"),
		r.URL.Query().Get("status"),
	)
	if err != nil {
		h.writeServiceError(w, err)
		return
	}

	fileName := fmt.Sprintf(
		"jedi-reimbursement-%s-%s.csv",
		filter.Period.StartDate.Format("20060102"),
		filter.Period.EndDate.Format("20060102"),
	)

	w.Header().Set(
		"Content-Type",
		"text/csv; charset=utf-8",
	)
	w.Header().Set(
		"Content-Disposition",
		mime.FormatMediaType(
			"attachment",
			map[string]string{"filename": fileName},
		),
	)
	w.Header().Set(
		"Content-Length",
		fmt.Sprintf("%d", len(data)),
	)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(data)
}

func (h *ReportHandler) writeServiceError(
	w http.ResponseWriter,
	err error,
) {
	var validationError *services.ReportValidationError

	switch {
	case errors.As(err, &validationError):
		responses.WriteJSON(
			w,
			http.StatusUnprocessableEntity,
			responses.APIResponse{
				Success: false,
				Message: "Filter laporan tidak valid",
				Errors:  validationError.Fields,
			},
		)
	case errors.Is(err, services.ErrInvalidReportFilter):
		responses.WriteJSON(
			w,
			http.StatusBadRequest,
			responses.APIResponse{
				Success: false,
				Message: "Filter laporan tidak valid",
			},
		)
	default:
		responses.WriteJSON(
			w,
			http.StatusInternalServerError,
			responses.APIResponse{
				Success: false,
				Message: "Laporan tidak dapat diproses",
			},
		)
	}
}
