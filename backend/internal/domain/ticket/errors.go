package ticket

import (
	"github.com/adriancf/demo-sistema-chamado/backend/internal/domain"
)

var (
	ErrTicketNotFound          = domain.NewError("ticket not found", domain.ErrNotFound)
	ErrInvalidTitle            = domain.NewError("ticket title must have between 3 and 150 characters", domain.ErrValidation)
	ErrInvalidDescription      = domain.NewError("ticket description must have between 5 and 2000 characters", domain.ErrValidation)
	ErrInvalidPriority         = domain.NewError("invalid priority: must be 'low', 'medium' or 'high'", domain.ErrValidation)
	ErrInvalidStatus           = domain.NewError("invalid status: must be 'open', 'in_progress', 'resolved' or 'closed'", domain.ErrValidation)
	ErrInvalidStatusTransition = domain.NewError("invalid status transition", domain.ErrValidation)
	ErrInvalidAssignee         = domain.NewError("assignee name cannot be empty", domain.ErrValidation)
	ErrInvalidTimestamp        = domain.NewError("updated_at cannot be earlier than created_at", domain.ErrValidation)
	ErrInvalidPagination       = domain.NewError("invalid pagination parameters", domain.ErrValidation)
)
