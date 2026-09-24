package assignee_test

import (
	"context"
	"testing"
	"time"

	appAssignee "github.com/adriancf/demo-sistema-chamado/backend/internal/application/assignee"
	domainAssignee "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/assignee"
	mockAssignee "github.com/adriancf/demo-sistema-chamado/backend/internal/util/test/mock/assignee"
	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateUseCase(t *testing.T) {
	ctx := context.Background()

	t.Run("should create a new assignee successfully", func(t *testing.T) {
		repo := mockAssignee.NewMockAssigneeRepository()
		repo.On("FindByName", ctx, "Carlos Silva").Return(nil, domainAssignee.ErrAssigneeNotFound)
		repo.On("Save", ctx, mock.MatchedBy(func(a domainAssignee.Assignee) bool {
			return a.Name == "Carlos Silva" && a.ID.Version() == uuid.Version(7)
		})).Return(nil)

		uc := appAssignee.NewCreateUseCase(repo)
		err := uc.Execute(ctx, appAssignee.CreateInput{Name: "Carlos Silva"})

		require.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("should fail when creating assignee with name too short", func(t *testing.T) {
		repo := mockAssignee.NewMockAssigneeRepository()
		uc := appAssignee.NewCreateUseCase(repo)

		err := uc.Execute(ctx, appAssignee.CreateInput{Name: "A"})
		require.ErrorIs(t, err, domainAssignee.ErrInvalidAssignee)
	})

	t.Run("should fail when creating duplicate assignee", func(t *testing.T) {
		repo := mockAssignee.NewMockAssigneeRepository()
		existing := domainAssignee.NewAssignee(uuid.New(), "Carlos Silva", time.Now().UTC())
		repo.On("FindByName", ctx, "Carlos Silva").Return(existing, nil)

		uc := appAssignee.NewCreateUseCase(repo)
		err := uc.Execute(ctx, appAssignee.CreateInput{Name: "Carlos Silva"})

		require.ErrorIs(t, err, domainAssignee.ErrAssigneeAlreadyExists)
		repo.AssertExpectations(t)
	})

	t.Run("should fail when repository returns error on Save", func(t *testing.T) {
		repo := mockAssignee.NewMockAssigneeRepository()
		repo.On("FindByName", ctx, "Carlos Silva").Return(nil, domainAssignee.ErrAssigneeNotFound)
		repo.On("Save", ctx, mock.Anything).Return(domainAssignee.ErrInvalidAssignee)

		uc := appAssignee.NewCreateUseCase(repo)
		err := uc.Execute(ctx, appAssignee.CreateInput{Name: "Carlos Silva"})

		require.Error(t, err)
		repo.AssertExpectations(t)
	})
}
