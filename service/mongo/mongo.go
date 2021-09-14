package mongo

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	NoFoundError string = "mongo: no documents in result"
)

// 查询参数
type Filter bson.M

type Client struct {
	database *mongo.Database // 数据库实例
}

func NewClient(uri string, db string) *Client {
	client := new(Client)
	if err := client.Connect(uri, db); err != nil {
		log.Panicf("[mongo] connect uri(%s) db(%s) err: %v", uri, db, err)
	}
	return client
}

func (o *Client) Connect(uri string, db string) error {
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		return err
	}

	o.database = client.Database(db)
	return nil
}

// 检查连接是否存在
func (o *Client) Ping() error {
	return o.database.Client().Ping(context.TODO(), nil)
}

// 查找数据
func (o *Client) FindOne(collection string, filter Filter, result interface{}) error {
	if err := o.Ping(); err != nil {
		return err
	}

	collect := o.database.Collection(collection)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	return collect.FindOne(ctx, filter).Decode(result)
}

// 插入数据
func (o *Client) InsertOne(collection string, document interface{}) error {
	if err := o.Ping(); err != nil {
		return err
	}

	collect := o.database.Collection(collection)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	if _, err := collect.InsertOne(ctx, document); err != nil {
		return err
	}
	return nil
}

// 更新数据
func (o *Client) UpdateOne(collection string, filter Filter, document interface{}) error {
	if err := o.Ping(); err != nil {
		return err
	}

	collect := o.database.Collection(collection)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	update := bson.D{{"$set", document}}
	if _, err := collect.UpdateOne(ctx, filter, update); err != nil {
		return err
	}
	return nil
}

// 更新数据，如果没有则插入数据
func (o *Client) UpsertOne(collection string, filter Filter, document interface{}) error {
	if err := o.Ping(); err != nil {
		return err
	}

	collect := o.database.Collection(collection)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	update := bson.D{{"$set", document}}
	opts := &options.UpdateOptions{}
	opts.SetUpsert(true)
	if _, err := collect.UpdateOne(ctx, filter, update, opts); err != nil {
		return err
	}
	return nil
}

// 删除数据
func (o *Client) DeleteOne(collection string, filter Filter) error {
	if err := o.Ping(); err != nil {
		return err
	}

	collect := o.database.Collection(collection)
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	if _, err := collect.DeleteOne(ctx, filter); err != nil {
		return err
	}
	return nil
}
