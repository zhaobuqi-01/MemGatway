package pkg

import (
	"fmt"
	"gateway/enity"
	"gateway/globals"
	"gateway/proxy/load_balance"
	"gateway/utils"
	"net"
	"net/http"
	"sync"
	"time"
)

// LoadBalanceAndTransport 接口组合了 GetLoadBalancer 和 GetTransportor 两个接口
type LoadBalanceAndTransport interface {
	GetLoadBalancer(service *enity.ServiceDetail) (load_balance.LoadBalance, error)
	GetTransportor(service *enity.ServiceDetail) (*http.Transport, error)
	Remove(serviceName string)
}

const (
	defaultUpstreamConnectTimeout = 30
	defaultUpstreamMaxIdle        = 100
	defaultUpstreamIdleTimeout    = 90
	defaultUpstreamHeaderTimeout  = 30
	defaultKeepAliveDuration      = 30 * time.Second
)

type loadBalanceAndTransport struct {
	loadBalanceMap sync.Map // 存储LoadBalancerItem的同步映射
	transportMap   sync.Map // 存储TransportItem的同步映射
}

// NewLoadBalancer 返回一个新的 LoadBalancer 实例
func NewLoadBalancerAndTransport() *loadBalanceAndTransport {
	return &loadBalanceAndTransport{
		loadBalanceMap: sync.Map{},
		transportMap:   sync.Map{},
	}
}

func (lbr *loadBalanceAndTransport) Remove(serviceName string) {
	lbr.loadBalanceMap.Delete(serviceName)
	lbr.transportMap.Delete(serviceName)
}

// GetLoadBalancer 获取LoadBalancer实例，如果不存在则创建一个新的实例并添加到映射中
func (lbr *loadBalanceAndTransport) GetLoadBalancer(service *enity.ServiceDetail) (load_balance.LoadBalance, error) {
	if service == nil || service.Info == nil || service.LoadBalance == nil {
		return nil, fmt.Errorf("service or service info or load balance is nil")
	}
	if service.GRPCRule == nil && service.HTTPRule == nil && service.TCPRule == nil {
		return nil, fmt.Errorf("grpc rule, http rule and tcp rule are all nil")
	}
	if ipList := utils.SplitStringByComma(service.LoadBalance.IpList); ipList == nil {
		return nil, fmt.Errorf("ip list is nil")
	}
	if weightList := utils.SplitStringByComma(service.LoadBalance.WeightList); weightList == nil {
		return nil, fmt.Errorf("weight list is nil")
	}

	if lbrItem, ok := lbr.loadBalanceMap.Load(service.Info.ServiceName); ok {
		return lbrItem.(load_balance.LoadBalance), nil
	}

	schema := "http://"
	if service.HTTPRule != nil && service.HTTPRule.NeedHttps == 1 {
		schema = "https://"
	}
	if service.Info.LoadType == globals.LoadTypeTCP || service.Info.LoadType == globals.LoadTypeGRPC {
		schema = ""
	}

	ipConf := make(map[string]string)
	if ipList := utils.SplitStringByComma(service.LoadBalance.IpList); ipList != nil {
		weightList := utils.SplitStringByComma(service.LoadBalance.WeightList)
		for i := 0; i < len(ipList); i++ {
			if i < len(weightList) {
				ipConf[ipList[i]] = weightList[i]
			} else {
				ipConf[ipList[i]] = "1"
			}
		}
	}

	mConf, err := load_balance.NewLoadBalanceCheckConf(fmt.Sprintf("%s%s", schema, "%s"), ipConf)
	if err != nil {
		return nil, err
	}

	lb := load_balance.LoadBanlanceFactorWithConf(load_balance.LbType(service.LoadBalance.RoundType), mConf)
	lbr.loadBalanceMap.Store(service.Info.ServiceName, lb)

	return lb, nil
}

// GetTransportor 根据服务详情获取Transportor实例，如果映射中不存在则创建一个新的实例并添加到映射中
func (t *loadBalanceAndTransport) GetTransportor(service *enity.ServiceDetail) (*http.Transport, error) {
	// 如果已经存在该服务的 TransportItem，则直接返回
	if transItem, ok := t.transportMap.Load(service.Info.ServiceName); ok {
		return transItem.(*http.Transport), nil
	}

	// 如果不存在该服务的 TransportItem，则创建一个新的实例并添加到映射中
	if service.LoadBalance.UpstreamConnectTimeout == 0 {
		service.LoadBalance.UpstreamConnectTimeout = defaultUpstreamConnectTimeout
	}
	if service.LoadBalance.UpstreamMaxIdle == 0 {
		service.LoadBalance.UpstreamMaxIdle = defaultUpstreamMaxIdle
	}
	if service.LoadBalance.UpstreamIdleTimeout == 0 {
		service.LoadBalance.UpstreamIdleTimeout = defaultUpstreamIdleTimeout
	}
	if service.LoadBalance.UpstreamHeaderTimeout == 0 {
		service.LoadBalance.UpstreamHeaderTimeout = defaultUpstreamHeaderTimeout
	}
	trans := &http.Transport{
		Proxy: http.ProxyFromEnvironment, // 从系统环境中获取代理信息（如果存在），并将其应用于请求
		DialContext: (&net.Dialer{ // 初始化 Dialer，控制如何建立与上游服务的连接，并绑定到 Transport 的 DialContext 字段上
			Timeout:   time.Duration(service.LoadBalance.UpstreamConnectTimeout) * time.Second, // 建立连接的超时时间
			KeepAlive: defaultKeepAliveDuration,                                                // 连接保持的时间
			DualStack: true,                                                                    // 是否启用 IPv6 和 IPv4
		}).DialContext,
		ForceAttemptHTTP2:     true,                                                                   // 启用 HTTP/2 支持
		MaxIdleConns:          service.LoadBalance.UpstreamMaxIdle,                                    // 最大空闲连接数，最多允许保持多少个空闲的连接
		IdleConnTimeout:       time.Duration(service.LoadBalance.UpstreamIdleTimeout) * time.Second,   // 空闲连接的超时时间
		TLSHandshakeTimeout:   10 * time.Second,                                                       // TLS 握手超时时间
		ResponseHeaderTimeout: time.Duration(service.LoadBalance.UpstreamHeaderTimeout) * time.Second, // 响应头部超时时间
	}

	// 将 TransportItem 添加到映射中并返回
	t.transportMap.Store(service.Info.ServiceName, trans)
	return trans, nil
}
