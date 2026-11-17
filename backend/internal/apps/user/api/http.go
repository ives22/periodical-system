package api

import (
	"fmt"
	"periodical/internal/apps/role"
	"periodical/internal/apps/token"
	"periodical/internal/apps/user"
	"periodical/internal/middleware"
	"periodical/internal/response"

	"github.com/gin-gonic/gin"
)

type UserApiHandler struct {
	svc     user.Service
	roleSvc role.Service
}

func NewUserApiHandler(svc user.Service, roleSvc ...role.Service) *UserApiHandler {
	var rSvc role.Service
	if len(roleSvc) > 0 {
		rSvc = roleSvc[0]
	}
	return &UserApiHandler{
		svc:     svc,
		roleSvc: rSvc,
	}
}

func (h *UserApiHandler) Registry(r gin.IRouter, tokenSvc token.Service) {
	v1 := r.Group("v1")
	v1.Use(middleware.Auth(tokenSvc))

	// 个人中心接口（所有已登录用户可用）
	v1.GET("/user/profile", h.GetProfile)
	v1.POST("/user/update", h.UpdatePassword)

	// 用户管理接口（按权限码控制）
	v1.POST("/user/list", middleware.RequirePermission("user:query"), h.QueryUsers)
	v1.POST("/user", middleware.RequirePermission("user:create"), h.CreateUser)
	v1.PUT("/user", middleware.RequirePermission("user:update"), h.UpdateUser)
	v1.POST("/user/delete", middleware.RequirePermission("user:delete"), h.DeleteUser)
}

// containsSuperAdminRole 检查角色ID列表中是否包含内置超级管理员角色
func (h *UserApiHandler) containsSuperAdminRole(c *gin.Context, roleIds []int64) bool {
	if h.roleSvc == nil {
		for _, id := range roleIds {
			if id == 3 {
				return true
			}
		}
		return false
	}
	for _, id := range roleIds {
		r, err := h.roleSvc.DescribeRole(c.Request.Context(), id)
		if err == nil && r.IsSuperAdmin() {
			return true
		}
	}
	return false
}

// GetProfile 获取当前登录用户信息、角色列表及权限码列表
func (h *UserApiHandler) GetProfile(c *gin.Context) {
	tk := middleware.GetTokenFromCtx(c)
	if tk == nil {
		response.Failed(c, response.CodeUnauthorized)
		return
	}

	u, err := h.svc.DescribeUser(c.Request.Context(), user.NewDescribeUserRequestById(fmt.Sprintf("%d", tk.UserId)))
	if err != nil {
		response.FailedWhitMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	u.Desensitize()
	response.Success(c, u)
}

// QueryUsers 分页查询用户列表
func (h *UserApiHandler) QueryUsers(c *gin.Context) {
	req := user.NewQueryUserRequest()
	if err := c.ShouldBindJSON(req); err != nil {
		response.Failed(c, response.CodeInvalidParam)
		return
	}

	set, err := h.svc.QueryUsers(c.Request.Context(), req)
	if err != nil {
		response.FailedWhitMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, set)
}

func (h *UserApiHandler) CreateUser(c *gin.Context) {
	req := user.NewCreateUserRequest()

	if err := c.ShouldBindJSON(req); err != nil {
		response.Failed(c, response.CodeInvalidParam)
		return
	}

	tk := middleware.GetTokenFromCtx(c)
	// 非超级管理员不可为新用户分配超级管理员角色
	if tk != nil && !tk.IsSuperAdmin && h.containsSuperAdminRole(c, req.RoleIds) {
		response.FailedWhitMsg(c, response.CodeForbidden, "非超级管理员不可分配超级管理员角色")
		return
	}

	u, err := h.svc.CreateUser(c.Request.Context(), req)
	if err != nil {
		response.FailedWhitMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	u.Desensitize()

	response.Success(c, u)
}

// UpdatePassword 修改个人密码（强制校验只能修改当前登录人自己的密码）
func (h *UserApiHandler) UpdatePassword(c *gin.Context) {
	req := user.NewUpdatePasswordRequest()

	if err := c.ShouldBindJSON(req); err != nil {
		response.Failed(c, response.CodeInvalidParam)
		return
	}

	tk := middleware.GetTokenFromCtx(c)
	if tk != nil {
		req.Username = tk.UserName
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

	tk := middleware.GetTokenFromCtx(c)
	if tk != nil && req.Username == tk.UserName {
		response.FailedWhitMsg(c, response.CodeInvalidParam, "不能删除当前正在登录的自身账号")
		return
	}

	targetUser, err := h.svc.DescribeUser(c.Request.Context(), user.NewDescribeUserRequestByUsername(req.Username))
	if err != nil {
		response.FailedWhitMsg(c, response.CodeInvalidParam, "目标用户不存在")
		return
	}

	if targetUser.IsSuperAdmin && (tk == nil || !tk.IsSuperAdmin) {
		response.FailedWhitMsg(c, response.CodeForbidden, "无权删除超级管理员账号")
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

	if err := req.Validate(); err != nil {
		response.FailedWhitMsg(c, response.CodeInvalidParam, err.Error())
		return
	}

	targetUser, err := h.svc.DescribeUser(c.Request.Context(), user.NewDescribeUserRequestByUsername(req.Username))
	if err != nil {
		response.FailedWhitMsg(c, response.CodeInvalidParam, "目标用户不存在")
		return
	}

	tk := middleware.GetTokenFromCtx(c)
	if tk != nil && !tk.IsSuperAdmin {
		if targetUser.IsSuperAdmin {
			response.FailedWhitMsg(c, response.CodeForbidden, "非超级管理员无权修改超级管理员账号")
			return
		}
		if h.containsSuperAdminRole(c, req.RoleIds) {
			response.FailedWhitMsg(c, response.CodeForbidden, "非超级管理员不可分配超级管理员角色")
			return
		}
	}

	u, err := h.svc.UpdateUser(c.Request.Context(), req)
	if err != nil {
		response.FailedWhitMsg(c, response.CodeServerBusy, err.Error())
		return
	}

	response.Success(c, u)
}
