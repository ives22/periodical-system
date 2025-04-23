package middleware

import (
	"periodical/internal/apps/token"
	"periodical/internal/response"

	"github.com/gin-gonic/gin"
)

func Auth(tokenService token.Service) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从cookie中获取token

		accessToken := c.Request.Header.Get("Authorization")
		if accessToken == "" {
			response.Failed(c, response.CodeUnauthorized)
			c.Abort()
			return
		}

		// 验证token
		tk, err := tokenService.ValidateToken(c.Request.Context(), &token.ValidateToken{
			AccessToken: accessToken,
		})
		if err != nil {
			response.FailedWhitMsg(c, response.CodeTokenExpired, err.Error())
			c.Abort()
			return
		}

		// 将用户信息存储到上下文中
		c.Set(token.TOKEN_GIN_KEY_NAME, tk)
		c.Next()

	}

}
