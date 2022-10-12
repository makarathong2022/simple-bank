package model

type Response struct {
	HasNext  bool         `json:"has_next"`
	PageSize int32        `json:"page_size"`
	Page     int32        `json:"page"`
	Total    int          `json:"total"`
	Result   any `json:"result"`
}
