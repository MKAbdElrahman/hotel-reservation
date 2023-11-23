package types

import (
	"errors"
	"time"
)

type Booking struct {
	ID            string        `json:"id" bson:"_id,omitempty"`
	UserID        string        `json:"user_id" bson:"user_id"`
	RoomID        string        `json:"room_id" bson:"room_id"`
	FromDate      time.Time     `json:"from_date" bson:"from_date"`
	TillDate      time.Time     `json:"till_date" bson:"till_date"`
	BookingStatus BookingStatus `json:"booking_status" bson:"booking_status"`
}

type NewBookingParams struct {
	UserID        string        `json:"user_id" bson:"user_id"`
	RoomID        string        `json:"room_id" bson:"room_id"`
	FromDate      time.Time     `json:"from_date" bson:"from_date"`
	TillDate      time.Time     `json:"till_date" bson:"till_date"`
	BookingStatus BookingStatus `json:"booking_status" bson:"booking_status"`
}

func (params NewBookingParams) Validate() error {
	if params.UserID == "" {
		return errors.New("UserID is required")
	}

	if params.RoomID == "" {
		return errors.New("RoomID is required")
	}

	if params.FromDate.IsZero() {
		return errors.New("FromDate is required and must be a valid time")
	}

	if params.TillDate.IsZero() {
		return errors.New("TillDate is required and must be a valid time")
	}

	if params.TillDate.Before(params.FromDate) {
		return errors.New("TillDate must be after FromDate")
	}

	now := time.Now()
	if params.FromDate.Before(now) {
		return errors.New("FromDate must be in the future")
	}
	return nil
}
func NewBookingFromParams(params NewBookingParams) *Booking {

	if params.BookingStatus == "" {
		params.BookingStatus = "Pending"
	}

	return &Booking{
		UserID:        params.UserID,
		RoomID:        params.RoomID,
		FromDate:      params.FromDate,
		TillDate:      params.TillDate,
		BookingStatus: params.BookingStatus,
	}
}

type BookingStatus string

const (
	StatusPending   BookingStatus = "Pending"
	StatusConfirmed BookingStatus = "Confirmed"
	StatusCanceled  BookingStatus = "Canceled"
)
