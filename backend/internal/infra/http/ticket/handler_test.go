package ticket_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	ticketApp "github.com/adriancf/demo-sistema-chamado/backend/internal/application/ticket"
	domainAssignee "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/assignee"
	domainTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/ticket"
	httpTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/infra/http/ticket"
	"github.com/adriancf/demo-sistema-chamado/backend/internal/util/pagination"
	mockAssignee "github.com/adriancf/demo-sistema-chamado/backend/internal/util/test/mock/assignee"
	mockTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/util/test/mock/ticket"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func setupTicketTest(t *testing.T) (*mockTicket.MockTicketRepository, *mockAssignee.MockAssigneeRepository, http.Handler) {
	ticketRepo := mockTicket.NewMockTicketRepository()
	assigneeRepo := mockAssignee.NewMockAssigneeRepository()

	resolveAssigneeUC := ticketApp.NewResolveAssigneeUseCase(assigneeRepo)
	ticketCreateUC := ticketApp.NewCreateUseCase(ticketRepo, resolveAssigneeUC)
	ticketUpdateUC := ticketApp.NewUpdateUseCase(ticketRepo, resolveAssigneeUC)
	ticketGetByIDUC := ticketApp.NewGetByIDUseCase(ticketRepo)
	ticketListUC := ticketApp.NewListUseCase(ticketRepo)

	handler := httpTicket.NewHandler(ticketCreateUC, ticketUpdateUC, ticketGetByIDUC, ticketListUC)

	r := chi.NewRouter()
	r.Route("/api/tickets", func(r chi.Router) {
		r.Post("/", handler.Create)
		r.Get("/", handler.List)
		r.Get("/{id}", handler.GetByID)
		r.Put("/{id}", handler.Update)
	})

	return ticketRepo, assigneeRepo, r
}

