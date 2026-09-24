package ticket

import (
	"time"

	appAssignee "github.com/adriancf/demo-sistema-chamado/backend/internal/application/assignee"
)

type CreateTicketInput struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"`
	Assignee    string `json:"assignee"`
	AutoAssign  bool   `json:"auto_assign"`
}

type UpdateTicketInput struct {
	Status     string `json:"status"`
	Assignee   string `json:"assignee"`
	AutoAssign bool   `json:"auto_assign"`
}

type ListTicketsInput struct {
	Page     int     `json:"page"`
	PageSize int     `json:"page_size"`
	SortBy   string  `json:"sort_by"`
	Order    string  `json:"order"`
	Status   *string `json:"status"`
	Assignee *string `json:"assignee"`
	Priority *string `json:"priority"`
	Search   *string `json:"search"`
}

type TicketOutput struct {
	ID          string                     `json:"id"`
	Title       string                     `json:"title"`
	Description string                     `json:"description"`
	Priority    string                     `json:"priority"`
	Status      string                     `json:"status"`
	AssigneeID  string                     `json:"assignee_id"`
	Assignee    string                     `json:"assignee"`
	AssigneeObj appAssignee.AssigneeOutput `json:"assignee_details,omitempty"`
	CreatedAt   time.Time                  `json:"created_at"`
	UpdatedAt   time.Time                  `json:"updated_at"`
}

type PaginatedTicketsOutput struct {
	Items       []TicketOutput `json:"items"`
	TotalItems  int            `json:"total_items"`
	TotalPages  int            `json:"total_pages"`
	CurrentPage int            `json:"current_page"`
	PageSize    int            `json:"page_size"`
}

// Aliases de conveniência e retrocompatibilidade
type CreateInput = CreateTicketInput
type UpdateInput = UpdateTicketInput
type ListInput = ListTicketsInput
type Output = TicketOutput
type PaginatedOutput = PaginatedTicketsOutput
