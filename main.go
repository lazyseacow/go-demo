package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/demo/config"
	"github.com/demo/database"
	_ "github.com/demo/docs" // Swagger 文档
	"github.com/demo/models"
	"github.com/demo/utils"

	"github.com/gin-gonic/gin"
)

// @title           Go-Demo API
// @version         2.0
// @description     这是一个基于 Go + Gin 的企业级 Web 框架 API 文档
// @description     支持 MySQL、MongoDB、Redis、JWT 认证
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@example.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8080
// @BasePath  /api/v1

// @securityDefinitions.apikey  ApiKeyAuth
// @in                          header
// @name                        X-Token
// @description                 JWT Token 认证，格式: Bearer {token} 或直接填写 token

// @securityDefinitions.apikey  BearerAuth
// @in                          header
// @name                        Authorization
// @description                 JWT Token 认证，格式: Bearer {token}

func main() {
	// 1. 加载配置
	// if err := config.LoadConfig("config.yaml"); err != nil {
	// 	log.Fatalf("❌ 加载配置失败: %v", err)
	// }

	// cfg := config.GetConfig()

	// // 2. 初始化日志系统
	// if err := utils.InitLogger(); err != nil {
	// 	log.Fatalf("❌ 初始化日志失败: %v", err)
	// }
	// defer utils.Sync()
	// utils.LogInfo("✅ 日志系统初始化成功")

	// // 3. 设置 Gin 模式
	// gin.SetMode(cfg.Server.Mode)

	// // 4. 初始化数据库连接
	// if err := database.InitMySQL(); err != nil {
	// 	utils.LogFatalf("❌ 初始化 MySQL 失败: %v", err)
	// }
	// defer database.CloseMySQL()

	// // 5. 初始化 MongoDB 连接（可选）
	// if err := database.InitMongoDB(); err != nil {
	// 	utils.LogWarnf("⚠️  初始化 MongoDB 失败: %v（跳过）", err)
	// } else {
	// 	defer database.CloseMongoDB()
	// }

	// // 6. 初始化 Redis 连接
	// if err := database.InitRedis(); err != nil {
	// 	utils.LogFatalf("❌ 初始化 Redis 失败: %v", err)
	// }
	// defer database.CloseRedis()

	// // 7. 自动迁移数据库表
	// if err := autoMigrate(); err != nil {
	// 	utils.LogFatalf("❌ 数据库迁移失败: %v", err)
	// }

	// // 8. 创建 Gin 引擎
	// r := gin.New()

	// // 9. 注册全局中间件
	// r.Use(middleware.RecoveryMiddleware())     // Panic 恢复
	// r.Use(middleware.CORSMiddleware())         // 跨域
	// r.Use(middleware.LoggerMiddleware())       // 日志
	// r.Use(middleware.RateLimitMiddleware(100)) // 限流：100 req/s

	// // 10. 注册 Swagger 文档
	// r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// // 11. 注册路由
	// routes.SetupRoutes(r)

	// 12. 启动服务器（支持优雅关闭）
	// startServer(r, cfg)

	main1()
}

// autoMigrate 自动迁移数据库表
func autoMigrate() error {
	db := database.GetDB()

	// 自动迁移表结构
	if err := db.AutoMigrate(
		&models.User{},
		// 在这里添加更多模型
	); err != nil {
		return err
	}

	utils.LogInfo("✅ 数据库表迁移成功")
	return nil
}

// startServer 启动服务器（支持优雅关闭）
func startServer(r *gin.Engine, cfg *config.Config) {
	srv := &http.Server{
		Addr:         cfg.Server.Port,
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	// 在 goroutine 中启动服务器
	go func() {
		utils.LogInfof("🚀 服务器启动成功，监听端口: %s", cfg.Server.Port)
		utils.LogInfof("📖 API 文档: http://localhost%s/swagger/index.html", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			utils.LogFatalf("❌ 服务器启动失败: %v", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	utils.LogInfo("⏳ 正在关闭服务器...")

	// 设置 5 秒的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		utils.LogFatalf("❌ 服务器强制关闭: %v", err)
	}

	utils.LogInfo("✅ 服务器已优雅退出")
}

type Server interface {
	ListenAndServe() error
	Shutdown() error
}

type ServerStarter struct {
	Server Server
}

func (s *ServerStarter) ListenAndServe() error {
	return s.Server.ListenAndServe()
}

func (s *ServerStarter) Shutdown() error {
	return s.Server.Shutdown()
}

type GinServer struct {
}

func (g *GinServer) ListenAndServe() error {
	fmt.Println("启动一个gin的服务")
	return nil
}

func (g *GinServer) Shutdown() error {
	fmt.Println("关闭一个gin的服务")
	return nil
}

type NginxServer struct {
}

func (n *NginxServer) ListenAndServe() error {
	fmt.Println("启动一个nginx的服务")
	return nil
}

func (n *NginxServer) Shutdown() error {
	fmt.Println("关闭一个nginx的服务")
	return nil
}

func main1() {
	starter1 := &ServerStarter{Server: &GinServer{}}
	starter2 := &ServerStarter{Server: &NginxServer{}}
	starter1.ListenAndServe()
	starter2.ListenAndServe()
	starter1.Shutdown()
	starter2.Shutdown()
}
