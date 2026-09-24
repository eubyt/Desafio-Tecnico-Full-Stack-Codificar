package ticket_test

import (
	"context"
	"errors"
	"testing"
	"time"

	appTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/application/ticket"
	domainAssignee "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/assignee"
	domainTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/ticket"
	mockAssignee "github.com/adriancf/demo-sistema-chamado/backend/internal/util/test/mock/assignee"
	mockTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/util/test/mock/ticket"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestUpdateUseCase(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()

	carlosID, _ := uuid.NewV7()
	carlos := domainAssignee.Assignee{ID: carlosID, Name: "Carlos Silva", CreatedAt: now, UpdatedAt: now}
	anaID, _ := uuid.NewV7()
	ana := domainAssignee.Assignee{ID: anaID, Name: "Ana Souza", CreatedAt: now, UpdatedAt: now}

	newTestTicket := func() (*domainTicket.Ticket, uuid.UUID) {
		id, _ := uuid.NewV7()
		ticket := domainTicket.NewTicket(id, "Original title", "Description with enough length", domainTicket.PriorityLow, carlos, now)
		return ticket, id
	}

	t.Run("should update ticket status and assignee successfully while keeping content intact", func(t *testing.T) {
		ticketRepo := mockTicket.NewMockTicketRepository()
		assigneeRepo := mockAssignee.NewMockAssigneeRepository()

		ticket, id := newTestTicket()
		ticketRepo.On("FindByID", ctx, id.String()).Return(ticket, nil).Once()
		assigneeRepo.On("FindByName", ctx, "Ana Souza").Return(&ana, nil).Once()
		ticketRepo.On("Update", ctx, mock.AnythingOfType("*ticket.Ticket")).Return(nil).Once()

		resolveUC := appTicket.NewResolveAssigneeUseCase(assigneeRepo)
		uc := appTicket.NewUpdateUseCase(ticketRepo, resolveUC)

		input := appTicket.UpdateInput{
			Status:   "in_progress",
			Assignee: "Ana Souza",
		}

		out, err := uc.Execute(ctx, id.String(), input)
		require.NoError(t, err)
		assert.Equal(t, "Original title", out.Title, "title must remain untouched")
		assert.Equal(t, "Description with enough length", out.Description, "description must remain untouched")
		assert.Equal(t, "low", out.Priority, "priority must remain untouched")
		assert.Equal(t, "in_progress", out.Status)
		assert.Equal(t, "Ana Souza", out.Assignee)

		ticketRepo.AssertExpectations(t)
		assigneeRepo.AssertExpectations(t)
	})

	t.Run("should update ticket with auto-assign (auto_assign: true)", func(t *testing.T) {
		ticketRepo := mockTicket.NewMockTicketRepository()
		assigneeRepo := mockAssignee.NewMockAssigneeRepository()

		ticket, id := newTestTicket()
		ticketRepo.On("FindByID", ctx, id.String()).Return(ticket, nil).Once()
		assigneeRepo.On("FindLeastLoaded", ctx).Return(&ana, nil).Once()
		ticketRepo.On("Update", ctx, mock.AnythingOfType("*ticket.Ticket")).Return(nil).Once()

		resolveUC := appTicket.NewResolveAssigneeUseCase(assigneeRepo)
		uc := appTicket.NewUpdateUseCase(ticketRepo, resolveUC)

		input := appTicket.UpdateInput{
			Status:     "in_progress",
			AutoAssign: true,
		}

		out, err := uc.Execute(ctx, id.String(), input)
		require.NoError(t, err)
		assert.Equal(t, "Ana Souza", out.Assignee)
		assert.Equal(t, "Original title", out.Title)

		ticketRepo.AssertExpectations(t)
		assigneeRepo.AssertExpectations(t)
	})

	t.Run("should fail when neither status nor assignee is provided", func(t *testing.T) {
		ticketRepo := mockTicket.NewMockTicketRepository()
		assigneeRepo := mockAssignee.NewMockAssigneeRepository()

		ticket, id := newTestTicket()
		ticketRepo.On("FindByID", ctx, id.String()).Return(ticket, nil).Once()

		resolveUC := appTicket.NewResolveAssigneeUseCase(assigneeRepo)
		uc := appTicket.NewUpdateUseCase(ticketRepo, resolveUC)

		input := appTicket.UpdateInput{}

		_, err := uc.Execute(ctx, id.String(), input)
		require.Error(t, err)
		assert.Equal(t, domainTicket.ErrInvalidStatus, err)
	})

	t.Run("should fail with invalid status", func(t *testing.T) {
		ticketRepo := mockTicket.NewMockTicketRepository()
		assigneeRepo := mockAssignee.NewMockAssigneeRepository()

		ticket, id := newTestTicket()
		ticketRepo.On("FindByID", ctx, id.String()).Return(ticket, nil).Once()

		resolveUC := appTicket.NewResolveAssigneeUseCase(assigneeRepo)
		uc := appTicket.NewUpdateUseCase(ticketRepo, resolveUC)

		input := appTicket.UpdateInput{
			Status:   "invalid_status",
			Assignee: "Ana Souza",
		}

		_, err := uc.Execute(ctx, id.String(), input)
		require.Error(t, err)
		assert.Equal(t, domainTicket.ErrInvalidStatus, err)
	})

	t.Run("should fail when updating non-existent ticket", func(t *testing.T) {
		ticketRepo := mockTicket.NewMockTicketRepository()
		assigneeRepo := mockAssignee.NewMockAssigneeRepository()

		ticketRepo.On("FindByID", ctx, "inexistente").Return(nil, domainTicket.ErrTicketNotFound).Once()

		resolveUC := appTicket.NewResolveAssigneeUseCase(assigneeRepo)
		uc := appTicket.NewUpdateUseCase(ticketRepo, resolveUC)

		input := appTicket.UpdateInput{
			Status:   "in_progress",
			Assignee: "Ana Souza",
		}

		_, err := uc.Execute(ctx, "inexistente", input)
		require.Error(t, err)
		assert.Equal(t, domainTicket.ErrTicketNotFound, err)
	})

	t.Run("should fail when repository returns error on Update", func(t *testing.T) {
		ticketRepo := mockTicket.NewMockTicketRepository()
		assigneeRepo := mockAssignee.NewMockAssigneeRepository()

		ticket, id := newTestTicket()
		ticketRepo.On("FindByID", ctx, id.String()).Return(ticket, nil).Once()
		assigneeRepo.On("FindByName", ctx, "Ana Souza").Return(&ana, nil).Once()
		ticketRepo.On("Update", ctx, mock.AnythingOfType("*ticket.Ticket")).Return(errors.New("db update error")).Once()

		resolveUC := appTicket.NewResolveAssigneeUseCase(assigneeRepo)
		uc := appTicket.NewUpdateUseCase(ticketRepo, resolveUC)

		input := appTicket.UpdateInput{
			Status:   "in_progress",
			Assignee: "Ana Souza",
		}

		_, err := uc.Execute(ctx, id.String(), input)
		require.Error(t, err)
		assert.Equal(t, "db update error", err.Error())
	})

	t.Run("should fail when assignee resolution fails on update", func(t *testing.T) {
		ticketRepo := mockTicket.NewMockTicketRepository()
		assigneeRepo := mockAssignee.NewMockAssigneeRepository()

		ticket, id := newTestTicket()
		ticketRepo.On("FindByID", ctx, id.String()).Return(ticket, nil).Once()
		assigneeRepo.On("FindLeastLoaded", ctx).Return(nil, domainAssignee.ErrNoAssigneesAvailable).Once()

		resolveUC := appTicket.NewResolveAssigneeUseCase(assigneeRepo)
		uc := appTicket.NewUpdateUseCase(ticketRepo, resolveUC)

		input := appTicket.UpdateInput{
			Status:     "in_progress",
			AutoAssign: true,
		}

		_, err := uc.Execute(ctx, id.String(), input)
		require.ErrorIs(t, err, domainAssignee.ErrNoAssigneesAvailable)
	})
}
