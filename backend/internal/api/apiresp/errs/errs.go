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

	// 群组相关
	ErrCodeGroupNotFound            = 13001
	ErrCodeGroupDismissed           = 13002
	ErrCodeGroupMemberNotFound      = 13003
	ErrCodeGroupPermissionDenied    = 13004
	ErrCodeGroupOwnerCannotQuit     = 13005
	ErrCodeGroupCannotRemoveOwner   = 13006
	ErrCodeGroupUseQuitGroup        = 13007
	ErrCodeGroupOnlyOwnerCanSetRole = 13008
	ErrCodeGroupRoleLevelTooHigh    = 13009
	ErrCodeGroupQuitSelfOnly        = 13010
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

	// 群组相关
	ErrGroupNotFound            = NewCodeError(ErrCodeGroupNotFound, "群组不存在")
	ErrGroupDismissed           = NewCodeError(ErrCodeGroupDismissed, "群组已解散")
	ErrGroupMemberNotFound      = NewCodeError(ErrCodeGroupMemberNotFound, "群成员不存在")
	ErrGroupPermissionDenied    = NewCodeError(ErrCodeGroupPermissionDenied, "群权限不足")
	ErrGroupOwnerCannotQuit     = NewCodeError(ErrCodeGroupOwnerCannotQuit, "群主无法直接退出，请先解散群组")
	ErrGroupCannotRemoveOwner   = NewCodeError(ErrCodeGroupCannotRemoveOwner, "无法移除群主")
	ErrGroupUseQuitGroup        = NewCodeError(ErrCodeGroupUseQuitGroup, "请使用退出群组接口")
	ErrGroupOnlyOwnerCanSetRole = NewCodeError(ErrCodeGroupOnlyOwnerCanSetRole, "只有群主可以调整角色等级")
	ErrGroupRoleLevelTooHigh    = NewCodeError(ErrCodeGroupRoleLevelTooHigh, "不能将角色设置为高于自身的等级")
	ErrGroupQuitSelfOnly        = NewCodeError(ErrCodeGroupQuitSelfOnly, "只能退出自己的群成员关系")
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
