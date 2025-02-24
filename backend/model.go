package main

import "time"

// Periodical  期刊信息
type Periodical struct {
	PeriodicalId          string      `json:"periodical_id"`          // 期刊ID
	Name                  string      `json:"name"`                   // 期刊名字
	PublicationTime       time.Time   `json:"publication_time"`       // 出刊时间
	UpdatedPeriodicalAt   time.Time   `json:"updated_periodical_at"`  // 更新时间
	DirectionList         interface{} `json:"direction_list"`         // 收稿方向   数据类型还么定
	CompetentUnit         string      `json:"competent_unit"`         // 主管单位
	Organization          string      `json:"organization"`           // 主办单位
	DomesticNumber        string      `json:"domestic_number"`        // 国内刊号
	InternationalNumber   string      `json:"international_number"`   // 国际刊号
	PublishingCycle       interface{} `json:"publishing_cycle"`       // 期刊周期（出版周期）
	FoundingTime          string      `json:"founding_time"`          // 创刊时间
	CitationRate          string      `json:"citation_rate"`          // 引用率
	PeriodicalDescription string      `json:"periodical_description"` // 期刊描述
	ColumnSetting         string      `json:"column_setting"`         // 栏目设置
	InfluenceFactor       int         `json:"influence_factor"`       // 影响因子
	AuditTime             string      `json:"audit_time"`             // 审核周期
	InvoiceReceipt        bool        `json:"invoice_receipt"`        // 是否可开杂志社发票
	AttentionMatter       string      `json:"attention_matter"`       // 文章注意事项， 注意事项
	ArticleNaming         string      `json:"article_naming"`         // 文章命名
	InternalProcesses     string      `json:"internal_processes"`     // 内部流程

}
