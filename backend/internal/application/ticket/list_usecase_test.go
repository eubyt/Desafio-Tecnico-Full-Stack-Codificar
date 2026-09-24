package ticket_test

import (
	"context"
	"testing"
	"time"

	appTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/application/ticket"
	domainAssignee "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/assignee"
	domainTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/ticket"
	domainPagination "github.com/adriancf/demo-sistema-chamado/backend/internal/util/pagination"
	mockTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/util/test/mock/ticket"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestListUseCase(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()

	carlosID, _ := uuid.NewV7()
	carlos := domainAssignee.Assignee{ID: carlosID, Name: "Carlos Silva", CreatedAt: now, UpdatedAt: now}

	t1ID, _ := uuid.NewV7()
	t2ID, _ := uuid.NewV7()
	tk1 := domainTicket.NewTicket(t1ID, "Ticket 1", "Description 12345", domainTicket.PriorityHigh, carlos, now)
	tk2 := domainTicket.NewTicket(t2ID, "Ticket 2", "Description 12345", domainTicket.PriorityLow, carlos, now)

	t.Run("should list tickets with pagination", func(t *testing.T) {
		ticketRepo := mockTicket.NewMockTicketRepository()
		uc := appTicket.NewListUseCase(ticketRepo)

		expectedResult := domainPagination.PaginatedResult[*domainTicket.Ticket]{
			Items:       []*domainTicket.Ticket{tk1, tk2},
			TotalItems:  2,
			TotalPages:  1,
			CurrentPage: 1,
			PageSize:    10,
		}
		ticketRepo.On("FindAll", ctx, mock.Anything, mock.Anything, mock.Anything).Return(expectedResult, nil).Once()

		status := "open"
		priority := "high"
		out, err := uc.Execute(ctx, appTicket.ListInput{
			Page:     1,
			PageSize: 10,
			Status:   &status,
			Priority: &priority,
			SortBy:   "priority",
			Order:    "asc",
		})

		require.NoError(t, err)
		assert.Equal(t, 2, out.TotalItems)
		assert.Equal(t, 1, out.CurrentPage)
		assert.Len(t, out.Items, 2)
		assert.Equal(t, "Ticket 1", out.Items[0].Title)
		ticketRepo.AssertExpectations(t)
	})

	t.Run("should list tickets with default sort and status", func(t *testing.T) {
		ticketRepo := mockTicket.NewMockTicketRepository()
		uc := appTicket.NewListUseCase(ticketRepo)

		expectedResult := domainPagination.PaginatedResult[*domainTicket.Ticket]{
			Items:       []*domainTicket.Ticket{tk1},
			TotalItems:  1,
			TotalPages:  1,
			CurrentPage: 1,
			PageSize:    10,
		}
		ticketRepo.On("FindAll", ctx, mock.Anything, mock.Anything, mock.Anything).Return(expectedResult, nil).Once()

		out, err := uc.Execute(ctx, appTicket.ListInput{
			SortBy: "status",
			Order:  "desc",
		})

		require.NoError(t, err)
		assert.Equal(t, 1, out.TotalItems)
		ticketRepo.AssertExpectations(t)
	})

	t.Run("should list tickets sorted by title", func(t *testing.T) {
		ticketRepo := mockTicket.NewMockTicketRepository()
		uc := appTicket.NewListUseCase(ticketRepo)

		expectedResult := domainPagination.PaginatedResult[*domainTicket.Ticket]{
			Items:       []*domainTicket.Ticket{},
			TotalItems:  0,
			TotalPages:  0,
			CurrentPage: 1,
			PageSize:    10,
		}
		ticketRepo.On("FindAll", ctx, mock.Anything, mock.Anything, mock.Anything).Return(expectedResult, nil).Once()

		out, err := uc.Execute(ctx, appTicket.ListInput{
			SortBy: "title",
		})

		require.NoError(t, err)
		assert.Equal(t, 0, out.TotalItems)
		ticketRepo.AssertExpectations(t)
	})

	t.Run("should return error for invalid status in filter", func(t *testing.T) {
		ticketRepo := mockTicket.NewMockTicketRepository()
		uc := appTicket.NewListUseCase(ticketRepo)

		invalidStatus := "invalido"
		_, err := uc.Execute(ctx, appTicket.ListInput{
			Status: &invalidStatus,
		})

		require.Error(t, err)
	})

	t.Run("should return error for invalid priority in filter", func(t *testing.T) {
		ticketRepo := mockTicket.NewMockTicketRepository()
		uc := appTicket.NewListUseCase(ticketRepo)

		invalidPriority := "invalida"
		_, err := uc.Execute(ctx, appTicket.ListInput{
			Priority: &invalidPriority,
		})

		require.Error(t, err)
	})

	t.Run("should return error when FindAll fails in repository", func(t *testing.T) {
		ticketRepo := mockTicket.NewMockTicketRepository()
		uc := appTicket.NewListUseCase(ticketRepo)

		ticketRepo.On("FindAll", ctx, mock.Anything, mock.Anything, mock.Anything).Return(domainPagination.PaginatedResult[*domainTicket.Ticket]{}, assert.AnError).Once()

		_, err := uc.Execute(ctx, appTicket.ListInput{})
		require.Error(t, err)
		ticketRepo.AssertExpectations(t)
	})
}
