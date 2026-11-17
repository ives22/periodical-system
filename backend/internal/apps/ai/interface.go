package ai

import (
	"context"
	"errors"
	"strings"
)

// Service 定义 AI Agent 与大模型管理服务接口
type Service interface {
	// ===== 大模型配置管理（超级管理员） =====
	QueryModelConfigs(ctx context.Context) ([]*AIModelConfig, error)
	GetActiveDefaultModel(ctx context.Context) (*AIModelConfig, error)
	SaveModelConfig(ctx context.Context, req *SaveModelConfigRequest) (*AIModelConfig, error)
	SetDefaultModel(ctx context.Context, id int64) error
	DeleteModelConfig(ctx context.Context, id int64) error
	TestModelConnection(ctx context.Context, id int64) (*TestConnectionResponse, error)
	FetchRemoteModels(ctx context.Context, req *FetchRemoteModelsRequest) (*FetchRemoteModelsResponse, error)

	// ===== 个人会话与对话管理 =====
	QueryUserSessions(ctx context.Context, userID int64, keyword string) ([]*AIChatSession, error)
	CreateSession(ctx context.Context, userID int64, username string, title string, scene string) (*AIChatSession, error)
	RenameSession(ctx context.Context, userID int64, sessionID string, title string, isPinned *int) error
	DeleteSession(ctx context.Context, userID int64, sessionID string) error
	QuerySessionMessages(ctx context.Context, userID int64, isSuperAdmin bool, sessionID string) ([]*AIChatMessage, error)

	// ===== 核心对话与人机协同写操作确认 =====
	SendMessage(ctx context.Context, req *SendMessageRequest) (*SendMessageResponse, error)
	ConfirmAction(ctx context.Context, req *ConfirmActionRequest) (*ConfirmActionResponse, error)

	// ===== 超级管理员后台全员对话审计 =====
	QueryAuditSessions(ctx context.Context, req *QueryAuditSessionsRequest) (*AuditSessionListResponse, error)
}

// SaveModelConfigRequest 新增或更新大模型配置请求
type SaveModelConfigRequest struct {
	Id           int64   `json:"id"`
	Name         string  `json:"name"`
	Provider     string  `json:"provider"`
	BaseURL      string  `json:"base_url"`
	APIKey       string  `json:"api_key"` // 为空或包含 **** 时在更新模式下保留原密钥
	ModelName    string  `json:"model_name"`
	SystemPrompt string  `json:"system_prompt"`
	Temperature  float64 `json:"temperature"`
	MaxTokens    int     `json:"max_tokens"`
	IsDefault    int     `json:"is_default"`
	Status       int     `json:"status"`
	Remark       string  `json:"remark"`
}

func (r *SaveModelConfigRequest) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	r.Provider = strings.TrimSpace(r.Provider)
	r.BaseURL = strings.TrimRight(strings.TrimSpace(r.BaseURL), "/")
	r.ModelName = strings.TrimSpace(r.ModelName)
	if r.Name == "" {
		return errors.New("渠道名称不能为空")
	}
	if r.BaseURL == "" {
		return errors.New("API Base URL 不能为空")
	}
	if r.ModelName == "" {
		return errors.New("模型标识 (Model Name) 不能为空")
	}
	if r.Temperature <= 0 || r.Temperature > 2 {
		r.Temperature = 0.2
	}
	if r.MaxTokens <= 0 {
		r.MaxTokens = 4096
	}
	return nil
}

// TestConnectionResponse 模型连通性测试结果
type TestConnectionResponse struct {
	Success   bool   `json:"success"`
	LatencyMs int64  `json:"latency_ms"`
	ModelName string `json:"model_name"`
	Reply     string `json:"reply"`
	Message   string `json:"message"`
}

// FetchRemoteModelsRequest 动态拉取渠道可用模型列表请求
type FetchRemoteModelsRequest struct {
	Id       int64  `json:"id"`
	Provider string `json:"provider"`
	BaseURL  string `json:"base_url"`
	APIKey   string `json:"api_key"`
}

// FetchRemoteModelsResponse 动态拉取渠道可用模型列表响应
type FetchRemoteModelsResponse struct {
	Models     []string `json:"models"`
	FromRemote bool     `json:"from_remote"`
	Message    string   `json:"message"`
}

// SendMessageRequest 发送 AI 对话请求
type SendMessageRequest struct {
	UserID          int64                  `json:"-"`
	Username        string                 `json:"-"`
	UserPermissions []string               `json:"-"`
	SessionID       string                 `json:"session_id"` // 为空则自动创建新会话
	Scene           string                 `json:"scene"`      // workspace | copilot
	Prompt          string                 `json:"prompt"`     // 用户输入的自然语言内容
	PageContext     map[string]interface{} `json:"page_context,omitempty"`
}

// SendMessageResponse AI 对话响应
type SendMessageResponse struct {
	Session          *AIChatSession `json:"session"`
	UserMessage      *AIChatMessage `json:"user_message"`
	AssistantMessage *AIChatMessage `json:"assistant_message"`
}

// ConfirmActionRequest 人机协同二次确认请求（确认或取消新增/修改/删除）
type ConfirmActionRequest struct {
	UserID          int64                  `json:"-"`
	Username        string                 `json:"-"`
	UserPermissions []string               `json:"-"`
	MessageID       int64                  `json:"message_id"`
	Decision        string                 `json:"decision"`                // confirm | cancel
	FormPayload     map[string]interface{} `json:"form_payload,omitempty"`  // 允许前端微调后的表单参数
	Diffs           []FieldDiff            `json:"diffs,omitempty"`         // 允许前端在确认单内增删/修改后的字段对比列表
	BatchUpdates    []BatchUpdateTarget    `json:"batch_updates,omitempty"` // 批量修改确认时前端提交的多本期刊微调数据
}

// ConfirmActionResponse 二次确认执行结果
type ConfirmActionResponse struct {
	UpdatedMessage *AIChatMessage `json:"updated_message"`
	FollowupNotice *AIChatMessage `json:"followup_notice,omitempty"`
}

// QueryAuditSessionsRequest 超管查询全员对话审计请求
type QueryAuditSessionsRequest struct {
	PageNumber int    `json:"page_number"`
	PageSize   int    `json:"page_size"`
	Username   string `json:"username"`
	Keyword    string `json:"keyword"`
	Scene      string `json:"scene"`
}

type AuditSessionListResponse struct {
	Total int64            `json:"total"`
	Items []*AIChatSession `json:"items"`
}
