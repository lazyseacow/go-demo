package routes

import (
	"github.com/demo/controllers"
	"github.com/demo/middleware"
	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine) {
	// 初始化控制器
	healthCtrl := controllers.NewHealthController()
	authCtrl := controllers.NewAuthController()
	userCtrl := controllers.NewUserController()
	articleCtrl := controllers.NewArticleController()

	// 健康检查接口（无需认证，用于监控和 K8s 探针）
	r.GET("/ping", healthCtrl.Ping)    // 简单健康检查
	r.GET("/health", healthCtrl.Check) // 完整健康检查（检查所有服务）
	r.GET("/ready", healthCtrl.Ready)  // 就绪检查（Kubernetes readiness probe）
	r.GET("/live", healthCtrl.Live)    // 存活检查（Kubernetes liveness probe）

	// API v1 路由组
	v1 := r.Group("/api/v1")
	{
		// 公开接口（不需要认证）
		public := v1.Group("/auth")
		{
			public.POST("/register", authCtrl.Register) // 注册
			public.POST("/login", authCtrl.Login)       // 登录
		}

		// 公开的文章接口
		publicArticles := v1.Group("/articles")
		{
			publicArticles.GET("", articleCtrl.GetArticleList)     // 获取文章列表
			publicArticles.GET("/:id", articleCtrl.GetArticleByID) // 获取文章详情
		}

		// 需要认证的接口
		authorized := v1.Group("")
		authorized.Use(middleware.AuthMiddleware())
		{
			// 认证相关
			auth := authorized.Group("/auth")
			{
				auth.POST("/logout", authCtrl.Logout) // 登出
			}

			// 用户管理
			user := authorized.Group("/users")
			{
				user.GET("", userCtrl.GetUserList)            // 获取用户列表
				user.GET("/:id", userCtrl.GetUserByID)        // 获取指定用户
				user.GET("/me", userCtrl.GetCurrentUser)      // 获取当前用户信息
				user.POST("/update", userCtrl.UpdateUser)     // 更新用户信息
				user.POST("/:id/delete", userCtrl.DeleteUser) // 删除用户
			}

			// 文章管理（需要认证）
			articles := authorized.Group("/articles")
			{
				articles.POST("", articleCtrl.CreateArticle)            // 创建文章
				articles.POST("/:id/update", articleCtrl.UpdateArticle) // 更新文章
				articles.POST("/:id/delete", articleCtrl.DeleteArticle) // 删除文章
				articles.POST("/:id/like", articleCtrl.LikeArticle)     // 点赞文章
			}
		}
	}

}
