package impl

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"periodical/conf"
	"periodical/internal/apps/ai"
	"periodical/internal/apps/periodical"
	"periodical/internal/common"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"
)

type AIImpl struct {
	db            *gorm.DB
	periodicalSvc periodical.Service
}

func NewAIImpl(periodicalSvc periodical.Service) *AIImpl {
	impl := &AIImpl{
		db:            conf.C().MySQL.GetConn().Debug(),
		periodicalSvc: periodicalSvc,
	}
	if err := impl.bootstrap(); err != nil {
		log.Printf("初始化AI模块数据表警告: %v", err)
	}
	return impl
}

func generateUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

func hasPerm(perms []string, target string) bool {
	for _, p := range perms {
		if p == "*" || p == target {
			return true
		}
	}
	return false
}

// bootstrap 幂等创建 AI 模块相关表并初始化默认主流模型渠道配置
func (i *AIImpl) bootstrap() error {
	tableDDLs := []string{
		`CREATE TABLE IF NOT EXISTS ai_model_configs (
			id int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
			name varchar(100) NOT NULL COMMENT '渠道显示名称',
			provider varchar(50) NOT NULL COMMENT '供应商标识：qwen/deepseek/openai/custom',
			base_url varchar(255) NOT NULL COMMENT 'API Base URL',
			api_key varchar(255) DEFAULT '' COMMENT 'API密钥',
			model_name varchar(100) NOT NULL COMMENT '模型标识',
			system_prompt text COMMENT '自定义系统提示词',
			temperature decimal(3,2) NOT NULL DEFAULT 0.20 COMMENT '采样温度',
			max_tokens int NOT NULL DEFAULT 4096 COMMENT '最大Token数',
			is_default tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否全局默认：1-是，0-否',
			status tinyint(1) NOT NULL DEFAULT 1 COMMENT '状态：1-启用，0-禁用',
			remark varchar(255) DEFAULT '' COMMENT '备注',
			created_at int NOT NULL DEFAULT 0 COMMENT '创建时间',
			updated_at int NOT NULL DEFAULT 0 COMMENT '更新时间',
			PRIMARY KEY (id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='AI大模型渠道配置表';`,
		`CREATE TABLE IF NOT EXISTS ai_chat_sessions (
			id int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
			session_id varchar(64) NOT NULL COMMENT '会话唯一UUID',
			user_id int unsigned NOT NULL COMMENT '所属用户ID',
			username varchar(64) NOT NULL DEFAULT '' COMMENT '所属用户名',
			title varchar(200) NOT NULL DEFAULT '新对话' COMMENT '会话标题',
			scene varchar(32) NOT NULL DEFAULT 'workspace' COMMENT '来源场景：workspace/copilot',
			model_name varchar(100) DEFAULT '' COMMENT '使用模型名',
			message_count int NOT NULL DEFAULT 0 COMMENT '消息总数',
			is_pinned tinyint(1) NOT NULL DEFAULT 0 COMMENT '是否置顶',
			created_at int NOT NULL DEFAULT 0 COMMENT '创建时间',
			updated_at int NOT NULL DEFAULT 0 COMMENT '更新时间',
			PRIMARY KEY (id),
			UNIQUE KEY idx_session_id (session_id),
			KEY idx_user_id (user_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='AI对话会话表';`,
		`CREATE TABLE IF NOT EXISTS ai_chat_messages (
			id int unsigned NOT NULL AUTO_INCREMENT COMMENT '主键ID',
			session_id varchar(64) NOT NULL COMMENT '会话UUID',
			user_id int unsigned NOT NULL COMMENT '所属用户ID',
			username varchar(64) NOT NULL DEFAULT '' COMMENT '所属用户名',
			role varchar(32) NOT NULL COMMENT '角色：user/assistant',
			content longtext COMMENT '消息文本内容',
			model_name varchar(100) DEFAULT '' COMMENT '响应模型名称',
			tool_steps json DEFAULT NULL COMMENT '工具调用步骤记录',
			action_card json DEFAULT NULL COMMENT '结构化交互卡片与确认状态',
			created_at int NOT NULL DEFAULT 0 COMMENT '创建时间',
			updated_at int NOT NULL DEFAULT 0 COMMENT '更新时间',
			PRIMARY KEY (id),
			KEY idx_msg_session_id (session_id),
			KEY idx_msg_user_id (user_id)
		) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='AI对话消息记录表';`,
	}

	for _, ddl := range tableDDLs {
		if err := i.db.Exec(ddl).Error; err != nil {
			return err
		}
	}

	// 初始化预设主流大模型配置（若为空）
	var count int64
	i.db.Model(&ai.AIModelConfig{}).Count(&count)
	if count == 0 {
		now := time.Now().Unix()
		presets := []*ai.AIModelConfig{
			{
				Name:         "通义千问 Qwen-Plus (阿里云百炼)",
				Provider:     "qwen",
				BaseURL:      "https://dashscope.aliyuncs.com/compatible-mode/v1",
				APIKey:       "",
				ModelName:    "qwen-plus",
				SystemPrompt: "你是专业的学术与商业期刊管理 AI 助手，精通期刊选刊推荐、征稿函信息提取、排期与报价维护。",
				Temperature:  0.2,
				MaxTokens:    4096,
				IsDefault:    1,
				Status:       1,
				Remark:       "兼容阿里云 DashScope OpenAI 协议，推荐日常选刊与信息提取使用",
				Meta:         &common.Meta{CreatedAt: now, UpdatedAt: now},
			},
			{
				Name:         "DeepSeek-V3 (深度求索官方)",
				Provider:     "deepseek",
				BaseURL:      "https://api.deepseek.com/v1",
				APIKey:       "",
				ModelName:    "deepseek-chat",
				SystemPrompt: "你是专业的学术与商业期刊管理 AI 助手，精通期刊选刊推荐、征稿函信息提取、排期与报价维护。",
				Temperature:  0.2,
				MaxTokens:    4096,
				IsDefault:    0,
				Status:       1,
				Remark:       "DeepSeek-V3 高性价比推理与工具调用模型",
				Meta:         &common.Meta{CreatedAt: now, UpdatedAt: now},
			},
			{
				Name:         "OpenAI GPT-4o",
				Provider:     "openai",
				BaseURL:      "https://api.openai.com/v1",
				APIKey:       "",
				ModelName:    "gpt-4o",
				SystemPrompt: "你是专业的学术与商业期刊管理 AI 助手，精通期刊选刊推荐、征稿函信息提取、排期与报价维护。",
				Temperature:  0.2,
				MaxTokens:    4096,
				IsDefault:    0,
				Status:       1,
				Remark:       "OpenAI 旗舰多模态与 Function Calling 模型",
				Meta:         &common.Meta{CreatedAt: now, UpdatedAt: now},
			},
		}
		for _, p := range presets {
			_ = i.db.Create(p).Error
		}
	}
	return nil
}

// ==================== 大模型配置管理 ====================

func (i *AIImpl) QueryModelConfigs(ctx context.Context) ([]*ai.AIModelConfig, error) {
	var list []*ai.AIModelConfig
	if err := i.db.WithContext(ctx).Order("is_default DESC, id ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	for _, item := range list {
		item.MaskKey()
	}
	return list, nil
}

func (i *AIImpl) GetActiveDefaultModel(ctx context.Context) (*ai.AIModelConfig, error) {
	var cfg ai.AIModelConfig
	err := i.db.WithContext(ctx).Where("is_default = 1 AND status = 1").First(&cfg).Error
	if err != nil {
		// 退化取第一个启用模型
		err = i.db.WithContext(ctx).Where("status = 1").Order("id ASC").First(&cfg).Error
		if err != nil {
			return nil, err
		}
	}
	cfg.MaskKey()
	return &cfg, nil
}

func (i *AIImpl) SaveModelConfig(ctx context.Context, req *ai.SaveModelConfigRequest) (*ai.AIModelConfig, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	now := time.Now().Unix()

	if req.Id > 0 {
		var existing ai.AIModelConfig
		if err := i.db.WithContext(ctx).First(&existing, req.Id).Error; err != nil {
			return nil, errors.New("未找到该模型配置")
		}

		apiKey := strings.TrimSpace(req.APIKey)
		if apiKey == "" || strings.Contains(apiKey, "****") {
			apiKey = existing.APIKey
		}

		err := i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			if req.IsDefault == 1 {
				if err := tx.Model(&ai.AIModelConfig{}).Where("id <> ?", req.Id).Update("is_default", 0).Error; err != nil {
					return err
				}
			}
			updates := map[string]interface{}{
				"name":          req.Name,
				"provider":      req.Provider,
				"base_url":      req.BaseURL,
				"api_key":       apiKey,
				"model_name":    req.ModelName,
				"system_prompt": req.SystemPrompt,
				"temperature":   req.Temperature,
				"max_tokens":    req.MaxTokens,
				"is_default":    req.IsDefault,
				"status":        req.Status,
				"remark":        req.Remark,
				"updated_at":    now,
			}
			return tx.Model(&ai.AIModelConfig{}).Where("id = ?", req.Id).Updates(updates).Error
		})
		if err != nil {
			return nil, err
		}
		var updated ai.AIModelConfig
		_ = i.db.WithContext(ctx).First(&updated, req.Id).Error
		updated.MaskKey()
		return &updated, nil
	}

	// 新增模式
	ins := &ai.AIModelConfig{
		Name:         req.Name,
		Provider:     req.Provider,
		BaseURL:      req.BaseURL,
		APIKey:       strings.TrimSpace(req.APIKey),
		ModelName:    req.ModelName,
		SystemPrompt: req.SystemPrompt,
		Temperature:  req.Temperature,
		MaxTokens:    req.MaxTokens,
		IsDefault:    req.IsDefault,
		Status:       req.Status,
		Remark:       req.Remark,
		Meta:         &common.Meta{CreatedAt: now, UpdatedAt: now},
	}

	err := i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if req.IsDefault == 1 {
			if err := tx.Model(&ai.AIModelConfig{}).Where("1 = 1").Update("is_default", 0).Error; err != nil {
				return err
			}
		}
		return tx.Create(ins).Error
	})
	if err != nil {
		return nil, err
	}
	ins.MaskKey()
	return ins, nil
}

func (i *AIImpl) SetDefaultModel(ctx context.Context, id int64) error {
	var target ai.AIModelConfig
	if err := i.db.WithContext(ctx).First(&target, id).Error; err != nil {
		return errors.New("目标模型配置不存在")
	}
	now := time.Now().Unix()
	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&ai.AIModelConfig{}).Where("1 = 1").Update("is_default", 0).Error; err != nil {
			return err
		}
		return tx.Model(&ai.AIModelConfig{}).Where("id = ?", id).Updates(map[string]interface{}{
			"is_default": 1,
			"status":     1,
			"updated_at": now,
		}).Error
	})
}

func (i *AIImpl) DeleteModelConfig(ctx context.Context, id int64) error {
	var target ai.AIModelConfig
	if err := i.db.WithContext(ctx).First(&target, id).Error; err != nil {
		return errors.New("模型配置不存在")
	}
	if target.IsDefault == 1 {
		return errors.New("当前默认模型不可直接删除，请先将其他模型设为默认")
	}
	return i.db.WithContext(ctx).Delete(&ai.AIModelConfig{}, id).Error
}

func (i *AIImpl) TestModelConnection(ctx context.Context, id int64) (*ai.TestConnectionResponse, error) {
	var cfg ai.AIModelConfig
	if err := i.db.WithContext(ctx).First(&cfg, id).Error; err != nil {
		return nil, errors.New("模型配置不存在")
	}
	if strings.TrimSpace(cfg.APIKey) == "" {
		return &ai.TestConnectionResponse{
			Success:   false,
			LatencyMs: 0,
			ModelName: cfg.ModelName,
			Message:   "尚未配置 API Key，请先点击【编辑】填写对应厂商的 API 密钥后再测试连接。",
		}, nil
	}

	start := time.Now()
	messages := []openAIMessage{
		{Role: "user", Content: "请仅回复“连接正常”四个字"},
	}
	respMsg, err := i.callOpenAICompatible(ctx, &cfg, messages, nil, 64)
	latency := time.Since(start).Milliseconds()
	if err != nil {
		return &ai.TestConnectionResponse{
			Success:   false,
			LatencyMs: latency,
			ModelName: cfg.ModelName,
			Message:   fmt.Sprintf("连接失败: %v", err),
		}, nil
	}

	return &ai.TestConnectionResponse{
		Success:   true,
		LatencyMs: latency,
		ModelName: cfg.ModelName,
		Reply:     strings.TrimSpace(respMsg.Content),
		Message:   fmt.Sprintf("连接成功 (%dms)", latency),
	}, nil
}

// getPresetModelsByProvider 当未配置密钥或远端不支持 /models 时返回预设推荐模型列表
func getPresetModelsByProvider(provider string) []string {
	switch provider {
	case "qwen":
		return []string{"qwen-plus", "qwen-max", "qwen-turbo", "qwen-long", "qwq-32b"}
	case "deepseek":
		return []string{"deepseek-chat", "deepseek-reasoner"}
	case "openai":
		return []string{"gpt-4o", "gpt-4o-mini", "gpt-4.1", "gpt-4.1-mini", "o3-mini"}
	default:
		return []string{"gpt-4o", "gpt-4o-mini", "claude-3-7-sonnet-20250219", "deepseek-chat", "qwen-plus"}
	}
}

