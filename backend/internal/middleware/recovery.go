package middleware

import (
	"fmt"
	"log"
	"net/http"

	"jedi-reimbursement-system/backend/internal/responses"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if recovered := recover(); recovered != nil {
				log.Printf(
					"request_id=%s panic=%s",
					RequestIDFromContext(r.Context()),
					fmt.Sprint(recovered),
				)

				responses.WriteJSON(
					w,
					http.StatusInternalServerError,
					responses.APIResponse{
						Success: false,
						Message: "Terjadi kesalahan pada server",
					},
				)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
