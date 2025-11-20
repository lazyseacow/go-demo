package middleware

import (
	"time"

	"github.com/demo/utils"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// LoggerMiddleware 基于 Zap 的日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery

		// 处理请求
		ctx.Next()

		// 计算耗时
		cost := time.Since(start)

		// 记录日志
		utils.LogInfo("HTTP Request",
			zap.String("method", ctx.Request.Method),
			zap.String("path", path),
			zap.String("query", query),
			zap.String("ip", ctx.ClientIP()),
			zap.String("user-agent", ctx.Request.UserAgent()),
			zap.Int("status", ctx.Writer.Status()),
			zap.Duration("latency", cost),
			zap.String("error", ctx.Errors.ByType(gin.ErrorTypePrivate).String()),
		)
	}
}

// RecoveryMiddleware Zap 版本的恢复中间件
func RecoveryMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				utils.LogError("Panic recovered",
					zap.Any("error", err),
					zap.String("method", ctx.Request.Method),
					zap.String("path", ctx.Request.URL.Path),
					zap.String("ip", ctx.ClientIP()),
				)

				ctx.AbortWithStatus(500)
			}
		}()
		ctx.Next()
	}
}
