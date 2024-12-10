package mongo

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/event"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"testing"
)

func TestMongo(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 10)
	defer cancel()
	monitor := &event.CommandMonitor{}
	opt := options.Client().ApplyURI("mongodb://localhost:27017/").
		SetMonitor(monitor)
	client, err := mongo.Connect(ctx, opt)
	assert.NoError(t, err)
	mdb := client.Database("webook")
	col := mdb.Collection("articles")

	res, err := col.InsertOne(ctx, Article{
		Id:      132,
		Title:   "test",
		Content: "test",
		Status:  1,
	})
	println(res.InsertedID)
	filter := bson.D{bson.E{Key: "id", Value: 132}}
	var art Article
	err = col.FindOne(ctx, filter).Decode(&art)
	assert.NoError(t, err)
	assert.Equal(t, "test", art.Title)
}

type Article struct {
	Id      int64  `bson:"id,omitempty"`
	Title   string `bson:"title,omitempty"`
	Content string `bson:"content,omitempty"`

	Status uint8 `bson:"status,omitempty"`
}
