package token

import (
	"encoding/json"
	"fmt"
	"github.com/rs/xid"
	"periodical/internal/apps/user"
	"time"
)

func NewToken() *Token {
	return &Token{
		AccessToken:           xid.New().String(),
		AccessTokenExpiredAt:  7200,
		RefreshToken:          xid.New().String(),
		RefreshTokenExpiredAt: 3600 * 24 * 7,
		CreatedAt:             time.Now().Unix(),
		IsEnable:              1,
	}
}

type Token struct {
	UserId                int64  `json:"user_id"`
	UserName              string `json:"username" gorm:"column:username"`
	AccessToken           string `json:"access_token"`             // 颁发给用户的访问令牌
	AccessTokenExpiredAt  int    `json:"access_token_expired_at"`  // 过期时间（2h） 单位秒
	RefreshToken          string `json:"refresh_token"`            // 刷新token（）
	RefreshTokenExpiredAt int    `json:"refresh_token_expired_at"` // 刷新token过期时间(7d)
	CreatedAt             int64  `json:"created_at"`               // 创建时间
	UpdatedAt             int64  `json:"updated_at"`               // 更新时间
	IsEnable              int64  `json:"is_enable"`                // 是否有效（1有效，0无效） 默认为1

	// 额外补充信息
	Role user.Role `json:"role" gorm:"-"`
}

func (t *Token) TableName() string {
	return "tokens"
}

func (t *Token) IsExpired() error {
	duration := time.Since(t.ExpiredTime())
	expiredSeconds := duration.Seconds()
	if expiredSeconds > 0 {
		return fmt.Errorf("token %s 过期了 %f秒", t.AccessToken, expiredSeconds)
	}
	return nil

}

// ExpiredTime 计算token过期时间
func (t *Token) ExpiredTime() time.Time {
	return time.Unix(t.CreatedAt, 0).Add(time.Duration(t.AccessTokenExpiredAt) * time.Second)
}

func (t *Token) String() string {
	dj, _ := json.Marshal(t)
	return string(dj)

}
