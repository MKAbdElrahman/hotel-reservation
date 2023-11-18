package business

import (
	"context"
	"errors"

	"github.com/mkabdelrahman/hotel-reservation/db"
	"github.com/mkabdelrahman/hotel-reservation/types"
)

type Manager struct {
	HotelStore   db.HotelStore
	RoomStore    db.RoomStore
	UserStore    db.UserStore
	BookingStore db.BookingStore
}

func NewManager(userStore db.UserStore, hotelStore db.HotelStore, roomStore db.RoomStore, bookingStore db.BookingStore) *Manager {
	return &Manager{
		UserStore:    userStore,
		HotelStore:   hotelStore,
		RoomStore:    roomStore,
		BookingStore: bookingStore,
	}
}

func (m *Manager) AddNewHotel(ctx context.Context, params types.HotelParams) (string, error) {

	hotel := types.NewHotelFromParams(params)
	insertedHotel, err := m.HotelStore.InsertHotel(ctx, hotel)
	if err != nil {
		return "", err
	}

	return insertedHotel.ID, nil
}

func (m *Manager) AddNewBooking(ctx context.Context, params types.NewBookingParams) (string, error) {

	booking := types.NewBookingFromParams(params)

	insertedBooking, err := m.BookingStore.InsertBooking(ctx, booking)
	if err != nil {
		return "", err
	}

	return insertedBooking.ID, nil
}

func (m *Manager) AddNewUser(ctx context.Context, params types.UserParams) (string, error) {
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return "", err
	}

	insertedUser, err := m.UserStore.InsertUser(ctx, user)
	if err != nil {
		return "", err
	}

	if insertedUser == nil {
		return "", errors.New("insertedUser is nil")
	}

	return insertedUser.ID.Hex(), nil
}

func (m *Manager) ListHotels(ctx context.Context) ([]*types.Hotel, error) {

	hotels, err := m.HotelStore.GetHotels(ctx)
	if err != nil {
		return nil, err
	}
	return hotels, nil
}

func (m *Manager) QueryHotels(ctx context.Context, criteria types.QueryCriteria) ([]*types.Hotel, error) {

	hotels, err := m.HotelStore.QueryHotels(ctx, criteria)
	if err != nil {
		return nil, err
	}
	return hotels, nil
}

func (m *Manager) AddNewRoom(ctx context.Context, params types.RoomParams, hotelID string) (string, error) {

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
		return "", err
	}

	hotel, err := m.HotelStore.GetHotel(ctx, hotelID)
	if err != nil {
		return "", err
	}
	hotel.Rooms = append(hotel.Rooms, insertedRoom.ID)
	err = m.HotelStore.UpdateHotel(ctx, hotel)
	if err != nil {
		return "", err
	}
	return room.ID, nil
}
func (m *Manager) ListRoomsForHotel(ctx context.Context, hotelID string) ([]types.Room, error) {
	return m.RoomStore.GetRoomsByHotelID(ctx, hotelID)
}
