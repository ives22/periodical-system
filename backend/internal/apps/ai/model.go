package ai

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"periodical/internal/common"
	"strings"
)

// AIModelConfig 大模型渠道配置表
type AIModelConfig struct {
	Id           int64   `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Name         string  `gorm:"column:name;type:varchar(100);not null" json:"name"`                  // 渠道显示名称，如：通义千问 Qwen-Plus
	Provider     string  `gorm:"column:provider;type:varchar(50);not null" json:"provider"`           // 供应商标识：qwen / deepseek / openai / custom
	BaseURL      string  `gorm:"column:base_url;type:varchar(255);not null" json:"base_url"`          // OpenAI兼容接口地址
	APIKey       string  `gorm:"column:api_key;type:varchar(255);default:''" json:"-"`                // 实际密钥（JSON序列化隐藏）
	APIKeyMasked string  `gorm:"-" json:"api_key_masked"`                                             // 前端脱敏展示密钥
	HasAPIKey    bool    `gorm:"-" json:"has_api_key"`                                                // 是否已配置密钥
	ModelName    string  `gorm:"column:model_name;type:varchar(100);not null" json:"model_name"`      // 模型标识，如：qwen-plus / deepseek-chat / gpt-4o
	SystemPrompt string  `gorm:"column:system_prompt;type:text" json:"system_prompt"`                 // 附加系统提示词
	Temperature  float64 `gorm:"column:temperature;type:decimal(3,2);default:0.20" json:"temperature"` // 采样温度
	MaxTokens    int     `gorm:"column:max_tokens;default:4096" json:"max_tokens"`                    // 最大输出Token
	IsDefault    int     `gorm:"column:is_default;default:0" json:"is_default"`                       // 是否全局默认模型：1-是，0-否
	Status       int     `gorm:"column:status;default:1" json:"status"`                               // 状态：1-启用，0-禁用
	Remark       string  `gorm:"column:remark;type:varchar(255);default:''" json:"remark"`            // 备注说明
	*common.Meta
}

func (*AIModelConfig) TableName() string {
	return "ai_model_configs"
}

// MaskKey 生成脱敏密钥信息
func (m *AIModelConfig) MaskKey() {
	key := strings.TrimSpace(m.APIKey)
	if key == "" {
		m.HasAPIKey = false
		m.APIKeyMasked = "未配置密钥"
		return
	}
	m.HasAPIKey = true
	if len(key) <= 8 {
		m.APIKeyMasked = "****"
		return
	}
	m.APIKeyMasked = key[:4] + "****" + key[len(key)-4:]
}

// AIChatSession AI 对话会话表
type AIChatSession struct {
	Id           int64  `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	SessionID    string `gorm:"column:session_id;type:varchar(64);uniqueIndex;not null" json:"session_id"`
	UserID       int64  `gorm:"column:user_id;index;not null" json:"user_id"`
	Username     string `gorm:"column:username;type:varchar(64);not null" json:"username"`
	Title        string `gorm:"column:title;type:varchar(200);not null" json:"title"`
	Scene        string `gorm:"column:scene;type:varchar(32);default:'workspace'" json:"scene"` // workspace | copilot
	ModelName    string `gorm:"column:model_name;type:varchar(100);default:''" json:"model_name"`
	MessageCount int    `gorm:"column:message_count;default:0" json:"message_count"`
	IsPinned     int    `gorm:"column:is_pinned;default:0" json:"is_pinned"`
	*common.Meta
}

func (*AIChatSession) TableName() string {
	return "ai_chat_sessions"
}

// ToolStep 记录一次 Agent 工具调用轨迹
type ToolStep struct {
	ToolName    string `json:"tool_name"`              // 工具名称
	Title       string `json:"title"`                  // 友好中文标题
	ArgsSummary string `json:"args_summary,omitempty"` // 调用入参摘要
	Summary     string `json:"summary"`                // 执行结果摘要
	DurationMs  int64  `json:"duration_ms,omitempty"`  // 执行耗时（毫秒）
	Status      string `json:"status"`                 // success | error | denied
}

type ToolStepSlice []ToolStep

func (t ToolStepSlice) Value() (driver.Value, error) {
	if len(t) == 0 {
		return "[]", nil
	}
	b, err := json.Marshal(t)
	return string(b), err
}

func (t *ToolStepSlice) Scan(value interface{}) error {
	if value == nil {
		*t = []ToolStep{}
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to scan ToolStepSlice: %v", value)
	}
	if len(bytes) == 0 {
		*t = []ToolStep{}
		return nil
	}
	return json.Unmarshal(bytes, t)
}

// FieldDiff 记录修改前后的字段对比
type FieldDiff struct {
	Field    string `json:"field"`
	Label    string `json:"label"`
	OldValue string `json:"old_value"`
	NewValue string `json:"new_value"`
}

