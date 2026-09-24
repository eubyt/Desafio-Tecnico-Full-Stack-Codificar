package ticket

import (
	"context"

	"github.com/adriancf/demo-sistema-chamado/backend/internal/util/pagination"
)

type SortField string

const (
	SortByCreatedAt SortField = "created_at"
	SortByPriority  SortField = "priority"
	SortByStatus    SortField = "status"
	SortByTitle     SortField = "title"
)

type Filter struct {
	Status   *Status   `json:"status,omitempty"`
	Assignee *string   `json:"assignee,omitempty"`
	Priority *Priority `json:"priority,omitempty"`
	Search   *string   `json:"search,omitempty"`
}

// Repository define o contrato de persistência para o agregado Ticket
type Repository interface {
	Save(ctx context.Context, ticket *Ticket) error
	Update(ctx context.Context, ticket *Ticket) error
	FindByID(ctx context.Context, id string) (*Ticket, error)
	FindAll(ctx context.Context, filter Filter, pagination pagination.PaginationRequest, sortOption pagination.SortOption) (pagination.PaginatedResult[*Ticket], error)
	CountOpenTicketsByAssignee(ctx context.Context) (map[string]int, error)
}
