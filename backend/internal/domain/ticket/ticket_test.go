package ticket_test

import (
	"strings"
	"testing"
	"time"

	domainAssignee "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/assignee"
	domainTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/ticket"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTicket(t *testing.T) {
	now := time.Now().UTC()
	id, err := uuid.NewV7()
	require.NoError(t, err)

	assigneeID, err := uuid.NewV7()
	require.NoError(t, err)
	assignee := domainAssignee.NewAssignee(assigneeID, "Carlos Silva", now)

	t.Run("should instantiate ticket with initial status 'open' and trimmed values", func(t *testing.T) {
		tk := domainTicket.NewTicket(
			id,
			"  Payment Gateway Timeout  ",
			"  Payment service returns 504 on checkout  ",
			domainTicket.PriorityHigh,
			*assignee,
			now,
		)

		require.NotNil(t, tk)
		assert.Equal(t, id, tk.ID)
		assert.Equal(t, uuid.Version(7), tk.ID.Version())
		assert.Equal(t, "Payment Gateway Timeout", tk.Title)
		assert.Equal(t, "Payment service returns 504 on checkout", tk.Description)
		assert.Equal(t, domainTicket.PriorityHigh, tk.Priority)
		assert.Equal(t, domainTicket.StatusOpen, tk.Status)
		assert.Equal(t, assignee.ID, tk.AssigneeID)
		assert.Equal(t, "Carlos Silva", tk.Assignee.Name)
		assert.Equal(t, now, tk.CreatedAt)
		assert.Equal(t, now, tk.UpdatedAt)
		assert.True(t, tk.IsOpen())
	})

	t.Run("should handle zero time by defaulting to non-zero UTC time", func(t *testing.T) {
		anaID, _ := uuid.NewV7()
		ana := domainAssignee.NewAssignee(anaID, "Ana Souza", now)
		tk := domainTicket.NewTicket(
			id,
			"Valid Title",
			"Valid description text",
			domainTicket.PriorityLow,
			*ana,
			time.Time{},
		)

		require.NotNil(t, tk)
		assert.False(t, tk.CreatedAt.IsZero())
		assert.False(t, tk.UpdatedAt.IsZero())
		assert.Equal(t, time.UTC, tk.CreatedAt.Location())
	})

	t.Run("should accept special characters, and unicode characters", func(t *testing.T) {
		title := "[PROD #99] Erro crítico na integração PIX & Cartão de Débito/Crédito!"
		description := "Payload com caracteres especiais: <xml><status>500</status></xml> & DROP TABLE 'users'; -- \nLinha 2: こんにちは."
		joseID, _ := uuid.NewV7()
		jose := domainAssignee.NewAssignee(joseID, "José da Silva & Cia (Equipe Técnica #01)", now)

		tk := domainTicket.NewTicket(
			id,
			title,
			description,
			domainTicket.PriorityHigh,
			*jose,
			now,
		)

		require.NotNil(t, tk)
		assert.Equal(t, title, tk.Title)
		assert.Equal(t, description, tk.Description)
		assert.Equal(t, *jose, tk.Assignee)
	})
}

