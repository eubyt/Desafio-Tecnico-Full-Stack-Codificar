package ticket

import (
	"context"

	domainTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/ticket"
	domainPagination "github.com/adriancf/demo-sistema-chamado/backend/internal/util/pagination"
	"github.com/stretchr/testify/mock"
)

type MockTicketRepository struct {
	mock.Mock
}

func NewMockTicketRepository() *MockTicketRepository {
	return &MockTicketRepository{}
}

func (m *MockTicketRepository) Save(ctx context.Context, t *domainTicket.Ticket) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

func (m *MockTicketRepository) Update(ctx context.Context, t *domainTicket.Ticket) error {
	args := m.Called(ctx, t)
	return args.Error(0)
}

func (m *MockTicketRepository) FindByID(ctx context.Context, id string) (*domainTicket.Ticket, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domainTicket.Ticket), args.Error(1)
}

func (m *MockTicketRepository) FindAll(ctx context.Context, filter domainTicket.Filter, pagination domainPagination.PaginationRequest, sortOption domainPagination.SortOption) (domainPagination.PaginatedResult[*domainTicket.Ticket], error) {
	args := m.Called(ctx, filter, pagination, sortOption)
	if args.Get(0) == nil {
		return domainPagination.PaginatedResult[*domainTicket.Ticket]{}, args.Error(1)
	}
	return args.Get(0).(domainPagination.PaginatedResult[*domainTicket.Ticket]), args.Error(1)
}

func (m *MockTicketRepository) CountOpenTicketsByAssignee(ctx context.Context) (map[string]int, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]int), args.Error(1)
}
