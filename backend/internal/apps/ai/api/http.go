package api

import (
	"periodical/internal/apps/ai"
	"periodical/internal/apps/token"
	"periodical/internal/middleware"
	"periodical/internal/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AIApiHandler struct {
	svc ai.Service
}

func NewAIApiHandler(svc ai.Service) *AIApiHandler {
	return &AIApiHandler{
		svc: svc,
	}
}

func (h *AIApiHandler) Registry(r *gin.Engine, tokenSvc token.Service) {
	v1 := r.Group("v1/ai")
	v1.Use(middleware.Auth(tokenSvc))

	// ===== 个人会话与 AI 对话接口（所有登录用户可用，内部工具执行校验具体期刊RBAC权限） =====
	v1.GET("/active-model", h.GetActiveDefaultModel)
	v1.GET("/sessions", h.QueryUserSessions)
	v1.POST("/sessions", h.CreateSession)
	v1.PUT("/sessions/:session_id", h.RenameSession)
	v1.DELETE("/sessions/:session_id", h.DeleteSession)
	v1.GET("/sessions/:session_id/messages", h.QuerySessionMessages)
	v1.POST("/chat", h.SendMessage)
	v1.POST("/action/confirm", h.ConfirmAction)

	// ===== 超级管理员：大模型渠道配置与全员对话审计 =====
	v1.GET("/models", middleware.RequirePermission("ai:config"), h.QueryModelConfigs)
	v1.POST("/models", middleware.RequirePermission("ai:config"), h.SaveModelConfig)
	v1.POST("/models/remote-list", middleware.RequirePermission("ai:config"), h.FetchRemoteModels)
	v1.PUT("/models/:id/default", middleware.RequirePermission("ai:config"), h.SetDefaultModel)
	v1.POST("/models/:id/test", middleware.RequirePermission("ai:config"), h.TestModelConnection)
	v1.DELETE("/models/:id", middleware.RequirePermission("ai:config"), h.DeleteModelConfig)
	v1.POST("/audit/sessions", middleware.RequirePermission("ai:audit"), h.QueryAuditSessions)
}

func (h *AIApiHandler) GetActiveDefaultModel(c *gin.Context) {
	cfg, err := h.svc.GetActiveDefaultModel(c.Request.Context())
	if err != nil {
		response.Success(c, gin.H{
			"name":         "内置本地智能引擎 (未配置外部模型)",
			"provider":     "local",
			"model_name":   "smart-agent-local",
			"has_api_key":  false,
			"is_connected": false,
		})
		return
	}
	response.Success(c, cfg)
}

func (h *AIApiHandler) QueryUserSessions(c *gin.Context) {
	tk := middleware.GetTokenFromCtx(c)
	keyword := c.Query("keyword")
	list, err := h.svc.QueryUserSessions(c.Request.Context(), tk.UserId, keyword)
	if err != nil {
		response.FailedWhitMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, list)
}

func (h *AIApiHandler) CreateSession(c *gin.Context) {
	tk := middleware.GetTokenFromCtx(c)
	var body struct {
		Title string `json:"title"`
		Scene string `json:"scene"`
	}
	_ = c.ShouldBindJSON(&body)
	session, err := h.svc.CreateSession(c.Request.Context(), tk.UserId, tk.UserName, body.Title, body.Scene)
	if err != nil {
		response.FailedWhitMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, session)
}

