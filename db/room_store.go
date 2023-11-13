package db

import (
	"context"
	"errors"
	"log"

	"github.com/mkabdelrahman/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error)
	GetRoom(ctx context.Context, roomID string) (*types.Room, error)
	DeleteRoom(ctx context.Context, roomID string) error
	UpdateRoom(ctx context.Context, room *types.Room) error
	GetRoomsByHotelID(ctx context.Context, hotelID string) ([]types.Room, error)
}

type MongoRoomStore struct {
	client   *mongo.Client
	collName string
	dbName   string
	coll     *mongo.Collection
}

func NewMongoRoomStore(client *mongo.Client, dbName string, collName string) *MongoRoomStore {

	return &MongoRoomStore{
		client:   client,
		collName: collName,
		dbName:   dbName,
		coll:     client.Database(dbName).Collection(collName),
	}
}

func (m *MongoRoomStore) InsertRoom(ctx context.Context, room *types.Room) (*types.Room, error) {
	result, err := m.coll.InsertOne(ctx, room)
	if err != nil {
		log.Printf("Error inserting room: %v\n", err)
		return nil, err
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("could not convert InsertedID to ObjectID")
	}
	room.ID = insertedID
	return room, nil
}

func (m *MongoRoomStore) GetRoom(ctx context.Context, roomID string) (*types.Room, error) {
	oid, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return nil, err
	}

	var room types.Room
	filter := bson.M{"_id": oid}
	err = m.coll.FindOne(ctx, filter).Decode(&room)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, types.ErrNotFound
		}
		log.Printf("Error getting room: %v\n", err)
		return nil, err
	}
	return &room, nil
}

func (m *MongoRoomStore) DeleteRoom(ctx context.Context, roomID string) error {
	oid, err := primitive.ObjectIDFromHex(roomID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	_, err = m.coll.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("Error deleting room: %v\n", err)
		return err
	}
	return nil
}

func (m *MongoRoomStore) UpdateRoom(ctx context.Context, room *types.Room) error {
	oid := room.ID

	filter := bson.M{"_id": oid}
	update := bson.M{"$set": room}
	_, err := m.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating room: %v\n", err)
		return err
	}
	return nil
}

func (m *MongoRoomStore) ListRooms(ctx context.Context) ([]*types.Room, error) {
	cursor, err := m.coll.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error listing rooms: %v\n", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var rooms []*types.Room
	for cursor.Next(ctx) {
		var hotel types.Room
		if err := cursor.Decode(&hotel); err != nil {
			log.Printf("Error decoding room: %v\n", err)
			return nil, err
		}
		rooms = append(rooms, &hotel)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v\n", err)
		return nil, err
	}

	return rooms, nil
}
func (r *MongoRoomStore) GetRoomsByHotelID(ctx context.Context, hotelID string) ([]types.Room, error) {

	oid, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"hotel_id": oid}

	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
			return nil, err
	}
	defer cursor.Close(ctx)

	var rooms []types.Room
	if err := cursor.All(ctx, &rooms); err != nil {
		return nil, err
	}

	return rooms, nil
}
