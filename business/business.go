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

