package ticket

import (
	"context"

	domainTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/ticket"
)

// GetByIDUseCase orquestra a busca de um chamado pelo seu identificador único
type GetByIDUseCase struct {
	ticketRepo domainTicket.Repository
}

// NewGetByIDUseCase instancia o caso de uso de busca por ID
func NewGetByIDUseCase(ticketRepo domainTicket.Repository) *GetByIDUseCase {
	return &GetByIDUseCase{
		ticketRepo: ticketRepo,
	}
}

// Execute recupera os dados de um chamado por ID
func (uc *GetByIDUseCase) Execute(ctx context.Context, id string) (Output, error) {
	t, err := uc.ticketRepo.FindByID(ctx, id)
	if err != nil {
		return Output{}, err
	}

	return toOutput(t), nil
}
