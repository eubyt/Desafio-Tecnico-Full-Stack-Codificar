package assignee

import (
	"github.com/adriancf/demo-sistema-chamado/backend/internal/domain"
)

var (
	ErrAssigneeNotFound      = domain.NewError("assignee not found", domain.ErrNotFound)
	ErrAssigneeAlreadyExists = domain.NewError("assignee already exists", domain.ErrConflict)
	ErrNoAssigneesAvailable  = domain.NewError("no assignees available for distribution", domain.ErrValidation)
	ErrInvalidAssignee       = domain.NewError("assignee name cannot be empty", domain.ErrValidation)
	ErrInvalidTimestamp      = domain.NewError("updated_at cannot be earlier than created_at", domain.ErrValidation)
)
