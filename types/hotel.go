package types

type Hotel struct {
	ID       string   `json:"id" bson:"_id,omitempty"`
	Name     string   `json:"name" bson:"name"`
	Location string   `json:"location" bson:"location"`
	Rooms    []string `json:"room_ids" bson:"rooms"`
}

type NewHotelParams struct {
	Name     string `json:"name" bson:"name"`
	Location string `json:"location" bson:"location"`
}
