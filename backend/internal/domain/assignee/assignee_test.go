package assignee_test

import (
	"strings"
	"testing"
	"time"

	domainAssignee "github.com/adriancf/demo-sistema-chamado/backend/internal/domain/assignee"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewAssignee(t *testing.T) {
	now := time.Now().UTC()
	id, err := uuid.NewV7()
	require.NoError(t, err)

	t.Run("should instantiate assignee with valid fields and trimmed name", func(t *testing.T) {
		a := domainAssignee.NewAssignee(id, "  Alice Smith  ", now)
		require.NotNil(t, a)
		assert.Equal(t, id, a.ID)
		assert.Equal(t, "Alice Smith", a.Name)
		assert.Equal(t, now, a.CreatedAt)
		assert.Equal(t, now, a.UpdatedAt)
	})

	t.Run("should handle zero time by defaulting to non-zero UTC time", func(t *testing.T) {
		a := domainAssignee.NewAssignee(id, "Bob Jones", time.Time{})
		require.NotNil(t, a)
		assert.False(t, a.CreatedAt.IsZero())
		assert.False(t, a.UpdatedAt.IsZero())
		assert.Equal(t, time.UTC, a.CreatedAt.Location())
	})
}

func TestAssigneeValidate(t *testing.T) {
	now := time.Now().UTC()
	validID, err := uuid.NewV7()
	require.NoError(t, err)

	baseAssignee := func() domainAssignee.Assignee {
		return domainAssignee.Assignee{
			ID:        validID,
			Name:      "Carlos Silva",
			CreatedAt: now,
			UpdatedAt: now,
		}
	}

	edgeCases := []struct {
		name        string
		modify      func(a *domainAssignee.Assignee)
		expectedErr error
	}{
		// Valid base case
		{
			name:        "valid assignee with standard name and matching timestamps",
			modify:      nil,
			expectedErr: nil,
		},
		// ID edge cases
		{
			name: "ID: nil UUID should fail with ErrInvalidAssignee",
			modify: func(a *domainAssignee.Assignee) {
				a.ID = uuid.Nil
			},
			expectedErr: domainAssignee.ErrInvalidAssignee,
		},
		{
			name: "ID: valid UUID v7 should succeed",
			modify: func(a *domainAssignee.Assignee) {
				v7ID, _ := uuid.NewV7()
				a.ID = v7ID
			},
			expectedErr: nil,
		},
		// Name boundary and content edge cases
		{
			name: "Name: empty string should fail with ErrInvalidAssignee",
			modify: func(a *domainAssignee.Assignee) {
				a.Name = ""
			},
			expectedErr: domainAssignee.ErrInvalidAssignee,
		},
		{
			name: "Name: only whitespace should fail with ErrInvalidAssignee",
			modify: func(a *domainAssignee.Assignee) {
				a.Name = "   \t\n  "
			},
			expectedErr: domainAssignee.ErrInvalidAssignee,
		},
		{
			name: "Name: shorter than 2 chars (1 char) should fail with ErrInvalidAssignee",
			modify: func(a *domainAssignee.Assignee) {
				a.Name = "A"
			},
			expectedErr: domainAssignee.ErrInvalidAssignee,
		},
		{
			name: "Name: trimmed length shorter than 2 chars should fail with ErrInvalidAssignee",
			modify: func(a *domainAssignee.Assignee) {
				a.Name = "   A   "
			},
			expectedErr: domainAssignee.ErrInvalidAssignee,
		},
		{
			name: "Name: exactly 2 chars (minimum boundary) should succeed",
			modify: func(a *domainAssignee.Assignee) {
				a.Name = "Ed"
			},
			expectedErr: nil,
		},
		{
			name: "Name: exactly 100 chars (maximum boundary) should succeed",
			modify: func(a *domainAssignee.Assignee) {
				a.Name = strings.Repeat("A", 100)
			},
			expectedErr: nil,
		},
		{
			name: "Name: exceeds 100 chars (101 chars) should fail with ErrInvalidAssignee",
			modify: func(a *domainAssignee.Assignee) {
				a.Name = strings.Repeat("A", 101)
			},
			expectedErr: domainAssignee.ErrInvalidAssignee,
		},
		{
			name: "Name: special characters, accents, numbers, and symbols should succeed",
			modify: func(a *domainAssignee.Assignee) {
				a.Name = "José & André (DevOps #1) [Tier-3]"
			},
			expectedErr: nil,
		},
		// Timestamp edge cases (UpdatedAt vs CreatedAt)
		{
			name: "Timestamp: UpdatedAt earlier than CreatedAt should fail with ErrInvalidTimestamp",
			modify: func(a *domainAssignee.Assignee) {
				a.CreatedAt = now
				a.UpdatedAt = now.Add(-10 * time.Minute)
			},
			expectedErr: domainAssignee.ErrInvalidTimestamp,
		},
		{
			name: "Timestamp: UpdatedAt equal to CreatedAt should succeed",
			modify: func(a *domainAssignee.Assignee) {
				a.CreatedAt = now
				a.UpdatedAt = now
			},
			expectedErr: nil,
		},
		{
			name: "Timestamp: UpdatedAt later than CreatedAt should succeed",
			modify: func(a *domainAssignee.Assignee) {
				a.CreatedAt = now
				a.UpdatedAt = now.Add(10 * time.Minute)
			},
			expectedErr: nil,
		},
	}

	for _, tc := range edgeCases {
		t.Run(tc.name, func(t *testing.T) {
			a := baseAssignee()
			if tc.modify != nil {
				tc.modify(&a)
			}

			err := a.Validate()
			if tc.expectedErr != nil {
				require.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedErr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
