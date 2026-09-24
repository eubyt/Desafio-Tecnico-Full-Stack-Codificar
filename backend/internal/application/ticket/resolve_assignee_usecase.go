package ticket

import (
	"context"
	"strings"

	domainAssignee "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/assignee"
	domainTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/ticket"
)

type ResolveAssigneeUseCase struct {
	assigneeRepo domainAssignee.Repository
}

func NewResolveAssigneeUseCase(assigneeRepo domainAssignee.Repository) *ResolveAssigneeUseCase {
	return &ResolveAssigneeUseCase{
		assigneeRepo: assigneeRepo,
	}
}

func (uc *ResolveAssigneeUseCase) Execute(ctx context.Context, assigneeName string, autoAssign bool) (*domainAssignee.Assignee, error) {
	if autoAssign {
		return uc.assigneeRepo.FindLeastLoaded(ctx)
	}

	trimmed := strings.TrimSpace(assigneeName)
	if trimmed == "" {
		return nil, domainTicket.ErrInvalidAssignee
	}

	return uc.assigneeRepo.FindByName(ctx, trimmed)
}
