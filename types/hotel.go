package types

type Hotel struct {
	ID       string   `json:"id" bson:"_id,omitempty"`
	Name     string   `json:"name" bson:"name"`
	Location string   `json:"location" bson:"location"`
	Rooms    []string `json:"room_ids" bson:"rooms"`
	Rating   Rating   `json:"rating" bson:"rating"`
}

type NewHotelParams struct {
	Name     string `json:"name" bson:"name"`
	Location string `json:"location" bson:"location"`
	Rating   Rating `json:"rating" bson:"rating"`
}

func NewHotelFromParams(params NewHotelParams) *Hotel {
	hotel := &Hotel{
		Name:     params.Name,
		Location: params.Location,
		Rating:   params.Rating,
	}
	return hotel
}
