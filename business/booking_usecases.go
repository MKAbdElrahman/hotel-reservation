package business

import (
	"context"
	"errors"

	"github.com/mkabdelrahman/hotel-reservation/types"
)

func (m *Manager) AddNewBooking(ctx context.Context, params types.NewBookingParams) (string, error) {
	// Make sure user ID in params exists
	user, err := m.UserStore.GetUserByID(ctx, params.UserID)
	if err != nil {
		return "", err
	}
	if user == nil {
		return "", errors.New("user not found")
	}

	// Make sure room ID in params exists
	room, err := m.RoomStore.GetRoomByID(ctx, params.RoomID)
	if err != nil {
		return "", err
	}
	if room == nil {
		return "", errors.New("room not found")
	}

	// Check if the room is already booked for the specified time range
	existingBooking, err := m.BookingStore.GetBookingByRoomAndTimeRange(ctx, params.RoomID, params.FromDate, params.TillDate)
	if err != nil {
		return "", err
	}
	if existingBooking != nil {
		// Check if the existing booking has been canceled
		if existingBooking.BookingStatus == types.StatusCanceled {
			// Proceed with the new booking
		} else {
			return "", errors.New("room is already booked for the specified time range")
		}
	}

	// Create a new booking from params
	newBooking := types.NewBookingFromParams(params)

	// Insert the new booking, the db will return a booking with the id field filled
	insertedBooking, err := m.BookingStore.InsertBooking(ctx, newBooking)
	if err != nil {
		return "", err
	}

	// Change the occupied field of the room to true and update
	room.Occupied = true
	err = m.RoomStore.UpdateRoom(ctx, room)
	if err != nil {
		// Rollback the booking if updating the room fails
		rollbackErr := m.BookingStore.DeleteBookingByID(ctx, insertedBooking.ID)
		if rollbackErr != nil {
			return "", err
		}

		return "", err
	}

	return insertedBooking.ID, nil
}

func (m *Manager) ListBookings(ctx context.Context) ([]*types.Booking, error) {

	bookings, err := m.BookingStore.GetBookings(ctx)
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (m *Manager) ListUserBookings(ctx context.Context, userID string) ([]*types.Booking, error) {

	bookings, err := m.BookingStore.GetBookingsByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return bookings, nil
}

func (m *Manager) CancelBooking(ctx context.Context, bookingID string) error {
	booking, err := m.BookingStore.GetBookingByID(ctx, bookingID)
	if err != nil {
		return err
	}

	if booking == nil {
		return errors.New("booking not found")
	}

	if booking.BookingStatus == types.StatusCanceled {
		return errors.New("booking is already canceled")
	}

	booking.BookingStatus = types.StatusCanceled

	err = m.BookingStore.UpdateBooking(ctx, booking)
	if err != nil {
		return err
	}

	return nil
}
