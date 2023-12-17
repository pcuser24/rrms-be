package responses

type SearchSortPaginationResponse struct {
	TotalCount int32         `json:"_totalCount"`
	Limit      int32         `json:"_limit"`
	Offset     int32         `json:"_offset"`
	SortBy     string        `json:"_sortby"`
	Order      string        `json:"_order"`
	Items      []interface{} `json:"items"`
}
