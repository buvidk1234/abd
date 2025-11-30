package errs

// 统一错误码定义
const (
	ErrCodeCommon         = 99999 // 通用错误
	ErrCodeSuccess        = 0
	ErrCodeInvalidParam   = 10001
	ErrCodeUnauthorized   = 10002
	ErrCodeNotFound       = 10003
	ErrCodeInternalServer = 20001

	// 用户相关
	ErrCodeUserNotFound      = 11001
	ErrCodeUserExists        = 11002
	ErrCodeUserPasswordWrong = 11003
	ErrCodeUserBanned        = 11004

	// 好友相关
	ErrCodeFriendNotFound        = 12001
	ErrCodeAlreadyFriends        = 12002
	ErrCodeFriendRequestSent     = 12003
	ErrCodeFriendRequestNotFound = 12004
	ErrCodeFriendBlocked         = 12005
)

// 常用错误变量
var (
	ErrCommon         = NewCodeError(ErrCodeCommon, "通用错误")
	ErrInvalidParam   = NewCodeError(ErrCodeInvalidParam, "参数错误")
	ErrUnauthorized   = NewCodeError(ErrCodeUnauthorized, "未授权")
	ErrNotFound       = NewCodeError(ErrCodeNotFound, "资源未找到")
	ErrInternalServer = NewCodeError(ErrCodeInternalServer, "服务器内部错误")

	// 用户相关
	ErrUserNotFound      = NewCodeError(ErrCodeUserNotFound, "用户不存在")
	ErrUserExists        = NewCodeError(ErrCodeUserExists, "用户已存在")
	ErrUserPasswordWrong = NewCodeError(ErrCodeUserPasswordWrong, "密码错误")
	ErrUserBanned        = NewCodeError(ErrCodeUserBanned, "用户已被封禁")

	// 好友相关
	ErrFriendNotFound        = NewCodeError(ErrCodeFriendNotFound, "好友不存在")
	ErrAlreadyFriends        = NewCodeError(ErrCodeAlreadyFriends, "已经是好友")
	ErrFriendRequestSent     = NewCodeError(ErrCodeFriendRequestSent, "好友请求已发送")
	ErrFriendRequestNotFound = NewCodeError(ErrCodeFriendRequestNotFound, "好友请求不存在")
	ErrFriendBlocked         = NewCodeError(ErrCodeFriendBlocked, "好友已被拉黑")
)

// CodeError 结构体和构造函数
type CodeError struct {
	Code int
	Msg  string
}

func (e *CodeError) Error() string {
	return e.Msg
}

func NewCodeError(code int, msg string) *CodeError {
	return &CodeError{Code: code, Msg: msg}
}

func (e *CodeError) WithDetail(detail string) *CodeError {
	return &CodeError{Code: e.Code, Msg: e.Msg + ": " + detail}
}
