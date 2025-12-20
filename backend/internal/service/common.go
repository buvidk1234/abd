package service

type PagedParams struct {
	Page     int `form:"page" json:"page"`
	PageSize int `form:"pageSize" json:"pageSize"`
}

type PagedResp[T any] struct {
	Page     int `json:"page"`
	Total    int `json:"total"`
	PageSize int `json:"pageSize"`
	Data     []T `json:"data"`
}
