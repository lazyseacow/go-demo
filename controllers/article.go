package controllers

import (
	"sync"
	"context"
	"net/http"

	"github.com/demo/common"
	"github.com/demo/database"
	"github.com/demo/middleware"
	"github.com/demo/models"
	"github.com/demo/service"
	"github.com/demo/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	_articleController *ArticleController
	_onceNewArticleController sync.Once
)

type ArticleController struct {
	*BaseController
	articleService *service.ArticleService
}

func NewArticleController() *ArticleController {
	_onceNewArticleController.Do(func() {
		_articleController = &ArticleController{
			BaseController: NewBaseController(),
			articleService: service.NewArticleService(),
		}
	})
	return _articleController
}

// CreateArticle 创建文章
// @Summary      创建文章
// @Description  创建新文章（存储在 MongoDB）
// @Tags         文章管理
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request  body      models.CreateArticleRequest  true  "文章内容"
// @Success      200      {object}  utils.Response{data=object{id=string}}  "创建成功"
// @Failure      14001    {object}  utils.Response  "参数错误"
// @Router       /articles [post]
func (ctrl *ArticleController) CreateArticle(ctx *gin.Context) {
	// 检查 MongoDB 是否可用
	if database.MongoDB == nil {
		utils.Fail(ctx, http.StatusServiceUnavailable, "文章服务暂时不可用，请稍后再试")
		return
	}

	// 使用专门的请求 DTO
	var req models.CreateArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Fail(ctx, common.CodeParamInvalid, "参数错误: "+err.Error())
		return
	}

	// 获取当前用户 ID
	userID := middleware.GetUserID(ctx)
	username := middleware.GetUsername(ctx)

	// 调用 Service 层处理业务逻辑
	article, err := ctrl.articleService.CreateArticle(req, userID, username)
	if err != nil {
		if customErr, ok := err.(*common.CustomError); ok {
			utils.Error(ctx, customErr)
			return
		}
		utils.Fail(ctx, common.CodeInternalError, err.Error())
		return
	}

	utils.Success(ctx, gin.H{
		"id": article.ID.Hex(),
	}, "创建成功")
}

// GetArticleList 获取文章列表（分页）
// @Summary      获取文章列表
// @Description  分页获取文章列表，支持搜索和筛选
// @Tags         文章管理
// @Accept       json
// @Produce      json
// @Param        page       query     int     false  "页码"  default(1)
// @Param        page_size  query     int     false  "每页数量"  default(10)
// @Param        keyword    query     string  false  "搜索关键词"
// @Param        status     query     int     false  "文章状态"
// @Success      200        {object}  utils.Response{data=models.PageResponse}  "获取成功"
// @Router       /articles [get]
func (ctrl *ArticleController) GetArticleList(ctx *gin.Context) {
	// 检查 MongoDB 是否可用
	if database.MongoDB == nil {
		utils.Fail(ctx, http.StatusServiceUnavailable, "文章服务暂时不可用，请稍后再试")
		return
	}

	var req models.ArticleQueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		req.Page = 1
		req.PageSize = 10
	}

	// 调用 Service 层处理业务逻辑
	result, err := ctrl.articleService.GetArticleList(req)
	if err != nil {
		if customErr, ok := err.(*common.CustomError); ok {
			utils.Error(ctx, customErr)
			return
		}
		utils.Fail(ctx, common.CodeInternalError, err.Error())
		return
	}

	utils.Success(ctx, result)
}

// GetArticleByID 根据 ID 获取文章
// @Summary      获取文章详情
// @Description  根据文章 ID 获取文章详细信息
// @Tags         文章管理
// @Accept       json
// @Produce      json
// @Param        id   path      string  true  "文章ID"
// @Success      200  {object}  utils.Response{data=models.Article}  "获取成功"
// @Failure      14001  {object}  utils.Response  "无效的文章 ID"
// @Failure      12001  {object}  utils.Response  "文章不存在"
// @Router       /articles/{id} [get]
func (ctrl *ArticleController) GetArticleByID(ctx *gin.Context) {
	// 检查 MongoDB 是否可用
	if database.MongoDB == nil {
		utils.Fail(ctx, http.StatusServiceUnavailable, "文章服务暂时不可用，请稍后再试")
		return
	}

	idStr := ctx.Param("id")

	// 调用 Service 层处理业务逻辑
	article, err := ctrl.articleService.GetArticleByID(idStr)
	if err != nil {
		if customErr, ok := err.(*common.CustomError); ok {
			utils.Error(ctx, customErr)
			return
		}
		utils.Fail(ctx, common.CodeInternalError, err.Error())
		return
	}

	utils.Success(ctx, article)
}

