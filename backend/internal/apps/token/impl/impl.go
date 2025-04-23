package impl

import (
	"context"
	"errors"
	"fmt"
	"log"
	"periodical/conf"
	"periodical/internal/apps/token"
	"periodical/internal/apps/user"

	"gorm.io/gorm"
)

type TokenImpl struct {
	db   *gorm.DB
	user user.Service
}

func NewTokenImpl(userSvc user.Service) *TokenImpl {
	return &TokenImpl{
		db:   conf.C().MySQL.GetConn().Debug(),
		user: userSvc,
	}
}

// Login 登录
func (i *TokenImpl) Login(ctx context.Context, req *token.LoginRequest) (*token.Token, error) {
	// 1. 查询用户
	uReq := user.NewDescribeUserRequestByUsername(req.Username)
	u, err := i.user.DescribeUser(ctx, uReq)
	if err != nil {
		return nil, err
	}

	// 2. 比对密码
	err = u.CheckPassword(req.Password)
	if err != nil {
		return nil, errors.New("用户名或者密码不正确")
	}

	// 3. 颁发token
	tk := token.NewToken()
	tk.UserId = u.Id
	tk.UserName = u.Username
	tk.Role = u.Role

	// 4. 保存token
	if err = i.db.WithContext(ctx).Create(tk).Error; err != nil {
		return nil, err
	}

	// 5. 颁发成功后，将之前的token标记为失效
	if err = i.db.Model(&tk).WithContext(ctx).
		Where("username = ? and access_token != ? and is_enable = 1", tk.UserName, tk.AccessToken).
		Update("is_enable", 0).Error; err != nil {
		return nil, err
	}

	return tk, nil
}

// Logout 退出登录
func (i *TokenImpl) Logout(ctx context.Context, req *token.LogoutRequest) error {
	tk := token.NewToken()
	err := i.db.Model(&tk).WithContext(ctx).
		Where("access_token = ? and refresh_token = ?", req.AccessToken, req.RefreshToken).
		Update("is_enable", "0").Error
	if err != nil {
		return err
	}

	return nil
}

// ValidateToken 校验token
func (i *TokenImpl) ValidateToken(ctx context.Context, req *token.ValidateToken) (*token.Token, error) {
	// 1. 查询token（是不是该系统颁发的）
	tk := token.NewToken()
	err := i.db.WithContext(ctx).Where("access_token = ?", req.AccessToken).First(tk).Error
	if err != nil {
		return nil, err
	}

	log.Println("校验token")
	log.Println(tk)
	// 2. 判断token是否被禁用
	if tk.IsEnable != 1 {
		return nil, errors.New("token已失效")
	}

	// 3. 判断token合法性
	// 3.1 判断AK是否过期
	if err = tk.IsExpired(); err != nil {
		return nil, err
	}

	// 3.2 补充用户信息，只是为了补充用户的角色
	uDesc := user.NewDescribeUserRequestById(fmt.Sprintf("%d", tk.UserId))
	u, err := i.user.DescribeUser(ctx, uDesc)
	if err != nil {
		return nil, err
	}

	tk.Role = u.Role

	return tk, nil
}
