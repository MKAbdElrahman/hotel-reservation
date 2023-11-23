package business

import (
	"github.com/mkabdelrahman/hotel-reservation/db"
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
