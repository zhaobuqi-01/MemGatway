package database

import (
	"fmt"
	"gateway/configs"
	"gateway/pkg/logger"
	"io/ioutil"
	"os"
	"strings"

	"go.uber.org/zap"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// ConnectMySQL 连接MySQL数据库 (Connect to MySQL database)
func ConnectMySQL() (*gorm.DB, error) {
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

func GetDB() *gorm.DB {
	return db
}

// Close 关闭数据库连接池
// Close closes the database connection pool
func CloseDB() error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// CreateDatabaseIfNotExists
func CreateDatabase() error {
	dbName := conf.DBName
	charset := conf.Charset
	collation := conf.Collation

	// Create the database if it doesn't exist
	createDbSQL := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s CHARACTER SET %s COLLATE %s;", dbName, charset, collation)
	if err := db.Exec(createDbSQL).Error; err != nil {
		return err
	}
	// Use the created database
	useDbSQL := fmt.Sprintf("USE %s;", dbName)
	if err := db.Exec(useDbSQL).Error; err != nil {
		return err
	}

	return nil
}

func ExecuteSQLFile(filePath string) error {
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		logger.Info("No .sql file provided, skipping execution.", zap.String("file_path", filePath))
		return nil
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	sqlQueries := strings.Split(string(content), ";")

	for _, query := range sqlQueries {
		if strings.TrimSpace(query) != "" {
			if err := db.Exec(query).Error; err != nil {
				return err
			}
		}
	}

	return nil
}

var (
	conf *configs.MySQLConfig
	db   *gorm.DB
)

func CheckDatabaseExists() (bool, error) {
	dbName := conf.DBName
	var count int64
	row := db.Raw("SELECT COUNT(*) FROM information_schema.schemata WHERE schema_name = ?", dbName).Row()
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func InitDB() {
	conf = configs.GetMysqlConfig()
	var err error
	db, err = ConnectMySQL()
	if err != nil {
		logger.Fatal("Failed to connect to MySQL: %v", zap.Error(err))
	}

	// Check if the database exists
	dbExists, err := CheckDatabaseExists()
	if err != nil {
		logger.Fatal("Failed to check if the database exists: %v", zap.Error(err))
	}

	// Create the database if it doesn't exist
	if !dbExists {
		if err := CreateDatabase(); err != nil {
			logger.Fatal("Failed to create the database: %v", zap.Error(err))
		}

		err = ExecuteSQLFile(conf.SqlFile)
		if err != nil {
			logger.Fatal("Failed to execute .sql file: %v", zap.Error(err))
		}
	}
	logger.Info("database already exists", zap.String("databaseName", conf.DBName))
	logger.Info("Start using an existing database")
}
