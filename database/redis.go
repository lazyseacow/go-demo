package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/demo/config"
	"github.com/redis/go-redis/v9"
)

var RDB *redis.Client

// InitRedis 初始化 Redis 连接
func InitRedis() error {
	cfg := config.GetConfig()
	redisCfg := cfg.Redis

	// 创建 Redis 客户端
	RDB = redis.NewClient(&redis.Options{
		Addr:         fmt.Sprintf("%s:%d", redisCfg.Host, redisCfg.Port),
		Password:     redisCfg.Password,
		DB:           redisCfg.DB,
		PoolSize:     redisCfg.PoolSize,
		MinIdleConns: redisCfg.MinIdleConns,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := RDB.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("连接 Redis 失败: %v", err)
	}

	log.Println("✅ Redis 连接成功")
	return nil
}

// GetRedis 获取 Redis 客户端
func GetRedis() *redis.Client {
	if RDB == nil {
		panic("Redis 未初始化，请先调用 InitRedis")
	}
	return RDB
}

// CloseRedis 关闭 Redis 连接
func CloseRedis() error {
	if RDB != nil {
		return RDB.Close()
	}
	return nil
}

// RedisHelper Redis 辅助方法
type RedisHelper struct{}

// Set 设置键值对
func (r *RedisHelper) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return RDB.Set(ctx, key, value, expiration).Err()
}

// Get 获取值
func (r *RedisHelper) Get(ctx context.Context, key string) (string, error) {
	return RDB.Get(ctx, key).Result()
}

// Del 删除键
func (r *RedisHelper) Del(ctx context.Context, keys ...string) error {
	return RDB.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func (r *RedisHelper) Exists(ctx context.Context, keys ...string) (int64, error) {
	return RDB.Exists(ctx, keys...).Result()
}

// Expire 设置过期时间
func (r *RedisHelper) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return RDB.Expire(ctx, key, expiration).Err()
}

// TTL 获取剩余过期时间
func (r *RedisHelper) TTL(ctx context.Context, key string) (time.Duration, error) {
	return RDB.TTL(ctx, key).Result()
}

var Redis = &RedisHelper{}
