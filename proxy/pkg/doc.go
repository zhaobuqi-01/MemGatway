// package pkg
//
// pkg包是反向代理服务器的公共包，封装在tcp，hhtp，grpc共用的函数，方法，全局变量
//
// 全局变量
//
// Cache 提供缓存功能，包括 AppCache 和 ServiceCache
// FlowLimiter 提供限流功能
// LoadBalanceTransport 提供负载均衡和传输功能
//
// 方法
// Init 函数用于初始化全局变量，它只会被执行一次
//
// # NewFlowLimiter 创建并返回一个新的flowLimiter实例
//
// GetLimiter 实现了Limiter接口中的GetLimiter方法，该方法首先检查映射中是否已经有对应服务的限流器，
// 如果有则返回，如果没有则新建一个限流器并存入sync.Map中
//
// # GetLoadBalancer 获取LoadBalancer实例，如果不存在则创建一个新的实例并添加到映射中
//
// # GetTransportor 获取Transportor实例，如果不存在则创建一个新的实例并添加到映射中
//
// GetApp 通过 appID 返回 app。
//
// LoadAppCache 将所有 app 数据加载到缓存中。
//
// LoadService 将所有 service 数据加载到缓存中。
//
// UpdateAppCache 通过appID和operation更新app缓存。
//
// UpdateServiceCache 通过serviceName和operation更新service缓存。
//
// # HTTPAccessMode 通过cxt获取服务的访问模式
//
// # GetGrpcServiceList 获取所有的grpc服务
//
// GetTcpServiceList 获取所有的tcp服务
package pkg
