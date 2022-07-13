package types

type PaginationReq struct {
	Limit   *int64  `json:"limit,omitempty"`
	Skip    *int64  `json:"skip,omitempty"`
	OrderBy *string `json:"order_by,omitempty"`
}
