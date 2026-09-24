package pagination_test

import (
	"testing"

	"github.com/adriancf/demo-sistema-chamado/backend/internal/util/pagination"
	"github.com/stretchr/testify/assert"
)

func TestPagination(t *testing.T) {
	t.Run("should parse sort orders correctly", func(t *testing.T) {
		assert.Equal(t, pagination.OrderAsc, pagination.ParseSortOrder("asc"))
		assert.Equal(t, pagination.OrderAsc, pagination.ParseSortOrder("ASC"))
		assert.Equal(t, pagination.OrderDesc, pagination.ParseSortOrder("desc"))
		assert.Equal(t, pagination.OrderDesc, pagination.ParseSortOrder("other"))
		assert.Equal(t, "ASC", pagination.OrderAsc.Direction())
		assert.Equal(t, "DESC", pagination.OrderDesc.Direction())
	})

	t.Run("should calculate offset and limit with defaults", func(t *testing.T) {
		req := pagination.PaginationRequest{Page: 0, PageSize: 0}
		assert.Equal(t, 0, req.Offset())
		assert.Equal(t, 10, req.Limit())

		req.Normalize()
		assert.Equal(t, 1, req.Page)
		assert.Equal(t, 10, req.PageSize)

		req2 := pagination.PaginationRequest{Page: 3, PageSize: 20}
		assert.Equal(t, 40, req2.Offset())
		assert.Equal(t, 20, req2.Limit())
	})

	t.Run("should create paginated result with correct pages calculation", func(t *testing.T) {
		items := []string{"item1", "item2"}
		result := pagination.NewPaginatedResult(items, 25, 1, 10)

		assert.Equal(t, 25, result.TotalItems)
		assert.Equal(t, 3, result.TotalPages)
		assert.Equal(t, 1, result.CurrentPage)
		assert.Equal(t, 10, result.PageSize)
		assert.Len(t, result.Items, 2)
	})

	t.Run("should handle zero total items in paginated result", func(t *testing.T) {
		result := pagination.NewPaginatedResult[string](nil, 0, 0, 0)
		assert.Equal(t, 0, result.TotalItems)
		assert.Equal(t, 0, result.TotalPages)
		assert.Equal(t, 1, result.CurrentPage)
		assert.Equal(t, 10, result.PageSize)
		assert.Empty(t, result.Items)
	})
}
