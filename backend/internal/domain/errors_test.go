package domain_test

import (
	"errors"
	"testing"

	"github.com/adriancf/demo-sistema-chamado/backend/internal/domain"
	"github.com/stretchr/testify/assert"
)

func TestDomainErrors(t *testing.T) {
	t.Run("should wrap base error and keep message and errors.Is compatibility", func(t *testing.T) {
		errNotFound := domain.NewError("custom item not found", domain.ErrNotFound)

		assert.Equal(t, "custom item not found", errNotFound.Error())
		assert.True(t, errors.Is(errNotFound, domain.ErrNotFound))
		assert.False(t, errors.Is(errNotFound, domain.ErrValidation))
		assert.True(t, errors.Is(errNotFound, errNotFound))
	})

	t.Run("should correctly match validation and conflict errors", func(t *testing.T) {
		errVal := domain.NewError("invalid field", domain.ErrValidation)
		errConflict := domain.NewError("duplicate record", domain.ErrConflict)

		assert.True(t, errors.Is(errVal, domain.ErrValidation))
		assert.False(t, errors.Is(errVal, domain.ErrConflict))

		assert.True(t, errors.Is(errConflict, domain.ErrConflict))
		assert.False(t, errors.Is(errConflict, domain.ErrValidation))
	})
}
