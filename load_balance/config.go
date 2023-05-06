package load_balance

// 配置主题
type LoadBalanceConf interface {
	Attach(o Observer)
	GetConf() []string
	WatchConf()
	UpdateConf(conf []string)
}

// Observer 观察者接口，用于实现观察者模式
type Observer interface {
	Update()
}
