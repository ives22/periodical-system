package periodical

import "context"

// Service 定义期刊的接口
type Service interface {
	// CreatePeriodical 创建期刊
	CreatePeriodical(context.Context, *CreatePeriodicalRequest) (*Periodical, error)
	// QueryPeriodical 查询期刊列表
	QueryPeriodical(context.Context, *QueryPeriodicalRequest) (*ResponsePeriodicalSet, error)
	// DescribePeriodical 查询期刊详情
	DescribePeriodical(context.Context, *DescribePeriodicalRequest) (*Periodical, error)
	// UpdatePeriodical 更新期刊
	UpdatePeriodical(context.Context, *UpdatePeriodicalRequest) (*Periodical, error)
	// DeletePeriodical 删除期刊
	DeletePeriodical(context.Context, *DeletePeriodicalRequest) error
	// QueryCategorize 查询分类列表
	QueryCategorize(context.Context, *QueryCategorizeRequest) (*ResponseCategorizeSet, error)
}

// CreatePeriodicalRequest 创建期刊的请求参数
type CreatePeriodicalRequest struct {
	*BasePeriodical         // 期刊基础信息
	CategorizeIDs   []int64 `json:"categorize_ids"` // 类型ID列表

	// 下面传的都是关联关系的字段
	//DirectionList  []int64 // 收稿方向列表 （传ID）
	//CycleList      []int64 // 出版周期列表 （传ID）
	//LevelList      []int64 // 期刊级别列表 （传ID）
	//CollectionList []int64 // 网站收录列表（传ID）
}

// QueryPeriodicalRequest 查询期刊列表的请求参数
type QueryPeriodicalRequest struct {
	PageSize                 int            `json:"page_size"`                  // 查询的页数
	PageNumber               int            `json:"page_number"`                // 每页显示的个数
	HasBaseTypeList          []*hasBaseType `json:"has_base_type_list"`         // 类型列表
	Invoice                  int            `json:"invoice"`                    // 是否可以开发票, 根据是否可开杂志社发票搜索，1 搜索 true，2 搜索 false，3 不加该条件搜索
	Name                     string         `json:"name"`                       // 期刊名字   根据期刊名字搜索
	PublishTimeStart         string         `json:"publish_time_start"`         // 出刊时间，根据出刊时间筛选   2024-12-01 （上旬）
	PublishTimeEnd           string         `json:"publish_time_end"`           // 出刊时间，根据出刊时间筛选   2024-12-10 （上旬）
	DomesticNumber           string         `json:"domestic_number"`            // 国内刊号   根据国内刊号搜索
	InternationalNumber      string         `json:"international_number"`       // 国际刊号  根据国际刊号搜索
	CompetentUnit            string         `json:"competent_unit"`             // 主管单位， 根据主管单位搜索
	Organization             string         `json:"organization"`               // 主办单位  根据主办单位搜索
	Works                    int            `json:"works"`                      // 是否是作品，根据是否是作品搜索 1 搜索 true，2 搜索 false，3 不加该条件搜索
	PeriodicalBatch          []int64        `json:"periodical_batch"`           // 期刊批次，1、2  第一批次；第二批次
	CompositeInfluenceFactor float64        `json:"composite_influence_factor"` // 复合影响因子，根据复合影响因子筛选 -99 表示0以上
	ColumnSetting            string         `json:"column_setting"`             // 栏目设置搜索
}

// hasBaseType 查询期刊时候，所携带的分类信息， 通用分类的ID和 类型（type）
type hasBaseType struct {
	BaseId int `json:"base_id"` // 对应categorizes的categorize_id字段，即通用分类的id
	Type   int `json:"type"`    // 对应categorizes的type字段, 即通用分类的所属type
}

// UpdatePeriodicalRequest 更新期刊时候传递的参数
type UpdatePeriodicalRequest struct {
	PeriodicalID string     `json:"periodical_id"`
	UpdateMode   UpdateMode `json:"update_mode"`
	*CreatePeriodicalRequest
}

// DescribePeriodicalRequest 查询期刊详情的请求参数
type DescribePeriodicalRequest struct {
	PeriodicalID string `json:"periodical_id"` // 期刊ID，根据期刊id查询期刊详情
}

type DeletePeriodicalRequest struct {
	PeriodicalID string `json:"periodical_id"` // 期刊ID，根据期刊id进行删除
}

func NewCreatePeriodicalRequest() *CreatePeriodicalRequest {
	return &CreatePeriodicalRequest{}
}

func NewUpdatePeriodicalRequest(id string) *UpdatePeriodicalRequest {
	return &UpdatePeriodicalRequest{
		PeriodicalID:            id,
		UpdateMode:              UPDATE_MODE_PUT,
		CreatePeriodicalRequest: NewCreatePeriodicalRequest(),
	}
}

func NewQueryPeriodicalRequest() *QueryPeriodicalRequest {
	return &QueryPeriodicalRequest{
		HasBaseTypeList: []*hasBaseType{},
	}
}

func NewDescribePeriodicalRequest(id string) *DescribePeriodicalRequest {
	return &DescribePeriodicalRequest{
		PeriodicalID: id,
	}
}

func NewDeletePeriodicalRequest(id string) *DeletePeriodicalRequest {
	return &DeletePeriodicalRequest{
		PeriodicalID: id,
	}
}

func (r *QueryPeriodicalRequest) Offset() int {
	return r.PageSize * (r.PageNumber - 1)
}

//=============== Categorize 相关请求参数结构体

// QueryCategorizeRequest 查询分类的请求参数
type QueryCategorizeRequest struct {
	PageSize   int `json:"page_size"`   // 查询的页数
	PageNumber int `json:"page_number"` // 每页显示的个数
	Status     int `json:"status"`      // 状态，默认1 启用
}

func NewQueryCategorizeRequest() *QueryCategorizeRequest {
	return &QueryCategorizeRequest{}
}

func (r *QueryCategorizeRequest) Offset() int {
	return r.PageSize * (r.PageNumber - 1)
}
