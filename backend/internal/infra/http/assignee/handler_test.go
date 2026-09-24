package assignee_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	assigneeApp "github.com/adriancf/demo-sistema-chamado/backend/internal/application/assignee"
	domainAssignee "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/assignee"
	httpAssignee "github.com/adriancf/demo-sistema-chamado/backend/internal/infra/http/assignee"
	mockAssignee "github.com/adriancf/demo-sistema-chamado/backend/internal/util/test/mock/assignee"
	mockTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/util/test/mock/ticket"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupAssigneeTest(t *testing.T) (*mockAssignee.MockAssigneeRepository, *mockTicket.MockTicketRepository, http.Handler) {
	assigneeRepo := mockAssignee.NewMockAssigneeRepository()
	ticketRepo := mockTicket.NewMockTicketRepository()

	createUC := assigneeApp.NewCreateUseCase(assigneeRepo)
	listUC := assigneeApp.NewListUseCase(assigneeRepo, ticketRepo)
	handler := httpAssignee.NewHandler(createUC, listUC)

	r := chi.NewRouter()
	r.Route("/api/assignees", func(r chi.Router) {
		r.Get("/", handler.List)
		r.Post("/", handler.Create)
	})

	return assigneeRepo, ticketRepo, r
}

func TestAssigneeHandler(t *testing.T) {
	now := time.Now().UTC()
	carlosID, _ := uuid.NewV7()
	carlos := domainAssignee.Assignee{ID: carlosID, Name: "Carlos Silva", CreatedAt: now, UpdatedAt: now}
	anaID, _ := uuid.NewV7()
	ana := domainAssignee.Assignee{ID: anaID, Name: "Ana Souza", CreatedAt: now, UpdatedAt: now}
	marianaID, _ := uuid.NewV7()
	mariana := domainAssignee.Assignee{ID: marianaID, Name: "Mariana Oliveira", CreatedAt: now, UpdatedAt: now}

	t.Run("should list existing assignees with open tickets count", func(t *testing.T) {
		assigneeRepo, ticketRepo, router := setupAssigneeTest(t)

		assignees := []domainAssignee.Assignee{carlos, ana}
		openCounts := map[string]int{"Carlos Silva": 2, "Ana Souza": 1}

		assigneeRepo.On("ListAll", mock.Anything).Return(assignees, nil).Once()
		ticketRepo.On("CountOpenTicketsByAssignee", mock.Anything).Return(openCounts, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/assignees", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var list []assigneeApp.Output
		err := json.Unmarshal(rec.Body.Bytes(), &list)
		require.NoError(t, err)
		assert.Len(t, list, 2)
		assert.Equal(t, "Carlos Silva", list[0].Name)
		assert.Equal(t, 2, list[0].OpenTicketsCount)

		assigneeRepo.AssertExpectations(t)
		ticketRepo.AssertExpectations(t)
	})

	t.Run("should return 500 when list use case fails", func(t *testing.T) {
		assigneeRepo, _, router := setupAssigneeTest(t)

		assigneeRepo.On("ListAll", mock.Anything).Return(nil, assert.AnError).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/assignees", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assigneeRepo.AssertExpectations(t)
	})

	t.Run("should create new assignee via POST", func(t *testing.T) {
		assigneeRepo, _, router := setupAssigneeTest(t)

		assigneeRepo.On("FindByName", mock.Anything, "Mariana Oliveira").Return(nil, domainAssignee.ErrAssigneeNotFound).Once()
		assigneeRepo.On("Save", mock.Anything, mock.AnythingOfType("assignee.Assignee")).Return(nil).Once()

		body := map[string]any{"name": "Mariana Oliveira"}
		raw, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/assignees", bytes.NewBuffer(raw))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assigneeRepo.AssertExpectations(t)
	})

	t.Run("should return 400 on malformed json body", func(t *testing.T) {
		_, _, router := setupAssigneeTest(t)

		req := httptest.NewRequest(http.MethodPost, "/api/assignees", bytes.NewBufferString("{invalid-json"))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("should return 400 when assignee name is invalid", func(t *testing.T) {
		assigneeRepo, _, router := setupAssigneeTest(t)

		body := map[string]any{"name": ""}
		raw, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/assignees", bytes.NewBuffer(raw))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assigneeRepo.AssertNotCalled(t, "Save", mock.Anything, mock.Anything)
	})

	t.Run("should return 409 when creating duplicate assignee", func(t *testing.T) {
		assigneeRepo, _, router := setupAssigneeTest(t)

		assigneeRepo.On("FindByName", mock.Anything, "Mariana Oliveira").Return(&mariana, nil).Once()

		body := map[string]any{"name": "Mariana Oliveira"}
		raw, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/assignees", bytes.NewBuffer(raw))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusConflict, rec.Code)
		assigneeRepo.AssertExpectations(t)
	})
}
