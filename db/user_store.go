package db

import (
	"context"
	"errors"

	"github.com/mkabdelrahman/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Dropper interface {
	Drop(c context.Context) error
}

type UserStore interface {
	Dropper
	GetUserByID(ctx context.Context, ID string) (*types.User, error)

	GetUserByEmail(ctx context.Context, email string) (*types.User, error)

	GetUsers(context.Context) ([]*types.User, error)

	GetUsersWithPagination(ctx context.Context, filter types.UsersPaginationFilter) ([]*types.User, error)

	InsertUser(ctx context.Context, user *types.User) (*types.User, error)

	DeleteUser(ctx context.Context, ID string) error

	UpdateUser(ctx context.Context, ID string, updateFields types.UpdateUserParams) (*types.User, error)
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

func (s *MongoUserStore) Drop(c context.Context) error {
	return s.coll.Drop(c)
}
func (s *MongoUserStore) GetUserByID(ctx context.Context, ID string) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}

	var u types.User
	err = s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&u)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("id not found")
		}
		return nil, err
	}
	return &u, nil
}

func (s *MongoUserStore) GetUserByEmail(ctx context.Context, email string) (*types.User, error) {

	var u types.User
	err := s.coll.FindOne(ctx, bson.M{"email": email}).Decode(&u)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("email not found")
		}
		return nil, err
	}

	return &u, nil
}
func (s *MongoUserStore) DeleteUser(ctx context.Context, ID string) error {

	oid, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return err
	}
	_, err = s.coll.DeleteOne(ctx, bson.M{"_id": oid})
	if err != nil {
		return err
	}

	return nil
}

func (s *MongoUserStore) GetUsers(ctx context.Context) ([]*types.User, error) {
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

func (s *MongoUserStore) InsertUser(ctx context.Context, user *types.User) (*types.User, error) {

	res, err := s.coll.InsertOne(ctx, user)

	if err != nil {
		return nil, err
	}

	user.ID = res.InsertedID.(primitive.ObjectID)

	return user, nil
}

func (s *MongoUserStore) UpdateUser(ctx context.Context, ID string, updateFields types.UpdateUserParams) (*types.User, error) {
	oid, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": oid}

	update := bson.M{"$set": updateFields.BSON()}

	_, err = s.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return nil, err
	}

	updatedUser, err := s.GetUserByID(ctx, ID)
	if err != nil {
		return nil, err
	}

	return updatedUser, nil
}

// GetUsersWithPagination retrieves users from the store with pagination using the provided filter.
func (s *MongoUserStore) GetUsersWithPagination(ctx context.Context, filter types.UsersPaginationFilter) ([]*types.User, error) {
	if err := filter.Validate(); err != nil {
		return nil, err
	}

	offset := (filter.Page - 1) * filter.PageSize

	options := options.Find().SetSkip(int64(offset)).SetLimit(int64(filter.PageSize))

	// Apply sorting options
	if filter.SortBy != "" {
		sortDir := 1
		if filter.SortDir == types.DescSort {
			sortDir = -1
		}

		options.SetSort(bson.D{{Key: filter.SortBy, Value: sortDir}})
	}

	cur, err := s.coll.Find(ctx, bson.M{}, options)
	if err != nil {
		return nil, err
	}

	var users []*types.User
	if err := cur.All(ctx, &users); err != nil {
		return nil, err
	}

	return users, nil
}
