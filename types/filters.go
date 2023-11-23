package types

import "errors"

const (
	MaxPageSize = 10
	MinPageSize = 1

	AscSort  = "asc"
	DescSort = "desc"

	DefaultPage     = 1
	DefaultPageSize = 5
	DefaultSort     = AscSort
	DefaultSortBy   = "firstName"
)

type UsersPaginationFilter struct {
	Page     int `form:"page"`
	PageSize int `form:"pageSize"`

	SortBy  string `form:"sortBy"`
	SortDir string `form:"sortDir"`
}

func NewUsersPaginationFilter() UsersPaginationFilter {
	return UsersPaginationFilter{
		Page:     DefaultPage,
		PageSize: DefaultPageSize,
		SortBy:   DefaultSortBy,
		SortDir:  DefaultSort,
	}
}

func (f *UsersPaginationFilter) Validate() error {
	if f.Page < 1 {
		return errors.New("page must be 1 or greater")
	}

	if f.PageSize < MinPageSize {
		return errors.New("pageSize must be greater than 0")
	}

	if f.PageSize > MaxPageSize {
		return errors.New("pageSize exceeds maximum limit")
	}

	if !(f.SortDir == AscSort || f.SortDir == DescSort) {
		return errors.New("sortDir must be 'asc' or 'desc'")
	}

	return nil
}
