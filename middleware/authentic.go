package middleware

import (
	"strings"

	"github.com/demo/common"
	"github.com/demo/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware JWT 认证中间件
func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从 Header 中获取 Token
		// 支持两种格式:
		// 1. Authorization: Bearer <token>
		// 2. X-Token: <token>
		token := ctx.GetHeader("X-Token")
		if token == "" {
			authHeader := ctx.GetHeader("Authorization")
			if authHeader != "" && strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		if token == "" {
			utils.Fail(ctx, common.CodeTokenMissing)
			ctx.Abort()
			return
		}

		// 解析 Token
		claims, err := utils.ParseJWT(token)
		if err != nil {
			utils.Fail(ctx, common.CodeTokenInvalid)
			ctx.Abort()
			return
		}

		// 设置用户信息到 ctx，后续业务可直接读取
		ctx.Set("user_id", claims.UserID)
		ctx.Set("username", claims.Username)
		ctx.Set("claims", claims)

		ctx.Next()
	}
}

// GetUserID 从 Context 中获取用户 ID
func GetUserID(ctx *gin.Context) int64 {
	if userID, exists := ctx.Get("user_id"); exists {
		if uid, ok := userID.(int64); ok {
			return uid
		}
	}
	return 0
}

// GetUsername 从 Context 中获取用户名
func GetUsername(ctx *gin.Context) string {
	if username, exists := ctx.Get("username"); exists {
		if name, ok := username.(string); ok {
			return name
		}
	}
	return ""
}

// GetClaims 从 Context 中获取完整的 Claims
func GetClaims(ctx *gin.Context) *utils.Claims {
	if claims, exists := ctx.Get("claims"); exists {
		if c, ok := claims.(*utils.Claims); ok {
			return c
		}
	}
	return nil
}
