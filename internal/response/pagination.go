package response

type Pagination struct {
	Page    int   `json:"page"`
	PerPage int   `json:"per_page"`
	Total   int64 `json:"total"`
}

type ListResult struct {
	Items      interface{} `json:"items"`
	Pagination Pagination  `json:"pagination"`
}

func List(items interface{}, page, per int, total int64) ListResult {
	return ListResult{
		Items: items,
		Pagination: Pagination{
			Page:    page,
			PerPage: per,
			Total:   total,
		},
	}
}
