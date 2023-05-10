package mysql

import (
	"fmt"
	"gateway/configs"
	"gateway/pkg/log"
	"time"

	"go.uber.org/zap"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	conf *configs.MySQLConfig
	db   *gorm.DB
)

// InitDB initializes the database
func Init() {
	// 初始化
	conf = configs.GetMysqlConfig()
	var err error
	db, err = connectMySQL()
	if err != nil {
		log.Fatal("Failed to connect to MySQL", zap.Error(err))
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
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Discard,
	})
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

// type GormLogger struct {
// 	level logger.LogLevel
// }

// func NewGormLogger(level logger.LogLevel) *GormLogger {
// 	return &GormLogger{level: level}
// }

// func (l *GormLogger) LogMode(level logger.LogLevel) logger.Interface {
// 	newlogger := *l
// 	newlogger.level = level
// 	return &newlogger
// }

// func (l GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
// 	log.Info(fmt.Sprintf(msg, data...))
// }

// func (l GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
// 	log.Warn(fmt.Sprintf(msg, data...))
// }

// func (l GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
// 	log.Error(fmt.Sprintf(msg, data...))
// }

// func (l GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
// 	elapsed := time.Since(begin)
// 	switch {
// 	case err != nil && l.level >= logger.Error:
// 		sql, rows := fc()
// 		log.Error("trace", zap.String("sql", sql), zap.Int64("rows", rows), zap.Duration("elapsed", elapsed), zap.Error(err))
// 	case elapsed > 200*time.Millisecond && l.level >= logger.Warn:
// 		sql, rows := fc()
// 		log.Warn("trace", zap.String("sql", sql), zap.Int64("rows", rows), zap.Duration("elapsed", elapsed))
// 	case l.level >= logger.Info:
// 		sql, rows := fc()
// 		log.Info("trace", zap.String("sql", sql), zap.Int64("rows", rows), zap.Duration("elapsed", elapsed))
// 	}
// }
