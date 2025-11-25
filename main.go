package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/demo/config"
	"github.com/demo/database"
	_ "github.com/demo/docs" // Swagger 文档
	"github.com/demo/middleware"
	"github.com/demo/models"
	"github.com/demo/routes"
	"github.com/demo/utils"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

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
	if err := config.LoadConfig("config.yaml"); err != nil {
		log.Fatalf("❌ 加载配置失败: %v", err)
	}

	cfg := config.GetConfig()

	// 2. 初始化日志系统
	if err := utils.InitLogger(); err != nil {
		log.Fatalf("❌ 初始化日志失败: %v", err)
	}
	defer utils.Sync()
	utils.LogInfo("✅ 日志系统初始化成功")

	// 3. 设置 Gin 模式
	gin.SetMode(cfg.Server.Mode)

	// 4. 初始化数据库连接
	if err := database.InitMySQL(); err != nil {
		utils.LogFatalf("❌ 初始化 MySQL 失败: %v", err)
	}
	defer database.CloseMySQL()

	// 5. 初始化 MongoDB 连接（可选）
	if err := database.InitMongoDB(); err != nil {
		utils.LogWarnf("⚠️  初始化 MongoDB 失败: %v（跳过）", err)
	} else {
		defer database.CloseMongoDB()
	}

	// 6. 初始化 Redis 连接
	if err := database.InitRedis(); err != nil {
		utils.LogFatalf("❌ 初始化 Redis 失败: %v", err)
	}
	defer database.CloseRedis()

	// 7. 自动迁移数据库表
	if err := autoMigrate(); err != nil {
		utils.LogFatalf("❌ 数据库迁移失败: %v", err)
	}

	// 8. 创建 Gin 引擎
	r := gin.New()

	// 9. 注册全局中间件
	r.Use(middleware.RecoveryMiddleware())     // Panic 恢复
	r.Use(middleware.CORSMiddleware())         // 跨域
	r.Use(middleware.LoggerMiddleware())       // 日志
	r.Use(middleware.RateLimitMiddleware(100)) // 限流：100 req/s

	// 10. 注册 Swagger 文档
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// 11. 注册路由
	routes.SetupRoutes(r)

	// 12. 启动服务器（支持优雅关闭）
	startServer(r, cfg)
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

// Server 接口定义了服务器的基本行为
// 任何实现了 ListenAndServe 和 Shutdown 方法的类型都可以被视为 Server
type Server interface {
	ListenAndServe() error // 启动服务器
	Shutdown() error       // 关闭服务器
}

// ServerStarter 是一个服务器启动器结构体
// 它使用了组合（composition）模式，内嵌了一个 Server 接口
// 这样可以灵活地传入不同的 Server 实现
type ServerStarter struct {
	Server Server // 持有一个 Server 接口类型的字段，实现了依赖注入
}

// ListenAndServe 是 ServerStarter 的方法
// 它委托（delegate）给内部的 Server 来执行实际的启动操作
// 这是装饰器模式的一种应用
func (s *ServerStarter) ListenAndServe() error {
	return s.Server.ListenAndServe()
}

// Shutdown 是 ServerStarter 的方法
// 它委托给内部的 Server 来执行实际的关闭操作
func (s *ServerStarter) Shutdown() error {
	return s.Server.Shutdown()
}

// GinServer 代表一个使用 Gin 框架的服务器
// 它是一个空结构体，因为这里只是演示，不需要存储状态
type GinServer struct {
}

// ListenAndServe 实现了 Server 接口的 ListenAndServe 方法
// 这使得 GinServer 成为了 Server 接口的一个实现
func (g *GinServer) ListenAndServe() error {
	fmt.Println("启动一个gin的服务")
	return nil
}

// Shutdown 实现了 Server 接口的 Shutdown 方法
func (g *GinServer) Shutdown() error {
	fmt.Println("关闭一个gin的服务")
	return nil
}

// NginxServer 代表一个 Nginx 服务器
// 它也实现了 Server 接口
type NginxServer struct {
}

// ListenAndServe 实现了 Server 接口的 ListenAndServe 方法
// 这使得 NginxServer 也成为了 Server 接口的一个实现
func (n *NginxServer) ListenAndServe() error {
	fmt.Println("启动一个nginx的服务")
	return nil
}

// Shutdown 实现了 Server 接口的 Shutdown 方法
func (n *NginxServer) Shutdown() error {
	fmt.Println("关闭一个nginx的服务")
	return nil
}

// NativeHTTPServer 代表一个 Go 原生的 HTTP 服务器
// 它使用 Go 标准库的 net/http 包实现
// 相比 GinServer 和 NginxServer 的简单示例，这是一个完整的实现
type NativeHTTPServer struct {
	server *http.Server // 存储 http.Server 实例，用于实际的服务器管理
	addr   string       // 监听地址，例如 ":8081"
}

// NewNativeHTTPServer 创建一个新的原生 HTTP 服务器实例
// addr: 监听地址，格式如 ":8081" 或 "localhost:8081"
// handler: HTTP 请求处理器，如果为 nil 则使用 http.DefaultServeMux
func NewNativeHTTPServer(addr string, handler http.Handler) *NativeHTTPServer {
	if handler == nil {
		// 如果没有提供处理器，使用默认的多路复用器
		handler = http.DefaultServeMux
	}

	return &NativeHTTPServer{
		addr: addr,
		server: &http.Server{
			Addr:         addr,
			Handler:      handler,
			ReadTimeout:  10 * time.Second, // 设置读取超时
			WriteTimeout: 10 * time.Second, // 设置写入超时
			IdleTimeout:  60 * time.Second, // 设置空闲连接超时
		},
	}
}

