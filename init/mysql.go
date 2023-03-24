package init

import (
	"gateway/pkg/database"
	log "gateway/pkg/logger"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"sync"
)

var DB *gorm.DB
var onceMySQL sync.Once

func init() {
	onceMySQL.Do(func() {
		db, err := database.ConnectMySQL()
		if err != nil {
			log.Fatal("Failed to connect to MySQL: %v", zap.Error(err))
		}
		DB = db
	})
}
