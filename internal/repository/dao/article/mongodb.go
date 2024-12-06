package article

import (
	"context"
	"errors"
	"github.com/bwmarrin/snowflake"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type MongoDBDAO struct {
	//client *mongo.Client
	// 代表 webook 的
	//database *mongo.Database
	// 代表的是制作库
	col *mongo.Collection
	// 代表的是线上库
	liveCol *mongo.Collection
	node    *snowflake.Node

	idGen IDGenerator
}

type IDGenerator func() int64

func (m MongoDBDAO) Insert(ctx context.Context, art Article) (int64, error) {
	now := time.Now().UnixMilli()
	art.Ctime = now
	art.Utime = now
	id := m.node.Generate().Int64()
	art.Id = id
	_, err := m.col.InsertOne(ctx, art)
	return id, err
}

func (m MongoDBDAO) UpdateById(ctx context.Context, art Article) error {
	filter := bson.M{"id": art.Id, "author_id": art.AuthorId}
	update := bson.D{bson.E{"$set", bson.M{
		"title":   art.Title,
		"content": art.Content,
		"utime":   time.Now().UnixMilli(),
		"status":  art.Status,
	}}}
	res, err := m.col.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}
	if res.ModifiedCount == 0 {
		return errors.New("没有修改任何数据")
	}
	return nil
}

func (m MongoDBDAO) GetByAuthor(ctx context.Context, author int64, offset, limit int) ([]Article, error) {
	//TODO implement me
	panic("implement me")
}

func (m MongoDBDAO) GetById(ctx context.Context, id int64) (Article, error) {
	//TODO implement me
	panic("implement me")
}

func (m MongoDBDAO) GetPubById(ctx context.Context, id int64) (PublishedArticle, error) {
	//TODO implement me
	panic("implement me")
}

func (m MongoDBDAO) Sync(ctx context.Context, art Article) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m MongoDBDAO) SyncStatus(ctx context.Context, author, id int64, status uint8) error {
	//TODO implement me
	panic("implement me")
}

func InitCollections(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	index := []mongo.IndexModel{
		{
			Keys:    bson.D{bson.E{Key: "id", Value: 1}},
			Options: options.Index().SetUnique(true),
		},
		{
			Keys: bson.D{bson.E{Key: "author_id", Value: 1},
				bson.E{Key: "ctime", Value: 1},
			},
			Options: options.Index(),
		},
	}
	_, err := db.Collection("articles").Indexes().
		CreateMany(ctx, index)
	if err != nil {
		return err
	}
	_, err = db.Collection("published_articles").Indexes().
		CreateMany(ctx, index)
	return err
}

func NewMongoDBDAOV1(db *mongo.Database, idGen IDGenerator) ArticleDAO {
	return &MongoDBDAO{
		col:     db.Collection("articles"),
		liveCol: db.Collection("published_articles"),
		//node:    node,
		idGen: idGen,
	}
}

func NewMongoDBDAO(db *mongo.Database, node *snowflake.Node) ArticleDAO {
	return &MongoDBDAO{
		col:     db.Collection("articles"),
		liveCol: db.Collection("published_articles"),
		node:    node,
	}
}
