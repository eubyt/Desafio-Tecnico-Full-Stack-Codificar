package assignee

import (
	"context"

	domainAssignee "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/assignee"
	domainTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/ticket"
)

type ListUseCase struct {
	assigneeRepo domainAssignee.Repository
	ticketRepo   domainTicket.Repository
}

func NewListUseCase(
	assigneeRepo domainAssignee.Repository,
	ticketRepo domainTicket.Repository,
) *ListUseCase {
	return &ListUseCase{
		assigneeRepo: assigneeRepo,
		ticketRepo:   ticketRepo,
	}
}

func (uc *ListUseCase) Execute(ctx context.Context) ([]Output, error) {
	assignees, err := uc.assigneeRepo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	openCounts, err := uc.ticketRepo.CountOpenTicketsByAssignee(ctx)
	if err != nil {
		return nil, err
	}

	output := make([]Output, 0, len(assignees))
	for _, a := range assignees {
		output = append(output, Output{
			ID:               a.ID.String(),
			Name:             a.Name,
			CreatedAt:        a.CreatedAt,
			UpdatedAt:        a.UpdatedAt,
			OpenTicketsCount: openCounts[a.Name],
		})
	}

	return output, nil
}
