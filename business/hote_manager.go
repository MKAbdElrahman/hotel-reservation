package business

import (
	"context"

	"github.com/mkabdelrahman/hotel-reservation/db"
	"github.com/mkabdelrahman/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Manager struct {
	HotelStore db.HotelStore
	RoomStore  db.RoomStore
}

func NewManager(hotelStore db.HotelStore, roomStore db.RoomStore) *Manager {
	return &Manager{
		HotelStore: hotelStore,
		RoomStore:  roomStore,
	}
}

func (m *Manager) AddNewHotel(ctx context.Context, params types.NewHotelParams) (primitive.ObjectID, error) {

	hotel := &types.Hotel{
		Name:     params.Name,
		Location: params.Location,
	}
	insertedHotel, err := m.HotelStore.InsertHotel(ctx, hotel)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	return insertedHotel.ID, nil
}

func (m *Manager) ListHotels(ctx context.Context) ([]*types.Hotel, error) {

	hotels, err := m.HotelStore.GetHotels(ctx)
	if err != nil {
		return nil, err
	}
	return hotels, nil
}

func (m *Manager) AddNewRoom(ctx context.Context, params types.NewRoomParams, hotelID primitive.ObjectID) (primitive.ObjectID, error) {

	room := &types.Room{
		HotelID:     hotelID,
		Number:      params.Number,
		Floor:       params.Floor,
		Type:        params.Type,
		Description: params.Description,
		Price:       params.Price,
		Occupied:    params.Occupied,
	}
	insertedRoom, err := m.RoomStore.InsertRoom(ctx, room)
	if err != nil {
		return primitive.ObjectID{}, err
	}

	hotel, err := m.HotelStore.GetHotel(ctx, hotelID.Hex())
	if err != nil {
		return primitive.ObjectID{}, err
	}
	hotel.Rooms = append(hotel.Rooms, insertedRoom.ID)
	err = m.HotelStore.UpdateHotel(ctx, hotel)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return room.ID, nil
}
func (m *Manager) ListRoomsForHotel(ctx context.Context, hotelID string) ([]types.Room, error) {
	return m.RoomStore.GetRoomsByHotelID(ctx, hotelID)
}
