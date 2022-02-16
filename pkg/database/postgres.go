package database

import (
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"time"
)

type PostgresConfig struct {
	Username              string        `mapstructure:"username"`
	Password              string        `mapstructure:"password"`
	Host                  string        `mapstructure:"host"`
	Port                  int           `mapstructure:"port"`
	Database              string        `mapstructure:"database"`
	SSL                   string        `mapstructure:"ssl"`
	Timezone              string        `mapstructure:"timezone"`
	MaxOpenConnections    int           `mapstructure:"max_open_connections"`
	MaxIdleConnections    int           `mapstructure:"max_idle_connections"`
	ConnectionMaxLifetime time.Duration `mapstructure:"connection_max_lifetime"`
}

var postgresConfig PostgresConfig

var db *gorm.DB = nil

func InitDB() {
	config := viper.Sub("postgres")
	if config == nil {
		log.Println("Postgres config not found")
		return
	}
	config.SetDefault("ssl", "disable")
	config.SetDefault("timezone", "Asia/Shanghai")
	config.SetDefault("max_open_connections", 16)
	config.SetDefault("max_idle_connections", 4)
	config.SetDefault("connection_max_lifetime", "4m")

	err := config.Unmarshal(&postgresConfig)
	if err != nil {
		panic(err)
	}

	db, err = newDB(postgresConfig.Database, logger.Info)
	if err != nil {
		panic(err)
	}
}

func GetDB() *gorm.DB {
	return db
}

func newDB(dbname string, logLevel logger.LogLevel) (gormDB *gorm.DB, err error) {
	switch dbType := "postgresql"; dbType {
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
	maxConn := postgresConfig.MaxOpenConnections
	maxIdleConn := postgresConfig.MaxIdleConnections
	maxLifetimeConn := postgresConfig.ConnectionMaxLifetime

	host := postgresConfig.Host
	port := postgresConfig.Port
	user := postgresConfig.Username
	password := postgresConfig.Password
	ssl := postgresConfig.SSL
	timezone := postgresConfig.Timezone
	// Define database connection for PostgreSQL.
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s", host, user, password, dbname, port, ssl, timezone)
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
