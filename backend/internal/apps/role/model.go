package role

import "periodical/internal/common"

const (
	SUPER_ADMIN_CODE = "super_admin"
	ADMIN_CODE       = "admin"
	AUTHOR_CODE      = "author"

	PERM_TYPE_MENU   = 1 // 菜单/目录权限
	PERM_TYPE_BUTTON = 2 // 按钮/接口操作权限
)

// Role 角色模型
type Role struct {
	Id          int64  `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	Name        string `json:"name" gorm:"column:name"`
	Code        string `json:"code" gorm:"column:code;uniqueIndex"`
	Description string `json:"description" gorm:"column:description"`
	Status      int64  `json:"status" gorm:"column:status;default:1"`       // 1-启用，0-禁用
	IsSystem    int64  `json:"is_system" gorm:"column:is_system;default:0"` // 1-系统内置，0-自定义
	Sort        int64  `json:"sort" gorm:"column:sort;default:0"`
	*common.Meta
	Permissions   []*Permission `json:"permissions,omitempty" gorm:"many2many:role_permissions;foreignKey:Id;joinForeignKey:RoleId;References:Id;joinReferences:PermissionId"`
	PermissionIds []int64       `json:"permission_ids,omitempty" gorm:"-"`
}

func (r *Role) TableName() string {
	return "roles"
}

func (r *Role) IsSuperAdmin() bool {
	return r.Code == SUPER_ADMIN_CODE || r.IsSystem == 1 && r.Code == SUPER_ADMIN_CODE
}

// Permission 权限资源模型
type Permission struct {
	Id       int64  `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	ParentId int64  `json:"parent_id" gorm:"column:parent_id;default:0"`
	Name     string `json:"name" gorm:"column:name"`
	Code     string `json:"code" gorm:"column:code;uniqueIndex"`
	Type     int64  `json:"type" gorm:"column:type;default:1"` // 1-菜单/目录，2-按钮/接口
	Sort     int64  `json:"sort" gorm:"column:sort;default:0"`
	*common.Meta
	Children []*Permission `json:"children,omitempty" gorm:"-"`
}

func (p *Permission) TableName() string {
	return "permissions"
}

// RolePermission 角色-权限关联表
type RolePermission struct {
	RoleId       int64 `json:"role_id" gorm:"primaryKey;column:role_id"`
	PermissionId int64 `json:"permission_id" gorm:"primaryKey;column:permission_id"`
}

func (rp *RolePermission) TableName() string {
	return "role_permissions"
}

// UserRole 用户-角色关联表（支持单用户绑定多个角色）
type UserRole struct {
	UserId int64 `json:"user_id" gorm:"primaryKey;column:user_id"`
	RoleId int64 `json:"role_id" gorm:"primaryKey;column:role_id"`
}

func (ur *UserRole) TableName() string {
	return "user_roles"
}
