package api

import (
	"log"
	"periodical/internal/apps/token"
	"periodical/internal/response"

	"github.com/gin-gonic/gin"
)

type TokenApiHandler struct {
	// 依赖控制器
	svc token.Service
}

func NewTokenApiHandler(tokenSvc token.Service) *TokenApiHandler {
	return &TokenApiHandler{
		svc: tokenSvc,
	}
}

func (h *TokenApiHandler) Registry(r gin.IRouter) {
	v1 := r.Group("v1")
	v1.POST("/auth/login", h.Login)
	v1.POST("/auth/logout", h.Logout)
}

func (h *TokenApiHandler) Login(c *gin.Context) {
	// 1. 获取用户的请求参数
	req := token.NewLoginRequest()
	if err := c.BindJSON(req); err != nil {
		response.Failed(c, response.CodeInvalidParam)
		return
	}

	// 2. 执行逻辑
	ins, err := h.svc.Login(c.Request.Context(), req)
	if err != nil {
		response.Failed(c, response.CodeServerBusy)
		return
	}

	c.SetCookie(token.TOKEN_COOKIE_NAME, ins.AccessToken, 0, "/", "localhost", false, true)

	// 3. 返回响应
	response.Success(c, ins)
}

func (h *TokenApiHandler) Logout(c *gin.Context) {
	req := token.NewLogoutRequest()

	if err := c.BindJSON(req); err != nil {
		response.Failed(c, response.CodeInvalidParam)
		return
	}
	log.Println("req", req)

	err := h.svc.Logout(c.Request.Context(), req)
	if err != nil {
		response.Failed(c, response.CodeServerBusy)
		return
	}

	response.Success(c, nil)
}
