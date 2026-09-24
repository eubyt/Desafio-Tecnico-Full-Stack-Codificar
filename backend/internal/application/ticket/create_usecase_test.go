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

func TestCreateUseCase(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()

	carlosID, _ := uuid.NewV7()
	carlos := domainAssignee.Assignee{ID: carlosID, Name: "Carlos Silva", CreatedAt: now, UpdatedAt: now}
	anaID, _ := uuid.NewV7()
	ana := domainAssignee.Assignee{ID: anaID, Name: "Ana Souza", CreatedAt: now, UpdatedAt: now}

	t.Run("should create ticket assigning manually to a valid assignee", func(t *testing.T) {
		ticketRepo := mockTicket.NewMockTicketRepository()
		assigneeRepo := mockAssignee.NewMockAssigneeRepository()

		assigneeRepo.On("FindByName", ctx, "Carlos Silva").Return(&carlos, nil).Once()
		ticketRepo.On("Save", ctx, mock.AnythingOfType("*ticket.Ticket")).Return(nil).Once()

		resolveUC := appTicket.NewResolveAssigneeUseCase(assigneeRepo)
		uc := appTicket.NewCreateUseCase(ticketRepo, resolveUC)

		input := appTicket.CreateInput{
			Title:       "Payment error",
			Description: "Communication failure with payment acquirer",
			Priority:    "high",
			Assignee:    "Carlos Silva",
			AutoAssign:  false,
		}

		out, err := uc.Execute(ctx, input)
		require.NoError(t, err)
		assert.NotEmpty(t, out.ID)
		assert.Equal(t, "Payment error", out.Title)
		assert.Equal(t, "high", out.Priority)
		assert.Equal(t, "open", out.Status)
		assert.Equal(t, "Carlos Silva", out.Assignee)

		ticketRepo.AssertExpectations(t)
		assigneeRepo.AssertExpectations(t)
	})

	t.Run("should fail when manually assigning to a non-existent assignee", func(t *testing.T) {
		ticketRepo := mockTicket.NewMockTicketRepository()
		assigneeRepo := mockAssignee.NewMockAssigneeRepository()

		assigneeRepo.On("FindByName", ctx, "Non Existent").Return(nil, domainAssignee.ErrAssigneeNotFound).Once()

		resolveUC := appTicket.NewResolveAssigneeUseCase(assigneeRepo)
		uc := appTicket.NewCreateUseCase(ticketRepo, resolveUC)

		input := appTicket.CreateInput{
			Title:       "Payment error",
			Description: "Communication failure with payment acquirer",
			Priority:    "high",
			Assignee:    "Non Existent",
			AutoAssign:  false,
		}

		_, err := uc.Execute(ctx, input)
		require.Error(t, err)
		assert.Equal(t, domainAssignee.ErrAssigneeNotFound, err)

		assigneeRepo.AssertExpectations(t)
	})

	t.Run("should create ticket with auto-assign to the least loaded assignee", func(t *testing.T) {
		ticketRepo := mockTicket.NewMockTicketRepository()
		assigneeRepo := mockAssignee.NewMockAssigneeRepository()

		assigneeRepo.On("FindLeastLoaded", ctx).Return(&ana, nil).Once()
		ticketRepo.On("Save", ctx, mock.AnythingOfType("*ticket.Ticket")).Return(nil).Once()

		resolveUC := appTicket.NewResolveAssigneeUseCase(assigneeRepo)
		uc := appTicket.NewCreateUseCase(ticketRepo, resolveUC)

		input := appTicket.CreateInput{
			Title:       "Question about invoice",
			Description: "How to issue supplementary return invoice?",
			Priority:    "low",
			AutoAssign:  true,
		}

		out, err := uc.Execute(ctx, input)
		require.NoError(t, err)
		assert.Equal(t, "Ana Souza", out.Assignee)

		ticketRepo.AssertExpectations(t)
		assigneeRepo.AssertExpectations(t)
	})

	t.Run("should fail if auto_assign is false and no assignee is provided", func(t *testing.T) {
		ticketRepo := mockTicket.NewMockTicketRepository()
		assigneeRepo := mockAssignee.NewMockAssigneeRepository()

		resolveUC := appTicket.NewResolveAssigneeUseCase(assigneeRepo)
		uc := appTicket.NewCreateUseCase(ticketRepo, resolveUC)

		input := appTicket.CreateInput{
			Title:       "Question about invoice",
			Description: "How to issue supplementary return invoice?",
			Priority:    "low",
			Assignee:    "",
			AutoAssign:  false,
		}

		_, err := uc.Execute(ctx, input)
		require.Error(t, err)
		assert.Equal(t, domainAssignee.ErrInvalidAssignee, err)
	})

	t.Run("should fail when creating ticket with invalid title", func(t *testing.T) {
		ticketRepo := mockTicket.NewMockTicketRepository()
		assigneeRepo := mockAssignee.NewMockAssigneeRepository()

		assigneeRepo.On("FindByName", ctx, "Carlos Silva").Return(&carlos, nil).Once()

		resolveUC := appTicket.NewResolveAssigneeUseCase(assigneeRepo)
		uc := appTicket.NewCreateUseCase(ticketRepo, resolveUC)

		input := appTicket.CreateInput{
			Title:       "ab",
			Description: "Valid description with sufficient length",
			Priority:    "low",
			Assignee:    "Carlos Silva",
		}

		_, err := uc.Execute(ctx, input)
		require.ErrorIs(t, err, domainTicket.ErrInvalidTitle)
		assigneeRepo.AssertExpectations(t)
	})

	t.Run("should fail when creating ticket with invalid description", func(t *testing.T) {
		ticketRepo := mockTicket.NewMockTicketRepository()
		assigneeRepo := mockAssignee.NewMockAssigneeRepository()

		assigneeRepo.On("FindByName", ctx, "Carlos Silva").Return(&carlos, nil).Once()

		resolveUC := appTicket.NewResolveAssigneeUseCase(assigneeRepo)
		uc := appTicket.NewCreateUseCase(ticketRepo, resolveUC)

		input := appTicket.CreateInput{
			Title:       "Valid Title",
			Description: "ops",
			Priority:    "low",
			Assignee:    "Carlos Silva",
		}

		_, err := uc.Execute(ctx, input)
		require.ErrorIs(t, err, domainTicket.ErrInvalidDescription)
		assigneeRepo.AssertExpectations(t)
	})

	t.Run("should fail when ticket repository returns error on Save", func(t *testing.T) {
		ticketRepo := mockTicket.NewMockTicketRepository()
		assigneeRepo := mockAssignee.NewMockAssigneeRepository()

		assigneeRepo.On("FindByName", ctx, "Carlos Silva").Return(&carlos, nil).Once()
		ticketRepo.On("Save", ctx, mock.AnythingOfType("*ticket.Ticket")).Return(errors.New("db error")).Once()

		resolveUC := appTicket.NewResolveAssigneeUseCase(assigneeRepo)
		uc := appTicket.NewCreateUseCase(ticketRepo, resolveUC)

		input := appTicket.CreateInput{
			Title:       "Payment error",
			Description: "Communication failure with payment acquirer",
			Priority:    "high",
			Assignee:    "Carlos Silva",
		}

		_, err := uc.Execute(ctx, input)
		require.Error(t, err)
		assert.Equal(t, "db error", err.Error())
		ticketRepo.AssertExpectations(t)
		assigneeRepo.AssertExpectations(t)
	})

	t.Run("should fail when creating ticket with invalid priority", func(t *testing.T) {
		ticketRepo := mockTicket.NewMockTicketRepository()
		assigneeRepo := mockAssignee.NewMockAssigneeRepository()

		resolveUC := appTicket.NewResolveAssigneeUseCase(assigneeRepo)
		uc := appTicket.NewCreateUseCase(ticketRepo, resolveUC)

		input := appTicket.CreateInput{
			Title:       "Valid Title",
			Description: "Valid description with sufficient length",
			Priority:    "invalid_priority",
			Assignee:    "Carlos Silva",
		}

		_, err := uc.Execute(ctx, input)
		require.ErrorIs(t, err, domainTicket.ErrInvalidPriority)
	})

	t.Run("should fail if assignee resolution fails", func(t *testing.T) {
		ticketRepo := mockTicket.NewMockTicketRepository()
		assigneeRepo := mockAssignee.NewMockAssigneeRepository()

		assigneeRepo.On("FindLeastLoaded", ctx).Return(nil, domainAssignee.ErrNoAssigneesAvailable).Once()

		resolveUC := appTicket.NewResolveAssigneeUseCase(assigneeRepo)
		uc := appTicket.NewCreateUseCase(ticketRepo, resolveUC)

		input := appTicket.CreateInput{
			Title:       "Valid Title",
			Description: "Valid description with sufficient length",
			Priority:    "low",
			AutoAssign:  true,
		}

		_, err := uc.Execute(ctx, input)
		require.ErrorIs(t, err, domainAssignee.ErrNoAssigneesAvailable)
		assigneeRepo.AssertExpectations(t)
	})
}
