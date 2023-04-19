package mysql

import (
	"fmt"
	"gateway/configs"
	"gateway/pkg/log"
	"time"

	"go.uber.org/zap"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	conf *configs.MySQLConfig
	db   *gorm.DB
)

// InitDB 初始化数据库
func InitDB() {
	conf = configs.GetMysqlConfig()
	var err error
	db, err = connectMySQL()
	if err != nil {
		log.Fatal("Failed to connect to MySQL: %v", zap.Error(err))
	}

	// 注册钩子函数
	db.Callback().Create().Before("gorm:before_create").Register("update_created_at", func(db *gorm.DB) {
		now := time.Now()
		if db.Statement.Schema != nil {
			if field, ok := db.Statement.Schema.FieldsByName["CreateAt"]; ok {
				field.Set(db.Statement.Context, db.Statement.ReflectValue, now)
			}
			if field, ok := db.Statement.Schema.FieldsByName["UpdateAt"]; ok {
				field.Set(db.Statement.Context, db.Statement.ReflectValue, now)
			}
		}
	})

	db.Callback().Update().Before("gorm:before_update").Register("update_updated_at", func(db *gorm.DB) {
		now := time.Now()
		if db.Statement.Schema != nil {
			if field, ok := db.Statement.Schema.FieldsByName["UpdateAt"]; ok {
				field.Set(db.Statement.Context, db.Statement.ReflectValue, now)
			}
		}
	})
}

// ConnectMySQL 连接MySQL数据库 (Connect to MySQL database)
func connectMySQL() (*gorm.DB, error) {
	// 构建DSN字符串 (Build the DSN string)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%t",
		conf.User, conf.Password, conf.Host, conf.Port, conf.DBName, conf.Charset, conf.ParseTime)

	// 连接数据库 (Connect to the database)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 设置连接池参数 (Set connection pool parameters)
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(conf.MaxIdleConns)
	sqlDB.SetMaxOpenConns(conf.MaxOpenConns)

	return db, nil
}

// GetDB 获取数据库连接
func GetDB() *gorm.DB {
	return db
}

// Close 关闭数据库连接池
func CloseDB() error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
