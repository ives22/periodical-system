package api

import (
	"periodical/internal/apps/role"
	"periodical/internal/apps/token"
	"periodical/internal/middleware"
	"periodical/internal/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

type RoleApiHandler struct {
	svc role.Service
}

func NewRoleApiHandler(svc role.Service) *RoleApiHandler {
	return &RoleApiHandler{
		svc: svc,
	}
}

func (h *RoleApiHandler) Registry(r gin.IRouter, tokenSvc token.Service) {
	v1 := r.Group("v1")
	v1.Use(middleware.Auth(tokenSvc))

	// 获取所有启用状态的角色（供用户管理下拉框使用）
	v1.GET("/role/all", h.QueryAllEnabledRoles)

	// 角色管理接口
	roleGroup := v1.Group("/role")
	roleGroup.POST("/list", middleware.RequirePermission("role:query"), h.QueryRoles)
	roleGroup.POST("", middleware.RequirePermission("role:create"), h.CreateRole)
	roleGroup.PUT("/:id", middleware.RequirePermission("role:update"), h.UpdateRole)
	roleGroup.DELETE("/:id", middleware.RequirePermission("role:delete"), h.DeleteRole)
	roleGroup.GET("/:id/permissions", middleware.RequirePermission("role:query"), h.GetRolePermissions)
	roleGroup.PUT("/:id/permissions", middleware.RequirePermission("role:update"), h.AssignRolePermissions)

	// 权限树查询接口
	v1.GET("/permission/tree", middleware.RequirePermission("role:query"), h.QueryPermissionTree)
}

func (h *RoleApiHandler) QueryAllEnabledRoles(c *gin.Context) {
	items, err := h.svc.QueryAllEnabledRoles(c.Request.Context())
	if err != nil {
		response.FailedWhitMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, items)
}

func (h *RoleApiHandler) QueryRoles(c *gin.Context) {
	req := role.NewQueryRoleRequest()
	if err := c.ShouldBindJSON(req); err != nil {
		response.Failed(c, response.CodeInvalidParam)
		return
	}
	set, err := h.svc.QueryRoles(c.Request.Context(), req)
	if err != nil {
		response.FailedWhitMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, set)
}

func (h *RoleApiHandler) CreateRole(c *gin.Context) {
	req := role.NewCreateRoleRequest()
	if err := c.ShouldBindJSON(req); err != nil {
		response.Failed(c, response.CodeInvalidParam)
		return
	}
	ins, err := h.svc.CreateRole(c.Request.Context(), req)
	if err != nil {
		response.FailedWhitMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, ins)
}

func (h *RoleApiHandler) UpdateRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Failed(c, response.CodeInvalidParam)
		return
	}

	req := role.NewUpdateRoleRequest(id)
	if err := c.ShouldBindJSON(req); err != nil {
		response.Failed(c, response.CodeInvalidParam)
		return
	}
	req.Id = id

	ins, err := h.svc.UpdateRole(c.Request.Context(), req)
	if err != nil {
		response.FailedWhitMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, ins)
}

func (h *RoleApiHandler) DeleteRole(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Failed(c, response.CodeInvalidParam)
		return
	}

	if err := h.svc.DeleteRole(c.Request.Context(), id); err != nil {
		response.FailedWhitMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *RoleApiHandler) GetRolePermissions(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Failed(c, response.CodeInvalidParam)
		return
	}

	ids, err := h.svc.GetRolePermissionIds(c.Request.Context(), id)
	if err != nil {
		response.FailedWhitMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, ids)
}

func (h *RoleApiHandler) AssignRolePermissions(c *gin.Context) {
	id, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil || id <= 0 {
		response.Failed(c, response.CodeInvalidParam)
		return
	}

	req := role.NewAssignPermissionsRequest(id)
	if err := c.ShouldBindJSON(req); err != nil {
		response.Failed(c, response.CodeInvalidParam)
		return
	}
	req.RoleId = id

	if err := h.svc.AssignRolePermissions(c.Request.Context(), req); err != nil {
		response.FailedWhitMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, nil)
}

func (h *RoleApiHandler) QueryPermissionTree(c *gin.Context) {
	tree, err := h.svc.QueryPermissionTree(c.Request.Context())
	if err != nil {
		response.FailedWhitMsg(c, response.CodeServerBusy, err.Error())
		return
	}
	response.Success(c, tree)
}
