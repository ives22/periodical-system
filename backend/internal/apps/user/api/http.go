package api

import (
	"github.com/gin-gonic/gin"
	"periodical/internal/apps/user"
	"periodical/internal/response"
)

type UserApiHandler struct {
	svc user.Service
}

func NewUserApiHandler(svc user.Service) *UserApiHandler {
	return &UserApiHandler{
		svc: svc,
	}
}

func (h *UserApiHandler) Registry(r gin.IRouter) {
	v1 := r.Group("v1")
	v1.POST("/user", h.CreateUser)
	v1.POST("/user/update", h.UpdatePassword)
	v1.POST("/user/delete", h.DeleteUser)
	v1.PUT("/user", h.UpdateUser)
}

func (h *UserApiHandler) CreateUser(c *gin.Context) {
	req := user.NewCreateUserRequest()

	if err := c.ShouldBindJSON(req); err != nil {
		response.Failed(c, response.CodeInvalidParam)
		return
	}

	u, err := h.svc.CreateUser(c.Request.Context(), req)
	if err != nil {
		response.FailedWhitMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	u.Password = ""

	response.Success(c, u)
}

func (h *UserApiHandler) UpdatePassword(c *gin.Context) {
	req := user.NewUpdatePasswordRequest()

	if err := c.ShouldBindJSON(req); err != nil {
		response.Failed(c, response.CodeInvalidParam)
		return
	}

	err := h.svc.UpdatePassword(c.Request.Context(), req)
	if err != nil {
		response.FailedWhitMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *UserApiHandler) DeleteUser(c *gin.Context) {
	req := user.NewDeleteUserRequest()

	if err := c.ShouldBindJSON(req); err != nil {
		response.Failed(c, response.CodeInvalidParam)
		return
	}

	if err := h.svc.DeleteUser(c.Request.Context(), req); err != nil {
		response.FailedWhitMsg(c, response.CodeInvalidParam, err.Error())
		return
	}

	response.Success(c, nil)
}

func (h *UserApiHandler) UpdateUser(c *gin.Context) {
	req := user.NewUpdateUserRequest()

	if err := c.ShouldBindJSON(req); err != nil {
		response.Failed(c, response.CodeInvalidParam)
		return
	}

	// 验证请求参数
	if err := req.Validate(); err != nil {
		response.FailedWhitMsg(c, response.CodeInvalidParam, err.Error())
		return
	}

	u, err := h.svc.UpdateUser(c.Request.Context(), req)
	if err != nil {
		response.FailedWhitMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, u)
}