// ListenAndServe 实现了 Server 接口的 ListenAndServe 方法
// 启动原生 HTTP 服务器并开始监听请求
// 这是一个阻塞调用，会一直运行直到服务器关闭
func (n *NativeHTTPServer) ListenAndServe() error {
	fmt.Printf("🚀 启动 Go 原生 HTTP 服务器，监听地址: %s\n", n.addr)

	// 注册一个简单的测试路由
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello from Native HTTP Server! 🎉\n")
		fmt.Fprintf(w, "Path: %s\n", r.URL.Path)
		fmt.Fprintf(w, "Method: %s\n", r.Method)
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"status": "healthy", "server": "native-http"}`)
	})

	// 启动服务器（阻塞式调用）
	if err := n.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		return fmt.Errorf("HTTP 服务器启动失败: %v", err)
	}

	return nil
}

// Shutdown 实现了 Server 接口的 Shutdown 方法
// 优雅地关闭服务器，等待现有连接处理完成
// 注意：这是一个阻塞调用，最多等待 5 秒
func (n *NativeHTTPServer) Shutdown() error {
	fmt.Println("⏳ 正在关闭 Go 原生 HTTP 服务器...")

	// 创建一个带超时的上下文，给服务器 5 秒时间完成现有请求
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := n.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("HTTP 服务器关闭失败: %v", err)
	}

	fmt.Println("✅ Go 原生 HTTP 服务器已成功关闭")
	return nil
}

// main1 演示了接口的多态性和组合模式的使用
// 支持优雅关闭：按 Ctrl+C 触发关闭流程
func main1() {
	// 创建第一个启动器，注入 GinServer 实例
	// 这里体现了依赖注入的思想：ServerStarter 不关心具体是什么服务器
	starter1 := &ServerStarter{Server: &GinServer{}}

	// 创建第二个启动器，注入 NginxServer 实例
	// 同样的 ServerStarter 可以启动不同类型的服务器
	starter2 := &ServerStarter{Server: &NginxServer{}}

	// 创建第三个启动器，注入 NativeHTTPServer 实例
	// 这是一个完整的 Go 原生 HTTP 服务器实现
	// 参数说明：":8081" 是监听地址，nil 表示使用默认的路由处理器
	nativeServer := NewNativeHTTPServer(":8081", nil)
	starter3 := &ServerStarter{Server: nativeServer}

	// 启动三个服务器，虽然调用的是相同的方法，但执行的是不同的实现
	// 这展示了接口的强大之处：同一个接口，不同的实现
	starter1.ListenAndServe() // 输出：启动一个gin的服务
	starter2.ListenAndServe() // 输出：启动一个nginx的服务

	// 演示：在 goroutine 中启动原生 HTTP 服务器
	// 因为 ListenAndServe 是阻塞调用，我们需要在单独的 goroutine 中运行
	go func() {
		if err := starter3.ListenAndServe(); err != nil {
			fmt.Printf("❌ 原生 HTTP 服务器错误: %v\n", err)
		}
	}()

	// 等待 2 秒，让服务器完全启动
	time.Sleep(2 * time.Second)
	fmt.Println("\n✅ 所有服务器已启动！")
	fmt.Println("📌 访问 http://localhost:8081 测试原生 HTTP 服务器")
	fmt.Println("📌 访问 http://localhost:8081/health 查看健康状态")
	fmt.Println("📌 按 Ctrl+C 可以优雅地关闭所有服务器")

	// ========== 优雅关闭机制 ==========
	// 创建一个信号通道，用于接收操作系统信号
	quit := make(chan os.Signal, 1)

	// 注册要监听的信号：
	// - syscall.SIGINT: 中断信号（Ctrl+C）
	// - syscall.SIGTERM: 终止信号（kill 命令默认信号）
	// 当收到这些信号时，会发送到 quit 通道
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// 阻塞等待关闭信号
	// 程序会一直运行，直到收到 SIGINT 或 SIGTERM 信号
	<-quit

	// 收到关闭信号后，开始优雅关闭流程
	fmt.Println("\n⏳ 收到关闭信号，开始优雅关闭所有服务器...")

	// 按顺序关闭所有服务器
	// 即使某个服务器关闭失败，也会继续关闭其他服务器
	if err := starter1.Shutdown(); err != nil {
		fmt.Printf("⚠️  关闭 GinServer 失败: %v\n", err)
	}

	if err := starter2.Shutdown(); err != nil {
		fmt.Printf("⚠️  关闭 NginxServer 失败: %v\n", err)
	}

	if err := starter3.Shutdown(); err != nil {
		fmt.Printf("⚠️  关闭 NativeHTTPServer 失败: %v\n", err)
	}

	fmt.Println("\n✅ 所有服务器已优雅关闭，程序退出")
}

func test(a, b, c int, d, e, f string) (aa, bb error) {
	aa = fmt.Errorf("error")
	bb = fmt.Errorf("error")
	return aa, bb
}
