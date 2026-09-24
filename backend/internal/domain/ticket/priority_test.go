package ticket_test

import (
	"testing"

	domainTicket "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/ticket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParsePriority(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expected    domainTicket.Priority
		expectError bool
	}{
		{"Should accept 'low'", "low", domainTicket.PriorityLow, false},
		{"Should accept 'medium'", "medium", domainTicket.PriorityMedium, false},
		{"Should accept 'high'", "high", domainTicket.PriorityHigh, false},
		{"Should accept with trailing spaces and uppercase '  HIGH '", "  HIGH ", domainTicket.PriorityHigh, false},
		{"Should accept mixed case 'Medium'", "Medium", domainTicket.PriorityMedium, false},
		{"Should reject invalid priority 'critical'", "critical", "", true},
		{"Should reject empty string", "", "", true},
		{"Should reject invalid numbers '1'", "1", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			priority, err := domainTicket.ParsePriority(tt.input)
			if tt.expectError {
				require.Error(t, err)
				assert.Equal(t, domainTicket.ErrInvalidPriority, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, priority)
				assert.True(t, priority.IsValid())
			}
		})
	}
}

func TestPriority_Weight(t *testing.T) {
	assert.Greater(t, domainTicket.PriorityHigh.Weight(), domainTicket.PriorityMedium.Weight())
	assert.Greater(t, domainTicket.PriorityMedium.Weight(), domainTicket.PriorityLow.Weight())
	assert.Equal(t, 0, domainTicket.Priority("unknown").Weight())
}
