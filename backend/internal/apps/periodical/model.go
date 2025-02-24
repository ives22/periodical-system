package periodical

import (
	"periodical/internal/common"
)

// BasePeriodical 期刊的基础信息
type BasePeriodical struct {
	Name            string `json:"name"`             // 期刊名字
	PublicationTime string `json:"publication_time"` // 出刊时间   (2024年-12月)
	//PublishTime           string `json:"publish_time"`           // 出刊时间 （对应上旬、中旬、下旬、月底等）
	UpdatedPeriodicalAt   int64  `json:"updated_periodical_at"`  // 更新时间   1732863779
	CompetentUnit         string `json:"competent_unit"`         // 主管单位    黑龙江省农垦总局
	Organization          string `json:"organization"`           // 主办单位    农垦日报社
	DomesticNumber        string `json:"domestic_number"`        // 国内刊号    23-1558/C
	InternationalNumber   string `json:"international_number"`   // 国际刊号    1674-7879
	FoundingTime          string `json:"founding_time"`          // 创刊时间    2009年
	CitationRate          string `json:"citation_rate"`          // 引用率      无
	PeriodicalDescription string `json:"periodical_description"` // 期刊描述    北大荒文化》省级旬刊，龙源新媒体xxxx
	ColumnSetting         string `json:"column_setting"`         // 栏目设置    文史在线、文学研究、文学文化、艺术设xxx）
	AuditTime             string `json:"audit_time"`             // 审核周期    当天
	InvoiceReceipt        bool   `json:"invoice_receipt"`        // 是否可开杂志社发票  false 0 否， true 1 可以
	AttentionMatter       string `json:"attention_matter"`       // 文章注意事项， 注意事项
	ArticleNaming         string `json:"article_naming"`         // 文章命名
	InternalProcesses     string `json:"internal_processes"`     // 内部流程

	Price                    float64 `json:"price"`                      // 期刊发表费用（服务费）
	IsWarp                   bool    `json:"is_warp"`                    // 是否全包（单选框） false 0 否， true 1 是
	PeriodicalPage           string  `json:"periodical_page"`            // 期刊页码  398页
	CompositeInfluenceFactor float64 `json:"composite_influence_factor"` // 复合影响因子
	SyntheticInfluenceFactor float64 `json:"synthetic_influence_factor"` // 综合影响因子
	PeriodicalBatch          int64   `json:"periodical_batch"`           //期刊批次  1 or 2
	Works                    bool    `json:"works"`                      // 可发作品    false 0 ， true 1 可发
}

// Periodical 期刊表
//type Periodical struct {
//	PeriodicalId int64 `json:"periodical_id"` // 期刊ID
//	*BasePeriodical
//	*common.Meta // 通用字段
//	//Categories   []*CommonCategories `gorm:"-"` // 使用 gorm:"-" 标记这个字段不需要持久化
//	//Categories []*CommonCategories `gorm:"many2many:periodical_common_categories;foreignKey:periodical_id;joinForeignKey:periodical_id;references:common_id;joinReferences:common_id"` // 使用 gorm:"-" 标记这个字段不需要持久化
//	Categories []*CommonCategories `gorm:"many2many:periodical_common_categories;joinForeignKey:periodical_id;References:common_id"` // 使用 gorm:"-" 标记这个字段不需要持久化
//}

// Periodical 期刊表
type Periodical struct {
	PeriodicalID    int           `gorm:"primaryKey" json:"periodical_id"`
	*BasePeriodical               // 加载上面定义个的期刊表字段
	*common.Meta                  // 通用信息（创建时间、更新时间），
	Categorizes     []*Categorize `gorm:"many2many:periodical_categorize;foreignKey:PeriodicalID;joinForeignKey:PeriodicalID;References:CategorizeID;JoinReferences:CategorizeID" json:"categorizes"`
}

// Categorize 通用分类表
type Categorize struct {
	CategorizeID int           `gorm:"primaryKey;autoIncrement;column:categorize_id" json:"categorize_id"`
	ParentId     int64         `json:"parent_id"`
	Name         string        `json:"name"`
	Type         int64         `json:"type"`
	Status       int64         `json:"status"`
	AppStatus    int64         `json:"app_status"`
	Sort         int64         `json:"sort"`
	Link         string        `json:"link"`
	ExtraInfo    string        `json:"extra_info"`
	*common.Meta               // 通用信息（创建时间、更新时间），
	Periodicals  []*Periodical `gorm:"many2many:periodical_categorize;foreignKey:CategorizeID;joinForeignKey:CategorizeID;References:PeriodicalID;JoinReferences:PeriodicalID" json:"periodicals"`
}

// PeriodicalCategorize 期刊表 和 通用分类表的关联表
// PeriodicalCategorize represents the `periodical_categorize` table
type PeriodicalCategorize struct {
	CategorizeID int `gorm:"primaryKey;column:categorize_id" json:"categorize_id"`
	PeriodicalID int `gorm:"primaryKey;column:periodical_id" json:"periodical_id"`
}

func (p *PeriodicalCategorize) TableName() string {
	return "periodical_categorize"
}

func NewPeriodical(req *CreatePeriodicalRequest) *Periodical {
	return &Periodical{
		BasePeriodical: req.BasePeriodical,
		Meta:           common.NewMeta(),
	}
}

//func NewPeriodical(req *CreatePeriodicalRequest) *Periodical {
//	return &Periodical{
//		BasePeriodical: req.BasePeriodical,
//		Meta:           common.NewMeta(),
//	}
//}

//// CommonCategories 存储了收稿方向、期刊级别、网络收录、出版周期等信息
//type CommonCategories struct {
//	CommonId  int64     `json:"common_id"`
//	ParentId  int64     `json:"parent_id"`
//	Name      string    `json:"name"`
//	Type      string    `json:"type"`
//	Status    int64     `json:"status"`
//	AppStatus int64     `json:"app_status"`
//	SortOrder int64     `json:"sort_order"`
//	Link      string    `json:"link"`
//	ExtraInfo string    `json:"extra_info"`
//	CreatedAt time.Time `json:"created_at"`
//	UpdatedAt time.Time `json:"updated_at"`
//}
//
//// PeriodicalCommonCategories 关联表模型
//type PeriodicalCommonCategories struct {
//	PeriodicalId int64 `json:"periodical_id"` // 期刊ID
//	CommonId     int64 `json:"common_id"`     // 通用分类ID
//}

//func NewCommonCategories() *CommonCategories {
//	return &CommonCategories{}
//}

//func (p *Periodical) TableName() string {
//	return "periodical"
//}
//func (c *CommonCategories) TableName() string {
//	return "common_categories"
//}
