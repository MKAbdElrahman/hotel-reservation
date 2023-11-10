package db

import (
	"context"

	"github.com/mkabdelrahman/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserStore interface {
	GetUserByID(ctx context.Context, ID string) (*types.User, error)
}

type MongoUserStore struct {
	client   *mongo.Client
	collName string
	dbName   string
	coll     *mongo.Collection
}

func NewMongoUserStore(client *mongo.Client, dbName string, collName string) *MongoUserStore {

	return &MongoUserStore{
		client:   client,
		collName: collName,
		dbName:   dbName,
		coll:     client.Database(dbName).Collection(collName),
	}
}
func (s MongoUserStore) GetUserByID(ctx context.Context, ID string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}

	var u types.User
	err = s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&u)
	if err != nil {
		return nil, err
	}
	return &u, nil
}
