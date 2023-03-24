package init

import (
	"context"
	"database/sql"
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

// GetConnection retrieves a connection from the connection pool
func GetConnection(ctx context.Context) (*sql.Conn, error) {
	sqlDB, err := DB.DB()
	if err != nil {
		return nil, err
	}

	conn, err := sqlDB.Conn(ctx)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// CloseConnection releases a connection back to the connection pool
func CloseConnection(conn *sql.Conn) error {
	if conn == nil {
		return nil
	}

	return conn.Close()
}

func CloseDB() error {
	sqlDB, err := DB.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
