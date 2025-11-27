package service

import (
	"context"
	"time"

	"github.com/demo/common"
	"github.com/demo/database"
	"github.com/demo/models"
	"github.com/demo/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ArticleService 文章服务（直接使用 MongoDB，不通过 Repository 层）
type ArticleService struct {
	*BaseService
}

// NewArticleService 创建文章服务实例
func NewArticleService() *ArticleService {
	return &ArticleService{
		BaseService: NewBaseService(),
	}
}

// getCollection 获取 MongoDB 集合（延迟初始化，检查 MongoDB 是否可用）
func (s *ArticleService) getCollection() (*mongo.Collection, error) {
	if database.MongoDB == nil {
		return nil, utils.NewError(common.CodeServiceUnavail, "MongoDB 服务不可用")
	}
	return database.GetMongoCollection("articles"), nil
}

// CreateArticle 创建文章
func (s *ArticleService) CreateArticle(req models.CreateArticleRequest, userID int64, username string) (*models.Article, error) {
	// 构建文章对象
	article := &models.Article{
		Title:     req.Title,
		Content:   req.Content,
		Tags:      req.Tags,
		Author:    username,
		UserID:    userID,
		Status:    req.Status,
		Views:     0,
		Likes:     0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 如果未指定状态，默认为已发布
	if article.Status == 0 {
		article.Status = 1
	}

	// 获取集合
	collection, err := s.getCollection()
	if err != nil {
		return nil, err
	}

	// 创建文章
	ctx := context.Background()
	result, err := database.Mongo.InsertOne(ctx, collection, article)
	if err != nil {
		return nil, common.NewError(common.CodeDBInsertFailed, "创建文章失败: "+err.Error())
	}

	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, common.NewError(common.CodeInternalError, "无效的 ObjectID")
	}

	article.ID = objectID
	return article, nil
}

// GetArticleList 获取文章列表
func (s *ArticleService) GetArticleList(req models.ArticleQueryRequest) (*models.PageResponse, error) {
	// 默认分页参数
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
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

	// 获取集合
	collection, err := s.getCollection()
	if err != nil {
		return nil, err
	}

	// 查询文章列表
	ctx := context.Background()
	var articles []models.Article

	// 计算总数
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, common.NewError(common.CodeDBQueryFailed, "查询文章列表失败: "+err.Error())
	}

	// 分页查询
	skip := int64(req.Page-1) * int64(req.PageSize)
	opts := options.Find().
		SetSkip(skip).
		SetLimit(int64(req.PageSize)).
		SetSort(bson.D{{Key: "created_at", Value: -1}}) // 按创建时间倒序

	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, common.NewError(common.CodeDBQueryFailed, "查询文章列表失败: "+err.Error())
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &articles); err != nil {
		return nil, common.NewError(common.CodeDBQueryFailed, "查询文章列表失败: "+err.Error())
	}

	return &models.PageResponse{
		Total:    total,
		Page:     req.Page,
		PageSize: req.PageSize,
		List:     articles,
	}, nil
}

// GetArticleByID 根据 ID 获取文章
func (s *ArticleService) GetArticleByID(id string) (*models.Article, error) {
	// 转换为 ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, common.NewError(common.CodeParamInvalid, "无效的文章 ID")
	}

	// 获取集合
	collection, err := s.getCollection()
	if err != nil {
		return nil, err
	}

	// 查询文章
	ctx := context.Background()
	var article models.Article
	filter := bson.M{"_id": objectID}
	if err := collection.FindOne(ctx, filter).Decode(&article); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, common.NewError(common.CodeArticleNotFound, "文章不存在")
		}
		return nil, common.NewError(common.CodeDBQueryFailed, "查询文章失败: "+err.Error())
	}

	// 增加浏览次数（异步，不阻塞）
	go func() {
		if coll, err := s.getCollection(); err == nil {
			update := bson.M{"$inc": bson.M{"views": 1}}
			_, _ = coll.UpdateOne(context.Background(), filter, update)
		}
	}()

	return &article, nil
}

// UpdateArticle 更新文章
func (s *ArticleService) UpdateArticle(id string, userID int64, req models.UpdateArticleRequest) error {
	// 转换为 ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return common.NewError(common.CodeParamInvalid, "无效的文章 ID")
	}

	// 构建更新内容
	updates := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	if req.Title != "" {
		updates["$set"].(bson.M)["title"] = req.Title
	}
	if req.Content != "" {
		updates["$set"].(bson.M)["content"] = req.Content
	}
	if len(req.Tags) > 0 {
		updates["$set"].(bson.M)["tags"] = req.Tags
	}
	if req.Status > 0 {
		updates["$set"].(bson.M)["status"] = req.Status
	}

	// 获取集合
	collection, err := s.getCollection()
	if err != nil {
		return err
	}

	// 更新文章（只能更新自己的文章）
	ctx := context.Background()
	filter := bson.M{
		"_id":     objectID,
		"user_id": userID,
	}

	result, err := collection.UpdateOne(ctx, filter, updates)
	if err != nil {
		return common.NewError(common.CodeDBUpdateFailed, "更新文章失败: "+err.Error())
	}

	if result.MatchedCount == 0 {
		return common.NewError(common.CodeArticleNotFound, "文章不存在或无权限")
	}

	return nil
}

// DeleteArticle 删除文章
func (s *ArticleService) DeleteArticle(id string, userID int64) error {
	// 转换为 ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return common.NewError(common.CodeParamInvalid, "无效的文章 ID")
	}

	// 获取集合
	collection, err := s.getCollection()
	if err != nil {
		return err
	}

	// 删除文章（只能删除自己的文章）
	ctx := context.Background()
	filter := bson.M{
		"_id":     objectID,
		"user_id": userID,
	}

	result, err := collection.DeleteOne(ctx, filter)
	if err != nil {
		return common.NewError(common.CodeDBDeleteFailed, "删除文章失败: "+err.Error())
	}

	if result.DeletedCount == 0 {
		return common.NewError(common.CodeArticleNotFound, "文章不存在或无权限")
	}

	return nil
}

// LikeArticle 点赞文章
func (s *ArticleService) LikeArticle(id string) error {
	// 转换为 ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return common.NewError(common.CodeParamInvalid, "无效的文章 ID")
	}

	// 获取集合
	collection, err := s.getCollection()
	if err != nil {
		return err
	}

	// 增加点赞数
	ctx := context.Background()
	filter := bson.M{"_id": objectID}
	update := bson.M{"$inc": bson.M{"likes": 1}}

	result, err := collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return common.NewError(common.CodeDBUpdateFailed, "点赞失败: "+err.Error())
	}

	if result.MatchedCount == 0 {
		return common.NewError(common.CodeArticleNotFound, "文章不存在")
	}

	return nil
}