// UpdateArticle 更新文章
// @Summary      更新文章
// @Description  更新指定文章的内容
// @Tags         文章管理
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id       path      string  true  "文章ID"
// @Param        request  body      models.UpdateArticleRequest  true  "更新参数"
// @Success      200      {object}  utils.Response  "更新成功"
// @Failure      14001    {object}  utils.Response  "无效的文章 ID"
// @Failure      12001    {object}  utils.Response  "文章不存在或无权限"
// @Router       /articles/{id}/update [post]
func (ctrl *ArticleController) UpdateArticle(ctx *gin.Context) {
	// 检查 MongoDB 是否可用
	if database.MongoDB == nil {
		utils.Fail(ctx, http.StatusServiceUnavailable, "文章服务暂时不可用，请稍后再试")
		return
	}

	idStr := ctx.Param("id")

	// 使用专门的更新请求 DTO
	var req models.UpdateArticleRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Fail(ctx, common.CodeParamInvalid, "参数错误: "+err.Error())
		return
	}

	// 获取当前用户 ID
	userID := middleware.GetUserID(ctx)

	// 调用 Service 层处理业务逻辑
	if err := ctrl.articleService.UpdateArticle(idStr, userID, req); err != nil {
		if customErr, ok := err.(*common.CustomError); ok {
			utils.Error(ctx, customErr)
			return
		}
		utils.Fail(ctx, common.CodeInternalError, err.Error())
		return
	}

	utils.Success(ctx, nil, "更新成功")
}

// DeleteArticle 删除文章
// @Summary      删除文章
// @Description  删除指定文章
// @Tags         文章管理
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      string  true  "文章ID"
// @Success      200  {object}  utils.Response  "删除成功"
// @Failure      14001  {object}  utils.Response  "无效的文章 ID"
// @Failure      12001  {object}  utils.Response  "文章不存在或无权限"
// @Router       /articles/{id}/delete [post]
func (ctrl *ArticleController) DeleteArticle(ctx *gin.Context) {
	// 检查 MongoDB 是否可用
	if database.MongoDB == nil {
		utils.Fail(ctx, http.StatusServiceUnavailable, "文章服务暂时不可用，请稍后再试")
		return
	}

	idStr := ctx.Param("id")

	// 获取当前用户 ID
	userID := middleware.GetUserID(ctx)

	// 调用 Service 层处理业务逻辑
	if err := ctrl.articleService.DeleteArticle(idStr, userID); err != nil {
		if customErr, ok := err.(*common.CustomError); ok {
			utils.Error(ctx, customErr)
			return
		}
		utils.Fail(ctx, common.CodeInternalError, err.Error())
		return
	}

	utils.Success(ctx, nil, "删除成功")
}

// LikeArticle 点赞文章
// @Summary      点赞文章
// @Description  为指定文章点赞（点赞数+1）
// @Tags         文章管理
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        id   path      string  true  "文章ID"
// @Success      200  {object}  utils.Response  "点赞成功"
// @Failure      14001  {object}  utils.Response  "无效的文章 ID"
// @Failure      12001  {object}  utils.Response  "文章不存在"
// @Router       /articles/{id}/like [post]
func (ctrl *ArticleController) LikeArticle(ctx *gin.Context) {
	// 检查 MongoDB 是否可用
	if database.MongoDB == nil {
		utils.Fail(ctx, http.StatusServiceUnavailable, "文章服务暂时不可用，请稍后再试")
		return
	}

	idStr := ctx.Param("id")

	// 调用 Service 层处理业务逻辑
	if err := ctrl.articleService.LikeArticle(idStr); err != nil {
		if customErr, ok := err.(*common.CustomError); ok {
			utils.Error(ctx, customErr)
			return
		}
		utils.Fail(ctx, common.CodeInternalError, err.Error())
		return
	}

	utils.Success(ctx, nil, "点赞成功")
}

// CreateArticleIndexes 创建文章索引（初始化时调用）
func CreateArticleIndexes() error {
	collection := database.GetMongoCollection("articles")

	// 创建索引
	indexes := []struct {
		keys    bson.D
		options *options.IndexOptions
	}{
		{
			keys:    bson.D{{Key: "user_id", Value: 1}},
			options: options.Index().SetName("idx_user_id"),
		},
		{
			keys:    bson.D{{Key: "author", Value: 1}},
			options: options.Index().SetName("idx_author"),
		},
		{
			keys:    bson.D{{Key: "status", Value: 1}},
			options: options.Index().SetName("idx_status"),
		},
		{
			keys:    bson.D{{Key: "tags", Value: 1}},
			options: options.Index().SetName("idx_tags"),
		},
		{
			keys:    bson.D{{Key: "created_at", Value: -1}},
			options: options.Index().SetName("idx_created_at"),
		},
		{
			// 文本索引用于全文搜索
			keys:    bson.D{{Key: "title", Value: "text"}, {Key: "content", Value: "text"}},
			options: options.Index().SetName("idx_text_search"),
		},
	}

	ctx := context.Background()
	for _, idx := range indexes {
		_, err := database.Mongo.CreateIndex(ctx, collection, idx.keys, idx.options)
		if err != nil {
			return err
		}
	}

	return nil
}