func (h *AIApiHandler) RenameSession(c *gin.Context) {
	tk := middleware.GetTokenFromCtx(c)
	sessionID := c.Param("session_id")
	var body struct {
		Title    string `json:"title"`
		IsPinned *int   `json:"is_pinned"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		response.Failed(c, response.CodeInvalidParam)
		return
	}
	if err := h.svc.RenameSession(c.Request.Context(), tk.UserId, sessionID, body.Title, body.IsPinned); err != nil {
		response.FailedWhitMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, gin.H{"session_id": sessionID})
}

func (h *AIApiHandler) DeleteSession(c *gin.Context) {
	tk := middleware.GetTokenFromCtx(c)
	sessionID := c.Param("session_id")
	if err := h.svc.DeleteSession(c.Request.Context(), tk.UserId, sessionID); err != nil {
		response.FailedWhitMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, gin.H{"session_id": sessionID})
}

func (h *AIApiHandler) QuerySessionMessages(c *gin.Context) {
	tk := middleware.GetTokenFromCtx(c)
	sessionID := c.Param("session_id")
	isSuperAdmin := tk.HasPermission("ai:audit") || tk.HasPermission("*")
	msgs, err := h.svc.QuerySessionMessages(c.Request.Context(), tk.UserId, isSuperAdmin, sessionID)
	if err != nil {
		response.FailedWhitMsg(c, response.CodeForbidden, err.Error())
		return
	}
	response.Success(c, msgs)
}

func (h *AIApiHandler) SendMessage(c *gin.Context) {
	tk := middleware.GetTokenFromCtx(c)
	req := &ai.SendMessageRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		response.Failed(c, response.CodeInvalidParam)
		return
	}
	req.UserID = tk.UserId
	req.Username = tk.UserName
	req.UserPermissions = tk.Permissions

	res, err := h.svc.SendMessage(c.Request.Context(), req)
	if err != nil {
		response.FailedWhitMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, res)
}

func (h *AIApiHandler) ConfirmAction(c *gin.Context) {
	tk := middleware.GetTokenFromCtx(c)
	req := &ai.ConfirmActionRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		response.Failed(c, response.CodeInvalidParam)
		return
	}
	req.UserID = tk.UserId
	req.Username = tk.UserName
	req.UserPermissions = tk.Permissions

	res, err := h.svc.ConfirmAction(c.Request.Context(), req)
	if err != nil {
		response.FailedWhitMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, res)
}

// ===== 超级管理员接口 =====

func (h *AIApiHandler) QueryModelConfigs(c *gin.Context) {
	list, err := h.svc.QueryModelConfigs(c.Request.Context())
	if err != nil {
		response.FailedWhitMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, list)
}

func (h *AIApiHandler) SaveModelConfig(c *gin.Context) {
	req := &ai.SaveModelConfigRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		response.Failed(c, response.CodeInvalidParam)
		return
	}
	cfg, err := h.svc.SaveModelConfig(c.Request.Context(), req)
	if err != nil {
		response.FailedWhitMsg(c, response.CodeInvalidParam, err.Error())
		return
	}
	response.Success(c, cfg)
}

func (h *AIApiHandler) SetDefaultModel(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Failed(c, response.CodeInvalidParam)
		return
	}
	if err := h.svc.SetDefaultModel(c.Request.Context(), id); err != nil {
		response.FailedWhitMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, gin.H{"id": id})
}

func (h *AIApiHandler) TestModelConnection(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Failed(c, response.CodeInvalidParam)
		return
	}
	res, err := h.svc.TestModelConnection(c.Request.Context(), id)
	if err != nil {
		response.FailedWhitMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, res)
}

func (h *AIApiHandler) FetchRemoteModels(c *gin.Context) {
	req := &ai.FetchRemoteModelsRequest{}
	if err := c.ShouldBindJSON(req); err != nil {
		response.Failed(c, response.CodeInvalidParam)
		return
	}
	res, err := h.svc.FetchRemoteModels(c.Request.Context(), req)
	if err != nil {
		response.FailedWhitMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, res)
}

func (h *AIApiHandler) DeleteModelConfig(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Failed(c, response.CodeInvalidParam)
		return
	}
	if err := h.svc.DeleteModelConfig(c.Request.Context(), id); err != nil {
		response.FailedWhitMsg(c, response.CodeInvalidParam, err.Error())
		return
	}
	response.Success(c, gin.H{"id": id})
}

func (h *AIApiHandler) QueryAuditSessions(c *gin.Context) {
	req := &ai.QueryAuditSessionsRequest{}
	_ = c.ShouldBindJSON(req)
	res, err := h.svc.QueryAuditSessions(c.Request.Context(), req)
	if err != nil {
		response.FailedWhitMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, res)
}
