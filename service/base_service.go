package service

import (
	"github.com/demo/database"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// BaseService 基础服务（包含数据库连接）
type BaseService struct {
	DB    *gorm.DB
	Redis *redis.Client
}

// NewBaseService 创建基础服务实例
func NewBaseService() *BaseService {
	return &BaseService{
		DB:    database.GetDB(),
		Redis: database.GetRedis(),
	}
}
