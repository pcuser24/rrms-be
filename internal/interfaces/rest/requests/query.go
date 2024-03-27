package requests

type SearchSortPaginationQuery struct {
	Limit  *int32   `query:"limit" validate:"omitempty,gte=0"`
	Offset *int32   `query:"offset" validate:"omitempty,gte=0"`
	SortBy []string `query:"sortby" validate:"omitempty"`
	Order  []string `query:"order" validate:"omitempty,dive,oneof=asc desc"`
}
