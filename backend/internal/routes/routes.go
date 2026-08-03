package routes

import (
	"database/sql"
	"net/http"

	"jedi-reimbursement-system/backend/internal/config"
	"jedi-reimbursement-system/backend/internal/handlers"
	"jedi-reimbursement-system/backend/internal/middleware"
	"jedi-reimbursement-system/backend/internal/models"
	"jedi-reimbursement-system/backend/internal/repositories"
	"jedi-reimbursement-system/backend/internal/responses"
	"jedi-reimbursement-system/backend/internal/security"
	"jedi-reimbursement-system/backend/internal/services"
	"jedi-reimbursement-system/backend/internal/storage"
)

func New(
	cfg config.Config,
	db *sql.DB,
	receiptStorage *storage.ReceiptStorage,
) http.Handler {
	mux := http.NewServeMux()

	tokenManager := security.NewTokenManager(
		cfg.JWTSecret,
		cfg.JWTExpiresIn,
	)

	userRepository := repositories.NewUserRepository(db)
	authService := services.NewAuthService(
		userRepository,
		tokenManager,
	)

	parkingClaimRepository := repositories.NewParkingClaimRepository(db)
	parkingClaimService := services.NewParkingClaimService(
		parkingClaimRepository,
	)

	overtimeClaimRepository := repositories.NewOvertimeClaimRepository(db)
	overtimeClaimService := services.NewOvertimeClaimService(
		overtimeClaimRepository,
	)

	claimReviewRepository := repositories.NewClaimReviewRepository(db)
	claimReviewService := services.NewClaimReviewService(
		claimReviewRepository,
	)

	reportRepository := repositories.NewReportRepository(db)
	reportService := services.NewReportService(reportRepository)

	healthHandler := handlers.NewHealthHandler(cfg, db)
	authHandler := handlers.NewAuthHandler(authService)
	accessHandler := handlers.NewAccessHandler()
	parkingClaimHandler := handlers.NewParkingClaimHandler(
		parkingClaimService,
		receiptStorage,
		cfg.ParkingReceiptMaxBytes,
	)
	overtimeClaimHandler := handlers.NewOvertimeClaimHandler(
		overtimeClaimService,
	)
	claimReviewHandler := handlers.NewClaimReviewHandler(
		claimReviewService,
		receiptStorage,
	)
	reportHandler := handlers.NewReportHandler(reportService)

	mux.HandleFunc("GET /api/v1/health", healthHandler.Check)
	mux.HandleFunc("POST /api/v1/auth/login", authHandler.Login)

	mux.Handle(
		"GET /api/v1/auth/me",
		authenticated(
			tokenManager,
			http.HandlerFunc(authHandler.Me),
		),
	)

	mux.Handle(
		"POST /api/v1/auth/logout",
		authenticated(
			tokenManager,
			http.HandlerFunc(authHandler.Logout),
		),
	)

	mux.Handle(
		"GET /api/v1/admin/ping",
		authenticatedForRoles(
			tokenManager,
			[]string{models.RoleAdmin},
			http.HandlerFunc(accessHandler.Admin),
		),
	)

	mux.Handle(
		"GET /api/v1/employee/ping",
		authenticatedForRoles(
			tokenManager,
			[]string{models.RoleEmployee},
			http.HandlerFunc(accessHandler.Employee),
		),
	)

	mux.Handle(
		"POST /api/v1/parking-claims",
		authenticatedForRoles(
			tokenManager,
			[]string{models.RoleEmployee},
			http.HandlerFunc(parkingClaimHandler.Create),
		),
	)

	mux.Handle(
		"GET /api/v1/parking-claims",
		authenticatedForRoles(
			tokenManager,
			[]string{models.RoleEmployee},
			http.HandlerFunc(parkingClaimHandler.List),
		),
	)

	mux.Handle(
		"GET /api/v1/parking-claims/{id}",
		authenticatedForRoles(
			tokenManager,
			[]string{models.RoleEmployee},
			http.HandlerFunc(parkingClaimHandler.Detail),
		),
	)

	mux.Handle(
		"GET /api/v1/parking-claims/{id}/receipt",
		authenticatedForRoles(
			tokenManager,
			[]string{models.RoleEmployee},
			http.HandlerFunc(parkingClaimHandler.Receipt),
		),
	)

	mux.Handle(
		"POST /api/v1/overtime-claims",
		authenticatedForRoles(
			tokenManager,
			[]string{models.RoleEmployee},
			http.HandlerFunc(overtimeClaimHandler.Create),
		),
	)

	mux.Handle(
		"GET /api/v1/overtime-claims",
		authenticatedForRoles(
			tokenManager,
			[]string{models.RoleEmployee},
			http.HandlerFunc(overtimeClaimHandler.List),
		),
	)

	mux.Handle(
		"GET /api/v1/overtime-claims/{id}",
		authenticatedForRoles(
			tokenManager,
			[]string{models.RoleEmployee},
			http.HandlerFunc(overtimeClaimHandler.Detail),
		),
	)

	mux.Handle(
		"GET /api/v1/admin/claims",
		authenticatedForRoles(
			tokenManager,
			[]string{models.RoleAdmin},
			http.HandlerFunc(claimReviewHandler.List),
		),
	)

	mux.Handle(
		"GET /api/v1/admin/claims/{type}/{id}",
		authenticatedForRoles(
			tokenManager,
			[]string{models.RoleAdmin},
			http.HandlerFunc(claimReviewHandler.Detail),
		),
	)

	mux.Handle(
		"PATCH /api/v1/admin/claims/{type}/{id}/status",
		authenticatedForRoles(
			tokenManager,
			[]string{models.RoleAdmin},
			http.HandlerFunc(claimReviewHandler.Review),
		),
	)

	mux.Handle(
		"GET /api/v1/admin/claims/{type}/{id}/history",
		authenticatedForRoles(
			tokenManager,
			[]string{models.RoleAdmin},
			http.HandlerFunc(claimReviewHandler.History),
		),
	)

	mux.Handle(
		"GET /api/v1/admin/claims/{type}/{id}/receipt",
		authenticatedForRoles(
			tokenManager,
			[]string{models.RoleAdmin},
			http.HandlerFunc(claimReviewHandler.Receipt),
		),
	)

	mux.Handle(
		"GET /api/v1/admin/dashboard",
		authenticatedForRoles(
			tokenManager,
			[]string{models.RoleAdmin},
			http.HandlerFunc(reportHandler.Dashboard),
		),
	)

	mux.Handle(
		"GET /api/v1/admin/reports",
		authenticatedForRoles(
			tokenManager,
			[]string{models.RoleAdmin},
			http.HandlerFunc(reportHandler.List),
		),
	)

	mux.Handle(
		"GET /api/v1/admin/reports/export",
		authenticatedForRoles(
			tokenManager,
			[]string{models.RoleAdmin},
			http.HandlerFunc(reportHandler.Export),
		),
	)

	mux.HandleFunc("/", notFound)

	var handler http.Handler = mux
	handler = middleware.CORS(cfg.FrontendURL, handler)
	handler = middleware.SecurityHeaders(handler)
	handler = middleware.Logger(handler)
	handler = middleware.RequestID(handler)
	handler = middleware.Recovery(handler)

	return handler
}

func authenticated(
	tokenManager *security.TokenManager,
	next http.Handler,
) http.Handler {
	return middleware.Auth(tokenManager, next)
}

func authenticatedForRoles(
	tokenManager *security.TokenManager,
	roles []string,
	next http.Handler,
) http.Handler {
	return middleware.Auth(
		tokenManager,
		middleware.RequireRoles(roles...)(next),
	)
}

func notFound(w http.ResponseWriter, _ *http.Request) {
	responses.WriteJSON(
		w,
		http.StatusNotFound,
		responses.APIResponse{
			Success: false,
			Message: "Endpoint tidak ditemukan",
		},
	)
}
