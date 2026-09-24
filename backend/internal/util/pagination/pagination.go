package pagination

import (
	"math"
	"strings"
)

type SortOrder string

const (
	OrderAsc  SortOrder = "asc"
	OrderDesc SortOrder = "desc"
)

func ParseSortOrder(s string) SortOrder {
	if strings.ToLower(strings.TrimSpace(s)) == "asc" {
		return OrderAsc
	}
	return OrderDesc
}

func (s SortOrder) Direction() string {
	if s == OrderAsc {
		return "ASC"
	}
	return "DESC"
}

type SortOption struct {
	Field string    `json:"field"`
	Order SortOrder `json:"order"`
}

type PaginationRequest struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
}

func (p *PaginationRequest) Normalize() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}
}

func (p PaginationRequest) Offset() int {
	page := p.Page
	if page <= 0 {
		page = 1
	}
	pageSize := p.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	return (page - 1) * pageSize
}

func (p PaginationRequest) Limit() int {
	if p.PageSize <= 0 {
		return 10
	}
	return p.PageSize
}

type PaginatedResult[T any] struct {
	Items       []T `json:"items"`
	TotalItems  int `json:"total_items"`
	TotalPages  int `json:"total_pages"`
	CurrentPage int `json:"current_page"`
	PageSize    int `json:"page_size"`
}

func NewPaginatedResult[T any](items []T, totalItems int, page, pageSize int) PaginatedResult[T] {
	if pageSize <= 0 {
		pageSize = 10
	}
	if page <= 0 {
		page = 1
	}
	totalPages := 0
	if totalItems > 0 {
		totalPages = int(math.Ceil(float64(totalItems) / float64(pageSize)))
	}
	if items == nil {
		items = []T{}
	}
	return PaginatedResult[T]{
		Items:       items,
		TotalItems:  totalItems,
		TotalPages:  totalPages,
		CurrentPage: page,
		PageSize:    pageSize,
	}
}
