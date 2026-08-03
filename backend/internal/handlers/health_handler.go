package handlers

import (
	"context"
	"net/http"
	"time"

	"jedi-reimbursement-system/backend/internal/config"
	"jedi-reimbursement-system/backend/internal/database"
	"jedi-reimbursement-system/backend/internal/responses"
)

type HealthHandler struct {
	config config.Config
	pinger database.Pinger
}

func NewHealthHandler(
	cfg config.Config,
	pinger database.Pinger,
) *HealthHandler {
	return &HealthHandler{
		config: cfg,
		pinger: pinger,
	}
}

func (h *HealthHandler) Check(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	statusCode := http.StatusOK
	databaseStatus := "connected"
	message := "API dan database berjalan dengan baik"

	if err := h.pinger.PingContext(ctx); err != nil {
		statusCode = http.StatusServiceUnavailable
		databaseStatus = "unavailable"
		message = "API berjalan, tetapi database tidak tersedia"
	}

	responses.WriteJSON(w, statusCode, responses.APIResponse{
		Success: statusCode == http.StatusOK,
		Message: message,
		Data: map[string]any{
			"application": h.config.AppName,
			"environment": h.config.AppEnv,
			"api":         "healthy",
			"database":    databaseStatus,
			"version":     "1.2.0",
			"timestamp":   time.Now().UTC().Format(time.RFC3339),
		},
	})
}
