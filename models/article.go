package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Article MongoDB 文章模型示例
type Article struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Title     string             `bson:"title" json:"title" binding:"required,min=1,max=200"`
	Content   string             `bson:"content" json:"content" binding:"required"`
	Author    string             `bson:"author" json:"author"`
	UserID    int64              `bson:"user_id" json:"user_id"`
	Tags      []string           `bson:"tags" json:"tags"`
	Views     int64              `bson:"views" json:"views"`
	Likes     int64              `bson:"likes" json:"likes"`
	Status    int                `bson:"status" json:"status"` // 1:已发布 0:草稿
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}

// CollectionName 集合名称
func (Article) CollectionName() string {
	return "articles"
}

// ArticleQueryRequest 文章查询请求
type ArticleQueryRequest struct {
	Page     int      `form:"page" json:"page" binding:"min=1"`
	PageSize int      `form:"page_size" json:"page_size" binding:"min=1,max=100"`
	Author   string   `form:"author" json:"author"`
	Status   int      `form:"status" json:"status"`
	Tags     []string `form:"tags" json:"tags"`
	Keyword  string   `form:"keyword" json:"keyword"` // 搜索关键词
}
