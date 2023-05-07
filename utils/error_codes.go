package utils

// 0-999 通用错误码
// 1000-1999 用户相关错误码
// 2000-2999 应用相关错误码
// 3000-3999 面板相关错误码
// 4000-4999 服务相关错误码
// 5000-5999 代理服务相关错误码

// 通用错误码
const (
	// SuccessCode 成功
	SuccessCode = iota

	// ParamBindingErrCode 参数绑定失败
	ParamBindingErrCode
	// InternalErrorCode 内部错误
	InternalErrorCode
	//IPMismatchErrCode IP不匹配
	IpMismatchErrCode
	// ClientIPInBlackListErrCode 客户端IP在黑名单中
	ClientIPInBlackListErrCode
	// ClientIPNotInWhiteListCode 客户端IP不在白名单中
	ClientIPNotInWhiteListCode
	// GetLoadBalancerErrCode 获取负载均衡器失败
	GetLoadBalancerErrCode
	// GetTransportorErrCode 获取传输器失败
	GetTransportorErrCode
	// GetLimiterErrCode 获取限流器失败
	GetLimiterErrCode
	// ServerLimiterAllowErrCode 服务端限流
	ServerLimiterAllowErrCode
	// ClientIPLimiterAllowErrCode 客户端IP限流
	ClientIPLimiterAllowErrCode
	// CircuitBreakerOpenErrCode 熔断器打开
	CircuitBreakerOpenErrCode
	// ReverseProxyErrCode 反向代理失败
	ReverseProxyErrCode
	// NoSuchHostErrCode 无此主机
	NoSuchHostErrCode
)

// 用户相关错误码
const (
	// UserNotLoggedInErrCode 用户未登录
	UserNotLoggedInErrCode = iota + 1001
	// UserLoginErrCode 用户登录失败
	UserLoginErrCode
	// UserLoginErrCode 用户退出登录失败
	UserLoginOutErrCode
	// UserInfoErrCode 获取用户信息失败
	UserInfoErrCode
	// UserChangePwdErrCode 修改密码失败
	UserChangePwdErrCode
)

// 应用相关错误码
const (
	// AppNotFoundErrCode 应用未找到
	AppNotFoundErrCode = iota + 2001
	// AppListErrCode 获取应用列表失败
	AppListErrCode
	// AppDetailErrCode 获取应用详情失败
	AppDetailErrCode
	// AppDeleteErrCode 删除应用失败
	AppDeleteErrCode
	// AppAddErrCode 添加应用失败
	AppAddErrCode
	// AppUpdateErrCode 更新应用失败
	AppUpdateErrCode
)

// 面板相关错误码
const (
	// PanelGroupDataErrCode 获取面板分组数据失败
	PanelGroupDataErrCode = iota + 3001
	// ServiceStatErrCode 获取服务统计数据失败
	ServiceStatErrCode
)

// 服务相关错误码
const (
	// ServiceNotFoundErrCode 服务未找到
	ServiceNotFoundErrCode = iota + 4001
	// ServiceListErrCode 获取服务列表失败
	ServiceListErrCode
	// ServiceDeleteErrCode 删除服务失败
	ServiceDeleteErrCode
	// ServiceDetailErrCode 获取服务详情失败
	ServiceDetailErrCode
	// AddHttpServiceErrCode 添加HTTP服务失败
	AddHttpServiceErrCode
	// UpdateHttpServiceErrCode 更新HTTP服务失败
	UpdateHttpServiceErrCode
	// AddTCPServiceErrCode 添加TCP服务失败
	AddTCPServiceErrCode
	// UpdateTCPServiceErrCode 更新TCP服务失败
	UpdateTCPServiceErrCode
	// AddGRPCServiceErrCode 添加UDP服务失败
	AddGRPCServiceErrCode
	// UpdateGRPCServiceErrCode 更新UDP服务失败
	UpdateGRPCServiceErrCode
)

// proxy相关错误码
const (
	// HttpAccessModeErrCode HTTP接入方式匹配失败
	HTTPAccessModeErrCode = iota + 5001
)
