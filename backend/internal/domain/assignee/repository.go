package assignee

import "context"

// Repository define o contrato de persistência para os responsáveis
type Repository interface {
	ListAll(ctx context.Context) ([]Assignee, error)
	FindByName(ctx context.Context, name string) (*Assignee, error)
	FindLeastLoaded(ctx context.Context) (*Assignee, error)
	Save(ctx context.Context, assignee Assignee) error
}
