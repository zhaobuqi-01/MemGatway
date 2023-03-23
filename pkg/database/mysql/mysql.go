package mysql

import (
	"fmt"
	"gateway/configs"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ConnectMySQL 连接MySQL数据库 (Connect to MySQL database)
func ConnectMySQL() (*gorm.DB, error) {
	conf := configs.GetMySQLConfig()

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
