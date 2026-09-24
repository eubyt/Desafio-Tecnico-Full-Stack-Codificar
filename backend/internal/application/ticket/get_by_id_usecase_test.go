package ticket_test

import (
	"context"
	"testing"
	"time"

	appTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/application/ticket"
	domainAssignee "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/assignee"
	domainTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/ticket"
	mockTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/util/test/mock/ticket"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetByIDUseCase(t *testing.T) {
	ctx := context.Background()
	ticketID, err := uuid.NewV7()
	require.NoError(t, err)

	carlosID, _ := uuid.NewV7()
	carlos := domainAssignee.Assignee{ID: carlosID, Name: "Carlos Silva", CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC()}
	tk := domainTicket.NewTicket(ticketID, "Login problem", "Cannot recover password", domainTicket.PriorityHigh, carlos, time.Now().UTC())

	t.Run("should return ticket when ID exists", func(t *testing.T) {
		ticketRepo := mockTicket.NewMockTicketRepository()
		ticketRepo.On("FindByID", ctx, ticketID.String()).Return(tk, nil)

		uc := appTicket.NewGetByIDUseCase(ticketRepo)
		out, err := uc.Execute(ctx, ticketID.String())

		require.NoError(t, err)
		assert.Equal(t, ticketID.String(), out.ID)
		assert.Equal(t, "Login problem", out.Title)
		assert.Equal(t, "Carlos Silva", out.Assignee)
		ticketRepo.AssertExpectations(t)
	})

	t.Run("should return error when ID does not exist", func(t *testing.T) {
		ticketRepo := mockTicket.NewMockTicketRepository()
		ticketRepo.On("FindByID", ctx, "t-999").Return(nil, domainTicket.ErrTicketNotFound)

		uc := appTicket.NewGetByIDUseCase(ticketRepo)
		_, err := uc.Execute(ctx, "t-999")

		require.ErrorIs(t, err, domainTicket.ErrTicketNotFound)
		ticketRepo.AssertExpectations(t)
	})
}
