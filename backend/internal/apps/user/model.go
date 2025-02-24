package user

import (
	"golang.org/x/crypto/bcrypt"
	"periodical/internal/common"
)

type User struct {
	Id                 int64 `json:"id"`
	*common.Meta             // 通用信息
	*CreateUserRequest       // 用户传递过来的信息
}

func NewUser(req *CreateUserRequest) *User {
	req.PasswordHash()

	return &User{
		Meta:              common.NewMeta(),
		CreateUserRequest: req,
	}
}

// CheckPassword 判断用户的密码是否正确
func (u *User) CheckPassword(password string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
}

func (u *User) TableName() string {
	return "users"
}