// BatchUpdateTarget 多本期刊批量修改时的单本目标项
type BatchUpdateTarget struct {
	PeriodicalID   int                    `json:"periodical_id"`
	PeriodicalName string                 `json:"periodical_name"`
	Level          string                 `json:"level,omitempty"`
	FormPayload    map[string]interface{} `json:"form_payload"`
	Diffs          []FieldDiff            `json:"diffs"`
}

// PeriodicalBriefCard 检索结果卡片中的期刊摘要项
type PeriodicalBriefCard struct {
	PeriodicalID          int      `json:"periodical_id"`
	Name                  string   `json:"name"`
	Level                 string   `json:"level"`
	Directions            string   `json:"directions"`
	Collections           []string `json:"collections"`
	Cycle                 string   `json:"cycle"`
	PublicationTime       string   `json:"publication_time"`
	Price                 string   `json:"price"`
	PeriodicalPage        string   `json:"periodical_page"`
	IsWarp                bool     `json:"is_warp"`
	InvoiceReceipt        bool     `json:"invoice_receipt"`
	Works                 bool     `json:"works"`
	IsHidden              bool     `json:"is_hidden"`
	Organization          string   `json:"organization"`
	CompetentUnit         string   `json:"competent_unit"`
	DomesticNumber        string   `json:"domestic_number"`
	InternationalNumber   string   `json:"international_number"`
	AuditTime             string   `json:"audit_time"`
	PeriodicalDescription string   `json:"periodical_description"`
}

// ActionCardData 人机协同交互卡片（查询结果卡片 / 新增确认卡片 / 修改Diff确认卡片 / 删除高危确认卡片）
type ActionCardData struct {
	CardType       string                 `json:"card_type"`                 // query_result | create_confirm | update_confirm | delete_confirm
	Status         string                 `json:"status"`                    // pending | confirmed | cancelled | readonly
	Title          string                 `json:"title"`                     // 卡片标题
	Description    string                 `json:"description"`               // 补充说明
	TotalCount     int64                  `json:"total_count,omitempty"`     // 查询总数（query_result）
	Items          []PeriodicalBriefCard  `json:"items,omitempty"`           // 查询列表（query_result）
	FilterPayload  map[string]interface{} `json:"filter_payload,omitempty"`  // 可同步给前端列表页的筛选条件（query_result）
	PeriodicalID   int                    `json:"periodical_id,omitempty"`   // 目标期刊 ID（update/delete）
	PeriodicalName string                 `json:"periodical_name,omitempty"` // 目标期刊名称
	FormPayload    map[string]interface{} `json:"form_payload,omitempty"`    // 待创建或待更新的完整表单数据（create/update）
	CategoryLabels map[string][]string    `json:"category_labels,omitempty"` // 匹配到的分类中文名展示（directions/collections/levels/cycle）
	Diffs          []FieldDiff            `json:"diffs,omitempty"`           // 修改前后字段对比（update_confirm）
	BatchUpdates   []BatchUpdateTarget    `json:"batch_updates,omitempty"`   // 多本期刊批量修改时的目标列表（update_confirm）
	Suggestions    []string               `json:"suggestions,omitempty"`     // 上下文智能追问建议胶囊
	ExecutedAt     int64                  `json:"executed_at,omitempty"`     // 确认执行时间戳
	ResultMessage  string                 `json:"result_message,omitempty"`  // 确认执行结果反馈
}

func (a *ActionCardData) Value() (driver.Value, error) {
	if a == nil {
		return nil, nil
	}
	b, err := json.Marshal(a)
	return string(b), err
}

func (a *ActionCardData) Scan(value interface{}) error {
	if value == nil {
		return nil
	}
	var bytes []byte
	switch v := value.(type) {
	case []byte:
		bytes = v
	case string:
		bytes = []byte(v)
	default:
		return fmt.Errorf("failed to scan ActionCardData: %v", value)
	}
	if len(bytes) == 0 {
		return nil
	}
	return json.Unmarshal(bytes, a)
}

// AIChatMessage AI 对话消息记录表
type AIChatMessage struct {
	Id         int64           `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	SessionID  string          `gorm:"column:session_id;type:varchar(64);index;not null" json:"session_id"`
	UserID     int64           `gorm:"column:user_id;index;not null" json:"user_id"`
	Username   string          `gorm:"column:username;type:varchar(64);not null" json:"username"`
	Role       string          `gorm:"column:role;type:varchar(32);not null" json:"role"` // user | assistant
	Content    string          `gorm:"column:content;type:longtext" json:"content"`
	ModelName  string          `gorm:"column:model_name;type:varchar(100);default:''" json:"model_name"`
	ToolSteps  ToolStepSlice   `gorm:"column:tool_steps;type:json" json:"tool_steps"`
	ActionCard *ActionCardData `gorm:"column:action_card;type:json" json:"action_card"`
	*common.Meta
}

func (*AIChatMessage) TableName() string {
	return "ai_chat_messages"
}
