package types

import "go.mongodb.org/mongo-driver/bson/primitive"

type Hotel struct {
	ID       primitive.ObjectID   `json:"id" bson:"_id,omitempty"`
	Name     string               `json:"name" bson:"name"`
	Location string               `json:"location" bson:"location"`
	Rooms    []primitive.ObjectID `json:"room_ids" bson:"rooms"`
}

type NewHotelParams struct {
	Name     string `json:"name" bson:"name"`
	Location string `json:"location" bson:"location"`
}
