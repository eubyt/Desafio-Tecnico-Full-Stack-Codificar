package assignee

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

type Assignee struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (a *Assignee) Validate() error {
	if a.ID == uuid.Nil {
		return ErrInvalidAssignee
	}
	trimmed := strings.TrimSpace(a.Name)
	if len(trimmed) < 2 || len(trimmed) > 100 {
		return ErrInvalidAssignee
	}
	if a.UpdatedAt.Before(a.CreatedAt) {
		return ErrInvalidTimestamp
	}
	return nil
}

func NewAssignee(id uuid.UUID, name string, now time.Time) *Assignee {
	if now.IsZero() {
		now = time.Now().UTC()
	}
	return &Assignee{
		ID:        id,
		Name:      strings.TrimSpace(name),
		CreatedAt: now.UTC(),
		UpdatedAt: now.UTC(),
	}
}
