package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

// Config 全局配置结构
type Config struct {
	Server   ServerConfig   `yaml:"server"`
	Database DatabaseConfig `yaml:"database"`
	Redis    RedisConfig    `yaml:"redis"`
	JWT      JWTConfig      `yaml:"jwt"`
	Log      LogConfig      `yaml:"log"`
	MongoDB  MongoDBConfig  `yaml:"mongodb"`
}

// ServerConfig 服务配置
type ServerConfig struct {
	Port         string        `yaml:"port"`
	Mode         string        `yaml:"mode"` // debug, release, test
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string `yaml:"driver"`
	Host            string `yaml:"host"`
	Port            int    `yaml:"port"`
	Database        string `yaml:"database"`
	Username        string `yaml:"username"`
	Password        string `yaml:"password"`
	Charset         string `yaml:"charset"`
	MaxIdleConns    int    `yaml:"max_idle_conns"`
	MaxOpenConns    int    `yaml:"max_open_conns"`
	ConnMaxLifetime int    `yaml:"conn_max_lifetime"` // 分钟
}

// MongoDBConfig MongoDB配置
type MongoDBConfig struct {
	Host        string `yaml:"host"`
	Port        int    `yaml:"port"`
	Database    string `yaml:"database"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	AuthSource  string `yaml:"auth_source"`   // 认证数据库
	MaxPoolSize int    `yaml:"max_pool_size"` // 最大连接池大小
	MinPoolSize int    `yaml:"min_pool_size"` // 最小连接池大小
}

// RedisConfig Redis配置
type RedisConfig struct {
	Host         string `yaml:"host"`
	Port         int    `yaml:"port"`
	Password     string `yaml:"password"`
	DB           int    `yaml:"db"`
	PoolSize     int    `yaml:"pool_size"`
	MinIdleConns int    `yaml:"min_idle_conns"`
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret      string `yaml:"secret"`
	ExpireHours int    `yaml:"expire_hours"`
	Issuer      string `yaml:"issuer"`
}

// LogConfig 日志配置
type LogConfig struct {
	Level      string `yaml:"level"` // debug, info, warn, error
	FilePath   string `yaml:"file_path"`
	MaxSize    int    `yaml:"max_size"`    // MB
	MaxBackups int    `yaml:"max_backups"` // 保留旧文件的个数
	MaxAge     int    `yaml:"max_age"`     // 保留旧文件的天数
	Compress   bool   `yaml:"compress"`
}

var Cfg *Config

// LoadConfig 加载配置文件
func LoadConfig(configPath string) error {
	// 读取配置文件
	data, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析配置
	Cfg = &Config{}
	if err := yaml.Unmarshal(data, Cfg); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 从环境变量覆盖敏感配置
	if dbPassword := os.Getenv("DB_PASSWORD"); dbPassword != "" {
		Cfg.Database.Password = dbPassword
	}
	if mongoPassword := os.Getenv("MONGO_PASSWORD"); mongoPassword != "" {
		Cfg.MongoDB.Password = mongoPassword
	}
	if redisPassword := os.Getenv("REDIS_PASSWORD"); redisPassword != "" {
		Cfg.Redis.Password = redisPassword
	}
	if jwtSecret := os.Getenv("JWT_SECRET"); jwtSecret != "" {
		Cfg.JWT.Secret = jwtSecret
	}

	log.Println("✅ 配置加载成功")
	return nil
}

// GetConfig 获取全局配置
func GetConfig() *Config {
	if Cfg == nil {
		panic("配置未初始化，请先调用 LoadConfig")
	}
	return Cfg
}
