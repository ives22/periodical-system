package user

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/crypto/bcrypt"
)

// Service 定义 User 包的能力，就是接口定义
type Service interface {
	// CreateUser 创建用户
	CreateUser(context.Context, *CreateUserRequest) (*User, error)
	// DeleteUser 删除用户
	DeleteUser(context.Context, *DeleteUserRequest) error
	// DescribeUser 查询用户
	DescribeUser(context.Context, *DescribeUserRequest) (*User, error)
	// UpdateUser 修改用户资料角色(重置密码、角色)
	UpdateUser(context.Context, *UpdateUserRequest) (*User, error)

	// UpdatePassword 添加修改密码方法
	UpdatePassword(context.Context, *UpdatePasswordRequest) error
}

// CreateUserRequest 创建用户时，用户传递的请求
type CreateUserRequest struct {
	Username string            `json:"username"`
	Password string            `json:"password"`
	Role     Role              `json:"role"`
	Label    map[string]string `json:"label" gorm:"serializer:json"`
	isHashed bool              // 用于校验密码是否被hash过
}

func NewCreateUserRequest() *CreateUserRequest {
	return &CreateUserRequest{
		Role:  ROLE_AUTHOR,
		Label: map[string]string{},
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
		DescribeValue: id,
	}
}

func NewDescribeUserRequestByUsername(username string) *DescribeUserRequest {
	return &DescribeUserRequest{
		DescribeBy:    DESCRIBE_BY_USERNAME,
		DescribeValue: username,
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
	Role        Role              `json:"role,omitempty"`         // 可选，用于更新角色
	Label       map[string]string `json:"label,omitempty"`        // 可选，用于更新标签
}

func NewUpdateUserRequest() *UpdateUserRequest {
	return &UpdateUserRequest{
		Label: map[string]string{},
	}
}

func (req *CreateUserRequest) Validate() error {
	if req.Username == "" || req.Password == "" {
		return errors.New("用户名或密码不能为空")
	}
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
	if req.Username == "" {
		return fmt.Errorf("用户名不能为空")
	}

	// 如果设置了角色，验证角色值是否有效
	if req.Role != ROLE_UNKNOWN && !req.Role.IsValid() {
		return fmt.Errorf("无效的角色值: %d", req.Role)
	}

	return nil
}
