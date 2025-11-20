package repository

import (
	"github.com/demo/database"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// BaseRepository 基础仓库
type BaseRepository struct {
	DB    *gorm.DB
	Redis *redis.Client
}

// NewBaseRepository 创建基础仓库实例
func NewBaseRepository() *BaseRepository {
	return &BaseRepository{
		DB:    database.GetDB(),
		Redis: database.GetRedis(),
	}
}
