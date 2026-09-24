package ticket

import (
	"context"
	"strings"

	domainTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/ticket"
	domainPagination "github.com/adriancf/demo-sistema-chamado/backend/internal/util/pagination"
)

// ListUseCase orquestra a listagem paginada, busca e ordenação de chamados
type ListUseCase struct {
	ticketRepo domainTicket.Repository
}

// NewListUseCase instancia o caso de uso de listagem de chamados
func NewListUseCase(ticketRepo domainTicket.Repository) *ListUseCase {
	return &ListUseCase{
		ticketRepo: ticketRepo,
	}
}

// Execute aplica filtros, ordenação e paginação para consulta de chamados
func (uc *ListUseCase) Execute(ctx context.Context, input ListInput) (PaginatedOutput, error) {
	var statusFilter *domainTicket.Status
	if input.Status != nil && *input.Status != "" {
		s, err := domainTicket.ParseStatus(*input.Status)
		if err != nil {
			return PaginatedOutput{}, err
		}
		statusFilter = &s
	}

	var priorityFilter *domainTicket.Priority
	if input.Priority != nil && *input.Priority != "" {
		p, err := domainTicket.ParsePriority(*input.Priority)
		if err != nil {
			return PaginatedOutput{}, err
		}
		priorityFilter = &p
	}

	filter := domainTicket.Filter{
		Status:   statusFilter,
		Priority: priorityFilter,
		Assignee: input.Assignee,
		Search:   input.Search,
	}

	paginationReq := domainPagination.PaginationRequest{
		Page:     input.Page,
		PageSize: input.PageSize,
	}
	paginationReq.Normalize()

	var sortField domainTicket.SortField
	switch strings.ToLower(input.SortBy) {
	case "priority":
		sortField = domainTicket.SortByPriority
	case "status":
		sortField = domainTicket.SortByStatus
	case "title":
		sortField = domainTicket.SortByTitle
	default:
		sortField = domainTicket.SortByCreatedAt
	}

	sortOption := domainPagination.SortOption{
		Field: string(sortField),
		Order: domainPagination.ParseSortOrder(input.Order),
	}

	result, err := uc.ticketRepo.FindAll(ctx, filter, paginationReq, sortOption)
	if err != nil {
		return PaginatedOutput{}, err
	}

	items := make([]Output, 0, len(result.Items))
	for _, t := range result.Items {
		items = append(items, toOutput(t))
	}

	return PaginatedOutput{
		Items:       items,
		TotalItems:  result.TotalItems,
		TotalPages:  result.TotalPages,
		CurrentPage: result.CurrentPage,
		PageSize:    result.PageSize,
	}, nil
}
