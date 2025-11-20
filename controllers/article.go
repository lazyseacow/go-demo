package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/demo/database"
	"github.com/demo/middleware"
	"github.com/demo/models"
	"github.com/demo/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ArticleController struct {
	*BaseController
}

func NewArticleController() *ArticleController {
	return &ArticleController{
		BaseController: NewBaseController(),
	}
}

// CreateArticle 创建文章（MongoDB 示例）
// @Summary      创建文章
// @Description  创建新文章（存储在 MongoDB）
// @Tags         文章管理
// @Accept       json
// @Produce      json
// @Security     ApiKeyAuth
// @Param        request  body      models.Article  true  "文章内容"
// @Success      200      {object}  utils.Response{data=object{id=string}}  "创建成功"
// @Failure      14001    {object}  utils.Response  "参数错误"
// @Router       /articles [post]
func (ctrl *ArticleController) CreateArticle(ctx *gin.Context) {
	var req models.Article
	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Fail(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 获取当前用户 ID
	userID := middleware.GetUserID(ctx)
	username := middleware.GetUsername(ctx)

	// 设置文章信息
	req.UserID = userID
	req.Author = username
	req.CreatedAt = time.Now()
	req.UpdatedAt = time.Now()
	req.Views = 0
	req.Likes = 0
	if req.Status == 0 {
		req.Status = 1 // 默认已发布
	}

	// 获取 MongoDB 集合
	collection := database.GetMongoCollection("articles")

	// 插入文档
	result, err := database.Mongo.InsertOne(context.Background(), collection, req)
	if err != nil {
		utils.Fail(ctx, http.StatusInternalServerError, "创建失败: "+err.Error())
		return
	}

	utils.Success(ctx, gin.H{
		"id": result.InsertedID,
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
	var req models.ArticleQueryRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		req.Page = 1
		req.PageSize = 10
	}

	// 构建查询条件
	filter := bson.M{}

	if req.Author != "" {
		filter["author"] = req.Author
	}
	if req.Status > 0 {
		filter["status"] = req.Status
	}
	if len(req.Tags) > 0 {
		filter["tags"] = bson.M{"$in": req.Tags}
	}
	if req.Keyword != "" {
		// 标题或内容包含关键词
		filter["$or"] = []bson.M{
			{"title": bson.M{"$regex": req.Keyword, "$options": "i"}},
			{"content": bson.M{"$regex": req.Keyword, "$options": "i"}},
		}
	}

	// 获取 MongoDB 集合
	collection := database.GetMongoCollection("articles")

	// 分页查询
	var articles []models.Article
	total, err := database.Mongo.Paginate(
		context.Background(),
		collection,
		filter,
		int64(req.Page),
		int64(req.PageSize),
		&articles,
		bson.D{{Key: "created_at", Value: -1}}, // 按创建时间倒序
	)

	if err != nil {
		utils.Fail(ctx, http.StatusInternalServerError, "查询失败: "+err.Error())
		return
	}

	utils.Success(ctx, models.PageResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     articles,
	})
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
	idStr := ctx.Param("id")

	// 转换为 ObjectID
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.Fail(ctx, http.StatusBadRequest, "无效的文章 ID")
		return
	}

	// 获取 MongoDB 集合
	collection := database.GetMongoCollection("articles")

	// 查询文档
	var article models.Article
	filter := bson.M{"_id": objectID}
	if err := database.Mongo.FindOne(context.Background(), collection, filter, &article); err != nil {
		utils.Fail(ctx, http.StatusNotFound, "文章不存在")
		return
	}

	// 增加浏览次数
	update := bson.M{"$inc": bson.M{"views": 1}}
	_, _ = database.Mongo.UpdateOne(context.Background(), collection, filter, update)

	utils.Success(ctx, article)
}

// UpdateArticle 更新文章
func (ctrl *ArticleController) UpdateArticle(ctx *gin.Context) {
	idStr := ctx.Param("id")

	// 转换为 ObjectID
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.Fail(ctx, http.StatusBadRequest, "无效的文章 ID")
		return
	}

	var req struct {
		Title   string   `json:"title" binding:"omitempty,min=1,max=200"`
		Content string   `json:"content"`
		Tags    []string `json:"tags"`
		Status  int      `json:"status"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		utils.Fail(ctx, http.StatusBadRequest, "参数错误: "+err.Error())
		return
	}

	// 获取当前用户 ID
	userID := middleware.GetUserID(ctx)

	// 获取 MongoDB 集合
	collection := database.GetMongoCollection("articles")

	// 检查文章是否属于当前用户
	filter := bson.M{
		"_id":     objectID,
		"user_id": userID,
	}

	// 构建更新内容
	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	if req.Title != "" {
		update["$set"].(bson.M)["title"] = req.Title
	}
	if req.Content != "" {
		update["$set"].(bson.M)["content"] = req.Content
	}
	if len(req.Tags) > 0 {
		update["$set"].(bson.M)["tags"] = req.Tags
	}
	if req.Status > 0 {
		update["$set"].(bson.M)["status"] = req.Status
	}

	// 更新文档
	result, err := database.Mongo.UpdateOne(context.Background(), collection, filter, update)
	if err != nil {
		utils.Fail(ctx, http.StatusInternalServerError, "更新失败: "+err.Error())
		return
	}

	if result.MatchedCount == 0 {
		utils.Fail(ctx, http.StatusNotFound, "文章不存在或无权限")
		return
	}

	utils.Success(ctx, nil, "更新成功")
}

// DeleteArticle 删除文章
func (ctrl *ArticleController) DeleteArticle(ctx *gin.Context) {
	idStr := ctx.Param("id")

	// 转换为 ObjectID
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.Fail(ctx, http.StatusBadRequest, "无效的文章 ID")
		return
	}

	// 获取当前用户 ID
	userID := middleware.GetUserID(ctx)

	// 获取 MongoDB 集合
	collection := database.GetMongoCollection("articles")

	// 只能删除自己的文章
	filter := bson.M{
		"_id":     objectID,
		"user_id": userID,
	}

	// 删除文档
	result, err := database.Mongo.DeleteOne(context.Background(), collection, filter)
	if err != nil {
		utils.Fail(ctx, http.StatusInternalServerError, "删除失败: "+err.Error())
		return
	}

	if result.DeletedCount == 0 {
		utils.Fail(ctx, http.StatusNotFound, "文章不存在或无权限")
		return
	}

	utils.Success(ctx, nil, "删除成功")
}

// LikeArticle 点赞文章
func (ctrl *ArticleController) LikeArticle(ctx *gin.Context) {
	idStr := ctx.Param("id")

	// 转换为 ObjectID
	objectID, err := primitive.ObjectIDFromHex(idStr)
	if err != nil {
		utils.Fail(ctx, http.StatusBadRequest, "无效的文章 ID")
		return
	}

	// 获取 MongoDB 集合
	collection := database.GetMongoCollection("articles")

	// 增加点赞数
	filter := bson.M{"_id": objectID}
	update := bson.M{"$inc": bson.M{"likes": 1}}

	result, err := database.Mongo.UpdateOne(context.Background(), collection, filter, update)
	if err != nil {
		utils.Fail(ctx, http.StatusInternalServerError, "点赞失败: "+err.Error())
		return
	}

	if result.MatchedCount == 0 {
		utils.Fail(ctx, http.StatusNotFound, "文章不存在")
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
