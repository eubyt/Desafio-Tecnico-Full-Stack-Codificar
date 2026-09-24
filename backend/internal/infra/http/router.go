package http

import (
	"net/http"

	httpAssignee "github.com/adriancf/demo-sistema-chamado/backend/internal/infra/http/assignee"
	httpResponse "github.com/adriancf/demo-sistema-chamado/backend/internal/infra/http/response"
	httpTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/infra/http/ticket"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Re-export common HTTP response helpers for convenience
var (
	JSON         = httpResponse.JSON
	RespondError = httpResponse.RespondError
)

type ErrorResponse = httpResponse.ErrorResponse

func NewRouter(ticketHandler *httpTicket.Handler, assigneeHandler *httpAssignee.Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Configuração de CORS
	// NOTE: Não use * em produção, especifique os domínios permitidos
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/health", HealthHandler)

	r.Route("/api", func(r chi.Router) {
		r.Route("/tickets", func(r chi.Router) {
			r.Post("/", ticketHandler.Create)
			r.Get("/", ticketHandler.List)
			r.Get("/{id}", ticketHandler.GetByID)
			r.Put("/{id}", ticketHandler.Update)
		})

		r.Route("/assignees", func(r chi.Router) {
			r.Get("/", assigneeHandler.List)
			r.Post("/", assigneeHandler.Create)
		})
	})

	return r
}

// HealthHandler godoc
// @Summary      Checagem de integridade da API
// @Description  Retorna o status de integridade do serviço.
// @Tags         Health
// @Produce      json
// @Success      200  {object}  map[string]string
// @Router       /health [get]
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	JSON(w, http.StatusOK, map[string]string{
		"status": "ok",
	})
}
