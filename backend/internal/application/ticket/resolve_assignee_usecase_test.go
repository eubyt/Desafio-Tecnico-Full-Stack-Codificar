package ticket_test

import (
	"context"
	"testing"
	"time"

	appTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/application/ticket"
	domainAssignee "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/assignee"
	domainTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/ticket"
	mockAssignee "github.com/adriancf/demo-sistema-chamado/backend/internal/util/test/mock/assignee"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResolveAssigneeUseCase(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()

	carlosID, _ := uuid.NewV7()
	carlos := domainAssignee.Assignee{ID: carlosID, Name: "Carlos Silva", CreatedAt: now, UpdatedAt: now}
	anaID, _ := uuid.NewV7()
	ana := domainAssignee.Assignee{ID: anaID, Name: "Ana Souza", CreatedAt: now, UpdatedAt: now}

	t.Run("should automatically select the least loaded assignee", func(t *testing.T) {
		repo := mockAssignee.NewMockAssigneeRepository()
		repo.On("FindLeastLoaded", ctx).Return(&ana, nil)

		uc := appTicket.NewResolveAssigneeUseCase(repo)
		selected, err := uc.Execute(ctx, "", true)

		require.NoError(t, err)
		assert.Equal(t, "Ana Souza", selected.Name)
		repo.AssertExpectations(t)
	})

	t.Run("should find assignee by name when autoAssign is false", func(t *testing.T) {
		repo := mockAssignee.NewMockAssigneeRepository()
		repo.On("FindByName", ctx, "Carlos Silva").Return(&carlos, nil)

		uc := appTicket.NewResolveAssigneeUseCase(repo)
		selected, err := uc.Execute(ctx, "Carlos Silva", false)

		require.NoError(t, err)
		assert.Equal(t, "Carlos Silva", selected.Name)
		repo.AssertExpectations(t)
	})

	t.Run("should return error when autoAssign is false and assignee name is empty", func(t *testing.T) {
		repo := mockAssignee.NewMockAssigneeRepository()
		uc := appTicket.NewResolveAssigneeUseCase(repo)

		_, err := uc.Execute(ctx, "   ", false)
		require.ErrorIs(t, err, domainTicket.ErrInvalidAssignee)
	})

	t.Run("should return error when assignee does not exist", func(t *testing.T) {
		repo := mockAssignee.NewMockAssigneeRepository()
		repo.On("FindByName", ctx, "Inexistente").Return(nil, domainAssignee.ErrAssigneeNotFound)

		uc := appTicket.NewResolveAssigneeUseCase(repo)
		_, err := uc.Execute(ctx, "Inexistente", false)

		require.ErrorIs(t, err, domainAssignee.ErrAssigneeNotFound)
		repo.AssertExpectations(t)
	})
}
