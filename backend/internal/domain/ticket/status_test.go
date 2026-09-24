package ticket_test

import (
	"testing"

	domainTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/ticket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseStatus(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    domainTicket.Status
		expectError bool
	}{
		{"Should accept 'open'", "open", domainTicket.StatusOpen, false},
		{"Should accept 'in_progress'", "in_progress", domainTicket.StatusInProgress, false},
		{"Should accept 'resolved'", "resolved", domainTicket.StatusResolved, false},
		{"Should accept 'closed'", "closed", domainTicket.StatusClosed, false},
		{"Should accept with trailing spaces and uppercase '  OPEN '", "  OPEN ", domainTicket.StatusOpen, false},
		{"Should reject invalid status 'pending'", "pending", "", true},
		{"Should reject empty string", "", "", true},
		{"Should reject numbers '123'", "123", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			status, err := domainTicket.ParseStatus(tt.input)
			if tt.expectError {
				require.Error(t, err)
				assert.Equal(t, domainTicket.ErrInvalidStatus, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, status)
				assert.True(t, status.IsValid())
			}
		})
	}
}

func TestStatus_CanTransitionTo(t *testing.T) {
	t.Run("Open transitions", func(t *testing.T) {
		assert.True(t, domainTicket.StatusOpen.CanTransitionTo(domainTicket.StatusOpen))
		assert.True(t, domainTicket.StatusOpen.CanTransitionTo(domainTicket.StatusInProgress))
		assert.True(t, domainTicket.StatusOpen.CanTransitionTo(domainTicket.StatusResolved))
		assert.True(t, domainTicket.StatusOpen.CanTransitionTo(domainTicket.StatusClosed))
		assert.False(t, domainTicket.StatusOpen.CanTransitionTo(domainTicket.Status("invalid")))
	})

	t.Run("InProgress transitions", func(t *testing.T) {
		assert.True(t, domainTicket.StatusInProgress.CanTransitionTo(domainTicket.StatusInProgress))
		assert.True(t, domainTicket.StatusInProgress.CanTransitionTo(domainTicket.StatusOpen))
		assert.True(t, domainTicket.StatusInProgress.CanTransitionTo(domainTicket.StatusResolved))
		assert.True(t, domainTicket.StatusInProgress.CanTransitionTo(domainTicket.StatusClosed))
	})

	t.Run("Resolved transitions", func(t *testing.T) {
		assert.True(t, domainTicket.StatusResolved.CanTransitionTo(domainTicket.StatusResolved))
		assert.True(t, domainTicket.StatusResolved.CanTransitionTo(domainTicket.StatusOpen))
		assert.True(t, domainTicket.StatusResolved.CanTransitionTo(domainTicket.StatusInProgress))
		assert.True(t, domainTicket.StatusResolved.CanTransitionTo(domainTicket.StatusClosed))
	})

	t.Run("Closed transitions", func(t *testing.T) {
		assert.True(t, domainTicket.StatusClosed.CanTransitionTo(domainTicket.StatusClosed))
		assert.True(t, domainTicket.StatusClosed.CanTransitionTo(domainTicket.StatusOpen))
		assert.True(t, domainTicket.StatusClosed.CanTransitionTo(domainTicket.StatusInProgress))
		assert.False(t, domainTicket.StatusClosed.CanTransitionTo(domainTicket.StatusResolved))
	})
}
