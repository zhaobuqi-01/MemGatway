package load_balance

import (
	"fmt"
	"net"
	"reflect"
	"sort"
	"time"
)

// 默认的负载均衡检查配置
const (
	DefaultCheckMethod    = 0
	DefaultCheckTimeout   = 5
	DefaultCheckMaxErrNum = 2
	DefaultCheckInterval  = 5
)

// LoadBalanceCheckConf 负载均衡检查配置
type LoadBalanceCheckConf struct {
	observers    []Observer        // 观察者列表
	confIpWeight map[string]string // 各服务器的IP和权重配置
	activeList   []string          // 当前可用的服务器列表
	format       string            // 配置格式
}

// Attach 注册观察者
func (s *LoadBalanceCheckConf) Attach(o Observer) {
	s.observers = append(s.observers, o)
}

// NotifyAllObservers 通知所有观察者
func (s *LoadBalanceCheckConf) NotifyAllObservers() {
	for _, obs := range s.observers {
		obs.Update()
	}
}

// GetConf 获取当前可用的服务器配置列表
func (s *LoadBalanceCheckConf) GetConf() []string {
	confList := []string{}
	for _, ip := range s.activeList {
		weight, ok := s.confIpWeight[ip]
		if !ok {
			weight = "50" // 默认权重
		}
		confList = append(confList, fmt.Sprintf(s.format, ip)+","+weight)
	}
	return confList
}

// WatchConf 监控配置变化
func (s *LoadBalanceCheckConf) WatchConf() {
	go func() {
		confIpErrNum := map[string]int{}
		for {
			changedList := []string{}
			// 检查每个服务器的可用性
			for item, _ := range s.confIpWeight {
				conn, err := net.DialTimeout("tcp", item, time.Duration(DefaultCheckTimeout)*time.Second)
				if err == nil {
					conn.Close()
					if _, ok := confIpErrNum[item]; ok {
						confIpErrNum[item] = 0
					}
				}
				if err != nil {
					if _, ok := confIpErrNum[item]; ok {
						confIpErrNum[item] += 1
					} else {
						confIpErrNum[item] = 1
					}
				}
				// 如果错误次数小于最大错误次数，将服务器添加到变更列表
				if confIpErrNum[item] < DefaultCheckMaxErrNum {
					changedList = append(changedList, item)
				}
			}
			sort.Strings(changedList)
			sort.Strings(s.activeList)
			// 如果服务器列表发生变化，则更新配置
			if !reflect.DeepEqual(changedList, s.activeList) {
				s.UpdateConf(changedList)
			}
			time.Sleep(time.Duration(DefaultCheckInterval) * time.Second)
		}
	}()
}

// Update
// 更新配置时，通知监听者也更新
func (s *LoadBalanceCheckConf) UpdateConf(conf []string) {
	s.activeList = conf
	for _, obs := range s.observers {
		obs.Update()
	}
}

// NewLoadBalanceCheckConf 创建一个新的负载均衡检查配置实例
func NewLoadBalanceCheckConf(format string, conf map[string]string) (*LoadBalanceCheckConf, error) {
	aList := []string{}
	// 默认初始化
	for item, _ := range conf {
		aList = append(aList, item)
	}
	mConf := &LoadBalanceCheckConf{format: format, activeList: aList, confIpWeight: conf}
	mConf.WatchConf()
	return mConf, nil
}