// FetchRemoteModels 动态调用渠道 GET {base_url}/models 接口拉取该渠道支持的模型列表
func (i *AIImpl) FetchRemoteModels(ctx context.Context, req *ai.FetchRemoteModelsRequest) (*ai.FetchRemoteModelsResponse, error) {
	baseURL := strings.TrimRight(strings.TrimSpace(req.BaseURL), "/")
	apiKey := strings.TrimSpace(req.APIKey)

	// 如果编辑已有配置且前端未重填密钥，自动从数据库读取已保存的 APIKey
	if (apiKey == "" || strings.Contains(apiKey, "****")) && req.Id > 0 {
		var existing ai.AIModelConfig
		if err := i.db.WithContext(ctx).First(&existing, req.Id).Error; err == nil {
			apiKey = strings.TrimSpace(existing.APIKey)
			if baseURL == "" {
				baseURL = strings.TrimRight(strings.TrimSpace(existing.BaseURL), "/")
			}
		}
	}

	if baseURL == "" {
		return &ai.FetchRemoteModelsResponse{
			Models:     getPresetModelsByProvider(req.Provider),
			FromRemote: false,
			Message:    "请先填写 API Base URL，当前为您展示预设模型列表",
		}, nil
	}

	if apiKey == "" {
		return &ai.FetchRemoteModelsResponse{
			Models:     getPresetModelsByProvider(req.Provider),
			FromRemote: false,
			Message:    "尚未填写 API Key，当前展示该供应商常用预设模型（填写 API Key 后可拉取渠道完整模型列表）",
		}, nil
	}

	endpoint := strings.TrimSuffix(baseURL, "/chat/completions")
	if !strings.HasSuffix(endpoint, "/models") {
		endpoint = endpoint + "/models"
	}

	httpCtx, cancel := context.WithTimeout(ctx, 15*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(httpCtx, http.MethodGet, endpoint, nil)
	if err != nil {
		return &ai.FetchRemoteModelsResponse{
			Models:     getPresetModelsByProvider(req.Provider),
			FromRemote: false,
			Message:    fmt.Sprintf("构建请求失败 (%v)，已回退展示预设模型", err),
		}, nil
	}
	httpReq.Header.Set("Authorization", "Bearer "+apiKey)
	httpReq.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return &ai.FetchRemoteModelsResponse{
			Models:     getPresetModelsByProvider(req.Provider),
			FromRemote: false,
			Message:    fmt.Sprintf("拉取远端模型失败 (%v)，已回退展示预设模型", err),
		}, nil
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode >= 400 {
		return &ai.FetchRemoteModelsResponse{
			Models:     getPresetModelsByProvider(req.Provider),
			FromRemote: false,
			Message:    fmt.Sprintf("远端返回 HTTP %d，已为您展示预设模型列表", resp.StatusCode),
		}, nil
	}

	// 解析标准 OpenAI {"data": [{"id": "xxx"}]} 格式
	var parsed struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}
	var models []string
	seen := make(map[string]bool)
	if err := json.Unmarshal(bodyBytes, &parsed); err == nil && len(parsed.Data) > 0 {
		for _, item := range parsed.Data {
			id := strings.TrimSpace(item.ID)
			if id != "" && !seen[id] {
				seen[id] = true
				models = append(models, id)
			}
		}
	}

	if len(models) == 0 {
		return &ai.FetchRemoteModelsResponse{
			Models:     getPresetModelsByProvider(req.Provider),
			FromRemote: false,
			Message:    "远端接口未返回模型数组，已为您展示预设模型列表",
		}, nil
	}

	sort.Strings(models)
	return &ai.FetchRemoteModelsResponse{
		Models:     models,
		FromRemote: true,
		Message:    fmt.Sprintf("已从渠道成功获取 %d 个可用模型", len(models)),
	}, nil
}

// ==================== 会话与历史消息管理 ====================

func (i *AIImpl) QueryUserSessions(ctx context.Context, userID int64, keyword string) ([]*ai.AIChatSession, error) {
	var list []*ai.AIChatSession
	tx := i.db.WithContext(ctx).Where("user_id = ?", userID)
	if strings.TrimSpace(keyword) != "" {
		tx = tx.Where("title LIKE ?", "%"+strings.TrimSpace(keyword)+"%")
	}
	err := tx.Order("is_pinned DESC, updated_at DESC").Limit(100).Find(&list).Error
	return list, err
}

func (i *AIImpl) CreateSession(ctx context.Context, userID int64, username string, title string, scene string) (*ai.AIChatSession, error) {
	now := time.Now().Unix()
	if strings.TrimSpace(title) == "" {
		title = "新建期刊智能对话"
	}
	if scene == "" {
		scene = "workspace"
	}
	defaultModel, _ := i.GetActiveDefaultModel(ctx)
	modelName := "qwen-plus"
	if defaultModel != nil {
		modelName = defaultModel.ModelName
	}

	session := &ai.AIChatSession{
		SessionID:    generateUUID(),
		UserID:       userID,
		Username:     username,
		Title:        title,
		Scene:        scene,
		ModelName:    modelName,
		MessageCount: 0,
		IsPinned:     0,
		Meta:         &common.Meta{CreatedAt: now, UpdatedAt: now},
	}
	if err := i.db.WithContext(ctx).Create(session).Error; err != nil {
		return nil, err
	}
	return session, nil
}

func (i *AIImpl) RenameSession(ctx context.Context, userID int64, sessionID string, title string, isPinned *int) error {
	updates := map[string]interface{}{
		"updated_at": time.Now().Unix(),
	}
	if strings.TrimSpace(title) != "" {
		updates["title"] = strings.TrimSpace(title)
	}
	if isPinned != nil {
		updates["is_pinned"] = *isPinned
	}
	return i.db.WithContext(ctx).Model(&ai.AIChatSession{}).
		Where("session_id = ? AND user_id = ?", sessionID, userID).
		Updates(updates).Error
}

func (i *AIImpl) DeleteSession(ctx context.Context, userID int64, sessionID string) error {
	return i.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("session_id = ? AND user_id = ?", sessionID, userID).Delete(&ai.AIChatSession{}).Error; err != nil {
			return err
		}
		return tx.Where("session_id = ? AND user_id = ?", sessionID, userID).Delete(&ai.AIChatMessage{}).Error
	})
}

func (i *AIImpl) QuerySessionMessages(ctx context.Context, userID int64, isSuperAdmin bool, sessionID string) ([]*ai.AIChatMessage, error) {
	var session ai.AIChatSession
	if err := i.db.WithContext(ctx).Where("session_id = ?", sessionID).First(&session).Error; err != nil {
		return nil, errors.New("会话不存在")
	}
	if !isSuperAdmin && session.UserID != userID {
		return nil, errors.New("无权查看他人对话记录")
	}

	var msgs []*ai.AIChatMessage
	err := i.db.WithContext(ctx).Where("session_id = ?", sessionID).Order("id ASC").Find(&msgs).Error
	return msgs, err
}

func (i *AIImpl) QueryAuditSessions(ctx context.Context, req *ai.QueryAuditSessionsRequest) (*ai.AuditSessionListResponse, error) {
	if req.PageNumber <= 0 {
		req.PageNumber = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}
	tx := i.db.WithContext(ctx).Model(&ai.AIChatSession{})
	if strings.TrimSpace(req.Username) != "" {
		tx = tx.Where("username LIKE ?", "%"+strings.TrimSpace(req.Username)+"%")
	}
	if strings.TrimSpace(req.Keyword) != "" {
		tx = tx.Where("title LIKE ?", "%"+strings.TrimSpace(req.Keyword)+"%")
	}
	if strings.TrimSpace(req.Scene) != "" {
		tx = tx.Where("scene = ?", strings.TrimSpace(req.Scene))
	}

	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, err
	}

	var items []*ai.AIChatSession
	offset := (req.PageNumber - 1) * req.PageSize
	if err := tx.Order("updated_at DESC").Offset(offset).Limit(req.PageSize).Find(&items).Error; err != nil {
		return nil, err
	}
	return &ai.AuditSessionListResponse{
		Total: total,
		Items: items,
	}, nil
}

// ==================== 核心 AI Agent 编排引擎 ====================

type openAIMessage struct {
	Role       string           `json:"role"`
	Content    string           `json:"content,omitempty"`
	ToolCalls  []openAIToolCall `json:"tool_calls,omitempty"`
	ToolCallID string           `json:"tool_call_id,omitempty"`
	Name       string           `json:"name,omitempty"`
}

type openAIToolCall struct {
	ID       string             `json:"id"`
	Type     string             `json:"type"`
	Function openAIFunctionCall `json:"function"`
}

type openAIFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type openAIChatRequest struct {
	Model       string          `json:"model"`
	Messages    []openAIMessage `json:"messages"`
	Tools       []interface{}   `json:"tools,omitempty"`
	ToolChoice  string          `json:"tool_choice,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
}

type openAIChatResponse struct {
	Choices []struct {
		Message      openAIMessage `json:"message"`
		FinishReason string        `json:"finish_reason"`
	} `json:"choices"`
	Error *struct {
		Message string `json:"message"`
		Type    string `json:"type"`
	} `json:"error,omitempty"`
}

func (i *AIImpl) callOpenAICompatible(ctx context.Context, cfg *ai.AIModelConfig, messages []openAIMessage, tools []interface{}, maxTokens int) (*openAIMessage, error) {
	endpoint := strings.TrimRight(strings.TrimSpace(cfg.BaseURL), "/")
	if !strings.HasSuffix(endpoint, "/chat/completions") {
		endpoint = endpoint + "/chat/completions"
	}

	if maxTokens <= 0 {
		maxTokens = cfg.MaxTokens
		if maxTokens <= 0 {
			maxTokens = 4096
		}
	}

	payload := openAIChatRequest{
		Model:       cfg.ModelName,
		Messages:    messages,
		Temperature: cfg.Temperature,
		MaxTokens:   maxTokens,
	}
	if len(tools) > 0 {
		payload.Tools = tools
		payload.ToolChoice = "auto"
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	httpCtx, cancel := context.WithTimeout(ctx, 45*time.Second)
	defer cancel()

	httpReq, err := http.NewRequestWithContext(httpCtx, http.MethodPost, endpoint, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+strings.TrimSpace(cfg.APIKey))

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(respBytes))
	}

	var parsed openAIChatResponse
	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		return nil, fmt.Errorf("解析模型响应失败: %v", err)
	}
	if parsed.Error != nil && parsed.Error.Message != "" {
		return nil, errors.New(parsed.Error.Message)
	}
	if len(parsed.Choices) == 0 {
		return nil, errors.New("模型未返回有效 Choices")
	}
	return &parsed.Choices[0].Message, nil
}

// SendMessage 处理用户自然语言指令，自动协调会话、工具调用与人机确认卡片
func (i *AIImpl) SendMessage(ctx context.Context, req *ai.SendMessageRequest) (*ai.SendMessageResponse, error) {
	prompt := strings.TrimSpace(req.Prompt)
	if prompt == "" {
		return nil, errors.New("请输入对话内容")
	}

	now := time.Now().Unix()

	// 1. 获取或自动创建会话
	var session *ai.AIChatSession
	if strings.TrimSpace(req.SessionID) != "" {
		var s ai.AIChatSession
		if err := i.db.WithContext(ctx).Where("session_id = ? AND user_id = ?", req.SessionID, req.UserID).First(&s).Error; err == nil {
			session = &s
		}
	}
	if session == nil {
		titleRunes := []rune(prompt)
		if len(titleRunes) > 22 {
			titleRunes = append(titleRunes[:22], []rune("...")...)
		}
		var err error
		session, err = i.CreateSession(ctx, req.UserID, req.Username, string(titleRunes), req.Scene)
		if err != nil {
			return nil, err
		}
	}

	// 2. 获取当前全局默认模型配置
	var activeModel ai.AIModelConfig
	hasValidRemoteModel := false
	if err := i.db.WithContext(ctx).Where("is_default = 1 AND status = 1").First(&activeModel).Error; err != nil {
		_ = i.db.WithContext(ctx).Where("status = 1").Order("id ASC").First(&activeModel).Error
	}
	modelLabel := activeModel.ModelName
	if modelLabel == "" {
		modelLabel = "smart-agent-local"
	}
	if strings.TrimSpace(activeModel.APIKey) != "" && strings.TrimSpace(activeModel.BaseURL) != "" {
		hasValidRemoteModel = true
	}

	// 3. 保存用户消息
	userMsg := &ai.AIChatMessage{
		SessionID: session.SessionID,
		UserID:    req.UserID,
		Username:  req.Username,
		Role:      "user",
		Content:   prompt,
		ModelName: modelLabel,
		ToolSteps: []ai.ToolStep{},
		Meta:      &common.Meta{CreatedAt: now, UpdatedAt: now},
	}
	if err := i.db.WithContext(ctx).Create(userMsg).Error; err != nil {
		return nil, err
	}

	// 4. 加载当前所有分类字典（收稿方向、级别、网站收录、周期）
	catSet, _ := i.periodicalSvc.QueryCategorize(ctx, &periodical.QueryCategorizeRequest{
		PageNumber: 1,
		PageSize:   500,
		Status:     1,
	})
	var allCats []*periodical.Categorize
	if catSet != nil {
		allCats = catSet.Items
	}

	// 5. 执行 Agent 推理与工具调用（优先调用配置的远程大模型，若未配Key或报错则无缝回退本地智能解析引擎）
	var replyContent string
	var toolSteps []ai.ToolStep
	var actionCard *ai.ActionCardData

	if hasValidRemoteModel {
		replyContent, toolSteps, actionCard = i.runLLMFunctionCallingAgent(ctx, &activeModel, session.SessionID, req, allCats)
	} else {
		replyContent, toolSteps, actionCard = i.runLocalSmartAgent(ctx, req, allCats, false, "")
	}

	// 6. 保存 Assistant 消息并更新会话计数与标题
	assistantNow := time.Now().Unix()
	assistantMsg := &ai.AIChatMessage{
		SessionID:  session.SessionID,
		UserID:     req.UserID,
		Username:   req.Username,
		Role:       "assistant",
		Content:    replyContent,
		ModelName:  modelLabel,
		ToolSteps:  toolSteps,
		ActionCard: actionCard,
		Meta:       &common.Meta{CreatedAt: assistantNow, UpdatedAt: assistantNow},
	}
	if err := i.db.WithContext(ctx).Create(assistantMsg).Error; err != nil {
		return nil, err
	}

	// 更新会话统计与标题
	sessionUpdates := map[string]interface{}{
		"model_name":    modelLabel,
		"message_count": gorm.Expr("message_count + 2"),
		"updated_at":    assistantNow,
	}
	if session.Title == "新建期刊智能对话" || session.Title == "新对话" {
		titleRunes := []rune(prompt)
		if len(titleRunes) > 22 {
			titleRunes = append(titleRunes[:22], []rune("...")...)
		}
		sessionUpdates["title"] = string(titleRunes)
		session.Title = string(titleRunes)
	}
	_ = i.db.WithContext(ctx).Model(&ai.AIChatSession{}).Where("session_id = ?", session.SessionID).Updates(sessionUpdates).Error
	session.MessageCount += 2
	session.UpdatedAt = assistantNow

	return &ai.SendMessageResponse{
		Session:          session,
		UserMessage:      userMsg,
		AssistantMessage: assistantMsg,
	}, nil
}

// buildAgentTools 定义提供给 OpenAI-Compatible 大模型的 Function Calling Tools
func buildAgentTools() []interface{} {
	return []interface{}{
		map[string]interface{}{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "query_periodicals",
				"description": "根据期刊名称、级别、网站收录、收稿方向、出版周期、出刊月份、是否开发票、是否全包、主管/主办单位等条件查询期刊列表",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name":            map[string]interface{}{"type": "string", "description": "期刊名称关键词（不含书名号）"},
						"competent_unit":  map[string]interface{}{"type": "string", "description": "主管单位关键词"},
						"organization":    map[string]interface{}{"type": "string", "description": "主办单位关键词"},
						"domestic_number": map[string]interface{}{"type": "string", "description": "国内刊号关键词"},
						"column_setting":  map[string]interface{}{"type": "string", "description": "栏目设置关键词"},
						"category_names": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "分类名称列表，如 ['国家级', '知网', '教育类', '月刊']",
						},
						"publish_month": map[string]interface{}{"type": "string", "description": "出刊月份，格式 YYYY-MM，如 2025-06"},
						"invoice":       map[string]interface{}{"type": "integer", "description": "是否可开杂志社发票：1-可开，2-不可开，0-不限"},
						"is_warp":       map[string]interface{}{"type": "integer", "description": "是否全包：1-全包，2-非全包，0-不限"},
						"works":         map[string]interface{}{"type": "integer", "description": "是否可发作品：1-可发，2-不可发，0-不限"},
						"hidden_status": map[string]interface{}{"type": "string", "description": "隐藏状态过滤：默认不传或传 visible_only（自动过滤已隐藏期刊）；仅当用户明确要求查看已隐藏期刊时传 hidden_only，明确要求包含隐藏期刊时传 all"},
						"order_by":      map[string]interface{}{"type": "string", "description": "排序字段：price | publication_time | periodical_page | updated_at"},
						"order_desc":    map[string]interface{}{"type": "string", "description": "排序方向：ascending | descending"},
					},
				},
			},
		},
		map[string]interface{}{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "prepare_create_periodical",
				"description": "当用户要求新增/录入/创建一本期刊，或粘贴了一段期刊征稿函/介绍文本要求解析入库时调用，生成待确认新增卡片",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"name":                       map[string]interface{}{"type": "string", "description": "期刊名称，建议带《》"},
						"publication_time":           map[string]interface{}{"type": "string", "description": "出刊时间，支持 YYYY-MM（如 2025-09）、YYYY-MM-DD 或 2025-09月底/中旬，系统会自动规范化，无需向用户反问具体哪一天"},
						"competent_unit":             map[string]interface{}{"type": "string", "description": "主管单位"},
						"organization":               map[string]interface{}{"type": "string", "description": "主办单位"},
						"domestic_number":            map[string]interface{}{"type": "string", "description": "国内刊号 CN"},
						"international_number":       map[string]interface{}{"type": "string", "description": "国际刊号 ISSN"},
						"founding_time":              map[string]interface{}{"type": "string", "description": "创刊年份，如 2009"},
						"citation_rate":              map[string]interface{}{"type": "string", "description": "引用率"},
						"audit_time":                 map[string]interface{}{"type": "string", "description": "审核周期"},
						"price":                      map[string]interface{}{"type": "string", "description": "期刊发表费用"},
						"periodical_page":            map[string]interface{}{"type": "string", "description": "期刊页码"},
						"is_warp":                    map[string]interface{}{"type": "boolean", "description": "是否全包"},
						"invoice_receipt":            map[string]interface{}{"type": "boolean", "description": "是否可开杂志社发票"},
						"works":                      map[string]interface{}{"type": "boolean", "description": "是否可发作品"},
						"is_hidden":                  map[string]interface{}{"type": "boolean", "description": "是否创建后设为隐藏，默认 false"},
						"composite_influence_factor": map[string]interface{}{"type": "string", "description": "复合影响因子"},
						"synthetic_influence_factor": map[string]interface{}{"type": "string", "description": "综合影响因子"},
						"column_setting":             map[string]interface{}{"type": "string", "description": "栏目设置"},
						"periodical_description":     map[string]interface{}{"type": "string", "description": "期刊综合描述"},
						"attention_matter":           map[string]interface{}{"type": "string", "description": "文章注意事项"},
						"author_info":                map[string]interface{}{"type": "string", "description": "作者数量与样刊说明"},
						"category_names": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "所属收稿方向、网站收录、期刊级别、出版周期的名称列表",
						},
					},
					"required": []string{"name"},
				},
			},
		},
		map[string]interface{}{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "prepare_update_periodical",
				"description": "当用户要求修改/更新/调整一本或多本已存在期刊的信息（如修改价格、出刊时间、页码、级别、周期、收录、方向、设为隐藏或取消隐藏等）时调用，生成修改前后Diff确认卡片。支持通过 target_names 一次性传入多本期刊进行批量修改或批量设置隐藏/取消隐藏",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"target_name": map[string]interface{}{"type": "string", "description": "要修改的目标期刊名称（单本修改时填写，或多本用顿号分隔）"},
						"target_names": map[string]interface{}{
							"type":        "array",
							"items":       map[string]interface{}{"type": "string"},
							"description": "当用户要求同时修改多本期刊时，传入目标期刊名称数组，如 ['河北农机', '文学艺术周刊', '考试周刊']",
						},
						"publication_time":           map[string]interface{}{"type": "string", "description": "新的出刊时间，支持 YYYY-MM（如 2025-09）、YYYY-MM-DD、1月 或 2025-09月底/中旬，系统会自动补全日期，切勿反问具体几号（未修改则留空）"},
						"price":                      map[string]interface{}{"type": "string", "description": "新的期刊发表费（未修改则留空）"},
						"periodical_page":            map[string]interface{}{"type": "string", "description": "新的期刊页码（未修改则留空）"},
						"level":                      map[string]interface{}{"type": "string", "description": "新的期刊级别，如：省级、国家级、国际期刊、学报（未修改则留空）"},
						"cycle":                      map[string]interface{}{"type": "string", "description": "新的出版周期，如：半月刊及以下、月刊及以上（未修改则留空）"},
						"collections":                map[string]interface{}{"type": "string", "description": "新的网站收录，多个用顿号分隔，如：知网、万方、维普（未修改则留空）"},
						"directions":                 map[string]interface{}{"type": "string", "description": "新的收稿方向，多个用顿号分隔，如：经济类、综合类（未修改则留空）"},
						"competent_unit":             map[string]interface{}{"type": "string", "description": "新的主管单位（未修改则留空）"},
						"organization":           map[string]interface{}{"type": "string", "description": "新的主办单位（未修改则留空）"},
						"domestic_number":        map[string]interface{}{"type": "string", "description": "新的国内刊号（未修改则留空）"},
						"international_number":   map[string]interface{}{"type": "string", "description": "新的国际刊号（未修改则留空）"},
						"audit_time":             map[string]interface{}{"type": "string", "description": "新的审核周期（未修改则留空）"},
						"column_setting":             map[string]interface{}{"type": "string", "description": "新的栏目设置（未修改则留空）"},
						"attention_matter":           map[string]interface{}{"type": "string", "description": "新的注意事项（未修改则留空）"},
						"periodical_description":     map[string]interface{}{"type": "string", "description": "新的期刊描述（未修改则留空）"},
						"founding_time":              map[string]interface{}{"type": "string", "description": "新的创刊年份（未修改则留空）"},
						"citation_rate":              map[string]interface{}{"type": "string", "description": "新的引用率（未修改则留空）"},
						"composite_influence_factor": map[string]interface{}{"type": "string", "description": "新的复合影响因子（未修改则留空）"},
						"synthetic_influence_factor": map[string]interface{}{"type": "string", "description": "新的综合影响因子（未修改则留空）"},
						"author_info":                map[string]interface{}{"type": "string", "description": "新的作者数量与样刊说明（未修改则留空）"},
						"periodical_final_proof":     map[string]interface{}{"type": "string", "description": "新的期刊清样说明（未修改则留空）"},
						"article_naming":             map[string]interface{}{"type": "string", "description": "新的文章命名规范（未修改则留空）"},
						"publication_process":        map[string]interface{}{"type": "string", "description": "新的发表流程（未修改则留空）"},
						"all_inclusive_process":      map[string]interface{}{"type": "string", "description": "新的全包流程（未修改则留空）"},
						"internal_processes":         map[string]interface{}{"type": "string", "description": "新的内部流程（未修改则留空）"},
						"is_warp":                    map[string]interface{}{"type": "boolean", "description": "是否全包"},
						"invoice_receipt":            map[string]interface{}{"type": "boolean", "description": "是否可开发票"},
						"works":                      map[string]interface{}{"type": "boolean", "description": "是否可发作品"},
						"is_hidden":                  map[string]interface{}{"type": "boolean", "description": "是否隐藏该期刊：true 表示设为隐藏/屏蔽/下架，false 表示取消隐藏/恢复显示/不隐藏"},
					},
					"required": []string{"target_name"},
				},
			},
		},
		map[string]interface{}{
			"type": "function",
			"function": map[string]interface{}{
				"name":        "prepare_delete_periodical",
				"description": "当用户要求删除/移除某本期刊时调用，生成高危删除二次确认卡片",
				"parameters": map[string]interface{}{
					"type": "object",
					"properties": map[string]interface{}{
						"target_name": map[string]interface{}{"type": "string", "description": "要删除的目标期刊名称"},
					},
					"required": []string{"target_name"},
				},
			},
		},
	}
}

func (i *AIImpl) buildSystemPrompt(cfg *ai.AIModelConfig, allCats []*periodical.Categorize) string {
	var dirs, levels, cols, cycles []string
	for _, c := range allCats {
		switch c.Type {
		case 1:
			dirs = append(dirs, c.Name)
		case 2:
			levels = append(levels, c.Name)
		case 3:
			cols = append(cols, c.Name)
		case 4:
			cycles = append(cycles, c.Name)
		}
	}

	basePrompt := `你是“期刊管理系统”的专属 AI 智能助手（Editorial & Publishing AI Agent）。
你的职责是协助编辑、业务顾问与管理员通过自然语言高效地对期刊库进行【选刊检索、征稿信息智能提取与新增、排期/价格修改、删除确认】以及【生成对外客户推荐报价文案】。

系统当前内置的标准分类字典如下（调用工具时请优先匹配以下标准名称）：
- 【收稿方向 (type=1)】：` + strings.Join(dirs, "、") + `
- 【期刊级别 (type=2)】：` + strings.Join(levels, "、") + `
- 【网站收录 (type=3)】：` + strings.Join(cols, "、") + `
- 【出版周期 (type=4)】：` + strings.Join(cycles, "、") + `

行为准则与输出排版规范（极其重要）：
1. 【真实库存优先】：当用户查询、筛选、对比或要求推荐期刊时，请调用 query_periodicals 工具获取真实库存数据后再回答，严禁编造不存在的期刊。
2. 【正文与卡片职责分离（禁止重复画表）】：当调用了 query_periodicals 工具后，前端界面已经会自动在你的回复下方渲染一个支持翻页、网格/表格切换与右侧分屏画布的【交互式期刊结果卡片 (ActionCard)】。因此，你在正文中**严禁再次输出完整的 Markdown 表格把返回的 8 本期刊机械地重复列一遍**（除非用户明确要求“把对比表格写在回复正文里”）。
3. 【正文输出顾问级数据洞察】：调用查询工具后，你的正文应当输出精炼的 **选刊数据洞察与顾问建议**（控制在 3~5 个要点内）：
   - **检索概况**：汇总匹配到的期刊总数及核心筛选维度；
   - **特征与价格分布**：简要总结当前结果集的网站收录分布（如知网/万方/维普占比）、发表费价格区间、全包或可开发票情况；
   - **精选推荐**：从返回结果中挑选 2~3 本最具代表性（如性价比最高、排期最快或综合权威度高）的刊物做一句话亮点点评，并引导用户在下方交互卡片中翻页浏览或点击【在右侧画布展开】。
4. 【新增与修改确认】：当用户要求新增期刊或粘贴了一段非结构化征稿函文本时，调用 prepare_create_periodical 提取字段；当要求修改或删除时，必须调用 prepare_update_periodical 或 prepare_delete_periodical 生成二次确认卡片，绝不能在未调用工具的情况下声称已直接修改或删除数据库。
5. 【连贯执行，严禁机械反问日期细节】：当用户要求修改或新增期刊排期时，若用户仅提供了月份或旬期（如“2025-09”、“9月中旬”、“2025-08月底”），请直接调用工具生成确认卡片！严禁停下来反问用户“请问具体是9月哪一天”，因为系统后端会自动规范化日期，且确认卡片本身支持用户在界面上行内直接编辑微调所有字段。
6. 【隐藏期刊管控规范】：日常选刊、筛刊、对比推荐时，query_periodicals 默认自动排除已隐藏期刊（请保持 hidden_status 为默认 visible_only，仅当用户明确说出“查看已隐藏的期刊”或“包含隐藏期刊”时才传 hidden_only 或 all）。当用户要求“把某本或多本期刊设为隐藏/屏蔽”或“把某本期刊取消隐藏/恢复显示/不隐藏”时，请调用 prepare_update_periodical 并传入 is_hidden=true 或 is_hidden=false 生成确认卡片。`

	if strings.TrimSpace(cfg.SystemPrompt) != "" {
		basePrompt = cfg.SystemPrompt + "\n\n" + basePrompt
	}
	return basePrompt
}

// runLLMFunctionCallingAgent 调用真实 OpenAI-Compatible 大模型并执行 Function Calling
func (i *AIImpl) runLLMFunctionCallingAgent(
	ctx context.Context,
	cfg *ai.AIModelConfig,
	sessionID string,
	req *ai.SendMessageRequest,
	allCats []*periodical.Categorize,
) (string, []ai.ToolStep, *ai.ActionCardData) {
	// 拉取最近 8 条历史消息构建多轮上下文
	var historyMsgs []*ai.AIChatMessage
	_ = i.db.WithContext(ctx).
		Where("session_id = ?", sessionID).
		Order("id DESC").
		Limit(8).
		Find(&historyMsgs)

	messages := []openAIMessage{
		{Role: "system", Content: i.buildSystemPrompt(cfg, allCats)},
	}
	for idx := len(historyMsgs) - 1; idx >= 0; idx-- {
		m := historyMsgs[idx]
		if m.Content == req.Prompt && idx == 0 {
			continue
		}
		if m.Role == "user" || m.Role == "assistant" {
			contentWithCtx := m.Content
			if m.Role == "assistant" && m.ActionCard != nil {
				contentWithCtx += "\n" + formatCardContextForLLM(m.ActionCard)
			}
			messages = append(messages, openAIMessage{
				Role:    m.Role,
				Content: contentWithCtx,
			})
		}
	}
	messages = append(messages, openAIMessage{
		Role:    "user",
		Content: req.Prompt,
	})

	tools := buildAgentTools()
	respMsg, err := i.callOpenAICompatible(ctx, cfg, messages, tools, cfg.MaxTokens)
	if err != nil {
		// 如果大模型调用报错（如 Key 过期、额度不足或网络超时），优雅回退到本地智能引擎并附带提示
		return i.runLocalSmartAgent(ctx, req, allCats, true, err.Error())
	}

	var toolSteps []ai.ToolStep
	var actionCard *ai.ActionCardData

	if len(respMsg.ToolCalls) > 0 {
		messages = append(messages, *respMsg)
		for _, tc := range respMsg.ToolCalls {
			var argsMap map[string]interface{}
			_ = json.Unmarshal([]byte(tc.Function.Arguments), &argsMap)

			step, card, toolOutputJSON := i.executeAgentTool(ctx, tc.Function.Name, argsMap, req.UserPermissions, allCats)
			toolSteps = append(toolSteps, step)
			if card != nil {
				actionCard = mergeActionCards(actionCard, card)
			}

			messages = append(messages, openAIMessage{
				Role:       "tool",
				ToolCallID: tc.ID,
				Name:       tc.Function.Name,
				Content:    toolOutputJSON,
			})
		}

		// 若工具已成功生成写操作确认卡片（新增/修改/删除），直接输出精炼的一体化引导文案，避免二次调用大模型输出重复的字段列表
		if actionCard != nil && (actionCard.CardType == "update_confirm" || actionCard.CardType == "create_confirm" || actionCard.CardType == "delete_confirm") {
			return i.buildSummaryFromCard(actionCard, toolSteps), toolSteps, actionCard
		}

		// 将查询工具执行结果回传给大模型生成最终自然语言选刊洞察
		finalMsg, err2 := i.callOpenAICompatible(ctx, cfg, messages, nil, cfg.MaxTokens)
		if err2 == nil && strings.TrimSpace(finalMsg.Content) != "" {
			return finalMsg.Content, toolSteps, actionCard
		}
		// 若二次总结超时，根据 actionCard 生成兜底专业回复
		return i.buildSummaryFromCard(actionCard, toolSteps), toolSteps, actionCard
	}

	return respMsg.Content, toolSteps, actionCard
}

// executeAgentTool 执行具体的业务工具并生成结构化卡片
func (i *AIImpl) executeAgentTool(
	ctx context.Context,
	toolName string,
	args map[string]interface{},
	userPerms []string,
	allCats []*periodical.Categorize,
) (ai.ToolStep, *ai.ActionCardData, string) {
	startTime := time.Now()
	calcMs := func() int64 {
		ms := time.Since(startTime).Milliseconds()
		if ms <= 0 {
			ms = 2
		}
		return ms
	}

	switch toolName {
	case "query_periodicals":
		if !hasPerm(userPerms, "periodical:view") && !hasPerm(userPerms, "periodical:query") {
			step := ai.ToolStep{
				ToolName:   toolName,
				Title:      "检索期刊库",
				Summary:    "权限不足：缺少期刊查询权限",
				DurationMs: calcMs(),
				Status:     "denied",
			}
			return step, nil, `{"error": "当前用户缺少期刊查询权限 (periodical:query)"}`
		}
		card, err := i.execQueryTool(ctx, args, userPerms, allCats)
		if err != nil {
			step := ai.ToolStep{
				ToolName:   toolName,
				Title:      "检索期刊库",
				Summary:    err.Error(),
				DurationMs: calcMs(),
				Status:     "error",
			}
			return step, nil, fmt.Sprintf(`{"error": %q}`, err.Error())
		}
		step := ai.ToolStep{
			ToolName:    toolName,
			Title:       "多维条件检索期刊库",
			ArgsSummary: card.Description,
			Summary:     fmt.Sprintf("共匹配到 %d 本期刊（首屏展示 %d 本）", card.TotalCount, len(card.Items)),
			DurationMs:  calcMs(),
			Status:      "success",
		}
		b, _ := json.Marshal(map[string]interface{}{
			"total_count": card.TotalCount,
			"items":       card.Items,
			"_ui_notice":  "前端已通过交互式卡片展示上述期刊列表（支持翻页与表格视图），正文中请勿重复输出完整 Markdown 表格逐条罗列这8本期刊，请直接输出【选刊数据洞察、价格与收录特征分析、2~3本代表刊物推荐亮点】。",
		})
		return step, card, string(b)

	case "prepare_create_periodical":
		if !hasPerm(userPerms, "periodical:create") {
			step := ai.ToolStep{
				ToolName:   toolName,
				Title:      "构建新增期刊入库单",
				Summary:    "权限拦截：当前账号无新增期刊权限 (periodical:create)",
				DurationMs: calcMs(),
				Status:     "denied",
			}
			return step, nil, `{"error": "权限不足：您当前绑定的角色没有【新增期刊 (periodical:create)】权限，无法录入新刊物。"}`
		}
		card := i.execPrepareCreateTool(args, allCats)
		step := ai.ToolStep{
			ToolName:    toolName,
			Title:       "智能提取期刊字段并生成入库确认单",
			ArgsSummary: fmt.Sprintf("目标刊名: %s", card.PeriodicalName),
			Summary:     fmt.Sprintf("已提取 %s 结构化字段，等待人工确认入库", card.PeriodicalName),
			DurationMs:  calcMs(),
			Status:      "success",
		}
		b, _ := json.Marshal(card.FormPayload)
		return step, card, fmt.Sprintf(`{"status": "pending_user_confirmation", "extracted": %s, "_ui_notice": "前端已生成入库确认卡片，请简要说明提取亮点并引导用户点击确认按钮。"}`, string(b))

	case "prepare_update_periodical":
		if !hasPerm(userPerms, "periodical:update") {
			step := ai.ToolStep{
				ToolName:   toolName,
				Title:      "构建期刊修改预览单",
				Summary:    "权限拦截：当前账号无编辑期刊权限 (periodical:update)",
				DurationMs: calcMs(),
				Status:     "denied",
			}
			return step, nil, `{"error": "权限不足：您当前绑定的角色没有【编辑期刊 (periodical:update)】权限，无法修改期刊信息。"}`
		}
		card, err := i.execPrepareUpdateTool(ctx, args, allCats)
		if err != nil {
			step := ai.ToolStep{
				ToolName:   toolName,
				Title:      "定位并构建期刊修改单",
				Summary:    err.Error(),
				DurationMs: calcMs(),
				Status:     "error",
			}
			return step, nil, fmt.Sprintf(`{"error": %q}`, err.Error())
		}
		argsSum := fmt.Sprintf("目标刊物: %s (ID: %d)", card.PeriodicalName, card.PeriodicalID)
		stepSum := fmt.Sprintf("已生成 %s 的 %d 项字段变更预览，等待确认", card.PeriodicalName, len(card.Diffs))
		if len(card.BatchUpdates) > 1 {
			var bNames []string
			totalDiffs := 0
			for _, bu := range card.BatchUpdates {
				bNames = append(bNames, bu.PeriodicalName)
				totalDiffs += len(bu.Diffs)
			}
			argsSum = fmt.Sprintf("批量目标 (%d本): %s", len(card.BatchUpdates), strings.Join(bNames, "、"))
			stepSum = fmt.Sprintf("已生成 %s 共 %d 本期刊的 %d 项字段变更预览，等待确认", strings.Join(bNames, "、"), len(card.BatchUpdates), totalDiffs)
		}
		step := ai.ToolStep{
			ToolName:    toolName,
			Title:       "对比期刊字段变更 (Diff)",
			ArgsSummary: argsSum,
			Summary:     stepSum,
			DurationMs:  calcMs(),
			Status:      "success",
		}
		b, _ := json.Marshal(card.Diffs)
		return step, card, fmt.Sprintf(`{"status": "pending_user_confirmation", "periodical_name": %q, "batch_count": %d, "diffs": %s, "_ui_notice": "前端已生成可交互修改确认卡片（用户可在卡片内直接微调数值并点击【确认执行修改】）。请简明告知已生成修改预览单，切勿怀疑字段未纳入。"}`, card.PeriodicalName, len(card.BatchUpdates), string(b))

	case "prepare_delete_periodical":
		if !hasPerm(userPerms, "periodical:delete") {
			step := ai.ToolStep{
				ToolName:   toolName,
				Title:      "构建期刊删除确认单",
				Summary:    "权限拦截：当前账号无删除期刊权限 (periodical:delete)",
				DurationMs: calcMs(),
				Status:     "denied",
			}
			return step, nil, `{"error": "权限不足：您当前绑定的角色没有【删除期刊 (periodical:delete)】权限，禁止删除操作。"}`
		}
		card, err := i.execPrepareDeleteTool(ctx, args)
		if err != nil {
			step := ai.ToolStep{
				ToolName:   toolName,
				Title:      "定位待删除期刊",
				Summary:    err.Error(),
				DurationMs: calcMs(),
				Status:     "error",
			}
			return step, nil, fmt.Sprintf(`{"error": %q}`, err.Error())
		}
		step := ai.ToolStep{
			ToolName:    toolName,
			Title:       "高危删除二次确认护栏",
			ArgsSummary: fmt.Sprintf("目标刊物: %s (ID: %d)", card.PeriodicalName, card.PeriodicalID),
			Summary:     fmt.Sprintf("已锁定 %s (ID: %d)，等待人工确认删除", card.PeriodicalName, card.PeriodicalID),
			DurationMs:  calcMs(),
			Status:      "success",
		}
		return step, card, fmt.Sprintf(`{"status": "pending_user_confirmation", "periodical_id": %d, "periodical_name": %q}`, card.PeriodicalID, card.PeriodicalName)
	}

	return ai.ToolStep{ToolName: toolName, Title: "未知工具", DurationMs: calcMs(), Status: "error"}, nil, `{"error": "unknown tool"}`
}

// ==================== 具体工具执行逻辑 ====================

func getStringArg(args map[string]interface{}, key string) string {
	if v, ok := args[key]; ok && v != nil {
		return strings.TrimSpace(fmt.Sprintf("%v", v))
	}
	return ""
}

func getIntArg(args map[string]interface{}, key string) int {
	if v, ok := args[key]; ok && v != nil {
		switch val := v.(type) {
		case float64:
			return int(val)
		case int:
			return val
		case string:
			n, _ := strconv.Atoi(val)
			return n
		}
	}
	return 0
}

func getBoolArg(args map[string]interface{}, key string) (bool, bool) {
	if v, ok := args[key]; ok && v != nil {
		switch val := v.(type) {
		case bool:
			return val, true
		case float64:
			return val == 1, true
		case string:
			s := strings.ToLower(strings.TrimSpace(val))
			if s == "" {
				return false, false
			}
			return s == "true" || s == "1" || s == "是" || s == "可" || s == "隐藏", true
		}
	}
	return false, false
}

func formatPeriodicalToBrief(p *periodical.Periodical) ai.PeriodicalBriefCard {
	var level, cycle string
	var dirs, cols []string
	for _, c := range p.Categorizes {
		switch c.Type {
		case 1:
			dirs = append(dirs, c.Name)
		case 2:
			level = c.Name
		case 3:
			cols = append(cols, c.Name)
		case 4:
			cycle = c.Name
		}
	}
	pubTime := p.PublicationTime
	if len(pubTime) >= 10 {
		pubTime = pubTime[:10]
	}
	return ai.PeriodicalBriefCard{
		PeriodicalID:          p.PeriodicalID,
		Name:                  p.Name,
		Level:                 level,
		Directions:            strings.Join(dirs, "、"),
		Collections:           cols,
		Cycle:                 cycle,
		PublicationTime:       pubTime,
		Price:                 p.Price,
		PeriodicalPage:        p.PeriodicalPage,
		IsWarp:                p.IsWarp,
		InvoiceReceipt:        p.InvoiceReceipt,
		Works:                 p.Works,
		IsHidden:              p.IsHidden,
		Organization:          p.Organization,
		CompetentUnit:         p.CompetentUnit,
		DomesticNumber:        p.DomesticNumber,
		InternationalNumber:   p.InternationalNumber,
		AuditTime:             p.AuditTime,
		PeriodicalDescription: p.PeriodicalDescription,
	}
}

func (i *AIImpl) execQueryTool(ctx context.Context, args map[string]interface{}, userPerms []string, allCats []*periodical.Categorize) (*ai.ActionCardData, error) {
	qReq := periodical.NewQueryPeriodicalRequest()
	qReq.PageNumber = 1
	qReq.PageSize = 8
	qReq.Name = strings.Trim(getStringArg(args, "name"), "《》")
	qReq.CompetentUnit = getStringArg(args, "competent_unit")
	qReq.Organization = getStringArg(args, "organization")
	qReq.DomesticNumber = getStringArg(args, "domestic_number")
	qReq.ColumnSetting = getStringArg(args, "column_setting")
	qReq.Invoice = getIntArg(args, "invoice")
	qReq.IsWarp = getIntArg(args, "is_warp")
	qReq.Works = getIntArg(args, "works")
	qReq.OrderBy = getStringArg(args, "order_by")
	qReq.OrderDesc = getStringArg(args, "order_desc")

	// 默认严格过滤已隐藏期刊（IsHidden = 2）；仅当拥有编辑权限且显式请求查询已隐藏或全部时放行
	qReq.IsHidden = 2
	if hasPerm(userPerms, "periodical:update") {
		hiddenStatus := strings.ToLower(strings.TrimSpace(getStringArg(args, "hidden_status")))
		if hiddenStatus == "hidden_only" {
			qReq.IsHidden = 1
		} else if hiddenStatus == "all" {
			qReq.IsHidden = 3
		}
	}

	publishMonth := getStringArg(args, "publish_month")
	if len(publishMonth) == 7 && strings.Contains(publishMonth, "-") {
		t, err := time.Parse("2006-01", publishMonth)
		if err == nil {
			qReq.PublishTimeStart = t.Format("2006-01-02")
			endOfMonth := t.AddDate(0, 1, -1)
			qReq.PublishTimeEnd = endOfMonth.Format("2006-01-02")
		}
	}

	// 匹配分类名称 -> HasBaseTypeList
	var matchedCatNames []string
	var baseTypes []map[string]int
	if rawNames, ok := args["category_names"].([]interface{}); ok {
		for _, rn := range rawNames {
			kw := strings.TrimSpace(fmt.Sprintf("%v", rn))
			if kw == "" {
				continue
			}
			for _, cat := range allCats {
				if strings.Contains(cat.Name, kw) || strings.Contains(kw, cat.Name) {
					baseTypes = append(baseTypes, map[string]int{
						"base_id": cat.CategorizeID,
						"type":    int(cat.Type),
					})
					matchedCatNames = append(matchedCatNames, cat.Name)
					break
				}
			}
		}
	}
	if len(baseTypes) > 0 {
		rawBytes, _ := json.Marshal(baseTypes)
		_ = json.Unmarshal(rawBytes, &qReq.HasBaseTypeList)
	}

	res, err := i.periodicalSvc.QueryPeriodical(ctx, qReq)
	if err != nil {
		return nil, err
	}

	var briefItems []ai.PeriodicalBriefCard
	for _, item := range res.Items {
		briefItems = append(briefItems, formatPeriodicalToBrief(item))
	}

	filterPayload := map[string]interface{}{
		"name":               qReq.Name,
		"competent_unit":     qReq.CompetentUnit,
		"organization":       qReq.Organization,
		"domestic_number":    qReq.DomesticNumber,
		"column_setting":     qReq.ColumnSetting,
		"publish_month":      publishMonth,
		"invoice":            qReq.Invoice,
		"is_warp":            qReq.IsWarp,
		"works":              qReq.Works,
		"is_hidden":          qReq.IsHidden,
		"has_base_type_list": baseTypes,
		"matched_cat_names":  matchedCatNames,
	}

	descParts := []string{}
	if qReq.IsHidden == 1 {
		descParts = append(descParts, "状态：仅已隐藏期刊")
	} else if qReq.IsHidden == 3 {
		descParts = append(descParts, "状态：含已隐藏期刊")
	}
	if qReq.Name != "" {
		descParts = append(descParts, "刊名含“"+qReq.Name+"”")
	}
	if len(matchedCatNames) > 0 {
		descParts = append(descParts, "分类："+strings.Join(matchedCatNames, " / "))
	}
	if publishMonth != "" {
		descParts = append(descParts, "出刊月份："+publishMonth)
	}
	if qReq.IsWarp == 1 {
		descParts = append(descParts, "全包")
	}
	if qReq.Invoice == 1 {
		descParts = append(descParts, "可开杂志社发票")
	}
	desc := "全部期刊检索"
	if len(descParts) > 0 {
		desc = "筛选条件：" + strings.Join(descParts, " · ")
	}

	// 动态生成上下文智能追问建议胶囊
	var suggestions []string
	if len(matchedCatNames) > 0 {
		catSummary := strings.Join(matchedCatNames, "")
		if qReq.IsWarp != 1 {
			suggestions = append(suggestions, fmt.Sprintf("进一步筛选其中支持全包且可开发票的%s期刊", catSummary))
		}
		if qReq.OrderBy != "price" {
			suggestions = append(suggestions, fmt.Sprintf("将上述%s期刊按发表费从低到高排序", catSummary))
		}
	} else {
		suggestions = append(suggestions, "筛选知网收录且支持全包的省级期刊", "查询所有国家级期刊并按价格从低到高排序")
	}
	if len(briefItems) > 0 {
		suggestions = append(suggestions, fmt.Sprintf("基于《%s》等前3本期刊生成一份对外客户选刊推荐报价方案", strings.Trim(briefItems[0].Name, "《》")))
	}

	return &ai.ActionCardData{
		CardType:      "query_result",
		Status:        "readonly",
		Title:         fmt.Sprintf("检索到 %d 本符合条件的期刊", res.Total),
		Description:   desc,
		TotalCount:    res.Total,
		Items:         briefItems,
		FilterPayload: filterPayload,
		Suggestions:   suggestions,
	}, nil
}

func (i *AIImpl) execPrepareCreateTool(args map[string]interface{}, allCats []*periodical.Categorize) *ai.ActionCardData {
	rawName := getStringArg(args, "name")
	if rawName == "" {
		rawName = "《未命名期刊》"
	} else if !strings.HasPrefix(rawName, "《") {
		rawName = "《" + strings.Trim(rawName, "《》") + "》"
	}

	var dirIDs, colIDs, levelIDs []int
	var cycleID int
	catLabels := map[string][]string{
		"directions":  {},
		"collections": {},
		"levels":      {},
		"cycle":       {},
	}

	var catKeywords []string
	if rawNames, ok := args["category_names"].([]interface{}); ok {
		for _, rn := range rawNames {
			catKeywords = append(catKeywords, strings.TrimSpace(fmt.Sprintf("%v", rn)))
		}
	}

	for _, kw := range catKeywords {
		if kw == "" {
			continue
		}
		for _, cat := range allCats {
			if strings.Contains(cat.Name, kw) || strings.Contains(kw, cat.Name) {
				switch cat.Type {
				case 1:
					dirIDs = append(dirIDs, cat.CategorizeID)
					catLabels["directions"] = append(catLabels["directions"], cat.Name)
				case 2:
					levelIDs = append(levelIDs, cat.CategorizeID)
					catLabels["levels"] = append(catLabels["levels"], cat.Name)
				case 3:
					colIDs = append(colIDs, cat.CategorizeID)
					catLabels["collections"] = append(catLabels["collections"], cat.Name)
				case 4:
					if cycleID == 0 {
						cycleID = cat.CategorizeID
						catLabels["cycle"] = []string{cat.Name}
					}
				}
			}
		}
	}

	// 若未匹配到任何收稿方向，默认兜底第一个收稿方向以便表单必填项通过
	if len(dirIDs) == 0 {
		for _, cat := range allCats {
			if cat.Type == 1 && strings.Contains(cat.Name, "综合") {
				dirIDs = append(dirIDs, cat.CategorizeID)
				catLabels["directions"] = append(catLabels["directions"], cat.Name)
				break
			}
		}
	}

	var allCatIDs []int
	allCatIDs = append(allCatIDs, dirIDs...)
	allCatIDs = append(allCatIDs, colIDs...)
	allCatIDs = append(allCatIDs, levelIDs...)
	if cycleID > 0 {
		allCatIDs = append(allCatIDs, cycleID)
	}

	isWarp, _ := getBoolArg(args, "is_warp")
	invoice, _ := getBoolArg(args, "invoice_receipt")
	works, _ := getBoolArg(args, "works")
	isHidden, _ := getBoolArg(args, "is_hidden")

	pubDate := normalizePublishDate(getStringArg(args, "publication_time"))
	cleanPrice := strings.TrimSuffix(strings.TrimSpace(getStringArg(args, "price")), "元")
	cleanPage := strings.TrimSuffix(strings.TrimSpace(getStringArg(args, "periodical_page")), "页")
	rawDesc := strings.TrimSpace(getStringArg(args, "periodical_description"))
	if len([]rune(rawDesc)) <= 16 && (strings.HasPrefix(rawDesc, "出版周期") || strings.HasPrefix(rawDesc, "级别") || strings.HasPrefix(rawDesc, "网站收录")) {
		rawDesc = ""
	}

	formPayload := map[string]interface{}{
		"name":                       rawName,
		"publication_time":           pubDate,
		"competent_unit":             getStringArg(args, "competent_unit"),
		"organization":               getStringArg(args, "organization"),
		"domestic_number":            getStringArg(args, "domestic_number"),
		"international_number":       getStringArg(args, "international_number"),
		"founding_time":              getStringArg(args, "founding_time"),
		"citation_rate":              getStringArg(args, "citation_rate"),
		"audit_time":                 getStringArg(args, "audit_time"),
		"price":                      cleanPrice,
		"periodical_page":            cleanPage,
		"is_warp":                    isWarp,
		"invoice_receipt":            invoice,
		"works":                      works,
		"is_hidden":                  isHidden,
		"composite_influence_factor": getStringArg(args, "composite_influence_factor"),
		"synthetic_influence_factor": getStringArg(args, "synthetic_influence_factor"),
		"column_setting":             getStringArg(args, "column_setting"),
		"periodical_description":     rawDesc,
		"attention_matter":           getStringArg(args, "attention_matter"),
		"author_info":                getStringArg(args, "author_info"),
		"categorize_ids":             allCatIDs,
		"directions":                 dirIDs,
		"collections":                colIDs,
		"levels":                     levelIDs,
		"cycle":                      cycleID,
	}

	return &ai.ActionCardData{
		CardType:       "create_confirm",
		Status:         "pending",
		Title:          "新增期刊入库确认单",
		Description:    "已为您自动提取期刊信息并匹配分类字典，您可在下方表单中直接修改、补充任意字段后点击【确认直接入库】。",
		PeriodicalName: rawName,
		FormPayload:    formPayload,
		CategoryLabels: catLabels,
		Suggestions: []string{
			fmt.Sprintf("查询库中与%s同级别的其他期刊", rawName),
			fmt.Sprintf("为%s撰写一段对外客户征稿推荐文案", rawName),
		},
	}
}

func (i *AIImpl) findPeriodicalByName(ctx context.Context, name string) (*periodical.Periodical, error) {
	clean := strings.Trim(strings.TrimSpace(name), "《》")
	if clean == "" {
		return nil, errors.New("请提供要操作的期刊名称")
	}
	qReq := periodical.NewQueryPeriodicalRequest()
	qReq.PageNumber = 1
	qReq.PageSize = 5
	qReq.Name = clean
	qReq.IsHidden = 3 // 允许定位全部状态（含已隐藏）的期刊，支持对已隐藏期刊执行取消隐藏或修改
	res, err := i.periodicalSvc.QueryPeriodical(ctx, qReq)
	if err != nil {
		return nil, err
	}
	if len(res.Items) == 0 {
		return nil, fmt.Errorf("未在期刊库中找到名称包含“%s”的刊物", clean)
	}
	// 优先精确匹配
	for _, item := range res.Items {
		if strings.Trim(item.Name, "《》") == clean {
			return i.periodicalSvc.DescribePeriodical(ctx, periodical.NewDescribePeriodicalRequest(strconv.Itoa(item.PeriodicalID)))
		}
	}
	return i.periodicalSvc.DescribePeriodical(ctx, periodical.NewDescribePeriodicalRequest(strconv.Itoa(res.Items[0].PeriodicalID)))
}

func extractTargetNamesFromArgs(args map[string]interface{}) []string {
	var result []string
	seen := make(map[string]bool)
	addClean := func(raw string) {
		c := strings.Trim(strings.TrimSpace(raw), "《》\"'")
		if c != "" && !seen[c] {
			seen[c] = true
			result = append(result, c)
		}
	}

	if rawList, ok := args["target_names"].([]interface{}); ok {
		for _, item := range rawList {
			addClean(fmt.Sprintf("%v", item))
		}
	} else if strList, ok := args["target_names"].([]string); ok {
		for _, item := range strList {
			addClean(item)
		}
	}

	tName := getStringArg(args, "target_name")
	if tName != "" {
		matches := bookTitleRegex.FindAllStringSubmatch(tName, -1)
		if len(matches) > 1 {
			for _, m := range matches {
				if len(m) > 1 {
					addClean(m[1])
				}
			}
		} else if strings.ContainsAny(tName, "、,，;；") {
			parts := strings.FieldsFunc(tName, func(r rune) bool {
				return r == '、' || r == ',' || r == '，' || r == ';' || r == '；'
			})
			for _, p := range parts {
				addClean(p)
			}
		} else {
			addClean(tName)
		}
	}
	return result
}

func normalizeCategoryFieldText(rawText string, catType int64, allCats []*periodical.Categorize) (string, []int64) {
	s := strings.TrimSpace(rawText)
	if s == "" {
		return "", nil
	}
	parts := strings.FieldsFunc(s, func(r rune) bool {
		return r == '、' || r == ',' || r == '，' || r == ';' || r == '；' || r == '/' || r == ' '
	})
	if len(parts) == 0 {
		parts = []string{s}
	}

	var matchedNames []string
	var matchedIDs []int64
	seenID := make(map[int]bool)

	for _, kw := range parts {
		kw = strings.TrimSpace(kw)
		if kw == "" {
			continue
		}
		if catType == 4 {
			if kw == "半月刊" || kw == "旬刊" || kw == "周刊" {
				kw = "半月刊及以下"
			} else if kw == "月刊" || kw == "双月刊" || kw == "季刊" || kw == "半年刊" || kw == "年刊" {
				kw = "月刊及以上"
			}
		}
		for _, c := range allCats {
			if c.Type != catType {
				continue
			}
			if strings.EqualFold(c.Name, kw) || strings.Contains(c.Name, kw) || strings.Contains(kw, c.Name) {
				if !seenID[c.CategorizeID] {
					seenID[c.CategorizeID] = true
					matchedNames = append(matchedNames, c.Name)
					matchedIDs = append(matchedIDs, int64(c.CategorizeID))
				}
				if catType == 2 || catType == 4 {
					break
				}
			}
		}
	}

	if len(matchedNames) > 0 {
		return strings.Join(matchedNames, "、"), matchedIDs
	}
	return s, nil
}

func syncCategoryFieldsToCategorizeIDs(m map[string]interface{}, allCats []*periodical.Categorize) {
	if m == nil || len(allCats) == 0 {
		return
	}
	catTypeByID := make(map[int64]int64)
	for _, c := range allCats {
		catTypeByID[int64(c.CategorizeID)] = c.Type
	}

	var currentIDs []int64
	if rawIDs, ok := m["categorize_ids"].([]interface{}); ok {
		for _, idVal := range rawIDs {
			switch v := idVal.(type) {
			case float64:
				currentIDs = append(currentIDs, int64(v))
			case int:
				currentIDs = append(currentIDs, int64(v))
			case int64:
				currentIDs = append(currentIDs, v)
			}
		}
	} else if int64IDs, ok := m["categorize_ids"].([]int64); ok {
		currentIDs = append(currentIDs, int64IDs...)
	}

	replaceTypeIDs := func(catType int64, fieldKey string) {
		val, exists := m[fieldKey]
		if !exists {
			return
		}
		strVal := strings.TrimSpace(fmt.Sprintf("%v", val))
		_, newTypeIDs := normalizeCategoryFieldText(strVal, catType, allCats)
		if strVal != "" && len(newTypeIDs) == 0 {
			return
		}
		var kept []int64
		for _, cid := range currentIDs {
			if catTypeByID[cid] != catType {
				kept = append(kept, cid)
			}
		}
		currentIDs = append(kept, newTypeIDs...)
	}

	replaceTypeIDs(2, "level")
	replaceTypeIDs(4, "cycle")
	replaceTypeIDs(3, "collections")
	replaceTypeIDs(1, "directions")

	m["categorize_ids"] = currentIDs
}

func (i *AIImpl) buildSingleUpdateTarget(ctx context.Context, targetName string, args map[string]interface{}, allCats []*periodical.Categorize) (*ai.BatchUpdateTarget, error) {
	existing, err := i.findPeriodicalByName(ctx, targetName)
	if err != nil {
		return nil, err
	}

	var existingCatIDs []int64
	var levelName string
	var cycleName string
	var collectionNames []string
	var directionNames []string
	for _, c := range existing.Categorizes {
		existingCatIDs = append(existingCatIDs, int64(c.CategorizeID))
		switch c.Type {
		case 1:
			directionNames = append(directionNames, c.Name)
		case 2:
			if levelName == "" {
				levelName = c.Name
			}
		case 3:
			collectionNames = append(collectionNames, c.Name)
		case 4:
			if cycleName == "" {
				cycleName = c.Name
			}
		}
	}
	oldCollectionsStr := strings.Join(collectionNames, "、")
	oldDirectionsStr := strings.Join(directionNames, "、")

	oldPubTime := existing.PublicationTime
	if len(oldPubTime) >= 10 {
		oldPubTime = oldPubTime[:10]
	}

	payload := map[string]interface{}{
		"periodical_id":              existing.PeriodicalID,
		"name":                       existing.Name,
		"publication_time":           oldPubTime,
		"competent_unit":             existing.CompetentUnit,
		"organization":               existing.Organization,
		"domestic_number":            existing.DomesticNumber,
		"international_number":       existing.InternationalNumber,
		"founding_time":              existing.FoundingTime,
		"citation_rate":              existing.CitationRate,
		"periodical_description":     existing.PeriodicalDescription,
		"column_setting":             existing.ColumnSetting,
		"audit_time":                 existing.AuditTime,
		"invoice_receipt":            existing.InvoiceReceipt,
		"attention_matter":           existing.AttentionMatter,
		"article_naming":             existing.ArticleNaming,
		"internal_processes":         existing.InternalProcesses,
		"publication_process":        existing.PublicationProcess,
		"all_inclusive_process":      existing.AllInclusiveProcess,
		"price":                      existing.Price,
		"is_warp":                    existing.IsWarp,
		"periodical_page":            existing.PeriodicalPage,
		"composite_influence_factor": existing.CompositeInfluenceFactor,
		"synthetic_influence_factor": existing.SyntheticInfluenceFactor,
		"periodical_batch":           existing.PeriodicalBatch,
		"works":                      existing.Works,
		"is_hidden":                  existing.IsHidden,
		"periodical_final_proof":     existing.PeriodicalFinalProof,
		"author_info":                existing.AuthorInfo,
		"level":                      levelName,
		"cycle":                      cycleName,
		"collections":                oldCollectionsStr,
		"directions":                 oldDirectionsStr,
		"categorize_ids":             existingCatIDs,
	}

	var diffs []ai.FieldDiff
	checkAndApply := func(field, label, oldVal, newVal string) {
		if newVal != "" {
			displayOld := oldVal
			if displayOld == "" {
				displayOld = "(空)"
			}
			diffs = append(diffs, ai.FieldDiff{
				Field:    field,
				Label:    label,
				OldValue: displayOld,
				NewValue: newVal,
			})
			payload[field] = newVal
		}
	}

	newPub := normalizePublishDateWithRef(getStringArg(args, "publication_time"), oldPubTime)
	checkAndApply("publication_time", "出刊时间", oldPubTime, newPub)
	checkAndApply("price", "期刊发表费", existing.Price, strings.TrimSuffix(getStringArg(args, "price"), "元"))
	checkAndApply("periodical_page", "期刊页码", existing.PeriodicalPage, strings.TrimSuffix(getStringArg(args, "periodical_page"), "页"))
	if rawLvl := getStringArg(args, "level"); rawLvl != "" {
		normLvl, _ := normalizeCategoryFieldText(rawLvl, 2, allCats)
		checkAndApply("level", "期刊级别", levelName, normLvl)
	}
	if rawCyc := getStringArg(args, "cycle"); rawCyc != "" {
		normCyc, _ := normalizeCategoryFieldText(rawCyc, 4, allCats)
		checkAndApply("cycle", "出版周期", cycleName, normCyc)
	}
	if rawCol := getStringArg(args, "collections"); rawCol != "" {
		normCol, _ := normalizeCategoryFieldText(rawCol, 3, allCats)
		checkAndApply("collections", "网站收录", oldCollectionsStr, normCol)
	}
	if rawDir := getStringArg(args, "directions"); rawDir != "" {
		normDir, _ := normalizeCategoryFieldText(rawDir, 1, allCats)
		checkAndApply("directions", "收稿方向", oldDirectionsStr, normDir)
	}
	checkAndApply("competent_unit", "主管单位", existing.CompetentUnit, getStringArg(args, "competent_unit"))
	checkAndApply("organization", "主办单位", existing.Organization, getStringArg(args, "organization"))
	checkAndApply("domestic_number", "国内刊号", existing.DomesticNumber, getStringArg(args, "domestic_number"))
	checkAndApply("international_number", "国际刊号", existing.InternationalNumber, getStringArg(args, "international_number"))
	checkAndApply("audit_time", "审核周期", existing.AuditTime, getStringArg(args, "audit_time"))
	checkAndApply("column_setting", "栏目设置", existing.ColumnSetting, getStringArg(args, "column_setting"))
	checkAndApply("attention_matter", "注意事项", existing.AttentionMatter, getStringArg(args, "attention_matter"))
	checkAndApply("periodical_description", "期刊描述", existing.PeriodicalDescription, getStringArg(args, "periodical_description"))
	checkAndApply("founding_time", "创刊时间", existing.FoundingTime, getStringArg(args, "founding_time"))
	checkAndApply("citation_rate", "引用率", existing.CitationRate, getStringArg(args, "citation_rate"))
	checkAndApply("composite_influence_factor", "复合影响因子", existing.CompositeInfluenceFactor, getStringArg(args, "composite_influence_factor"))
	checkAndApply("synthetic_influence_factor", "综合影响因子", existing.SyntheticInfluenceFactor, getStringArg(args, "synthetic_influence_factor"))
	checkAndApply("author_info", "作者与样刊说明", existing.AuthorInfo, getStringArg(args, "author_info"))
	checkAndApply("periodical_final_proof", "期刊清样", existing.PeriodicalFinalProof, getStringArg(args, "periodical_final_proof"))
	checkAndApply("article_naming", "文章命名", existing.ArticleNaming, getStringArg(args, "article_naming"))
	checkAndApply("publication_process", "发表流程", existing.PublicationProcess, getStringArg(args, "publication_process"))
	checkAndApply("all_inclusive_process", "全包流程", existing.AllInclusiveProcess, getStringArg(args, "all_inclusive_process"))
	checkAndApply("internal_processes", "内部流程", existing.InternalProcesses, getStringArg(args, "internal_processes"))

	syncCategoryFieldsToCategorizeIDs(payload, allCats)

	if bVal, ok := getBoolArg(args, "is_warp"); ok && bVal != existing.IsWarp {
		diffs = append(diffs, ai.FieldDiff{
			Field:    "is_warp",
			Label:    "是否全包",
			OldValue: formatBoolCN(existing.IsWarp),
			NewValue: formatBoolCN(bVal),
		})
		payload["is_warp"] = bVal
	}
	if bVal, ok := getBoolArg(args, "invoice_receipt"); ok && bVal != existing.InvoiceReceipt {
		diffs = append(diffs, ai.FieldDiff{
			Field:    "invoice_receipt",
			Label:    "可开发票",
			OldValue: formatBoolCN(existing.InvoiceReceipt),
			NewValue: formatBoolCN(bVal),
		})
		payload["invoice_receipt"] = bVal
	}
	if bVal, ok := getBoolArg(args, "works"); ok && bVal != existing.Works {
		diffs = append(diffs, ai.FieldDiff{
			Field:    "works",
			Label:    "可发作品",
			OldValue: formatBoolCN(existing.Works),
			NewValue: formatBoolCN(bVal),
		})
		payload["works"] = bVal
	}

	hasExplicitHidden := false
	if bVal, ok := getBoolArg(args, "is_hidden"); ok {
		hasExplicitHidden = true
		if bVal != existing.IsHidden {
			diffs = append(diffs, ai.FieldDiff{
				Field:    "is_hidden",
				Label:    "是否隐藏",
				OldValue: formatBoolCN(existing.IsHidden),
				NewValue: formatBoolCN(bVal),
			})
			payload["is_hidden"] = bVal
		}
	}

	// 若用户仅说了“修改《某刊》”而未在句子里带具体新值，自动预置核心常改项供用户在确认卡片内直接输入修改
	if len(diffs) == 0 {
		if hasExplicitHidden {
			diffs = []ai.FieldDiff{
				{
					Field:    "is_hidden",
					Label:    "是否隐藏",
					OldValue: formatBoolCN(existing.IsHidden),
					NewValue: formatBoolCN(existing.IsHidden),
				},
			}
		} else {
			diffs = []ai.FieldDiff{
				{Field: "publication_time", Label: "出刊时间", OldValue: oldPubTime, NewValue: oldPubTime},
				{Field: "price", Label: "期刊发表费", OldValue: existing.Price, NewValue: existing.Price},
				{Field: "periodical_page", Label: "期刊页码", OldValue: existing.PeriodicalPage, NewValue: existing.PeriodicalPage},
			}
		}
	}

	return &ai.BatchUpdateTarget{
		PeriodicalID:   existing.PeriodicalID,
		PeriodicalName: existing.Name,
		Level:          levelName,
		FormPayload:    payload,
		Diffs:          diffs,
	}, nil
}

func (i *AIImpl) execPrepareUpdateTool(ctx context.Context, args map[string]interface{}, allCats []*periodical.Categorize) (*ai.ActionCardData, error) {
	targetNames := extractTargetNamesFromArgs(args)
	if len(targetNames) == 0 {
		return nil, errors.New("请提供要修改的期刊名称")
	}

	var targets []ai.BatchUpdateTarget
	var lastErr error
	for _, tName := range targetNames {
		t, err := i.buildSingleUpdateTarget(ctx, tName, args, allCats)
		if err != nil {
			lastErr = err
			continue
		}
		// 避免同名重复加入
		exists := false
		for _, prev := range targets {
			if prev.PeriodicalID == t.PeriodicalID {
				exists = true
				break
			}
		}
		if !exists {
			targets = append(targets, *t)
		}
	}

	if len(targets) == 0 {
		if lastErr != nil {
			return nil, lastErr
		}
		return nil, errors.New("未找到匹配的目标期刊")
	}

	first := targets[0]
	if len(targets) == 1 {
		return &ai.ActionCardData{
			CardType:       "update_confirm",
			Status:         "pending",
			Title:          fmt.Sprintf("修改期刊确认单 · %s", first.PeriodicalName),
			Description:    fmt.Sprintf("即将更新 %s 的 %d 项属性，请核对以下变更对比（Diff）后点击【确认执行修改】。", first.PeriodicalName, len(first.Diffs)),
			PeriodicalID:   first.PeriodicalID,
			PeriodicalName: first.PeriodicalName,
			FormPayload:    first.FormPayload,
			Diffs:          first.Diffs,
			BatchUpdates:   targets,
			Suggestions: []string{
				fmt.Sprintf("查看 %s 修改后的完整详情信息", first.PeriodicalName),
				"继续查询其他需要调整排期的期刊",
			},
		}, nil
	}

	var names []string
	for _, t := range targets {
		names = append(names, t.PeriodicalName)
	}
	return &ai.ActionCardData{
		CardType:       "update_confirm",
		Status:         "pending",
		Title:          fmt.Sprintf("批量修改期刊确认单 · 共 %d 本期刊", len(targets)),
		Description:    fmt.Sprintf("即将批量更新 %s 共 %d 本期刊的属性，您可在下方对各刊物独立微调、增删字段或一键同步新值，核对无误后点击【确认批量修改】。", strings.Join(names, "、"), len(targets)),
		PeriodicalID:   first.PeriodicalID,
		PeriodicalName: first.PeriodicalName,
		FormPayload:    first.FormPayload,
		Diffs:          first.Diffs,
		BatchUpdates:   targets,
		Suggestions: []string{
			fmt.Sprintf("查询 %s 的最新详情", strings.Join(names, "、")),
			"继续查询其他需要调整排期的期刊",
		},
	}, nil
}

// mergeActionCards 当大模型在单轮对话中并发调用多次 prepare_update_periodical 时，将其合并为一张批量修改确认卡片，避免后一次覆盖前一次
func mergeActionCards(prev *ai.ActionCardData, next *ai.ActionCardData) *ai.ActionCardData {
	if prev == nil {
		return next
	}
	if next == nil {
		return prev
	}
	if prev.CardType == "update_confirm" && next.CardType == "update_confirm" {
		extractTargets := func(c *ai.ActionCardData) []ai.BatchUpdateTarget {
			if len(c.BatchUpdates) > 0 {
				return c.BatchUpdates
			}
			if c.PeriodicalID > 0 {
				return []ai.BatchUpdateTarget{
					{
						PeriodicalID:   c.PeriodicalID,
						PeriodicalName: c.PeriodicalName,
						FormPayload:    c.FormPayload,
						Diffs:          c.Diffs,
					},
				}
			}
			return nil
		}

		merged := extractTargets(prev)
		nextTargets := extractTargets(next)

		for _, nt := range nextTargets {
			foundIdx := -1
			for idx, pt := range merged {
				if pt.PeriodicalID == nt.PeriodicalID {
					foundIdx = idx
					break
				}
			}
			if foundIdx == -1 {
				merged = append(merged, nt)
			} else {
				// 同一本期刊多次调用时合并字段变更
				if merged[foundIdx].FormPayload == nil {
					merged[foundIdx].FormPayload = make(map[string]interface{})
				}
				for _, nd := range nt.Diffs {
					merged[foundIdx].FormPayload[nd.Field] = nt.FormPayload[nd.Field]
					diffReplaced := false
					for dIdx, ed := range merged[foundIdx].Diffs {
						if ed.Field == nd.Field {
							merged[foundIdx].Diffs[dIdx].NewValue = nd.NewValue
							diffReplaced = true
							break
						}
					}
					if !diffReplaced {
						merged[foundIdx].Diffs = append(merged[foundIdx].Diffs, nd)
					}
				}
			}
		}

		prev.BatchUpdates = merged
		if len(merged) > 0 {
			prev.PeriodicalID = merged[0].PeriodicalID
			prev.PeriodicalName = merged[0].PeriodicalName
			prev.FormPayload = merged[0].FormPayload
			prev.Diffs = merged[0].Diffs
		}
		if len(merged) > 1 {
			var names []string
			for _, m := range merged {
				names = append(names, m.PeriodicalName)
			}
			prev.Title = fmt.Sprintf("批量修改期刊确认单 · 共 %d 本期刊", len(merged))
			prev.Description = fmt.Sprintf("即将批量更新 %s 共 %d 本期刊的属性，您可在下方对各刊物独立微调、增删字段或一键同步新值，核对无误后点击【确认批量修改】。", strings.Join(names, "、"), len(merged))
			prev.Suggestions = []string{
				fmt.Sprintf("查询 %s 的最新详情", strings.Join(names, "、")),
				"继续查询其他需要调整排期的期刊",
			}
		}
		return prev
	}
	return next
}

func (i *AIImpl) execPrepareDeleteTool(ctx context.Context, args map[string]interface{}) (*ai.ActionCardData, error) {
	targetName := getStringArg(args, "target_name")
	existing, err := i.findPeriodicalByName(ctx, targetName)
	if err != nil {
		return nil, err
	}
	brief := formatPeriodicalToBrief(existing)
	return &ai.ActionCardData{
		CardType:       "delete_confirm",
		Status:         "pending",
		Title:          fmt.Sprintf("高危操作：删除期刊 %s", existing.Name),
		Description:    "删除后该期刊及其分类关联将从数据库移除，请再次确认是否执行删除操作。",
		PeriodicalID:   existing.PeriodicalID,
		PeriodicalName: existing.Name,
		Items:          []ai.PeriodicalBriefCard{brief},
	}, nil
}

func formatBoolCN(b bool) string {
	if b {
		return "是"
	}
	return "否"
}

var pureMonthOnlyRe = regexp.MustCompile(`^(0?[1-9]|1[0-2])月(初|中旬|下旬|月底|月末)?$`)
var refYearRe = regexp.MustCompile(`^(20\d{2})-\d{2}`)

// normalizePublishDateWithRef 支持当用户仅输入 "1月"、"9月中旬"（不含年份）且原排期为标准 YYYY-MM-DD 时自动继承原排期年份
func normalizePublishDateWithRef(raw string, refDate string) string {
	norm := normalizePublishDate(raw)
	if norm == "" {
		return ""
	}
	if m := pureMonthOnlyRe.FindStringSubmatch(norm); len(m) > 1 {
		if ym := refYearRe.FindStringSubmatch(strings.TrimSpace(refDate)); len(ym) > 1 {
			mon, _ := strconv.Atoi(m[1])
			day := "01"
			if strings.Contains(norm, "中旬") {
				day = "15"
			} else if strings.Contains(norm, "下旬") {
				day = "21"
			} else if strings.Contains(norm, "月底") || strings.Contains(norm, "月末") {
				day = "28"
			}
			return fmt.Sprintf("%s-%02d-%s", ym[1], mon, day)
		}
	}
	return norm
}

// normalizePublishDate 将 "2025-09"、"2025年9月中旬"、"2025-09月底" 等自然语言排期统一规范化为 YYYY-MM-DD
func normalizePublishDate(raw string) string {
	s := strings.TrimSpace(raw)
	if s == "" {
		return ""
	}
	// 若已是标准 YYYY-MM-DD
	if matched, _ := regexp.MatchString(`^\d{4}-\d{2}-\d{2}$`, s); matched {
		return s
	}
	// 提取年、月、日
	fullDateRe := regexp.MustCompile(`(20\d{2})[年\-/\.]+(0?[1-9]|1[0-2])[月\-/\.]+(0?[1-9]|[12]\d|3[01])`)
	if m := fullDateRe.FindStringSubmatch(s); len(m) > 3 {
		mon, _ := strconv.Atoi(m[2])
		day, _ := strconv.Atoi(m[3])
		return fmt.Sprintf("%s-%02d-%02d", m[1], mon, day)
	}
	// 提取年、月 + 旬期描述
	if m := monthRegex.FindStringSubmatch(s); len(m) > 2 {
		mon, _ := strconv.Atoi(m[2])
		day := "01"
		if strings.Contains(s, "中旬") {
			day = "15"
		} else if strings.Contains(s, "下旬") {
			day = "21"
		} else if strings.Contains(s, "月底") || strings.Contains(s, "月末") {
			day = "28"
		}
		return fmt.Sprintf("%s-%02d-%s", m[1], mon, day)
	}
	return s
}

// formatCardContextForLLM 将上一轮交互卡片的结构化上下文注入多轮历史，确保多轮指代与连续追问丝滑连贯
func formatCardContextForLLM(card *ai.ActionCardData) string {
	if card == nil {
		return ""
	}
	switch card.CardType {
	case "query_result":
		var itemDescs []string
		for idx, it := range card.Items {
			if idx >= 8 {
				break
			}
			itemDescs = append(itemDescs, fmt.Sprintf("%d.%s(ID:%d,级别:%s,价格:%s,排期:%s)", idx+1, it.Name, it.PeriodicalID, it.Level, it.Price, it.PublicationTime))
		}
		return fmt.Sprintf("[系统上下文：本条回复下方附带了查询结果卡片，共匹配 %d 本，首屏列表：%s]", card.TotalCount, strings.Join(itemDescs, "；"))
	case "update_confirm":
		if len(card.BatchUpdates) > 1 {
			var names []string
			for _, bu := range card.BatchUpdates {
				names = append(names, fmt.Sprintf("%s(ID:%d)", bu.PeriodicalName, bu.PeriodicalID))
			}
			return fmt.Sprintf("[系统上下文：本条回复已生成批量修改确认卡片，目标期刊共 %d 本：%s，卡片状态：%s]", len(card.BatchUpdates), strings.Join(names, "、"), card.Status)
		}
		b, _ := json.Marshal(card.Diffs)
		return fmt.Sprintf("[系统上下文：本条回复已生成修改确认卡片，目标期刊：%s(ID:%d)，卡片状态：%s，变更项：%s]", card.PeriodicalName, card.PeriodicalID, card.Status, string(b))
	case "create_confirm":
		return fmt.Sprintf("[系统上下文：本条回复已生成新增期刊入库确认卡片，目标期刊：%s，卡片状态：%s]", card.PeriodicalName, card.Status)
	case "delete_confirm":
		return fmt.Sprintf("[系统上下文：本条回复已生成删除确认卡片，目标期刊：%s(ID:%d)，卡片状态：%s]", card.PeriodicalName, card.PeriodicalID, card.Status)
	}
	return ""
}

func (i *AIImpl) buildSummaryFromCard(card *ai.ActionCardData, steps []ai.ToolStep) string {
	if card == nil {
		if len(steps) > 0 {
			return steps[0].Summary
		}
		return "操作已完成。"
	}
	switch card.CardType {
	case "query_result":
		if card.TotalCount == 0 {
			return "未检索到符合该条件的期刊，您可以尝试放宽筛选条件（如调整出刊月份或收录网站）。"
		}
		var warpCount, invoiceCount int
		var sampleNames []string
		for idx, it := range card.Items {
			if it.IsWarp {
				warpCount++
			}
			if it.InvoiceReceipt {
				invoiceCount++
			}
			if idx < 3 {
				sampleNames = append(sampleNames, it.Name)
			}
		}
		summaryText := fmt.Sprintf("### 📊 选刊检索洞察\n- **匹配规模**：共为您从期刊库中检索到 **%d** 本符合条件的刊物（%s）。\n- **首屏样本特征**：前 %d 本刊物中，支持全包共 **%d** 本，可开杂志社发票共 **%d** 本。\n", card.TotalCount, card.Description, len(card.Items), warpCount, invoiceCount)
		if len(sampleNames) > 0 {
			summaryText += fmt.Sprintf("- **重点关注刊物**：如 **%s** 等，您可在下方交互卡片中切换【网格/表格视图】、直接翻页，或点击【在右侧画布展开】进行深度对比与生成客户报价单。", strings.Join(sampleNames, "、"))
		}
		return summaryText
	case "create_confirm":
		return fmt.Sprintf("已为您从文本中提取出 **%s** 的结构化信息并自动匹配了系统分类字典。请核对下方入库确认单，点击 **【确认直接入库】** 或 **【填入新增抽屉微调】**：", card.PeriodicalName)
	case "update_confirm":
		var errNotes []string
		for _, st := range steps {
			if st.Status == "error" && st.Summary != "" {
				errNotes = append(errNotes, st.Summary)
			}
		}
		errSuffix := ""
		if len(errNotes) > 0 {
			errSuffix = fmt.Sprintf("（注：%s）", strings.Join(errNotes, "；"))
		}
		if len(card.BatchUpdates) > 1 {
			var boldNames []string
			totalDiffs := 0
			for _, bu := range card.BatchUpdates {
				boldNames = append(boldNames, "**"+bu.PeriodicalName+"**")
				totalDiffs += len(bu.Diffs)
			}
			return fmt.Sprintf("已定位到 %s 共 **%d 本**期刊，并为您生成了批量修改前后的字段对比（共 **%d** 项变更）%s。请核对下方卡片并点击 **【确认批量修改】**：", strings.Join(boldNames, "、"), len(card.BatchUpdates), totalDiffs, errSuffix)
		}
		return fmt.Sprintf("已定位到 **%s**，并为您生成了修改前后的字段对比（共 **%d** 项变更）%s。请核对下方 Diff 卡片并点击 **【确认执行修改】**：", card.PeriodicalName, len(card.Diffs), errSuffix)
	case "delete_confirm":
		return fmt.Sprintf("已定位到目标刊物 **%s**（ID: `%d`）。由于删除属于高危不可逆操作，请在下方卡片中点击 **【确认删除】** 完成最终执行：", card.PeriodicalName, card.PeriodicalID)
	default:
		return "请查看下方操作卡片。"
	}
}

// ==================== 本地智能解析回退引擎 (Smart Local Agent Fallback) ====================

var (
	bookTitleRegex = regexp.MustCompile(`《([^》]+)》`)
	priceRegex     = regexp.MustCompile(`(?:价格|发表费|费用|版面费|全包价?|¥|￥)\s*[:：为是调到改到]*\s*(\d+(?:\.\d+)?(?:[^\s，,。；;]{0,10}))`)
	pageRegex      = regexp.MustCompile(`(\d+)\s*页`)
	monthRegex     = regexp.MustCompile(`(202\d)[年\-/\.]+(0?[1-9]|1[0-2])`)
	shortMonthRe   = regexp.MustCompile(`(?:改为|改成|调到|调整为|设为|时间[为是]?)\s*((?:0?[1-9]|1[0-2])月(?:初|中旬|下旬|月底|月末)?)`)
	cnNumberRegex  = regexp.MustCompile(`(?:国内刊号|CN)[:：\s]*([0-9]{2}-[0-9]{3,5}/[A-Za-z0-9]+)`)
	issnRegex      = regexp.MustCompile(`(?:国际刊号|ISSN)[:：\s]*([0-9]{4}-[0-9]{3}[0-9Xx])`)
)

func (i *AIImpl) runLocalSmartAgent(
	ctx context.Context,
	req *ai.SendMessageRequest,
	allCats []*periodical.Categorize,
	isFallbackFromError bool,
	llmErrMsg string,
) (string, []ai.ToolStep, *ai.ActionCardData) {
	prompt := req.Prompt
	var matchedCats []string
	for _, c := range allCats {
		cleanCat := strings.TrimSuffix(c.Name, "类")
		if strings.Contains(prompt, c.Name) || (len([]rune(cleanCat)) >= 2 && strings.Contains(prompt, cleanCat)) {
			matchedCats = append(matchedCats, c.Name)
		}
	}

	var bookNames []string
	seenBooks := make(map[string]bool)
	for _, m := range bookTitleRegex.FindAllStringSubmatch(prompt, -1) {
		if len(m) > 1 {
			bn := strings.TrimSpace(m[1])
			if bn != "" && !seenBooks[bn] {
				seenBooks[bn] = true
				bookNames = append(bookNames, bn)
			}
		}
	}
	var bookName string
	if len(bookNames) > 0 {
		bookName = bookNames[0]
	}

	// 判断意图：1. 删除意图  2. 修改意图（含隐藏/取消隐藏）  3. 新增/文本提取意图  4. 查询/选刊意图
	wantUnhide := strings.Contains(prompt, "取消隐藏") || strings.Contains(prompt, "不隐藏") || strings.Contains(prompt, "解除隐藏") || strings.Contains(prompt, "恢复显示") || strings.Contains(prompt, "取消屏蔽")
	wantHide := !wantUnhide && (strings.Contains(prompt, "设为隐藏") || strings.Contains(prompt, "设置隐藏") || strings.Contains(prompt, "改成隐藏") || strings.Contains(prompt, "改为隐藏") || (len(bookNames) > 0 && (strings.Contains(prompt, "隐藏") || strings.Contains(prompt, "屏蔽"))))

	isDelete := strings.Contains(prompt, "删除") || strings.Contains(prompt, "删掉") || strings.Contains(prompt, "移除")
	isUpdate := wantUnhide || wantHide || strings.Contains(prompt, "修改") || strings.Contains(prompt, "改为") || strings.Contains(prompt, "改成") || strings.Contains(prompt, "更新") || strings.Contains(prompt, "调价") || strings.Contains(prompt, "调到")
	isCreate := strings.Contains(prompt, "新增") || strings.Contains(prompt, "录入") || strings.Contains(prompt, "添加") || strings.Contains(prompt, "创建") || (strings.Contains(prompt, "主管单位") && strings.Contains(prompt, "主办单位"))

	var step ai.ToolStep
	var card *ai.ActionCardData

	if isDelete && bookName != "" {
		step, card, _ = i.executeAgentTool(ctx, "prepare_delete_periodical", map[string]interface{}{
			"target_name": bookName,
		}, req.UserPermissions, allCats)
	} else if isUpdate && len(bookNames) > 0 {
		args := map[string]interface{}{
			"target_name":  bookName,
			"target_names": toInterfaceSlice(bookNames),
		}
		if wantHide {
			args["is_hidden"] = true
		} else if wantUnhide {
			args["is_hidden"] = false
		}
		if m := priceRegex.FindStringSubmatch(prompt); len(m) > 1 {
			args["price"] = m[1]
		}
		if m := pageRegex.FindStringSubmatch(prompt); len(m) > 1 {
			args["periodical_page"] = m[1]
		}
		if m := monthRegex.FindStringSubmatch(prompt); len(m) > 2 {
			mon, _ := strconv.Atoi(m[2])
			day := "01"
			if strings.Contains(prompt, "中旬") {
				day = "15"
			} else if strings.Contains(prompt, "下旬") {
				day = "21"
			} else if strings.Contains(prompt, "月底") {
				day = "28"
			}
			args["publication_time"] = fmt.Sprintf("%s-%02d-%s", m[1], mon, day)
		} else if sm := shortMonthRe.FindStringSubmatch(prompt); len(sm) > 1 {
			args["publication_time"] = sm[1]
		}
		step, card, _ = i.executeAgentTool(ctx, "prepare_update_periodical", args, req.UserPermissions, allCats)
	} else if isCreate {
		if bookName == "" {
			bookName = "新录入期刊"
		}
		args := map[string]interface{}{
			"name":                   "《" + bookName + "》",
			"category_names":         toInterfaceSlice(matchedCats),
			"is_warp":                strings.Contains(prompt, "全包"),
			"invoice_receipt":        strings.Contains(prompt, "发票"),
			"works":                  strings.Contains(prompt, "作品"),
			"periodical_description": prompt,
		}
		if m := priceRegex.FindStringSubmatch(prompt); len(m) > 1 {
			args["price"] = m[1]
		}
		if m := pageRegex.FindStringSubmatch(prompt); len(m) > 1 {
			args["periodical_page"] = m[1]
		}
		if m := monthRegex.FindStringSubmatch(prompt); len(m) > 2 {
			mon, _ := strconv.Atoi(m[2])
			args["publication_time"] = fmt.Sprintf("%s-%02d-01", m[1], mon)
		}
		if m := cnNumberRegex.FindStringSubmatch(prompt); len(m) > 1 {
			args["domestic_number"] = m[1]
		}
		if m := issnRegex.FindStringSubmatch(prompt); len(m) > 1 {
			args["international_number"] = m[1]
		}
		if comp := extractKeyVal(prompt, "主管单位"); comp != "" {
			args["competent_unit"] = comp
		}
		if org := extractKeyVal(prompt, "主办单位"); org != "" {
			args["organization"] = org
		}
		if col := extractKeyVal(prompt, "栏目设置"); col != "" {
			args["column_setting"] = col
		}
		step, card, _ = i.executeAgentTool(ctx, "prepare_create_periodical", args, req.UserPermissions, allCats)
	} else {
		// 默认执行多维智能查询
		args := map[string]interface{}{
			"category_names": toInterfaceSlice(matchedCats),
		}
		if bookName != "" {
			args["name"] = bookName
		}
		if m := monthRegex.FindStringSubmatch(prompt); len(m) > 2 {
			mon, _ := strconv.Atoi(m[2])
			args["publish_month"] = fmt.Sprintf("%s-%02d", m[1], mon)
		}
		if strings.Contains(prompt, "全包") {
			args["is_warp"] = 1
		}
		if strings.Contains(prompt, "发票") {
			args["invoice"] = 1
		}
		if strings.Contains(prompt, "作品") {
			args["works"] = 1
		}
		if strings.Contains(prompt, "已隐藏") || strings.Contains(prompt, "隐藏的期刊") || strings.Contains(prompt, "隐藏期刊") || strings.Contains(prompt, "哪些隐藏") {
			args["hidden_status"] = "hidden_only"
		}
		if strings.Contains(prompt, "便宜") || strings.Contains(prompt, "低到高") {
			args["order_by"] = "price"
			args["order_desc"] = "ascending"
		}
		step, card, _ = i.executeAgentTool(ctx, "query_periodicals", args, req.UserPermissions, allCats)
	}

	reply := i.buildSummaryFromCard(card, []ai.ToolStep{step})
	if step.Status == "denied" {
		reply = "⚠️ **" + step.Summary + "**\n\n您当前登录的账号没有执行该写操作的权限。如需操作，请联系超级管理员在【角色与权限管理】中为您授权。"
	}
	if !isFallbackFromError {
		reply += "\n\n> 💡 *提示：当前使用的是系统内置的【本地智能语义解析引擎】。超级管理员可随时前往 **【系统管理 -> AI 模型与审计管理】** 填入通义千问 (Qwen)、DeepSeek 或 OpenAI 的 API Key，开启完整的大模型深度对话与推理能力。*"
	} else if llmErrMsg != "" {
		reply += fmt.Sprintf("\n\n> ⚠️ *远端大模型调用异常（%s），已自动为您无缝切换至本地智能引擎完成处理。*", llmErrMsg)
	}
	return reply, []ai.ToolStep{step}, card
}

func toInterfaceSlice(strs []string) []interface{} {
	res := make([]interface{}, 0, len(strs))
	for _, s := range strs {
		res = append(res, s)
	}
	return res
}

func extractKeyVal(text, key string) string {
	re := regexp.MustCompile(key + `\s*[:：为是]\s*([^\s，,。；;\n]+)`)
	if m := re.FindStringSubmatch(text); len(m) > 1 {
		return strings.TrimSpace(m[1])
	}
	return ""
}

// ==================== 二次确认落库执行器 (ConfirmAction) ====================

func (i *AIImpl) ConfirmAction(ctx context.Context, req *ai.ConfirmActionRequest) (*ai.ConfirmActionResponse, error) {
	var msg ai.AIChatMessage
	if err := i.db.WithContext(ctx).First(&msg, req.MessageID).Error; err != nil {
		return nil, errors.New("未找到对应的待确认消息")
	}
	if msg.UserID != req.UserID {
		return nil, errors.New("无权操作他人的确认卡片")
	}
	if msg.ActionCard == nil {
		return nil, errors.New("该消息不包含可执行的操作卡片")
	}
	if msg.ActionCard.Status != "pending" {
		return nil, fmt.Errorf("该操作卡片状态已变更为 %s，请勿重复操作", msg.ActionCard.Status)
	}

	now := time.Now().Unix()

	// 若用户点击取消
	if req.Decision == "cancel" {
		msg.ActionCard.Status = "cancelled"
		msg.ActionCard.ExecutedAt = now
		msg.ActionCard.ResultMessage = "已取消操作"
		msg.UpdatedAt = now
		if err := i.db.WithContext(ctx).Save(&msg).Error; err != nil {
			return nil, err
		}
		return &ai.ConfirmActionResponse{UpdatedMessage: &msg}, nil
	}

	// 允许前端传入卡片内行内微调后的 FormPayload 与增删后的 Diffs
	payload := msg.ActionCard.FormPayload
	if len(req.FormPayload) > 0 {
		for k, v := range req.FormPayload {
			if payload == nil {
				payload = make(map[string]interface{})
			}
			payload[k] = v
		}
		msg.ActionCard.FormPayload = payload
	}
	if len(req.Diffs) > 0 {
		msg.ActionCard.Diffs = req.Diffs
	} else if len(req.FormPayload) > 0 {
		// 同步更新 Diffs 中的展示新值
		for idx, d := range msg.ActionCard.Diffs {
			if newVal, ok := payload[d.Field]; ok {
				if b, isBool := newVal.(bool); isBool {
					msg.ActionCard.Diffs[idx].NewValue = formatBoolCN(b)
				} else {
					msg.ActionCard.Diffs[idx].NewValue = fmt.Sprintf("%v", newVal)
				}
			}
		}
	}

	sanitizeCell := func(s string, maxRunes int) string {
		clean := strings.ReplaceAll(s, "\r\n", " ")
		clean = strings.ReplaceAll(clean, "\n", " ")
		clean = strings.ReplaceAll(clean, "|", "｜")
		clean = strings.TrimSpace(clean)
		if clean == "" {
			return "(空)"
		}
		runes := []rune(clean)
		if maxRunes > 0 && len(runes) > maxRunes {
			return string(runes[:maxRunes]) + "..."
		}
		return clean
	}

	timeStr := time.Unix(now, 0).Format("2006-01-02 15:04:05")

	switch msg.ActionCard.CardType {
	case "create_confirm":
		if !hasPerm(req.UserPermissions, "periodical:create") {
			return nil, errors.New("权限不足：缺少新增期刊权限 (periodical:create)")
		}
		createReq := buildCreateRequestFromMap(payload)
		created, err := i.periodicalSvc.CreatePeriodical(ctx, createReq)
		if err != nil {
			return nil, fmt.Errorf("创建期刊失败: %v", err)
		}
		cleanName := "《" + strings.Trim(created.Name, "《》") + "》"
		fullCreated, _ := i.periodicalSvc.DescribePeriodical(ctx, periodical.NewDescribePeriodicalRequest(strconv.Itoa(created.PeriodicalID)))
		var briefItems []ai.PeriodicalBriefCard
		if fullCreated != nil {
			briefItems = []ai.PeriodicalBriefCard{formatPeriodicalToBrief(fullCreated)}
		}

		// 原地更新当前会话消息与卡片（一体化展示，不再额外分裂出第二条 AI 消息）
		msg.Content = fmt.Sprintf("✅ 新刊 **%s**（数据库分配 ID：`#%d`）已于 `%s` 成功写入数据库生效（操作人：`%s`）：", cleanName, created.PeriodicalID, timeStr, req.Username)
		msg.ActionCard.Status = "confirmed"
		msg.ActionCard.Title = fmt.Sprintf("新刊入库已生效 · %s", cleanName)
		msg.ActionCard.Description = ""
		msg.ActionCard.PeriodicalID = created.PeriodicalID
		msg.ActionCard.ExecutedAt = now
		msg.ActionCard.Items = briefItems
		msg.ActionCard.ResultMessage = fmt.Sprintf("已于 %s 完成入库持久化 (ID: #%d)", timeStr, created.PeriodicalID)
		msg.ActionCard.Suggestions = []string{
			fmt.Sprintf("继续修改%s的栏目设置或注意事项", cleanName),
			"查询库中所有同级别期刊",
		}

	case "update_confirm":
		if !hasPerm(req.UserPermissions, "periodical:update") {
			return nil, errors.New("权限不足：缺少编辑期刊权限 (periodical:update)")
		}

		var allCats []*periodical.Categorize
		if catRes, err := i.periodicalSvc.QueryCategorize(ctx, periodical.NewQueryCategorizeRequest()); err == nil && catRes != nil {
			allCats = catRes.Items
		}

		// 判断是否为多本批量修改（或者前端传了 batch_updates）
		batchTargets := req.BatchUpdates
		if len(batchTargets) == 0 && len(msg.ActionCard.BatchUpdates) > 1 {
			batchTargets = msg.ActionCard.BatchUpdates
		}

		if len(batchTargets) > 1 {
			var briefItems []ai.PeriodicalBriefCard
			var updatedNames []string
			var boldNames []string
			totalDiffs := 0

			for idx := range batchTargets {
				bt := &batchTargets[idx]
				for dIdx, d := range bt.Diffs {
					if bt.FormPayload == nil {
						bt.FormPayload = make(map[string]interface{})
					}
					if newVal, ok := bt.FormPayload[d.Field]; ok {
						if b, isBool := newVal.(bool); isBool {
							bt.Diffs[dIdx].NewValue = formatBoolCN(b)
						} else {
							bt.Diffs[dIdx].NewValue = fmt.Sprintf("%v", newVal)
						}
					} else if d.NewValue != "" {
						bt.FormPayload[d.Field] = d.NewValue
					}
				}
				syncCategoryFieldsToCategorizeIDs(bt.FormPayload, allCats)
				pid := strconv.Itoa(bt.PeriodicalID)
				updateReq := periodical.NewUpdatePeriodicalRequest(pid)
				updateReq.CreatePeriodicalRequest = buildCreateRequestFromMap(bt.FormPayload)
				updated, err := i.periodicalSvc.UpdatePeriodical(ctx, updateReq)
				if err != nil {
					return nil, fmt.Errorf("更新期刊 %s 失败: %v", bt.PeriodicalName, err)
				}
				updatedNames = append(updatedNames, updated.Name)
				boldNames = append(boldNames, "**"+updated.Name+"**")
				totalDiffs += len(bt.Diffs)

				fullUpdated, _ := i.periodicalSvc.DescribePeriodical(ctx, periodical.NewDescribePeriodicalRequest(pid))
				if fullUpdated != nil {
					briefItems = append(briefItems, formatPeriodicalToBrief(fullUpdated))
				}
			}

			msg.Content = fmt.Sprintf("✅ %s 共 **%d 本**期刊（合计 **%d 项**属性变更）已于 `%s` 批量写入数据库生效（操作人：`%s`）：", strings.Join(boldNames, "、"), len(batchTargets), totalDiffs, timeStr, req.Username)
			msg.ActionCard.Status = "confirmed"
			msg.ActionCard.Title = fmt.Sprintf("批量修改已生效 · 共 %d 本期刊", len(batchTargets))
			msg.ActionCard.Description = ""
			msg.ActionCard.ExecutedAt = now
			msg.ActionCard.BatchUpdates = batchTargets
			msg.ActionCard.Items = briefItems
			msg.ActionCard.Suggestions = []string{
				fmt.Sprintf("查询 %s 的最新详情", strings.Join(updatedNames, "、")),
				"继续查询其他需要调整排期的期刊",
			}
			msg.ActionCard.ResultMessage = fmt.Sprintf("已于 %s 批量持久化写入数据库（共更新 %d 本期刊 · %d 项变更）", timeStr, len(batchTargets), totalDiffs)
		} else {
			if len(batchTargets) == 1 {
				payload = batchTargets[0].FormPayload
				msg.ActionCard.FormPayload = payload
				msg.ActionCard.Diffs = batchTargets[0].Diffs
				msg.ActionCard.PeriodicalID = batchTargets[0].PeriodicalID
				msg.ActionCard.PeriodicalName = batchTargets[0].PeriodicalName
			}
			for _, d := range msg.ActionCard.Diffs {
				if payload == nil {
					payload = make(map[string]interface{})
				}
				if _, exists := payload[d.Field]; !exists && d.NewValue != "" {
					payload[d.Field] = d.NewValue
				}
			}
			syncCategoryFieldsToCategorizeIDs(payload, allCats)
			pid := strconv.Itoa(msg.ActionCard.PeriodicalID)
			updateReq := periodical.NewUpdatePeriodicalRequest(pid)
			updateReq.CreatePeriodicalRequest = buildCreateRequestFromMap(payload)
			updated, err := i.periodicalSvc.UpdatePeriodical(ctx, updateReq)
			if err != nil {
				return nil, fmt.Errorf("更新期刊失败: %v", err)
			}

			var rollbackParts []string
			for _, d := range msg.ActionCard.Diffs {
				if d.Field == "is_hidden" {
					if d.OldValue == "否" {
						rollbackParts = append(rollbackParts, "取消隐藏")
					} else if d.OldValue == "是" {
						rollbackParts = append(rollbackParts, "设为隐藏")
					}
				} else if d.OldValue != "" && d.OldValue != "(空)" && len([]rune(d.OldValue)) <= 30 {
					rollbackParts = append(rollbackParts, fmt.Sprintf("%s改回 %s", d.Label, sanitizeCell(d.OldValue, 30)))
				}
			}

			fullUpdated, _ := i.periodicalSvc.DescribePeriodical(ctx, periodical.NewDescribePeriodicalRequest(pid))
			var briefItems []ai.PeriodicalBriefCard
			if fullUpdated != nil {
				briefItems = []ai.PeriodicalBriefCard{formatPeriodicalToBrief(fullUpdated)}
			}
			var sugs []string
			if len(rollbackParts) > 0 {
				sugs = append(sugs, fmt.Sprintf("撤销刚才的修改：把%s的%s", updated.Name, strings.Join(rollbackParts, "，")))
			}
			sugs = append(sugs, fmt.Sprintf("查询%s同级别的其他期刊", updated.Name))

			// 原地将当前消息与卡片演进为“已确认生效 + 最新期刊档案快照”一体化视图，不再分裂出第二条重复消息
			msg.Content = fmt.Sprintf("✅ **%s** 的 **%d 项**属性变更已于 `%s` 写入数据库生效（操作人：`%s`）：", updated.Name, len(msg.ActionCard.Diffs), timeStr, req.Username)
			msg.ActionCard.Status = "confirmed"
			msg.ActionCard.Title = fmt.Sprintf("修改已生效 · %s", updated.Name)
			msg.ActionCard.Description = ""
			msg.ActionCard.ExecutedAt = now
			msg.ActionCard.Items = briefItems
			msg.ActionCard.BatchUpdates = []ai.BatchUpdateTarget{
				{
					PeriodicalID:   msg.ActionCard.PeriodicalID,
					PeriodicalName: updated.Name,
					FormPayload:    payload,
					Diffs:          msg.ActionCard.Diffs,
				},
			}
			msg.ActionCard.Suggestions = sugs
			msg.ActionCard.ResultMessage = fmt.Sprintf("已于 %s 持久化写入数据库（共生效 %d 项变更）", timeStr, len(msg.ActionCard.Diffs))
		}

	case "delete_confirm":
		if !hasPerm(req.UserPermissions, "periodical:delete") {
			return nil, errors.New("权限不足：缺少删除期刊权限 (periodical:delete)")
		}
		pid := strconv.Itoa(msg.ActionCard.PeriodicalID)
		if err := i.periodicalSvc.DeletePeriodical(ctx, periodical.NewDeletePeriodicalRequest(pid)); err != nil {
			return nil, fmt.Errorf("删除期刊失败: %v", err)
		}
		msg.Content = fmt.Sprintf("🗑️ 期刊 **%s**（原 ID：`#%d`）已于 `%s` 从数据库永久移除（操作人：`%s`）。", msg.ActionCard.PeriodicalName, msg.ActionCard.PeriodicalID, timeStr, req.Username)
		msg.ActionCard.Status = "confirmed"
		msg.ActionCard.ExecutedAt = now
		msg.ActionCard.ResultMessage = fmt.Sprintf("已于 %s 永久移除 %s (ID: #%d)", timeStr, msg.ActionCard.PeriodicalName, msg.ActionCard.PeriodicalID)

	default:
		return nil, errors.New("不支持确认的卡片类型")
	}

	msg.UpdatedAt = now
	if err := i.db.WithContext(ctx).Save(&msg).Error; err != nil {
		return nil, err
	}

	return &ai.ConfirmActionResponse{
		UpdatedMessage: &msg,
	}, nil
}

