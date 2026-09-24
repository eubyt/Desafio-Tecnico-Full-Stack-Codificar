package postgres

import (
	"time"

	"github.com/adriancf/demo-sistema-chamado/backend/internal/domain/assignee"
	"github.com/adriancf/demo-sistema-chamado/backend/internal/domain/ticket"
	"github.com/adriancf/demo-sistema-chamado/backend/internal/infra/postgres/sqlc"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

func uuidToPg(u uuid.UUID) pgtype.UUID {
	return pgtype.UUID{Bytes: u, Valid: true}
}

func pgToUUID(p pgtype.UUID) uuid.UUID {
	if !p.Valid {
		return uuid.Nil
	}
	return uuid.UUID(p.Bytes)
}

func timeToPg(t time.Time) pgtype.Timestamptz {
	if t.IsZero() {
		return pgtype.Timestamptz{Valid: false}
	}
	return pgtype.Timestamptz{Time: t.UTC(), Valid: true}
}

func pgToTime(p pgtype.Timestamptz) time.Time {
	if !p.Valid {
		return time.Time{}
	}
	return p.Time.UTC()
}

func textToPg(s *string) pgtype.Text {
	if s == nil || *s == "" {
		return pgtype.Text{Valid: false}
	}
	return pgtype.Text{String: *s, Valid: true}
}

func toDomainAssignee(m sqlc.Assignee) assignee.Assignee {
	return assignee.Assignee{
		ID:        pgToUUID(m.ID),
		Name:      m.Name,
		CreatedAt: pgToTime(m.CreatedAt),
		UpdatedAt: pgToTime(m.UpdatedAt),
	}
}

func toDomainTicketFromFindRow(row sqlc.FindTicketByIDRow) *ticket.Ticket {
	return &ticket.Ticket{
		ID:          pgToUUID(row.ID),
		Title:       row.Title,
		Description: row.Description,
		Priority:    ticket.Priority(row.Priority),
		Status:      ticket.Status(row.Status),
		AssigneeID:  pgToUUID(row.AssigneeID),
		Assignee: assignee.Assignee{
			ID:        pgToUUID(row.AssigneeRefID),
			Name:      row.AssigneeName,
			CreatedAt: pgToTime(row.AssigneeCreatedAt),
			UpdatedAt: pgToTime(row.AssigneeUpdatedAt),
		},
		CreatedAt: pgToTime(row.CreatedAt),
		UpdatedAt: pgToTime(row.UpdatedAt),
	}
}

func toDomainTicketFromListRow(row sqlc.ListTicketsWithFilterAndSortRow) *ticket.Ticket {
	return &ticket.Ticket{
		ID:          pgToUUID(row.ID),
		Title:       row.Title,
		Description: row.Description,
		Priority:    ticket.Priority(row.Priority),
		Status:      ticket.Status(row.Status),
		AssigneeID:  pgToUUID(row.AssigneeID),
		Assignee: assignee.Assignee{
			ID:        pgToUUID(row.AssigneeRefID),
			Name:      row.AssigneeName,
			CreatedAt: pgToTime(row.AssigneeCreatedAt),
			UpdatedAt: pgToTime(row.AssigneeUpdatedAt),
		},
		CreatedAt: pgToTime(row.CreatedAt),
		UpdatedAt: pgToTime(row.UpdatedAt),
	}
}
