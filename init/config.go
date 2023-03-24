// init/config.go
package init

import (
	"gateway/configs"
	"log"
	"sync"
)

// 用于确保仅初始化一次
// Used to ensure that initialization occurs only once
var onceConfig sync.Once

func init() {
	// 使用 sync.Once 仅执行一次初始化
	// Use sync.Once to initialize only once
	onceConfig.Do(func() {
		err := configs.LoadConfigurations()
		if err != nil {
			// 如果配置解析失败，则打印错误并退出
			// If the configuration parsing fails, print the error and exit
			log.Fatalf("Failed to load configurations: %v", err)
		}
	})
}
