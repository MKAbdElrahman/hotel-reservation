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

type HotelStore interface {
	InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error)

	GetHotel(ctx context.Context, hotelID string) (*types.Hotel, error)

	UpdateHotel(ctx context.Context, hotel *types.Hotel) error

	DeleteHotel(ctx context.Context, hotelID string) error

	GetHotels(ctx context.Context) ([]*types.Hotel, error)
}

type MongoHotelStore struct {
	client   *mongo.Client
	collName string
	dbName   string
	coll     *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client, dbName string, collName string) *MongoHotelStore {

	return &MongoHotelStore{
		client:   client,
		collName: collName,
		dbName:   dbName,
		coll:     client.Database(dbName).Collection(collName),
	}
}
func (m *MongoHotelStore) InsertHotel(ctx context.Context, hotel *types.Hotel) (*types.Hotel, error) {
	result, err := m.coll.InsertOne(ctx, hotel)
	if err != nil {
		log.Printf("Error inserting hotel: %v\n", err)
		return nil, err
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("could not convert InsertedID to ObjectID")
	}
	hotel.ID = insertedID
	return hotel, nil
}

func (m *MongoHotelStore) GetHotel(ctx context.Context, hotelID string) (*types.Hotel, error) {
	oid, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		return nil, err
	}

	var hotel types.Hotel
	filter := bson.M{"_id": oid}
	err = m.coll.FindOne(ctx, filter).Decode(&hotel)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, types.ErrNotFound
		}
		log.Printf("Error getting hotel: %v\n", err)
		return nil, err
	}
	return &hotel, nil
}

func (m *MongoHotelStore) UpdateHotel(ctx context.Context, hotel *types.Hotel) error {
	oid := hotel.ID

	filter := bson.M{"_id": oid}
	update := bson.M{"$set": hotel}
	_, err := m.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating hotel: %v\n", err)
		return err
	}
	return nil
}

func (m *MongoHotelStore) DeleteHotel(ctx context.Context, hotelID string) error {
	oid, err := primitive.ObjectIDFromHex(hotelID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}
	_, err = m.coll.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("Error deleting hotel: %v\n", err)
		return err
	}
	return nil
}

func (m *MongoHotelStore) GetHotels(ctx context.Context) ([]*types.Hotel, error) {
	cursor, err := m.coll.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error listing hotels: %v\n", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var hotels []*types.Hotel
	for cursor.Next(ctx) {
		var hotel types.Hotel
		if err := cursor.Decode(&hotel); err != nil {
			log.Printf("Error decoding hotel: %v\n", err)
			return nil, err
		}
		hotels = append(hotels, &hotel)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v\n", err)
		return nil, err
	}

	return hotels, nil
}
