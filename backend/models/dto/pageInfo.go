package dto

type PageInfo struct {
	Size  int64       `json:"size"`
	Total int64       `json:"total"`
	List  interface{} `json:"list"`
}

// PageQuery 请求一个页面数据的dto对象
type PageQuery struct {
	Query        interface{}
	Page         int
	PageSize     int
	SortProperty string
	SortRule     string
}
