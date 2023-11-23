package types

import (
	"fmt"
)

type Room struct {
	ID          string   `json:"id" bson:"_id,omitempty"`
	HotelID     string   `json:"hotel_id" bson:"hotel_id,omitempty"`
	Number      string   `json:"number" bson:"number"`
	Floor       int      `json:"floor" bson:"floor"`
	Type        RoomType `json:"type" bson:"type"`
	Description string   `json:"description" bson:"description"`
	Price       float64  `json:"price" bson:"price"`
	Occupied    bool     `json:"occupied" bson:"occupied"`
}

type NewRoomParams struct {
	Number      string   `json:"number" bson:"number"`
	Floor       int      `json:"floor" bson:"floor"`
	Type        RoomType `json:"type" bson:"type"`
	Description string   `json:"description" bson:"description"`
	Price       float64  `json:"price" bson:"price"`
	Occupied    bool     `json:"occupied" bson:"occupied"`
}

func NewRoomFromParams(params NewRoomParams) *Room {
	room := &Room{
		Number:      params.Number,
		Floor:       params.Floor,
		Type:        params.Type,
		Description: params.Description,
		Price:       params.Price,
		Occupied:    params.Occupied,
	}
	return room
}

type RoomType int

const (
	_ RoomType = iota
	StandardRoom
	DeluxeRoom
	SuiteRoom
)

func (rt RoomType) String() string {
	switch rt {
	case StandardRoom:
		return "Standard Room"
	case DeluxeRoom:
		return "Deluxe Room"
	case SuiteRoom:
		return "Suite Room"
	default:
		return fmt.Sprintf("Unknown RoomType: %d", rt)
	}
}
