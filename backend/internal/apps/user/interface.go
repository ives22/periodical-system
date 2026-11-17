package user

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Service 定义 User 包的能力，就是接口定义
type Service interface {
	// CreateUser 创建用户
	CreateUser(context.Context, *CreateUserRequest) (*User, error)
	// DeleteUser 删除用户
	DeleteUser(context.Context, *DeleteUserRequest) error
	// DescribeUser 查询单个用户（含角色与权限集）
	DescribeUser(context.Context, *DescribeUserRequest) (*User, error)
	// QueryUsers 分页查询用户列表
	QueryUsers(context.Context, *QueryUserRequest) (*UserSet, error)
	// UpdateUser 修改用户资料与多角色(重置密码、角色列表、标签)
	UpdateUser(context.Context, *UpdateUserRequest) (*User, error)
	// UpdatePassword 修改个人密码
	UpdatePassword(context.Context, *UpdatePasswordRequest) error
}

// CreateUserRequest 创建用户时，用户传递的请求
type CreateUserRequest struct {
	Username string            `json:"username" gorm:"column:username"`
	Password string            `json:"password,omitempty" gorm:"column:password"`
	Role     Role              `json:"role" gorm:"column:role"`
	RoleIds  []int64           `json:"role_ids" gorm:"-"`
	Label    map[string]string `json:"label" gorm:"column:label;serializer:json"`
	isHashed bool              `gorm:"-"` // 用于校验密码是否被hash过
}

func NewCreateUserRequest() *CreateUserRequest {
	return &CreateUserRequest{
		Role:    ROLE_AUTHOR,
		RoleIds: []int64{},
		Label:   map[string]string{},
	}
}

type DeleteUserRequest struct {
	Username string `json:"username"`
}

func NewDeleteUserRequest() *DeleteUserRequest {
	return &DeleteUserRequest{}
}

type DescribeUserRequest struct {
	DescribeBy    DescribeBy `json:"describe_by"`
	DescribeValue string     `json:"describe_value"`
}

func NewDescribeUserRequestById(id string) *DescribeUserRequest {
	return &DescribeUserRequest{
		DescribeBy:    DESCRIBE_BY_ID,
		DescribeValue: id,
	}
}

func NewDescribeUserRequestByUsername(username string) *DescribeUserRequest {
	return &DescribeUserRequest{
		DescribeBy:    DESCRIBE_BY_USERNAME,
		DescribeValue: username,
	}
}

type QueryUserRequest struct {
	PageSize   int    `json:"page_size"`
	PageNumber int    `json:"page_number"`
	Username   string `json:"username"`
	RoleId     int64  `json:"role_id"`
}

func NewQueryUserRequest() *QueryUserRequest {
	return &QueryUserRequest{
		PageSize:   20,
		PageNumber: 1,
	}
}

func (r *QueryUserRequest) Offset() int {
	if r.PageNumber <= 0 {
		r.PageNumber = 1
	}
	if r.PageSize <= 0 {
		r.PageSize = 20
	}
	return (r.PageNumber - 1) * r.PageSize
}

type UserSet struct {
	Total int64   `json:"total"`
	Items []*User `json:"items"`
}

func NewUserSet() *UserSet {
	return &UserSet{
		Items: []*User{},
	}
}

type UpdatePasswordRequest struct {
	Username    string `json:"username"`
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

func NewUpdatePasswordRequest() *UpdatePasswordRequest {
	return &UpdatePasswordRequest{}
}

type UpdateUserRequest struct {
	Username    string            `json:"username"`
	NewPassword string            `json:"new_password,omitempty"` // 可选，用于重置密码
	Role        Role              `json:"role,omitempty"`         // 兼容旧单角色字段
	RoleIds     []int64           `json:"role_ids,omitempty"`     // 支持绑定多个角色ID
	Label       map[string]string `json:"label,omitempty"`        // 可选，用于更新标签
}

func NewUpdateUserRequest() *UpdateUserRequest {
	return &UpdateUserRequest{
		RoleIds: []int64{},
		Label:   map[string]string{},
	}
}

func (req *CreateUserRequest) Validate() error {
	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" || req.Password == "" {
		return errors.New("用户名或密码不能为空")
	}
	// 如果未传 RoleIds，则回退使用 Role 字段或默认普通用户角色(2)
	if len(req.RoleIds) == 0 {
		if req.Role > 0 {
			req.RoleIds = []int64{int64(req.Role)}
		} else {
			req.RoleIds = []int64{int64(ROLE_AUTHOR)}
		}
	}
	req.Role = Role(req.RoleIds[0])
	return nil
}

func (req *CreateUserRequest) SetIsHashed() {
	req.isHashed = true
}

func (req *CreateUserRequest) SetNoHashed() {
	req.isHashed = false
}

func (req *CreateUserRequest) PasswordHash() {
	if req.isHashed {
		return
	}

	b, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	req.Password = string(b)
	req.isHashed = true
}

// Validate 验证更新请求参数
func (req *UpdateUserRequest) Validate() error {
	req.Username = strings.TrimSpace(req.Username)
	if req.Username == "" {
		return fmt.Errorf("用户名不能为空")
	}
	if len(req.RoleIds) == 0 && req.Role > 0 {
		req.RoleIds = []int64{int64(req.Role)}
	}
	return nil
}
