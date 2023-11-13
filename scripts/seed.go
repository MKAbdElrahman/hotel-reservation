package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

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

type Config struct {
	MONGODB_URI string `conf:"default:mongodb://localhost:27017,flag:dburi,env:DB_URI"`
}

func main() {
	// Config
	var cfg Config
	help, err := conf.Parse("", &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return
		}
		log.Fatalf("Error parsing configuration: %v\n", err)
		return
	}

	// DB Client
	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MONGODB_URI))
	if err != nil {
		log.Fatal(err)
	}

	// SCRIPT

	hotelStore := db.NewMongoHotelStore(client, dbName, hotelColl)
	roomStore := db.NewMongoRoomStore(client, dbName, roomColl)

	manager := business.NewManager(hotelStore, roomStore)

	hotelID, err := manager.AddNewHotel(ctx, types.NewHotelParams{
		Name:     "Dolco",
		Location: "Cairo",
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

	printHotelWithRooms(manager, ctx, hotelID.Hex())
}

func printHotelWithRooms(manager *business.Manager, ctx context.Context, hotelID string) {
	// Get hotel information
	hotel, err := manager.HotelStore.GetHotel(ctx, hotelID)
	if err != nil {
		log.Printf("Error fetching hotel %s: %v", hotelID, err)
		return
	}

	// Print hotel information
	fmt.Printf("Hotel: %s (ID: %s, Location: %s)\n", hotel.Name, hotel.ID, hotel.Location)

	// Get rooms for the hotel
	rooms, err := manager.RoomStore.GetRoomsByHotelID(ctx, hotelID)
	if err != nil {
		log.Printf("Error fetching rooms for hotel %s: %v", hotelID, err)
		return
	}

	// Print rooms for the hotel
	fmt.Println("Rooms:")
	for _, room := range rooms {
		fmt.Printf("  Room %s\n", room.Number)
		fmt.Printf("    Type: %s\n", room.Type)
		fmt.Printf("    Floor: %d\n", room.Floor)
		fmt.Printf("    Price: %.2f\n", room.Price)
		fmt.Printf("    Occupied: %t\n", room.Occupied)
		fmt.Printf("    Description: %s\n", room.Description)
	}

	fmt.Println(strings.Repeat("-", 40)) // Separator line
}
