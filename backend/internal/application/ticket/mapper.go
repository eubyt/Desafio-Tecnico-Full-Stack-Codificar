package ticket

import (
	appAssignee "github.com/adriancf/demo-sistema-chamado/backend/internal/application/assignee"
	domainTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/ticket"
)

func toOutput(t *domainTicket.Ticket) Output {
	return Output{
		ID:          t.ID.String(),
		Title:       t.Title,
		Description: t.Description,
		Priority:    string(t.Priority),
		Status:      string(t.Status),
		AssigneeID:  t.AssigneeID.String(),
		Assignee:    t.Assignee.Name,
		AssigneeObj: appAssignee.Output{
			ID:        t.Assignee.ID.String(),
			Name:      t.Assignee.Name,
			CreatedAt: t.Assignee.CreatedAt,
			UpdatedAt: t.Assignee.UpdatedAt,
		},
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
