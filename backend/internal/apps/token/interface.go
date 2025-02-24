package token

import "context"

type Service interface {
	// 登录接口
	Login(context.Context, *LoginRequest) (*Token, error)
	// 退出接口
	Logout(context.Context, *LogoutRequest) error
	// 校验token，内部中间层使用
	ValidateToken(context.Context, *ValidateToken) (*Token, error)
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LogoutRequest struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ValidateToken struct {
	AccessToken string `json:"access_token"`
}

func NewLoginRequest() *LoginRequest {
	return &LoginRequest{}
}

func NewLogoutRequest() *LogoutRequest {
	return &LogoutRequest{}
}

func NewValidateToken(at string) *ValidateToken {
	return &ValidateToken{
		AccessToken: at,
	}
}
