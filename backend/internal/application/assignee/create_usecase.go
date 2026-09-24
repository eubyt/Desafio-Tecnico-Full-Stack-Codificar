package assignee

import (
	"context"
	"time"

	domainAssignee "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/assignee"
	"github.com/google/uuid"
)

type CreateUseCase struct {
	assigneeRepo domainAssignee.Repository
}

func NewCreateUseCase(assigneeRepo domainAssignee.Repository) *CreateUseCase {
	return &CreateUseCase{
		assigneeRepo: assigneeRepo,
	}
}
func (uc *CreateUseCase) Execute(ctx context.Context, input CreateInput) error {
	id, err := uuid.NewV7()
	if err != nil {
		return err
	}

	a := domainAssignee.NewAssignee(id, input.Name, time.Now().UTC())
	if err := a.Validate(); err != nil {
		return err
	}

	existing, err := uc.assigneeRepo.FindByName(ctx, a.Name)
	if err == nil && existing != nil {
		return domainAssignee.ErrAssigneeAlreadyExists
	}

	return uc.assigneeRepo.Save(ctx, *a)
}
