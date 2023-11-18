package types

import (
	"time"
)

type Booking struct {
	ID       string    `json:"id" bson:"_id,omitempty"`
	UserID   string    `json:"user_id" bson:"user_id,omitempty"`
	RoomID   string    `json:"room_id" bson:"room_id,omitempty"`
	FromDate time.Time `json:"from_date" bson:"from_date,omitempty"`
	TillDate time.Time `json:"till_date" bson:"till_date,omitempty"`
}

type NewBookingParams struct {
	UserID   string    `json:"user_id" bson:"user_id,omitempty"`
	RoomID   string    `json:"room_id" bson:"room_id,omitempty"`
	FromDate time.Time `json:"from_date" bson:"from_date,omitempty"`
	TillDate time.Time `json:"till_date" bson:"till_date,omitempty"`
}

func NewBookingFromParams(params NewBookingParams) *Booking {

	return &Booking{
		RoomID:   params.RoomID,
		FromDate: params.FromDate,
		TillDate: params.TillDate,
	}
}

type BookingStatus int

const (
	StatusPending BookingStatus = iota
	StatusConfirmed
	StatusCanceled
)
