// init/logger.go
package init

import (
	"gateway/pkg/logger"
	"log"
	"sync"
)

// 用于确保仅初始化一次
// Used to ensure that initialization occurs only once
var onceLogger sync.Once

func init() {
	// 使用 sync.Once 仅执行一次初始化
	// Use sync.Once to initialize only once
	onceLogger.Do(func() {
		err := logger.InitLogger()
		if err != nil {
			// 如果日志初始化失败，则打印错误并退出
			// If the logger initialization fails, print the error and exit
			log.Fatalf("Failed to initialize logger: %v", err)
		}
	})
}
