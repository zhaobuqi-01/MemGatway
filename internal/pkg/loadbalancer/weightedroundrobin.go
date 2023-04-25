// loadbalancer/weightedroundrobin.go

package loadbalancer

import (
	"errors"
	"fmt"
	"strconv"
)

type WeightRoundRobinBalance struct {
	curIndex int
	rss      []*WeightNode // 服务器列表
	rsw      []int         // 节点权重（未使用）
}

type WeightNode struct {
	addr            string // 节点地址
	weight          int    // 节点权重（用于有效权重计算）
	currentWeight   int    // 节点当前权重（用于负载均衡）
	effectiveWeight int    // 节点有效权重（用于有效权重计算）
}

// Add 添加一个节点到服务器列表，包括地址和权重
func (r *WeightRoundRobinBalance) Add(params ...string) error {
	if len(params) != 2 {
		return errors.New("param len need 2")
	}
	parInt, err := strconv.ParseInt(params[1], 10, 64)
	if err != nil {
		return err
	}
	node := &WeightNode{addr: params[0], weight: int(parInt)}
	node.effectiveWeight = node.weight // 将有效权重初始化为节点权重
	r.rss = append(r.rss, node)        // 将节点添加到服务器列表
	return nil
}

// 使用加权轮询算法选择下一个节点
func (r *WeightRoundRobinBalance) get() string {
	total := 0
	var best *WeightNode // 当前权重最高的节点（即最佳选择）
	for i := 0; i < len(r.rss); i++ {
		w := r.rss[i]
		total += w.effectiveWeight           // 计算有效权重之和
		w.currentWeight += w.effectiveWeight // 将节点当前权重加上节点的有效权重
		if w.effectiveWeight < w.weight {
			w.effectiveWeight++ // 如果有效权重小于节点权重，则增加有效权重
		}
		if best == nil || w.currentWeight > best.currentWeight {
			best = w // 如果当前节点的当前权重比当前最佳节点的当前权重更高，则更新最佳选择
		}
	}
	if best == nil {
		return "" // 如果没有可供选择的节点，则返回空字符串
	}
	best.currentWeight -= total // 将最佳选择的节点的当前权重减去有效权重之和
	return best.addr            // 返回最佳选择的节点的地址
}

// Remove 从服务器列表中移除指定地址的节点
func (r *WeightRoundRobinBalance) Remove(addr string) error {
	for i := 0; i < len(r.rss); i++ {
		if r.rss[i].addr == addr {
			r.rss = append(r.rss[:i], r.rss[i+1:]...) // 使用切片append方法从列表中移除节点
			return nil
		}
	}
	return fmt.Errorf("node %s not found", addr) // 如果没有找到节点，则返回错误
}
