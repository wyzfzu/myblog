package utils

type ErrorCode struct {
	Code string
	Msg  string
}

var (
	SUCCESS          = ErrorCode{Code: "0", Msg: "成功"}
	PARAM_ERROR      = ErrorCode{Code: "1001", Msg: "请求参数错误"}
	FAIL             = ErrorCode{Code: "1002", Msg: "服务器错误"}
	TOKEN_EXPIRED    = ErrorCode{Code: "1003", Msg: "令牌已失效"}
	NO_AUTH          = ErrorCode{Code: "1004", Msg: "没有权限"}
	DATA_NOT_EXISTS  = ErrorCode{Code: "1005", Msg: "数据不存在"}
	USER_NAME_EXISTS = ErrorCode{Code: "2001", Msg: "用户名已存在"}
	USER_NOT_EXISTS  = ErrorCode{Code: "2002", Msg: "用户不存在"}
)

var (
	USER_ID_CTX_KEY = "UserID"
	POST_CTX_KEY    = "Post"
)
