package user

import (
	"golang.org/x/crypto/bcrypt"
	"periodical/internal/apps/role"
	"periodical/internal/common"
)

type User struct {
	Id                 int64 `json:"id" gorm:"primaryKey;autoIncrement;column:id"`
	*common.Meta             // 通用信息
	*CreateUserRequest       // 用户传递过来的信息

	// 多角色与权限扩展信息
	Roles        []*role.Role `json:"roles" gorm:"many2many:user_roles;foreignKey:Id;joinForeignKey:UserId;References:Id;joinReferences:RoleId"`
	Permissions  []string     `json:"permissions,omitempty" gorm:"-"`
	IsSuperAdmin bool         `json:"is_super_admin" gorm:"-"`
}

func NewUser(req *CreateUserRequest) *User {
	req.PasswordHash()

	return &User{
		Meta:              common.NewMeta(),
		CreateUserRequest: req,
		Roles:             []*role.Role{},
		Permissions:       []string{},
	}
}

// CheckPassword 判断用户的密码是否正确
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

// Desensitize 脱敏密码并整理RoleIds
func (u *User) Desensitize() {
	if u.CreateUserRequest != nil {
		u.Password = ""
		if len(u.Roles) > 0 {
			ids := make([]int64, 0, len(u.Roles))
			for _, r := range u.Roles {
				ids = append(ids, r.Id)
				if r.IsSuperAdmin() {
					u.IsSuperAdmin = true
				}
			}
			u.RoleIds = ids
		} else if u.RoleIds == nil {
			u.RoleIds = []int64{}
		}
	}
	if u.Roles == nil {
		u.Roles = []*role.Role{}
	}
}

func (u *User) TableName() string {
	return "users"
}
