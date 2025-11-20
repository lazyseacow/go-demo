package middleware

import (
	"sync"
	"time"

	"github.com/demo/common"
	"github.com/demo/utils"
	"github.com/gin-gonic/gin"
)

// RateLimiter 简单的内存限流器
type RateLimiter struct {
	rate     int                 // 每秒允许的请求数
	visitors map[string]*visitor // IP -> 访问者
	mu       sync.RWMutex
}

type visitor struct {
	lastSeen time.Time
	count    int
}

// NewRateLimiter 创建限流器
func NewRateLimiter(rate int) *RateLimiter {
	rl := &RateLimiter{
		rate:     rate,
		visitors: make(map[string]*visitor),
	}

	// 定期清理过期的访问者
	go rl.cleanupVisitors()

	return rl
}

// cleanupVisitors 清理过期的访问者
func (rl *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)
		rl.mu.Lock()
		for ip, v := range rl.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(rl.visitors, ip)
			}
		}
		rl.mu.Unlock()
	}
}

// Allow 检查是否允许访问
func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if !exists {
		rl.visitors[ip] = &visitor{
			lastSeen: time.Now(),
			count:    1,
		}
		return true
	}

	// 如果距离上次请求超过 1 秒，重置计数
	if time.Since(v.lastSeen) > time.Second {
		v.count = 1
		v.lastSeen = time.Now()
		return true
	}

	// 检查是否超过限制
	if v.count >= rl.rate {
		return false
	}

	v.count++
	v.lastSeen = time.Now()
	return true
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(rate int) gin.HandlerFunc {
	limiter := NewRateLimiter(rate)

	return func(ctx *gin.Context) {
		ip := ctx.ClientIP()

		if !limiter.Allow(ip) {
			utils.Fail(ctx, common.CodeTooManyRequests)
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