func TestTicketValidate(t *testing.T) {
	now := time.Now().UTC()
	validID, _ := uuid.NewV7()
	assigneeID, _ := uuid.NewV7()
	validAssignee := domainAssignee.Assignee{
		ID:        assigneeID,
		Name:      "Carlos Silva",
		CreatedAt: now,
		UpdatedAt: now,
	}

	baseTicket := func() domainTicket.Ticket {
		return domainTicket.Ticket{
			ID:          validID,
			Title:       "Valid Title Here",
			Description: "Valid description text for support ticket",
			Priority:    domainTicket.PriorityMedium,
			Status:      domainTicket.StatusOpen,
			AssigneeID:  assigneeID,
			Assignee:    validAssignee,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
	}

	edgeCases := []struct {
		name        string
		modify      func(t *domainTicket.Ticket)
		expectedErr error
	}{
		// Valid base case
		{
			name:        "valid ticket with all standard fields",
			modify:      nil,
			expectedErr: nil,
		},
		// ID edge cases
		{
			name: "ID: nil UUID should fail with ErrTicketNotFound",
			modify: func(t *domainTicket.Ticket) {
				t.ID = uuid.Nil
			},
			expectedErr: domainTicket.ErrTicketNotFound,
		},
		{
			name: "ID: valid UUID v7 should succeed",
			modify: func(t *domainTicket.Ticket) {
				v7ID, _ := uuid.NewV7()
				t.ID = v7ID
			},
			expectedErr: nil,
		},
		// Title edge cases
		{
			name: "Title: exactly 3 chars (minimum allowed boundary) should succeed",
			modify: func(t *domainTicket.Ticket) {
				t.Title = "Bug"
			},
			expectedErr: nil,
		},
		{
			name: "Title: exactly 150 chars (maximum allowed boundary) should succeed",
			modify: func(t *domainTicket.Ticket) {
				t.Title = strings.Repeat("A", 150)
			},
			expectedErr: nil,
		},
		{
			name: "Title: fewer than 3 chars (2 chars) should fail with ErrInvalidTitle",
			modify: func(t *domainTicket.Ticket) {
				t.Title = "ab"
			},
			expectedErr: domainTicket.ErrInvalidTitle,
		},
		{
			name: "Title: exceeds 150 chars (151 chars) should fail with ErrInvalidTitle",
			modify: func(t *domainTicket.Ticket) {
				t.Title = strings.Repeat("A", 151)
			},
			expectedErr: domainTicket.ErrInvalidTitle,
		},
		{
			name: "Title: only whitespace should fail with ErrInvalidTitle",
			modify: func(t *domainTicket.Ticket) {
				t.Title = "   \t\r\n   "
			},
			expectedErr: domainTicket.ErrInvalidTitle,
		},
		{
			name: "Title: trimmed length less than 3 chars should fail with ErrInvalidTitle",
			modify: func(t *domainTicket.Ticket) {
				t.Title = "   ab   "
			},
			expectedErr: domainTicket.ErrInvalidTitle,
		},
		{
			name: "Title: special characters, symbols, and query string should succeed",
			modify: func(t *domainTicket.Ticket) {
				t.Title = "[CRITICAL] Gateway 502: API/v1?user_id=123&env=prod & <Test>"
			},
			expectedErr: nil,
		},
		// Description edge cases
		{
			name: "Description: exactly 5 chars (minimum allowed boundary) should succeed",
			modify: func(t *domainTicket.Ticket) {
				t.Description = "12345"
			},
			expectedErr: nil,
		},
		{
			name: "Description: exactly 2000 chars (maximum allowed boundary) should succeed",
			modify: func(t *domainTicket.Ticket) {
				t.Description = strings.Repeat("D", 2000)
			},
			expectedErr: nil,
		},
		{
			name: "Description: fewer than 5 chars (4 chars) should fail with ErrInvalidDescription",
			modify: func(t *domainTicket.Ticket) {
				t.Description = "1234"
			},
			expectedErr: domainTicket.ErrInvalidDescription,
		},
		{
			name: "Description: exceeds 2000 chars (2001 chars) should fail with ErrInvalidDescription",
			modify: func(t *domainTicket.Ticket) {
				t.Description = strings.Repeat("D", 2001)
			},
			expectedErr: domainTicket.ErrInvalidDescription,
		},
		{
			name: "Description: only whitespace should fail with ErrInvalidDescription",
			modify: func(t *domainTicket.Ticket) {
				t.Description = "   \n\t   "
			},
			expectedErr: domainTicket.ErrInvalidDescription,
		},
		{
			name: "Description: multi-line, code blocks, JSON, and HTML should succeed",
			modify: func(t *domainTicket.Ticket) {
				t.Description = "```json\n{\"error\": \"failed\", \"code\": 500}\n```\n<p>HTML payload test</p>"
			},
			expectedErr: nil,
		},
		// Priority edge cases
		{
			name: "Priority: low should succeed",
			modify: func(t *domainTicket.Ticket) {
				t.Priority = domainTicket.PriorityLow
			},
			expectedErr: nil,
		},
		{
			name: "Priority: medium should succeed",
			modify: func(t *domainTicket.Ticket) {
				t.Priority = domainTicket.PriorityMedium
			},
			expectedErr: nil,
		},
		{
			name: "Priority: high should succeed",
			modify: func(t *domainTicket.Ticket) {
				t.Priority = domainTicket.PriorityHigh
			},
			expectedErr: nil,
		},
		{
			name: "Priority: empty string should fail with ErrInvalidPriority",
			modify: func(t *domainTicket.Ticket) {
				t.Priority = domainTicket.Priority("")
			},
			expectedErr: domainTicket.ErrInvalidPriority,
		},
		{
			name: "Priority: unexpected string 'critical' should fail with ErrInvalidPriority",
			modify: func(t *domainTicket.Ticket) {
				t.Priority = domainTicket.Priority("critical")
			},
			expectedErr: domainTicket.ErrInvalidPriority,
		},
		{
			name: "Priority: uppercase string 'HIGH' should fail with ErrInvalidPriority",
			modify: func(t *domainTicket.Ticket) {
				t.Priority = domainTicket.Priority("HIGH")
			},
			expectedErr: domainTicket.ErrInvalidPriority,
		},
		{
			name: "Priority: trailing space 'low ' should fail with ErrInvalidPriority",
			modify: func(t *domainTicket.Ticket) {
				t.Priority = domainTicket.Priority("low ")
			},
			expectedErr: domainTicket.ErrInvalidPriority,
		},
		// Status edge cases
		{
			name: "Status: open should succeed",
			modify: func(t *domainTicket.Ticket) {
				t.Status = domainTicket.StatusOpen
			},
			expectedErr: nil,
		},
		{
			name: "Status: in_progress should succeed",
			modify: func(t *domainTicket.Ticket) {
				t.Status = domainTicket.StatusInProgress
			},
			expectedErr: nil,
		},
		{
			name: "Status: resolved should succeed",
			modify: func(t *domainTicket.Ticket) {
				t.Status = domainTicket.StatusResolved
			},
			expectedErr: nil,
		},
		{
			name: "Status: closed should succeed",
			modify: func(t *domainTicket.Ticket) {
				t.Status = domainTicket.StatusClosed
			},
			expectedErr: nil,
		},
		{
			name: "Status: empty string should fail with ErrInvalidStatus",
			modify: func(t *domainTicket.Ticket) {
				t.Status = domainTicket.Status("")
			},
			expectedErr: domainTicket.ErrInvalidStatus,
		},
		{
			name: "Status: unexpected string 'pending' should fail with ErrInvalidStatus",
			modify: func(t *domainTicket.Ticket) {
				t.Status = domainTicket.Status("pending")
			},
			expectedErr: domainTicket.ErrInvalidStatus,
		},
		{
			name: "Status: uppercase string 'OPEN' should fail with ErrInvalidStatus",
			modify: func(t *domainTicket.Ticket) {
				t.Status = domainTicket.Status("OPEN")
			},
			expectedErr: domainTicket.ErrInvalidStatus,
		},
		// Assignee edge cases
		{
			name: "Assignee: nil AssigneeID should fail with ErrInvalidAssignee",
			modify: func(t *domainTicket.Ticket) {
				t.AssigneeID = uuid.Nil
			},
			expectedErr: domainTicket.ErrInvalidAssignee,
		},
		{
			name: "Assignee: empty name should fail with ErrInvalidAssignee",
			modify: func(t *domainTicket.Ticket) {
				t.Assignee.Name = ""
			},
			expectedErr: domainTicket.ErrInvalidAssignee,
		},
		{
			name: "Assignee: only whitespace and tabs in name should fail with ErrInvalidAssignee",
			modify: func(t *domainTicket.Ticket) {
				t.Assignee.Name = "   \t   "
			},
			expectedErr: domainTicket.ErrInvalidAssignee,
		},
		{
			name: "Assignee: accented characters and special symbols should succeed",
			modify: func(t *domainTicket.Ticket) {
				t.Assignee.Name = "Maria da Conceição & André"
			},
			expectedErr: nil,
		},
		// Timestamp edge cases (UpdatedAt vs CreatedAt)
		{
			name: "Timestamp: UpdatedAt earlier than CreatedAt should fail with ErrInvalidTimestamp",
			modify: func(t *domainTicket.Ticket) {
				t.CreatedAt = now
				t.UpdatedAt = now.Add(-1 * time.Minute)
			},
			expectedErr: domainTicket.ErrInvalidTimestamp,
		},
		{
			name: "Timestamp: UpdatedAt equal to CreatedAt should succeed",
			modify: func(t *domainTicket.Ticket) {
				t.CreatedAt = now
				t.UpdatedAt = now
			},
			expectedErr: nil,
		},
		{
			name: "Timestamp: UpdatedAt later than CreatedAt should succeed",
			modify: func(t *domainTicket.Ticket) {
				t.CreatedAt = now
				t.UpdatedAt = now.Add(2 * time.Hour)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			tk := baseTicket()
			if tc.modify != nil {
				tc.modify(&tk)
			}
			err := tk.Validate()
			if tc.expectedErr != nil {
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestTicketUpdate(t *testing.T) {
	now := time.Now().UTC()
	id, _ := uuid.NewV7()
	assigneeID, _ := uuid.NewV7()
	carlos := domainAssignee.Assignee{ID: assigneeID, Name: "Carlos Silva", CreatedAt: now, UpdatedAt: now}
	anaID, _ := uuid.NewV7()
	ana := domainAssignee.Assignee{ID: anaID, Name: "Ana Souza", CreatedAt: now, UpdatedAt: now}

	tk := domainTicket.NewTicket(
		id,
		"Initial Title",
		"Initial description of the ticket",
		domainTicket.PriorityLow,
		carlos,
		now,
	)

	t.Run("should update fields successfully with valid data", func(t *testing.T) {
		later := now.Add(10 * time.Minute)
		err := tk.Update(
			"  Updated Title With Details  ",
			"  Updated description with more than 5 characters  ",
			domainTicket.PriorityHigh,
			domainTicket.StatusInProgress,
			ana,
			later,
		)

		require.NoError(t, err)
		assert.Equal(t, "Updated Title With Details", tk.Title)
		assert.Equal(t, "Updated description with more than 5 characters", tk.Description)
		assert.Equal(t, domainTicket.PriorityHigh, tk.Priority)
		assert.Equal(t, domainTicket.StatusInProgress, tk.Status)
		assert.Equal(t, ana.ID, tk.AssigneeID)
		assert.Equal(t, "Ana Souza", tk.Assignee.Name)
		assert.Equal(t, later, tk.UpdatedAt)
		assert.True(t, tk.IsOpen())
	})

	t.Run("should handle zero time by defaulting to non-zero UTC time", func(t *testing.T) {
		err := tk.Update(
			"Title",
			"Description",
			domainTicket.PriorityMedium,
			domainTicket.StatusOpen,
			carlos,
			time.Time{},
		)

		require.NoError(t, err)
		assert.False(t, tk.UpdatedAt.IsZero())
	})

	t.Run("should reject invalid status transition from closed to resolved", func(t *testing.T) {
		later := now.Add(20 * time.Minute)
		tk.Status = domainTicket.StatusClosed
		err := tk.Update(
			tk.Title,
			tk.Description,
			tk.Priority,
			domainTicket.StatusResolved,
			tk.Assignee,
			later,
		)
		require.Error(t, err)
		assert.Equal(t, domainTicket.ErrInvalidStatusTransition, err)
	})

	t.Run("should allow reopening from closed to open or in_progress", func(t *testing.T) {
		tk.Status = domainTicket.StatusClosed

		err := tk.Update(
			tk.Title,
			tk.Description,
			tk.Priority,
			domainTicket.StatusOpen,
			tk.Assignee,
			now,
		)
		require.NoError(t, err)
		assert.Equal(t, domainTicket.StatusOpen, tk.Status)

		tk.Status = domainTicket.StatusClosed
		err = tk.Update(
			tk.Title,
			tk.Description,
			tk.Priority,
			domainTicket.StatusInProgress,
			tk.Assignee,
			now,
		)
		require.NoError(t, err)
		assert.Equal(t, domainTicket.StatusInProgress, tk.Status)
	})
}

func TestTicketOperations(t *testing.T) {
	now := time.Now().UTC()
	id, _ := uuid.NewV7()
	assigneeID, _ := uuid.NewV7()
	carlos := domainAssignee.Assignee{ID: assigneeID, Name: "Carlos Silva", CreatedAt: now, UpdatedAt: now}
	brunoID, _ := uuid.NewV7()
	bruno := domainAssignee.Assignee{ID: brunoID, Name: "Bruno Santos", CreatedAt: now, UpdatedAt: now}
	anaID, _ := uuid.NewV7()
	ana := domainAssignee.Assignee{ID: anaID, Name: "Ana Souza", CreatedAt: now, UpdatedAt: now}

	tk := domainTicket.NewTicket(
		id,
		"Login Problem",
		"User cannot login with Google SSO",
		domainTicket.PriorityMedium,
		carlos,
		now,
	)

	t.Run("AssignTo should successfully reassign ticket", func(t *testing.T) {
		later := now.Add(5 * time.Minute)
		err := tk.AssignTo(bruno, later)
		require.NoError(t, err)
		assert.Equal(t, bruno.ID, tk.AssigneeID)
		assert.Equal(t, "Bruno Santos", tk.Assignee.Name)
		assert.Equal(t, later, tk.UpdatedAt)
	})

	t.Run("AssignTo should reject empty or whitespace assignee", func(t *testing.T) {
		err := tk.AssignTo(domainAssignee.Assignee{ID: uuid.Nil, Name: "   "}, now)
		require.ErrorIs(t, err, domainTicket.ErrInvalidAssignee)
	})

	t.Run("AssignTo should handle zero time", func(t *testing.T) {
		err := tk.AssignTo(ana, time.Time{})
		require.NoError(t, err)
		assert.False(t, tk.UpdatedAt.IsZero())
	})

	t.Run("ChangeStatus should transition through valid lifecycle", func(t *testing.T) {
		later := now.Add(15 * time.Minute)
		err := tk.ChangeStatus(domainTicket.StatusResolved, later)
		require.NoError(t, err)
		assert.Equal(t, domainTicket.StatusResolved, tk.Status)
		assert.False(t, tk.IsOpen(), "Resolved status should not be considered open")
	})

	t.Run("ChangeStatus should reject invalid status transition", func(t *testing.T) {
		tk.Status = domainTicket.StatusClosed
		err := tk.ChangeStatus(domainTicket.StatusResolved, now)
		require.ErrorIs(t, err, domainTicket.ErrInvalidStatusTransition)
	})

	t.Run("ChangeStatus should handle zero time", func(t *testing.T) {
		tk.Status = domainTicket.StatusOpen
		err := tk.ChangeStatus(domainTicket.StatusInProgress, time.Time{})
		require.NoError(t, err)
		assert.False(t, tk.UpdatedAt.IsZero())
	})

	t.Run("IsOpen status checks", func(t *testing.T) {
		tk.Status = domainTicket.StatusOpen
		assert.True(t, tk.IsOpen())

		tk.Status = domainTicket.StatusInProgress
		assert.True(t, tk.IsOpen())

		tk.Status = domainTicket.StatusResolved
		assert.False(t, tk.IsOpen())

		tk.Status = domainTicket.StatusClosed
		assert.False(t, tk.IsOpen())
	})
}
