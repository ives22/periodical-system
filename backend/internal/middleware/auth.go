package middleware

import (
	"periodical/internal/apps/token"
	"periodical/internal/response"

	"github.com/gin-gonic/gin"
)

// Auth 认证中间件：校验Token有效性并挂载用户角色与权限集
func Auth(tokenService token.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.Request.Header.Get("Authorization")
		if accessToken == "" {
			response.Failed(c, response.CodeUnauthorized)
			c.Abort()
			return
		}

		// 验证token并获取用户最新角色与权限列表
		tk, err := tokenService.ValidateToken(c.Request.Context(), &token.ValidateToken{
			AccessToken: accessToken,
		})
		if err != nil {
			response.FailedWhitMsg(c, response.CodeTokenExpired, err.Error())
			c.Abort()
			return
		}

		// 将Token与权限信息存储到上下文中
		c.Set(token.TOKEN_GIN_KEY_NAME, tk)
		c.Next()
	}
}

// GetTokenFromCtx 从Gin上下文中获取当前登录用户的Token信息
func GetTokenFromCtx(c *gin.Context) *token.Token {
	val, exists := c.Get(token.TOKEN_GIN_KEY_NAME)
	if !exists {
		return nil
	}
	tk, ok := val.(*token.Token)
	if !ok {
		return nil
	}
	return tk
}

// RequirePermission 动态权限校验中间件：校验当前用户是否拥有指定的权限标识码
func RequirePermission(permCode string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tk := GetTokenFromCtx(c)
		if tk == nil {
			response.Failed(c, response.CodeUnauthorized)
			c.Abort()
			return
		}

		if !tk.HasPermission(permCode) {
			response.FailedWhitMsg(c, response.CodeForbidden, "无权限执行此操作（缺少权限: "+permCode+"）")
			c.Abort()
			return
		}

		c.Next()
	}
}