func buildCreateRequestFromMap(m map[string]interface{}) *periodical.CreatePeriodicalRequest {
	if m != nil {
		for _, bf := range []string{"is_warp", "invoice_receipt", "works", "is_hidden"} {
			if bVal, ok := getBoolArg(m, bf); ok {
				m[bf] = bVal
			}
		}
	}
	req := periodical.NewCreatePeriodicalRequest()
	b, _ := json.Marshal(m)
	var base periodical.BasePeriodical
	_ = json.Unmarshal(b, &base)
	req.BasePeriodical = &base

	if rawIDs, ok := m["categorize_ids"].([]interface{}); ok {
		for _, idVal := range rawIDs {
			switch v := idVal.(type) {
			case float64:
				req.CategorizeIDs = append(req.CategorizeIDs, int64(v))
			case int:
				req.CategorizeIDs = append(req.CategorizeIDs, int64(v))
			case int64:
				req.CategorizeIDs = append(req.CategorizeIDs, v)
			}
		}
	} else if int64IDs, ok := m["categorize_ids"].([]int64); ok {
		req.CategorizeIDs = append(req.CategorizeIDs, int64IDs...)
	} else if intIDs, ok := m["categorize_ids"].([]int); ok {
		for _, idVal := range intIDs {
			req.CategorizeIDs = append(req.CategorizeIDs, int64(idVal))
		}
	}
	return req
}
