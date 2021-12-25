package database

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"strconv"
	"time"
)

var db *gorm.DB = nil

func init() {
	var err error

	db, err = newDB(os.Getenv("DB_SERVER_DB_NAME"), logger.Info)
	if err != nil {
		panic(err)
	}
}

func GetDB() *gorm.DB {
	return db
}

func newDB(dbname string, logLevel logger.LogLevel) (gormDB *gorm.DB, err error) {
	switch dbType := os.Getenv("DB_SERVER_TYPE"); dbType {
	case "postgresql":
		fmt.Println("使用Postgresql")
		gormDB, err = newPostgreSQLConnection(dbname, logLevel)
	default:
		return nil, fmt.Errorf("当前仅支持Postresql,不支持%s", dbType)
	}
	if err != nil {
		return nil, err
	}
	return gormDB, nil
}

func newPostgreSQLConnection(dbname string, logLevel logger.LogLevel) (*gorm.DB, error) {
	// Define database connection settings.
	maxConn, _ := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	maxIdleConn, _ := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	maxLifetimeConn, _ := strconv.Atoi(os.Getenv("DB_MAX_LIFETIME_CONNECTIONS"))

	host := os.Getenv("DB_SERVER_HOST")
	port := os.Getenv("DB_SERVER_PORT")
	user := os.Getenv("DB_SERVER_USER")
	password := os.Getenv("DB_SERVER_PASSWORD")
	ssl := os.Getenv("DB_SERVER_SSL")
	timezone := os.Getenv("DB_SERVER_TIMEZONE")
	// Define database connection for PostgreSQL.
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s", host, user, password, dbname, port, ssl, timezone)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	// Set database connection settings.
	sqlDB.SetMaxOpenConns(maxConn)                                                              // the default is 0 (unlimited)
	sqlDB.SetMaxIdleConns(maxIdleConn)                                                          // defaultMaxIdleConns = 2
	sqlDB.SetConnMaxLifetime(time.Duration(int64(maxLifetimeConn) * time.Minute.Nanoseconds())) // 0, connections are reused forever

	// Try to ping database.
	if err := sqlDB.Ping(); err != nil {
		defer sqlDB.Close() // close database connection
		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}

	return db, nil
}
