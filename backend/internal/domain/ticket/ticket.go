package ticket

import (
	"strings"
	"time"

	"github.com/adriancf/demo-sistema-chamado/backend/internal/domain/assignee"
	"github.com/google/uuid"
)

type Ticket struct {
	ID          uuid.UUID         `json:"id"`
	Title       string            `json:"title"`
	Description string            `json:"description"`
	Priority    Priority          `json:"priority"`
	Status      Status            `json:"status"`
	AssigneeID  uuid.UUID         `json:"assignee_id"`
	Assignee    assignee.Assignee `json:"assignee"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

func (t *Ticket) Validate() error {
	if t.ID == uuid.Nil {
		return ErrTicketNotFound
	}

	titleTrimmed := strings.TrimSpace(t.Title)
	if len(titleTrimmed) < 3 || len(titleTrimmed) > 150 {
		return ErrInvalidTitle
	}

	descTrimmed := strings.TrimSpace(t.Description)
	if len(descTrimmed) < 5 || len(descTrimmed) > 2000 {
		return ErrInvalidDescription
	}

	if !t.Priority.IsValid() {
		return ErrInvalidPriority
	}

	if !t.Status.IsValid() {
		return ErrInvalidStatus
	}

	if t.AssigneeID == uuid.Nil {
		return ErrInvalidAssignee
	}

	if strings.TrimSpace(t.Assignee.Name) == "" {
		return ErrInvalidAssignee
	}

	if t.UpdatedAt.Before(t.CreatedAt) {
		return ErrInvalidTimestamp
	}

	return nil
}

func NewTicket(id uuid.UUID, title, description string, priority Priority, a assignee.Assignee, now time.Time) *Ticket {
	if now.IsZero() {
		now = time.Now().UTC()
	}

	return &Ticket{
		ID:          id,
		Title:       strings.TrimSpace(title),
		Description: strings.TrimSpace(description),
		Priority:    priority,
		Status:      StatusOpen,
		AssigneeID:  a.ID,
		Assignee:    a,
		CreatedAt:   now.UTC(),
		UpdatedAt:   now.UTC(),
	}
}

func (t *Ticket) AssignTo(a assignee.Assignee) error {
	if a.ID == uuid.Nil || strings.TrimSpace(a.Name) == "" {
		return ErrInvalidAssignee
	}
	t.AssigneeID = a.ID
	t.Assignee = a
	t.UpdatedAt = time.Now().UTC()
	return nil
}

func (t *Ticket) ChangeStatus(newStatus Status) error {
	if !t.Status.CanTransitionTo(newStatus) {
		return ErrInvalidStatusTransition
	}
	t.Status = newStatus
	t.UpdatedAt = time.Now().UTC()
	return nil
}

func (t *Ticket) IsOpen() bool {
	return t.Status.IsOpen()
}
