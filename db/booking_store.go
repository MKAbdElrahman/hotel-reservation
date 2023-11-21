package db

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/mkabdelrahman/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingStore interface {
	InsertBooking(ctx context.Context, hotel *types.Booking) (*types.Booking, error)

	GetBookingByID(ctx context.Context, bookingID string) (*types.Booking, error)

	GetBookingByRoomAndTimeRange(ctx context.Context, roomID string, fromDate, tillDate time.Time) (*types.Booking, error)

	GetBookings(ctx context.Context) ([]*types.Booking, error)

	UpdateBooking(ctx context.Context, booking *types.Booking) error

	// DeleteBooking(ctx context.Context, bookigID string) error

	// QueryBooking(ctx context.Context, criteria types.BookingQueryCriteria) ([]*types.Booking, error)
}

type MongoBookingStore struct {
	client   *mongo.Client
	collName string
	dbName   string
	coll     *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client, dbName string, collName string) *MongoBookingStore {

	return &MongoBookingStore{
		client:   client,
		collName: collName,
		dbName:   dbName,
		coll:     client.Database(dbName).Collection(collName),
	}
}

func (m *MongoBookingStore) InsertBooking(ctx context.Context, booking *types.Booking) (*types.Booking, error) {
	result, err := m.coll.InsertOne(ctx, booking)
	if err != nil {
		log.Printf("Error inserting booking: %v\n", err)
		return nil, err
	}
	insertedID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, errors.New("could not convert InsertedID to ObjectID")
	}
	booking.ID = insertedID.Hex()
	return booking, nil
}

func (s *MongoBookingStore) Drop(c context.Context) error {
	return s.coll.Drop(c)
}

func (s *MongoBookingStore) GetBookingByID(ctx context.Context, ID string) (*types.Booking, error) {
	oid, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}

	var b types.Booking
	err = s.coll.FindOne(ctx, bson.M{"_id": oid}).Decode(&b)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.New("id not found")
		}
		return nil, err
	}
	return &b, nil
}

func (m *MongoBookingStore) GetBookingByRoomAndTimeRange(ctx context.Context, roomID string, fromDate, tillDate time.Time) (*types.Booking, error) {
	// Check if the room is booked for the specified time range
	filter := bson.M{
		"room_id":   roomID,
		"from_date": bson.M{"$lt": tillDate},
		"till_date": bson.M{"$gt": fromDate},
	}

	var existingBooking types.Booking
	err := m.coll.FindOne(ctx, filter).Decode(&existingBooking)
	if err == mongo.ErrNoDocuments {
		// No booking found for the specified time range, the room is available
		return nil, nil
	} else if err != nil {
		// Handle other errors
		return nil, err
	}

	// A booking exists for the specified time range, indicating the room is already booked
	return &existingBooking, nil
}

func (m *MongoBookingStore) GetBookings(ctx context.Context) ([]*types.Booking, error) {
	cursor, err := m.coll.Find(ctx, bson.M{})
	if err != nil {
		log.Printf("Error listing bookings: %v\n", err)
		return nil, err
	}
	defer cursor.Close(ctx)

	var bookings []*types.Booking
	for cursor.Next(ctx) {
		var booking types.Booking
		if err := cursor.Decode(&booking); err != nil {
			log.Printf("Error decoding booking: %v\n", err)
			return nil, err
		}
		bookings = append(bookings, &booking)
	}

	if err := cursor.Err(); err != nil {
		log.Printf("Cursor error: %v\n", err)
		return nil, err
	}

	return bookings, nil
}

func (m *MongoBookingStore) UpdateBooking(ctx context.Context, booking *types.Booking) error {
	oid, err := primitive.ObjectIDFromHex(booking.ID)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": oid}

	update := bson.M{
		"$set": bson.M{
			"user_id":        booking.UserID,
			"room_id":        booking.RoomID,
			"from_date":      booking.FromDate,
			"till_date":      booking.TillDate,
			"booking_status": booking.BookingStatus,
		},
	}

	result, err := m.coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.ModifiedCount == 0 {
		return errors.New("no matching booking found for update")
	}

	return nil
}
