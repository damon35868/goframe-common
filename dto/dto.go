package dto

type PageReqDto struct {
	Page     int `json:"page" v:"required|min:1#必须输入分页"`
	PageSize int `json:"pageSize" v:"required|min:1#必须输入分页量"`
}
type PageResDto[T any] struct {
	Items       []T  `json:"items"`       // 分页数据
	HasNextPage bool `json:"hasNextPage"` // 当前条件下是否还有更多分页
	TotalCount  int  `json:"totalCount"`  // 满足条件下的全部数据总条数
}
