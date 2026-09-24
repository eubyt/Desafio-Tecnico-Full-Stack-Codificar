package ticket

import (
	"context"

	domainTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/ticket"
)

type UpdateUseCase struct {
	ticketRepo        domainTicket.Repository
	resolveAssigneeUC *ResolveAssigneeUseCase
}

func NewUpdateUseCase(
	ticketRepo domainTicket.Repository,
	resolveAssigneeUC *ResolveAssigneeUseCase,
) *UpdateUseCase {
	return &UpdateUseCase{
		ticketRepo:        ticketRepo,
		resolveAssigneeUC: resolveAssigneeUC,
	}
}

func (uc *UpdateUseCase) Execute(ctx context.Context, id string, input UpdateInput) (Output, error) {
	t, err := uc.ticketRepo.FindByID(ctx, id)
	if err != nil {
		return Output{}, err
	}

	if input.Status == "" && input.Assignee == "" && !input.AutoAssign {
		return Output{}, domainTicket.ErrInvalidStatus
	}

	if input.Status != "" {
		status, err := domainTicket.ParseStatus(input.Status)
		if err != nil {
			return Output{}, err
		}
		if err := t.ChangeStatus(status); err != nil {
			return Output{}, err
		}
	}

	if input.Assignee != "" || input.AutoAssign {
		selectedAssignee, err := uc.resolveAssigneeUC.Execute(ctx, input.Assignee, input.AutoAssign)
		if err != nil {
			return Output{}, err
		}
		if err := t.AssignTo(*selectedAssignee); err != nil {
			return Output{}, err
		}
	}

	if err := t.Validate(); err != nil {
		return Output{}, err
	}

	if err := uc.ticketRepo.Update(ctx, t); err != nil {
		return Output{}, err
	}

	return toOutput(t), nil
}
