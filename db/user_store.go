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
	GetUsers(context.Context) ([]*types.User, error)

	InsertUser(ctx context.Context, user *types.User) (*types.User, error)
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

func (s MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
	cur, err := s.coll.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}

	var users []*types.User
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (s MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {

	res, err := s.coll.InsertOne(ctx, user)

	if err != nil {
		return nil, err
	}

	user.ID = res.InsertedID.(primitive.ObjectID)

	return user, nil
}
