package controllers

import (
	"github.com/demo/database"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// BaseController 基础控制器
type BaseController struct {
	DB    *gorm.DB
	Redis *redis.Client
}

// NewBaseController 创建基础控制器
func NewBaseController() *BaseController {
	return &BaseController{
		DB:    database.GetDB(),
		Redis: database.GetRedis(),
	}
}
