package assignee

import "time"

type CreateAssigneeInput struct {
	Name string `json:"name"`
}

type AssigneeOutput struct {
	ID               string    `json:"id"`
	Name             string    `json:"name"`
	CreatedAt        time.Time `json:"created_at,omitempty"`
	UpdatedAt        time.Time `json:"updated_at,omitempty"`
	OpenTicketsCount int       `json:"open_tickets_count"`
}

// Aliases para conveniência e retrocompatibilidade
type CreateInput = CreateAssigneeInput
type Output = AssigneeOutput
