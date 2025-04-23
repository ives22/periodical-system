package response

type ResCode int64

const (
	CodeSuccess ResCode = 200 + iota
	CodeInvalidParam
	CodeServerBusy
	CodeUnauthorized = 401 // 未授权
	CodeTokenExpired = 402 // token 过期
)

var CodeMsgMap = map[ResCode]string{
	CodeSuccess:      "success",
	CodeInvalidParam: "请求参数错误",
	CodeServerBusy:   "服务器错误",
	CodeUnauthorized: "未授权",
	CodeTokenExpired: "Token已经过期",
}

func (c ResCode) Msg() string {
	msg, ok := CodeMsgMap[c]
	if !ok {
		msg = CodeMsgMap[CodeServerBusy]
	}
	return msg
}
