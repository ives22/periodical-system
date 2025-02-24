package periodical

// 存放返Response的数据格式，统一放在这里方便管理

// ResponsePeriodicalSet 查询期刊列列表返回参数
type ResponsePeriodicalSet struct {
	Total int64         `json:"total"` // 总数
	Items []*Periodical `json:"items"` // 期刊列表
}

func NewResponsePeriodicalSet() *ResponsePeriodicalSet {
	return &ResponsePeriodicalSet{
		Items: []*Periodical{},
	}
}

type ResponseCategorizeSet struct {
	Total int64         `json:"total"`
	Items []*Categorize `json:"items"`
}

func NewResponseCategorizeSet() *ResponseCategorizeSet {
	return &ResponseCategorizeSet{
		Items: []*Categorize{},
	}
}
