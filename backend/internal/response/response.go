package response

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Data struct {
	Code ResCode     `json:"code"`
	Msg  interface{} `json:"msg"`
	Data interface{} `json:"data"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, &Data{
		Code: CodeSuccess,
		Msg:  CodeSuccess.Msg(),
		Data: data,
	})
}

func Failed(c *gin.Context, code ResCode) {
	c.JSON(http.StatusOK, &Data{
		Code: code,
		Msg:  code.Msg(),
		Data: nil,
	})
}

func FailedWhitMsg(c *gin.Context, code ResCode, msg interface{}) {
	c.JSON(http.StatusOK, &Data{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
}
