package role

import (
	"context"
	"errors"
	"strings"
)

type Service interface {
	// 角色管理
	CreateRole(ctx context.Context, req *CreateRoleRequest) (*Role, error)
	UpdateRole(ctx context.Context, req *UpdateRoleRequest) (*Role, error)
	DeleteRole(ctx context.Context, id int64) error
	DescribeRole(ctx context.Context, id int64) (*Role, error)
	QueryRoles(ctx context.Context, req *QueryRoleRequest) (*RoleSet, error)
	QueryAllEnabledRoles(ctx context.Context) ([]*Role, error)

	// 权限管理与角色授权
	QueryPermissionTree(ctx context.Context) ([]*Permission, error)
	GetRolePermissionIds(ctx context.Context, roleId int64) ([]int64, error)
	AssignRolePermissions(ctx context.Context, req *AssignPermissionsRequest) error

	// 用户多角色关联与权限并集计算
	GetUserRolesAndPermissions(ctx context.Context, userId int64) (roles []*Role, permCodes []string, isSuperAdmin bool, err error)
	SetUserRoles(ctx context.Context, userId int64, roleIds []int64) error
}

type CreateRoleRequest struct {
	Name          string  `json:"name"`
	Code          string  `json:"code"`
	Description   string  `json:"description"`
	Status        int64   `json:"status"`
	Sort          int64   `json:"sort"`
	PermissionIds []int64 `json:"permission_ids"`
}

func NewCreateRoleRequest() *CreateRoleRequest {
	return &CreateRoleRequest{
		Status:        1,
		Sort:          10,
		PermissionIds: []int64{},
	}
}

func (r *CreateRoleRequest) Validate() error {
	r.Name = strings.TrimSpace(r.Name)
	r.Code = strings.TrimSpace(r.Code)
	if r.Name == "" {
		return errors.New("角色名称不能为空")
	}
	if r.Code == "" {
		return errors.New("角色编码不能为空")
	}
	return nil
}

type UpdateRoleRequest struct {
	Id          int64  `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	Status      *int64 `json:"status"`
	Sort        *int64 `json:"sort"`
}

func NewUpdateRoleRequest(id int64) *UpdateRoleRequest {
	return &UpdateRoleRequest{
		Id: id,
	}
}

func (r *UpdateRoleRequest) Validate() error {
	if r.Id <= 0 {
		return errors.New("无效的角色ID")
	}
	r.Name = strings.TrimSpace(r.Name)
	if r.Name == "" {
		return errors.New("角色名称不能为空")
	}
	return nil
}

type QueryRoleRequest struct {
	PageSize   int    `json:"page_size"`
	PageNumber int    `json:"page_number"`
	Keyword    string `json:"keyword"`
	Status     *int64 `json:"status"`
}

func NewQueryRoleRequest() *QueryRoleRequest {
	return &QueryRoleRequest{
		PageSize:   20,
		PageNumber: 1,
	}
}

func (r *QueryRoleRequest) Offset() int {
	if r.PageNumber <= 0 {
		r.PageNumber = 1
	}
	if r.PageSize <= 0 {
		r.PageSize = 20
	}
	return (r.PageNumber - 1) * r.PageSize
}

type RoleSet struct {
	Total int64   `json:"total"`
	Items []*Role `json:"items"`
}

func NewRoleSet() *RoleSet {
	return &RoleSet{
		Items: []*Role{},
	}
}

type AssignPermissionsRequest struct {
	RoleId        int64   `json:"role_id"`
	PermissionIds []int64 `json:"permission_ids"`
}

func NewAssignPermissionsRequest(roleId int64) *AssignPermissionsRequest {
	return &AssignPermissionsRequest{
		RoleId:        roleId,
		PermissionIds: []int64{},
	}
}
