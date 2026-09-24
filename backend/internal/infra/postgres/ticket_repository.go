package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/adriancf/demo-sistema-chamado/backend/internal/domain/ticket"
	"github.com/adriancf/demo-sistema-chamado/backend/internal/infra/postgres/sqlc"
	"github.com/adriancf/demo-sistema-chamado/backend/internal/util/pagination"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TicketRepository struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

func NewTicketRepository(pool *pgxpool.Pool) *TicketRepository {
	return &TicketRepository{
		pool:    pool,
		queries: sqlc.New(pool),
	}
}

func (r *TicketRepository) Save(ctx context.Context, t *ticket.Ticket) error {
	return r.queries.CreateTicket(ctx, sqlc.CreateTicketParams{
		ID:          uuidToPg(t.ID),
		Title:       t.Title,
		Description: t.Description,
		Priority:    string(t.Priority),
		Status:      string(t.Status),
		AssigneeID:  uuidToPg(t.AssigneeID),
		CreatedAt:   timeToPg(t.CreatedAt),
		UpdatedAt:   timeToPg(t.UpdatedAt),
	})
}

func (r *TicketRepository) Update(ctx context.Context, t *ticket.Ticket) error {
	rowsAffected, err := r.queries.UpdateTicket(ctx, sqlc.UpdateTicketParams{
		ID:          uuidToPg(t.ID),
		Title:       t.Title,
		Description: t.Description,
		Priority:    string(t.Priority),
		Status:      string(t.Status),
		AssigneeID:  uuidToPg(t.AssigneeID),
		UpdatedAt:   timeToPg(t.UpdatedAt),
	})
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return ticket.ErrTicketNotFound
	}
	return nil
}

func (r *TicketRepository) FindByID(ctx context.Context, id string) (*ticket.Ticket, error) {
	parsedID, err := uuid.Parse(strings.TrimSpace(id))
	if err != nil {
		return nil, ticket.ErrTicketNotFound
	}

	row, err := r.queries.FindTicketByID(ctx, uuidToPg(parsedID))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, ticket.ErrTicketNotFound
		}
		return nil, err
	}

	return toDomainTicketFromFindRow(row), nil
}

func (r *TicketRepository) FindAll(ctx context.Context, filter ticket.Filter, paginationReq pagination.PaginationRequest, sortOption pagination.SortOption) (pagination.PaginatedResult[*ticket.Ticket], error) {
	var filterStatus, filterPriority, filterAssignee, filterSearch *string
	if filter.Status != nil && *filter.Status != "" {
		s := string(*filter.Status)
		filterStatus = &s
	}
	if filter.Priority != nil && *filter.Priority != "" {
		p := string(*filter.Priority)
		filterPriority = &p
	}
	if filter.Assignee != nil && strings.TrimSpace(*filter.Assignee) != "" {
		a := strings.TrimSpace(*filter.Assignee)
		filterAssignee = &a
	}
	if filter.Search != nil && strings.TrimSpace(*filter.Search) != "" {
		s := strings.TrimSpace(*filter.Search)
		filterSearch = &s
	}

	totalCount, err := r.queries.CountTicketsWithFilter(ctx, sqlc.CountTicketsWithFilterParams{
		Status:   textToPg(filterStatus),
		Priority: textToPg(filterPriority),
		Assignee: textToPg(filterAssignee),
		Search:   textToPg(filterSearch),
	})
	if err != nil {
		return pagination.PaginatedResult[*ticket.Ticket]{}, err
	}

	rows, err := r.queries.ListTicketsWithFilterAndSort(ctx, sqlc.ListTicketsWithFilterAndSortParams{
		Status:     textToPg(filterStatus),
		Priority:   textToPg(filterPriority),
		Assignee:   textToPg(filterAssignee),
		Search:     textToPg(filterSearch),
		SortBy:     strings.ToLower(string(sortOption.Field)),
		SortOrder:  strings.ToLower(string(sortOption.Order)),
		PageOffset: int32(paginationReq.Offset()),
		PageLimit:  int32(paginationReq.Limit()),
	})
	if err != nil {
		return pagination.PaginatedResult[*ticket.Ticket]{}, err
	}

	tickets := make([]*ticket.Ticket, len(rows))
	for i, row := range rows {
		tickets[i] = toDomainTicketFromListRow(row)
	}

	return pagination.NewPaginatedResult(tickets, int(totalCount), paginationReq.Page, paginationReq.Limit()), nil
}

func (r *TicketRepository) CountOpenTicketsByAssignee(ctx context.Context) (map[string]int, error) {
	rows, err := r.queries.CountOpenTicketsByAssignee(ctx)
	if err != nil {
		return nil, err
	}

	counts := make(map[string]int)
	for _, r := range rows {
		counts[r.AssigneeName] = int(r.OpenCount)
	}
	return counts, nil
}
