package pkg

const (
	SuccessCode       ResponseCode = iota // 成功
	UndefErrorCode                        // 未定义的错误
	ValidErrorCode                        // 校验错误
	InternalErrorCode                     // 内部错误

	InvalidRequestErrorCode ResponseCode = 401  // 请求未经授权
	CustomizeCode           ResponseCode = 1000 // 自定义错误码
)
