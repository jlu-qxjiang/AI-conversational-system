package controller

type ResCode int64

// 常量
const (
	CodeSuccess ResCode = 1000 + iota
	CodeInvalidParam
	CodeUserExist
	CodeUserNotExist
	CodeInvalidPassword
	CodeServerBusy
	CodeInvalidtoken
	CodeNeedLogin
	CodeFileNotExist
	CodeFileExist
	CodeFileUploadSuccess
	CodeFileNotUploaded
	CodeFileUploadFailed
)

var codeMsgMap = map[ResCode]string{
	CodeSuccess:           "OK",
	CodeInvalidParam:      "请求参数错误",
	CodeUserExist:         "用户已存在",
	CodeUserNotExist:      "用户不存在",
	CodeInvalidPassword:   "用户名或密码错误",
	CodeServerBusy:        "服务繁忙",
	CodeInvalidtoken:      "无效的token",
	CodeNeedLogin:         "需要登录",
	CodeFileNotExist:      "文件不存在",
	CodeFileUploadSuccess: "文件上传成功",
	CodeFileExist:         "文件已存在",
	CodeFileNotUploaded:   "文件未上传",
	CodeFileUploadFailed:  "文件上传失败",
}

func (c ResCode) Msg() string {
	msg, ok := codeMsgMap[c]
	if !ok {
		msg = codeMsgMap[CodeServerBusy]
	}
	return msg
}
