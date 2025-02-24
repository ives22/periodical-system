package response

type ResCode int64

const (
	CodeSuccess ResCode = 200 + iota
	CodeInvalidParam
	CodeServerBusy
)

var CodeMsgMap = map[ResCode]string{
	CodeSuccess:      "success",
	CodeInvalidParam: "请求参数错误",
	CodeServerBusy:   "服务器错误",
}

func (c ResCode) Msg() string {
	msg, ok := CodeMsgMap[c]
	if !ok {
		msg = CodeMsgMap[CodeServerBusy]
	}
	return msg
}
