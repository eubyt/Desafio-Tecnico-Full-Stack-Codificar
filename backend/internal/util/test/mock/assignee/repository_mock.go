package assignee

import (
	"context"

	domainAssignee "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/assignee"
	"github.com/stretchr/testify/mock"
)

type MockAssigneeRepository struct {
	mock.Mock
}

func NewMockAssigneeRepository() *MockAssigneeRepository {
	return &MockAssigneeRepository{}
}

func (m *MockAssigneeRepository) ListAll(ctx context.Context) ([]domainAssignee.Assignee, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]domainAssignee.Assignee), args.Error(1)
}

func (m *MockAssigneeRepository) FindByName(ctx context.Context, name string) (*domainAssignee.Assignee, error) {
	args := m.Called(ctx, name)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainAssignee.Assignee), args.Error(1)
}

func (m *MockAssigneeRepository) FindLeastLoaded(ctx context.Context) (*domainAssignee.Assignee, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainAssignee.Assignee), args.Error(1)
}

func (m *MockAssigneeRepository) Save(ctx context.Context, a domainAssignee.Assignee) error {
	args := m.Called(ctx, a)
	return args.Error(0)
}
