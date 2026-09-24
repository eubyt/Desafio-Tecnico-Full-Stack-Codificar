package assignee_test

import (
	"context"
	"testing"
	"time"

	appAssignee "github.com/adriancf/demo-sistema-chamado/backend/internal/application/assignee"
	domainAssignee "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/assignee"
	mockAssignee "github.com/adriancf/demo-sistema-chamado/backend/internal/util/test/mock/assignee"
	mockTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/util/test/mock/ticket"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListUseCase(t *testing.T) {
	ctx := context.Background()
	now := time.Now().UTC()

	id1, _ := uuid.NewV7()
	id2, _ := uuid.NewV7()

	assignees := []domainAssignee.Assignee{
		*domainAssignee.NewAssignee(id1, "Carlos Silva", now),
		*domainAssignee.NewAssignee(id2, "Ana Souza", now),
	}

	openCounts := map[string]int{
		"Carlos Silva": 3,
		"Ana Souza":    1,
	}

	t.Run("should list assignees with correct open ticket count", func(t *testing.T) {
		repoA := mockAssignee.NewMockAssigneeRepository()
		repoA.On("ListAll", ctx).Return(assignees, nil).Once()

		repoT := mockTicket.NewMockTicketRepository()
		repoT.On("CountOpenTicketsByAssignee", ctx).Return(openCounts, nil).Once()

		ucList := appAssignee.NewListUseCase(repoA, repoT)
		list, err := ucList.Execute(ctx)
		require.NoError(t, err)
		assert.Len(t, list, 2)

		counts := make(map[string]int)
		for _, item := range list {
			counts[item.Name] = item.OpenTicketsCount
		}
		assert.Equal(t, 3, counts["Carlos Silva"])
		assert.Equal(t, 1, counts["Ana Souza"])

		repoA.AssertExpectations(t)
		repoT.AssertExpectations(t)
	})

	t.Run("should fail when assignee repository fails on list", func(t *testing.T) {
		repoA := mockAssignee.NewMockAssigneeRepository()
		repoA.On("ListAll", ctx).Return(nil, assert.AnError).Once()

		repoT := mockTicket.NewMockTicketRepository()
		ucList := appAssignee.NewListUseCase(repoA, repoT)

		_, err := ucList.Execute(ctx)
		require.Error(t, err)
		repoA.AssertExpectations(t)
	})

	t.Run("should fail when ticket repository fails on counting open tickets", func(t *testing.T) {
		repoA := mockAssignee.NewMockAssigneeRepository()
		repoA.On("ListAll", ctx).Return(assignees, nil).Once()

		repoT := mockTicket.NewMockTicketRepository()
		repoT.On("CountOpenTicketsByAssignee", ctx).Return(nil, assert.AnError).Once()

		ucList := appAssignee.NewListUseCase(repoA, repoT)

		_, err := ucList.Execute(ctx)
		require.Error(t, err)
		repoA.AssertExpectations(t)
		repoT.AssertExpectations(t)
	})
}
