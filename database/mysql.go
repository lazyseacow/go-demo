package database

import (
	"fmt"
	"log"
	"time"

	"github.com/demo/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// InitMySQL 初始化 MySQL 连接
func InitMySQL() error {
	cfg := config.GetConfig()
	dbCfg := cfg.Database

	// 构建 DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		dbCfg.Username,
		dbCfg.Password,
		dbCfg.Host,
		dbCfg.Port,
		dbCfg.Database,
		dbCfg.Charset,
	)

	// 设置 GORM 日志级别
	var logLevel logger.LogLevel
	switch cfg.Log.Level {
	case "debug":
		logLevel = logger.Info
	case "info":
		logLevel = logger.Warn
	default:
		logLevel = logger.Error
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	})

	if err != nil {
		return fmt.Errorf("连接 MySQL 失败: %v", err)
	}

	// 获取底层的 sql.DB 并设置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库实例失败: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(dbCfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(dbCfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(dbCfg.ConnMaxLifetime) * time.Minute)

	// 测试连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("MySQL 连接测试失败: %v", err)
	}

	DB = db
	log.Println("✅ MySQL 连接成功")
	return nil
}

// GetDB 获取数据库实例
func GetDB() *gorm.DB {
	if DB == nil {
		panic("数据库未初始化，请先调用 InitMySQL")
	}
	return DB
}

// CloseMySQL 关闭数据库连接
func CloseMySQL() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
