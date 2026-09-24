package postgres

import (
	"context"
	"errors"
	"strings"

	"github.com/adriancf/demo-sistema-chamado/backend/internal/domain/assignee"
	"github.com/adriancf/demo-sistema-chamado/backend/internal/infra/postgres/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AssigneeRepository struct {
	pool    *pgxpool.Pool
	queries *sqlc.Queries
}

func NewAssigneeRepository(pool *pgxpool.Pool) *AssigneeRepository {
	return &AssigneeRepository{
		pool:    pool,
		queries: sqlc.New(pool),
	}
}

func (r *AssigneeRepository) ListAll(ctx context.Context) ([]assignee.Assignee, error) {
	rows, err := r.queries.ListAssignees(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]assignee.Assignee, len(rows))
	for i, row := range rows {
		result[i] = toDomainAssignee(row)
	}
	return result, nil
}

func (r *AssigneeRepository) FindByName(ctx context.Context, name string) (*assignee.Assignee, error) {
	trimmed := strings.TrimSpace(name)
	row, err := r.queries.FindAssigneeByName(ctx, trimmed)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, assignee.ErrAssigneeNotFound
		}
		return nil, err
	}

	res := toDomainAssignee(row)
	return &res, nil
}

func (r *AssigneeRepository) FindLeastLoaded(ctx context.Context) (*assignee.Assignee, error) {
	row, err := r.queries.FindLeastLoadedAssignee(ctx)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, assignee.ErrNoAssigneesAvailable
		}
		return nil, err
	}

	res := toDomainAssignee(row)
	return &res, nil
}

func (r *AssigneeRepository) Save(ctx context.Context, a assignee.Assignee) error {
	err := r.queries.UpsertAssignee(ctx, sqlc.UpsertAssigneeParams{
		ID:        uuidToPg(a.ID),
		Name:      a.Name,
		CreatedAt: timeToPg(a.CreatedAt),
		UpdatedAt: timeToPg(a.UpdatedAt),
	})
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // unique_violation
			return assignee.ErrAssigneeAlreadyExists
		}
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") || strings.Contains(strings.ToLower(err.Error()), "unique") {
			return assignee.ErrAssigneeAlreadyExists
		}
		return err
	}
	return nil
}
