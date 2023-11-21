package business

import (
	"context"
	"errors"
	"fmt"

	"github.com/mkabdelrahman/hotel-reservation/types"
)

func (m *Manager) AddNewBooking(ctx context.Context, booking *types.Booking) (string, error) {

	existingBooking, err := m.BookingStore.GetBookingByRoomAndTimeRange(ctx, booking.RoomID, booking.FromDate, booking.TillDate)

	if err != nil {
		return "", err
	}
	if existingBooking != nil {
		return "", errors.New("room is already booked for the specified time range")
	}

	fmt.Printf("%+v", booking)

	insertedBooking, err := m.BookingStore.InsertBooking(ctx, booking)

	if err != nil {
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
