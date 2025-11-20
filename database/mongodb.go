package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/demo/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var MongoDB *mongo.Client
var MongoDBName string

// InitMongoDB 初始化 MongoDB 连接
func InitMongoDB() error {
	cfg := config.GetConfig()
	mongoCfg := cfg.MongoDB

	// 构建连接 URI
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%d",
		mongoCfg.Username,
		mongoCfg.Password,
		mongoCfg.Host,
		mongoCfg.Port,
	)

	// 设置客户端选项
	clientOptions := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(uint64(mongoCfg.MaxPoolSize)).
		SetMinPoolSize(uint64(mongoCfg.MinPoolSize)).
		SetMaxConnIdleTime(5 * time.Minute).
		SetConnectTimeout(10 * time.Second).
		SetSocketTimeout(30 * time.Second)

	// 如果配置了认证源，设置认证数据库
	if mongoCfg.AuthSource != "" {
		clientOptions.SetAuth(options.Credential{
			AuthSource: mongoCfg.AuthSource,
			Username:   mongoCfg.Username,
			Password:   mongoCfg.Password,
		})
	}

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 连接 MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("连接 MongoDB 失败: %v", err)
	}

	// 测试连接
	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return fmt.Errorf("MongoDB 连接测试失败: %v", err)
	}

	MongoDB = client
	MongoDBName = mongoCfg.Database

	log.Println("✅ MongoDB 连接成功")
	return nil
}

// GetMongoDB 获取 MongoDB 客户端
func GetMongoDB() *mongo.Client {
	if MongoDB == nil {
		panic("MongoDB 未初始化，请先调用 InitMongoDB")
	}
	return MongoDB
}

// GetMongoDatabase 获取指定的数据库
func GetMongoDatabase(name ...string) *mongo.Database {
	dbName := MongoDBName
	if len(name) > 0 && name[0] != "" {
		dbName = name[0]
	}
	return MongoDB.Database(dbName)
}

// GetMongoCollection 获取指定的集合
func GetMongoCollection(collectionName string, dbName ...string) *mongo.Collection {
	db := GetMongoDatabase(dbName...)
	return db.Collection(collectionName)
}

// CloseMongoDB 关闭 MongoDB 连接
func CloseMongoDB() error {
	if MongoDB != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return MongoDB.Disconnect(ctx)
	}
	return nil
}

// MongoHelper MongoDB 辅助方法
type MongoHelper struct{}

// InsertOne 插入单个文档
func (m *MongoHelper) InsertOne(ctx context.Context, collection *mongo.Collection, document any) (*mongo.InsertOneResult, error) {
	return collection.InsertOne(ctx, document)
}

// InsertMany 插入多个文档
func (m *MongoHelper) InsertMany(ctx context.Context, collection *mongo.Collection, documents []any) (*mongo.InsertManyResult, error) {
	return collection.InsertMany(ctx, documents)
}

// FindOne 查询单个文档
func (m *MongoHelper) FindOne(ctx context.Context, collection *mongo.Collection, filter any, result any) error {
	return collection.FindOne(ctx, filter).Decode(result)
}

// FindMany 查询多个文档
func (m *MongoHelper) FindMany(ctx context.Context, collection *mongo.Collection, filter any, results any, opts ...*options.FindOptions) error {
	cursor, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	return cursor.All(ctx, results)
}

// UpdateOne 更新单个文档
func (m *MongoHelper) UpdateOne(ctx context.Context, collection *mongo.Collection, filter any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return collection.UpdateOne(ctx, filter, update, opts...)
}

// UpdateMany 更新多个文档
func (m *MongoHelper) UpdateMany(ctx context.Context, collection *mongo.Collection, filter any, update any, opts ...*options.UpdateOptions) (*mongo.UpdateResult, error) {
	return collection.UpdateMany(ctx, filter, update, opts...)
}

// DeleteOne 删除单个文档
func (m *MongoHelper) DeleteOne(ctx context.Context, collection *mongo.Collection, filter any) (*mongo.DeleteResult, error) {
	return collection.DeleteOne(ctx, filter)
}

// DeleteMany 删除多个文档
func (m *MongoHelper) DeleteMany(ctx context.Context, collection *mongo.Collection, filter any) (*mongo.DeleteResult, error) {
	return collection.DeleteMany(ctx, filter)
}

// Count 统计文档数量
func (m *MongoHelper) Count(ctx context.Context, collection *mongo.Collection, filter any) (int64, error) {
	return collection.CountDocuments(ctx, filter)
}

// Aggregate 聚合查询
func (m *MongoHelper) Aggregate(ctx context.Context, collection *mongo.Collection, pipeline any, results any) error {
	cursor, err := collection.Aggregate(ctx, pipeline)
	if err != nil {
		return err
	}
	defer cursor.Close(ctx)
	return cursor.All(ctx, results)
}

// CreateIndex 创建索引
func (m *MongoHelper) CreateIndex(ctx context.Context, collection *mongo.Collection, keys any, opts ...*options.IndexOptions) (string, error) {
	indexModel := mongo.IndexModel{
		Keys: keys,
	}
	if len(opts) > 0 {
		indexModel.Options = opts[0]
	}
	return collection.Indexes().CreateOne(ctx, indexModel)
}

// CreateIndexes 创建多个索引
func (m *MongoHelper) CreateIndexes(ctx context.Context, collection *mongo.Collection, indexes []mongo.IndexModel) ([]string, error) {
	return collection.Indexes().CreateMany(ctx, indexes)
}

// DropIndex 删除索引
func (m *MongoHelper) DropIndex(ctx context.Context, collection *mongo.Collection, indexName string) error {
	_, err := collection.Indexes().DropOne(ctx, indexName)
	return err
}

// Transaction 执行事务（需要 MongoDB 副本集或分片集群）
func (m *MongoHelper) Transaction(ctx context.Context, fn func(sessCtx mongo.SessionContext) error) error {
	session, err := MongoDB.StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	// 在事务中执行操作
	_, err = session.WithTransaction(ctx, func(sessCtx mongo.SessionContext) (any, error) {
		return nil, fn(sessCtx)
	})

	return err
}

// Paginate 分页查询
func (m *MongoHelper) Paginate(ctx context.Context, collection *mongo.Collection, filter any, page, pageSize int64, results any, sort ...any) (int64, error) {
	// 计算总数
	total, err := collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	// 计算跳过的文档数
	skip := (page - 1) * pageSize

	// 构建查询选项
	opts := options.Find().
		SetSkip(skip).
		SetLimit(pageSize)

	// 设置排序
	if len(sort) > 0 {
		opts.SetSort(sort[0])
	}

	// 执行查询
	cursor, err := collection.Find(ctx, filter, opts)
	if err != nil {
		return total, err
	}
	defer cursor.Close(ctx)

	// 解码结果
	if err := cursor.All(ctx, results); err != nil {
		return total, err
	}

	return total, nil
}

// BulkWrite 批量写入操作
func (m *MongoHelper) BulkWrite(ctx context.Context, collection *mongo.Collection, models []mongo.WriteModel, opts ...*options.BulkWriteOptions) (*mongo.BulkWriteResult, error) {
	return collection.BulkWrite(ctx, models, opts...)
}

// Distinct 获取去重后的值
func (m *MongoHelper) Distinct(ctx context.Context, collection *mongo.Collection, fieldName string, filter any) ([]any, error) {
	return collection.Distinct(ctx, fieldName, filter)
}

var Mongo = &MongoHelper{}
