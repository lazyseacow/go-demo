package controllers

import (
	"context"
	"time"

	"github.com/demo/database"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// HealthController 健康检查控制器
type HealthController struct {
	*BaseController
}

// NewHealthController 创建健康检查控制器实例
func NewHealthController() *HealthController {
	return &HealthController{
		BaseController: NewBaseController(),
	}
}

// ServiceStatus 服务状态
type ServiceStatus struct {
	Status  string `json:"status"`  // healthy, unhealthy, unknown
	Message string `json:"message"` // 状态描述
	Latency string `json:"latency"` // 响应延迟
}

// HealthResponse 健康检查响应
type HealthResponse struct {
	Status    string                    `json:"status"`    // healthy, degraded, unhealthy
	Timestamp int64                     `json:"timestamp"` // 时间戳
	Services  map[string]ServiceStatus  `json:"services"`  // 各服务状态
	Version   string                    `json:"version"`   // 版本号
}

// Check 健康检查
// @Summary      健康检查
// @Description  检查系统和所有依赖服务的健康状态
// @Tags         系统
// @Accept       json
// @Produce      json
// @Success      200  {object}  HealthResponse  "系统健康"
// @Failure      503  {object}  HealthResponse  "系统异常"
// @Router       /health [get]
func (ctrl *HealthController) Check(ctx *gin.Context) {
	services := make(map[string]ServiceStatus)

	// 检查 MySQL
	mysqlStatus := ctrl.checkMySQL()
	services["mysql"] = mysqlStatus

	// 检查 MongoDB
	mongoStatus := ctrl.checkMongoDB()
	services["mongodb"] = mongoStatus

	// 检查 Redis
	redisStatus := ctrl.checkRedis()
	services["redis"] = redisStatus

	// 判断整体状态
	overallStatus := ctrl.getOverallStatus(services)

	response := HealthResponse{
		Status:    overallStatus,
		Timestamp: time.Now().Unix(),
		Services:  services,
		Version:   "2.0.0",
	}

	// 根据整体状态返回不同的 HTTP 状态码
	if overallStatus == "healthy" {
		ctx.JSON(200, response)
	} else if overallStatus == "degraded" {
		ctx.JSON(200, response) // 部分服务异常，但系统可用
	} else {
		ctx.JSON(503, response) // 关键服务异常，系统不可用
	}
}

// Ping 简单健康检查
// @Summary      简单健康检查
// @Description  返回简单的 pong 响应，用于快速检查服务是否在线
// @Tags         系统
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string  "服务在线"
// @Router       /ping [get]
func (ctrl *HealthController) Ping(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"message": "pong",
		"time":    time.Now().Format("2006-01-02 15:04:05"),
	})
}

// checkMySQL 检查 MySQL 连接
func (ctrl *HealthController) checkMySQL() ServiceStatus {
	start := time.Now()

	// 使用 defer + recover 捕获可能的 panic
	defer func() {
		if r := recover(); r != nil {
			// 数据库未初始化时可能 panic
		}
	}()

	// 直接访问 DB 变量，避免 panic
	db := database.DB
	if db == nil {
		return ServiceStatus{
			Status:  "unknown",
			Message: "数据库未初始化",
			Latency: "0ms",
		}
	}

	sqlDB, err := db.DB()
	if err != nil {
		return ServiceStatus{
			Status:  "unhealthy",
			Message: "获取数据库实例失败: " + err.Error(),
			Latency: time.Since(start).String(),
		}
	}

	// Ping 数据库
	if err := sqlDB.Ping(); err != nil {
		return ServiceStatus{
			Status:  "unhealthy",
			Message: "数据库连接失败: " + err.Error(),
			Latency: time.Since(start).String(),
		}
	}

	// 执行简单查询测试
	if err := db.Raw("SELECT 1").Error; err != nil {
		return ServiceStatus{
			Status:  "unhealthy",
			Message: "数据库查询失败: " + err.Error(),
			Latency: time.Since(start).String(),
		}
	}

	latency := time.Since(start)
	return ServiceStatus{
		Status:  "healthy",
		Message: "MySQL 运行正常",
		Latency: latency.String(),
	}
}

