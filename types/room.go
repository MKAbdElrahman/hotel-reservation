package types

import (
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	ID          primitive.ObjectID `json:"id" bson:"_id,omitempty"`
	HotelID     primitive.ObjectID `json:"hotel_id" bson:"hotel_id,omitempty"`
	Number      string             `json:"number" bson:"number"`
	Floor       int                `json:"floor" bson:"floor"`
	Type        RoomType           `json:"type" bson:"type"`
	Description string             `json:"description" bson:"description"`
	Price       float64            `json:"price" bson:"price"`
	Occupied    bool               `json:"occupied" bson:"occupied"`
}

type NewRoomParams struct {
	Number      string   `json:"number" bson:"number"`
	Floor       int      `json:"floor" bson:"floor"`
	Type        RoomType `json:"type" bson:"type"`
	Description string   `json:"description" bson:"description"`
	Price       float64  `json:"price" bson:"price"`
	Occupied    bool     `json:"occupied" bson:"occupied"`
}

type RoomType int

const (
	StandardRoom RoomType = iota
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
