package main

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/ardanlabs/conf/v3"
	"github.com/mkabdelrahman/hotel-reservation/business"
	"github.com/mkabdelrahman/hotel-reservation/db"
	"github.com/mkabdelrahman/hotel-reservation/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName    = "hotel-reservation"
	hotelColl = "hotels"
	roomColl  = "rooms"
)

var (
	cfg    Config
	client *mongo.Client

	hotelStore *db.MongoHotelStore

	roomStore *db.MongoRoomStore

	manager *business.Manager
)

type Config struct {
	MONGODB_URI string `conf:"default:mongodb://localhost:27017,flag:dburi,env:DB_URI"`
}

func Init(ctx context.Context) {

	help, err := conf.Parse("", &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return
		}
		log.Fatalf("Error parsing configuration: %v\n", err)
		return
	}

	client, err = mongo.Connect(ctx, options.Client().ApplyURI(cfg.MONGODB_URI))
	if err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client, dbName, hotelColl)
	roomStore = db.NewMongoRoomStore(client, dbName, roomColl)

	hotelStore.Drop(ctx)
	roomStore.Drop(ctx)

	manager = business.NewManager(hotelStore, roomStore)

}
func main() {

	ctx := context.Background()

	Init(ctx)

	seedHotel(ctx, "Dolcica", "Madrid")
	seedHotel(ctx, "Lapache", "Paris")

}

func seedHotel(ctx context.Context, name string, location string) {

	hotelID, err := manager.AddNewHotel(ctx, types.NewHotelParams{
		Name:     name,
		Location: location,
	})

	if err != nil {
		log.Fatal(err)
	}

	rooms := []types.NewRoomParams{
		{
			Number:      "101",
			Floor:       1,
			Type:        types.DeluxeRoom,
			Price:       150.0,
			Occupied:    false,
			Description: "Spacious room with a city view.",
		},
		{
			Number:      "202",
			Floor:       2,
			Type:        types.StandardRoom,
			Price:       100.0,
			Occupied:    false,
			Description: "Cozy room with modern amenities.",
		},
		{
			Number:      "305",
			Floor:       3,
			Type:        types.SuiteRoom,
			Price:       200.0,
			Occupied:    false,
			Description: "Luxurious suite with a balcony and sea view.",
		},
		{
			Number:      "410",
			Floor:       4,
			Type:        types.DeluxeRoom,
			Price:       160.0,
			Occupied:    false,
			Description: "Elegant room with premium furnishings.",
		},
	}

	for _, room := range rooms {
		_, err := manager.AddNewRoom(ctx, room, hotelID)
		if err != nil {
			log.Fatal(err)
		}
	}
}