func TestTicketHandler_CreateTicket(t *testing.T) {
	now := time.Now().UTC()
	carlosID, _ := uuid.NewV7()
	carlos := domainAssignee.Assignee{ID: carlosID, Name: "Carlos Silva", CreatedAt: now, UpdatedAt: now}
	anaID, _ := uuid.NewV7()
	ana := domainAssignee.Assignee{ID: anaID, Name: "Ana Souza", CreatedAt: now, UpdatedAt: now}

	t.Run("should create ticket via POST with manual assignee", func(t *testing.T) {
		ticketRepo, assigneeRepo, router := setupTicketTest(t)

		assigneeRepo.On("FindByName", mock.Anything, "Carlos Silva").Return(&carlos, nil).Once()
		ticketRepo.On("Save", mock.Anything, mock.AnythingOfType("*ticket.Ticket")).Return(nil).Once()

		body := map[string]any{
			"title":       "Store checkout error",
			"description": "Clicking pay closes modal without feedback",
			"priority":    "high",
			"assignee":    "Carlos Silva",
			"auto_assign": false,
		}
		raw, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/tickets", bytes.NewBuffer(raw))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var res ticketApp.Output
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		require.NoError(t, err)
		assert.NotEmpty(t, res.ID)
		assert.Equal(t, "Store checkout error", res.Title)
		assert.Equal(t, "high", res.Priority)
		assert.Equal(t, "open", res.Status)
		assert.Equal(t, "Carlos Silva", res.Assignee)

		ticketRepo.AssertExpectations(t)
		assigneeRepo.AssertExpectations(t)
	})

	t.Run("should create ticket via POST with auto assignment", func(t *testing.T) {
		ticketRepo, assigneeRepo, router := setupTicketTest(t)

		assigneeRepo.On("FindLeastLoaded", mock.Anything).Return(&ana, nil).Once()
		ticketRepo.On("Save", mock.Anything, mock.AnythingOfType("*ticket.Ticket")).Return(nil).Once()

		body := map[string]any{
			"title":       "Password reset support",
			"description": "User is not receiving verification email",
			"priority":    "medium",
			"auto_assign": true,
		}
		raw, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/tickets", bytes.NewBuffer(raw))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusCreated, rec.Code)

		var res ticketApp.Output
		err := json.Unmarshal(rec.Body.Bytes(), &res)
		require.NoError(t, err)
		assert.NotEmpty(t, res.ID)
		assert.Equal(t, "Ana Souza", res.Assignee)

		ticketRepo.AssertExpectations(t)
		assigneeRepo.AssertExpectations(t)
	})

	t.Run("should return 400 when body payload is invalid JSON", func(t *testing.T) {
		_, _, router := setupTicketTest(t)

		req := httptest.NewRequest(http.MethodPost, "/api/tickets", bytes.NewBufferString("{invalid-json"))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("should return 400 for invalid domain validation data", func(t *testing.T) {
		_, _, router := setupTicketTest(t)

		body := map[string]any{
			"title":       "Hi",
			"description": "short",
			"priority":    "invalid",
		}
		raw, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/api/tickets", bytes.NewBuffer(raw))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestTicketHandler_ListAndGet(t *testing.T) {
	now := time.Now().UTC()
	carlosID, _ := uuid.NewV7()
	carlos := domainAssignee.Assignee{ID: carlosID, Name: "Carlos Silva", CreatedAt: now, UpdatedAt: now}
	ticketID, _ := uuid.NewV7()
	testTicket := domainTicket.NewTicket(ticketID, "Slow admin dashboard loading", "Dashboard charts take more than 10s to load", domainTicket.PriorityLow, carlos, now)

	t.Run("should list tickets with pagination and filters", func(t *testing.T) {
		ticketRepo, _, router := setupTicketTest(t)

		paginatedResult := pagination.NewPaginatedResult([]*domainTicket.Ticket{testTicket}, 1, 1, 10)
		ticketRepo.On("FindAll", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(paginatedResult, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/tickets?page=1&page_size=10&sort_by=created_at&order=desc&status=open&priority=low&assignee=Carlos%20Silva&search=dashboard", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var listRes ticketApp.PaginatedOutput
		err := json.Unmarshal(rec.Body.Bytes(), &listRes)
		require.NoError(t, err)
		assert.Equal(t, 1, listRes.TotalItems)
		assert.Equal(t, 1, listRes.CurrentPage)
		assert.Len(t, listRes.Items, 1)

		ticketRepo.AssertExpectations(t)
	})

	t.Run("should get ticket by ID", func(t *testing.T) {
		ticketRepo, _, router := setupTicketTest(t)

		ticketRepo.On("FindByID", mock.Anything, ticketID.String()).Return(testTicket, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/tickets/"+ticketID.String(), nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var found ticketApp.Output
		err := json.Unmarshal(rec.Body.Bytes(), &found)
		require.NoError(t, err)
		assert.Equal(t, ticketID.String(), found.ID)

		ticketRepo.AssertExpectations(t)
	})

	t.Run("should return 404 for non-existent ticket ID", func(t *testing.T) {
		ticketRepo, _, router := setupTicketTest(t)

		ticketRepo.On("FindByID", mock.Anything, "non-existent-id").Return(nil, domainTicket.ErrTicketNotFound).Once()

		req := httptest.NewRequest(http.MethodGet, "/api/tickets/non-existent-id", nil)
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusNotFound, rec.Code)
		ticketRepo.AssertExpectations(t)
	})
}

func TestTicketHandler_UpdateTicket(t *testing.T) {
	now := time.Now().UTC()
	carlosID, _ := uuid.NewV7()
	carlos := domainAssignee.Assignee{ID: carlosID, Name: "Carlos Silva", CreatedAt: now, UpdatedAt: now}
	anaID, _ := uuid.NewV7()
	ana := domainAssignee.Assignee{ID: anaID, Name: "Ana Souza", CreatedAt: now, UpdatedAt: now}
	ticketID, _ := uuid.NewV7()

	newTicket := func() *domainTicket.Ticket {
		return domainTicket.NewTicket(ticketID, "Ticket to be updated", "Original description of the ticket", domainTicket.PriorityLow, carlos, now)
	}

	t.Run("should update ticket via PUT with manual assignee", func(t *testing.T) {
		ticketRepo, assigneeRepo, router := setupTicketTest(t)

		ticketRepo.On("FindByID", mock.Anything, ticketID.String()).Return(newTicket(), nil).Once()
		assigneeRepo.On("FindByName", mock.Anything, "Ana Souza").Return(&ana, nil).Once()
		ticketRepo.On("Update", mock.Anything, mock.AnythingOfType("*ticket.Ticket")).Return(nil).Once()

		updateBody := map[string]any{
			"status":   "in_progress",
			"assignee": "Ana Souza",
		}
		updateRaw, _ := json.Marshal(updateBody)
		req := httptest.NewRequest(http.MethodPut, "/api/tickets/"+ticketID.String(), bytes.NewBuffer(updateRaw))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var updated ticketApp.Output
		err := json.Unmarshal(rec.Body.Bytes(), &updated)
		require.NoError(t, err)
		assert.Equal(t, "Ticket to be updated", updated.Title, "title must remain unchanged")
		assert.Equal(t, "low", updated.Priority, "priority must remain unchanged")
		assert.Equal(t, "in_progress", updated.Status)
		assert.Equal(t, "Ana Souza", updated.Assignee)

		ticketRepo.AssertExpectations(t)
		assigneeRepo.AssertExpectations(t)
	})

	t.Run("should return 400 on malformed json body for update", func(t *testing.T) {
		_, _, router := setupTicketTest(t)

		req := httptest.NewRequest(http.MethodPut, "/api/tickets/"+ticketID.String(), bytes.NewBufferString("{invalid-json"))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("should update ticket via PUT with auto assignment", func(t *testing.T) {
		ticketRepo, assigneeRepo, router := setupTicketTest(t)

		ticketRepo.On("FindByID", mock.Anything, ticketID.String()).Return(newTicket(), nil).Once()
		assigneeRepo.On("FindLeastLoaded", mock.Anything).Return(&ana, nil).Once()
		ticketRepo.On("Update", mock.Anything, mock.AnythingOfType("*ticket.Ticket")).Return(nil).Once()

		updateBody := map[string]any{
			"status":      "in_progress",
			"auto_assign": true,
		}
		updateRaw, _ := json.Marshal(updateBody)
		req := httptest.NewRequest(http.MethodPut, "/api/tickets/"+ticketID.String(), bytes.NewBuffer(updateRaw))
		rec := httptest.NewRecorder()

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusOK, rec.Code)

		var updated ticketApp.Output
		err := json.Unmarshal(rec.Body.Bytes(), &updated)
		require.NoError(t, err)
		assert.Equal(t, "Ticket to be updated", updated.Title, "title must remain unchanged")
		assert.Equal(t, "Ana Souza", updated.Assignee)

		ticketRepo.AssertExpectations(t)
		assigneeRepo.AssertExpectations(t)
	})
}