// checkMongoDB 检查 MongoDB 连接
func (ctrl *HealthController) checkMongoDB() ServiceStatus {
	start := time.Now()

	// 使用 defer + recover 捕获 panic
	defer func() {
		if r := recover(); r != nil {
			// MongoDB 未初始化时会 panic，这里捕获
		}
	}()

	// 直接访问 MongoDB 变量，避免 panic
	client := database.MongoDB
	if client == nil {
		return ServiceStatus{
			Status:  "unknown",
			Message: "MongoDB 未初始化（可选服务）",
			Latency: "0ms",
		}
	}

	// Ping MongoDB
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return ServiceStatus{
			Status:  "unhealthy",
			Message: "MongoDB 连接失败: " + err.Error(),
			Latency: time.Since(start).String(),
		}
	}

	latency := time.Since(start)
	return ServiceStatus{
		Status:  "healthy",
		Message: "MongoDB 运行正常",
		Latency: latency.String(),
	}
}

// checkRedis 检查 Redis 连接
func (ctrl *HealthController) checkRedis() ServiceStatus {
	start := time.Now()

	// 使用 defer + recover 捕获可能的 panic
	defer func() {
		if r := recover(); r != nil {
			// Redis 未初始化时可能 panic
		}
	}()

	// 直接访问 RDB 变量，避免 panic
	rdb := database.RDB
	if rdb == nil {
		return ServiceStatus{
			Status:  "unknown",
			Message: "Redis 未初始化",
			Latency: "0ms",
		}
	}

	// Ping Redis
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return ServiceStatus{
			Status:  "unhealthy",
			Message: "Redis 连接失败: " + err.Error(),
			Latency: time.Since(start).String(),
		}
	}

	latency := time.Since(start)
	return ServiceStatus{
		Status:  "healthy",
		Message: "Redis 运行正常",
		Latency: latency.String(),
	}
}

// getOverallStatus 获取整体健康状态
func (ctrl *HealthController) getOverallStatus(services map[string]ServiceStatus) string {
	healthyCount := 0
	unhealthyCount := 0
	unknownCount := 0

	// 统计各服务状态
	for _, service := range services {
		switch service.Status {
		case "healthy":
			healthyCount++
		case "unhealthy":
			unhealthyCount++
		case "unknown":
			unknownCount++
		}
	}

	// MySQL 和 Redis 是关键服务，必须健康
	mysqlStatus := services["mysql"].Status
	redisStatus := services["redis"].Status

	// 如果关键服务异常，系统不可用
	if mysqlStatus == "unhealthy" || redisStatus == "unhealthy" {
		return "unhealthy"
	}

	// 如果有非关键服务异常，系统降级
	if unhealthyCount > 0 || unknownCount > 0 {
		return "degraded"
	}

	// 所有服务都正常
	return "healthy"
}

// Ready 就绪检查（Kubernetes readiness probe）
// @Summary      就绪检查
// @Description  检查服务是否已准备好接收流量
// @Tags         系统
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string  "服务就绪"
// @Failure      503  {object}  map[string]string  "服务未就绪"
// @Router       /ready [get]
func (ctrl *HealthController) Ready(ctx *gin.Context) {
	// 检查关键服务
	mysqlStatus := ctrl.checkMySQL()
	redisStatus := ctrl.checkRedis()

	// MySQL 和 Redis 必须健康才能提供服务
	if mysqlStatus.Status == "healthy" && redisStatus.Status == "healthy" {
		ctx.JSON(200, gin.H{
			"status":  "ready",
			"message": "服务已就绪",
		})
	} else {
		ctx.JSON(503, gin.H{
			"status":  "not_ready",
			"message": "服务未就绪",
			"mysql":   mysqlStatus.Status,
			"redis":   redisStatus.Status,
		})
	}
}

// Live 存活检查（Kubernetes liveness probe）
// @Summary      存活检查
// @Description  检查服务是否存活
// @Tags         系统
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]string  "服务存活"
// @Router       /live [get]
func (ctrl *HealthController) Live(ctx *gin.Context) {
	ctx.JSON(200, gin.H{
		"status":  "alive",
		"message": "服务运行中",
		"time":    time.Now().Unix(),
	})
}

