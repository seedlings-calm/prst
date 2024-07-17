package db

import (
	"fmt"
	"log"
	"time"

	cfg "github.com/seedlings-calm/prst/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

var Db *gorm.DB

func GormMysql() {
	config := cfg.GetGlobalConf().Mysql
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.User, config.Password, config.Host, config.Port, config.DBName)

	DB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get sql.DB from gorm.DB: %v", err)
	}

	// Set database connection pool parameters
	sqlDB.SetMaxIdleConns(config.IdleConn) //  设置空闲连接池中的最大连接数。
	sqlDB.SetMaxOpenConns(config.OpenConn) //  设置数据库的最大打开连接数。
	sqlDB.SetConnMaxLifetime(time.Hour)
	Db = DB
}
