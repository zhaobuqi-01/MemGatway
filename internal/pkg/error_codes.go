package pkg

const (
	SuccessCode = iota // 成功

	ParamBindingErrCode    = 1001 // 参数绑定错误
	InternalErrorCode      = 2001 // 内部错误(业务逻辑返回的错误)
	UserNotLoggedInErrCode = 3001
	IpMismatchErrCode      = 3002
)
