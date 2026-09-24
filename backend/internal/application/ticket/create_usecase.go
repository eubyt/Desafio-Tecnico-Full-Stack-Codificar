package ticket

import (
	"context"
	"time"

	domainTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/ticket"
	"github.com/google/uuid"
)

// CreateUseCase orquestra a abertura de chamados com suporte a atribuição manual ou automática
type CreateUseCase struct {
	ticketRepo        domainTicket.Repository
	resolveAssigneeUC *ResolveAssigneeUseCase
}

// NewCreateUseCase instancia o caso de uso para criação de chamados
func NewCreateUseCase(
	ticketRepo domainTicket.Repository,
	resolveAssigneeUC *ResolveAssigneeUseCase,
) *CreateUseCase {
	return &CreateUseCase{
		ticketRepo:        ticketRepo,
		resolveAssigneeUC: resolveAssigneeUC,
	}
}

// Execute cria um novo chamado validando dados e persistindo no repositório
func (uc *CreateUseCase) Execute(ctx context.Context, input CreateInput) (Output, error) {
	priority, err := domainTicket.ParsePriority(input.Priority)
	if err != nil {
		return Output{}, err
	}

	selectedAssignee, err := uc.resolveAssigneeUC.Execute(ctx, input.Assignee, input.AutoAssign)
	if err != nil {
		return Output{}, err
	}

	ticketID, err := uuid.NewV7()
	if err != nil {
		return Output{}, err
	}
	t := domainTicket.NewTicket(
		ticketID,
		input.Title,
		input.Description,
		priority,
		*selectedAssignee,
		time.Now().UTC(),
	)

	if err := t.Validate(); err != nil {
		return Output{}, err
	}

	if err := uc.ticketRepo.Save(ctx, t); err != nil {
		return Output{}, err
	}

	return toOutput(t), nil
}
