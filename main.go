package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/gin-gonic/gin"
	"github.com/mkabdelrahman/hotel-reservation/api"
	"github.com/mkabdelrahman/hotel-reservation/business"
	"github.com/mkabdelrahman/hotel-reservation/db"
	"github.com/mkabdelrahman/hotel-reservation/middleware"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName    = "hotel-reservation"
	userColl  = "users"
	hotelColl = "hotels"
	roomColl  = "rooms"
)

type Config struct {
	MONGODB_URI string `conf:"default:mongodb://localhost:27017,flag:dburi,env:DB_URI"`
	Port        int    `conf:"default:8080,env:PORT"`
}

func main() {

	// CONFIG
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

	// Mongodb

	// fmt.Println(cfg.MONGODB_URI)

	ctx := context.Background()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MONGODB_URI))

	if err != nil {
		log.Fatal(err)
	}

	userStore := db.NewMongoUserStore(client, dbName, userColl)
	hotelStore := db.NewMongoHotelStore(client, dbName, hotelColl)
	roomStore := db.NewMongoRoomStore(client, dbName, roomColl)

	hotelManager := business.NewManager(hotelStore, roomStore)

	// Handlers Initialization

	userHandler := api.NewUserHandler(userStore)
	hotelHandler := api.NewHotelHandler(hotelManager)

	// Router
	engine := gin.New()

	engine.Use(middleware.DefaultStructuredLogger())
	engine.Use(gin.Recovery())

	v1 := engine.Group("/api/v1")

	v1.GET("/user/:id", userHandler.HandleGetUser)
	v1.DELETE("/user/:id", userHandler.HandleDeleteUser)

	v1.GET("/user", userHandler.HandleGetUsers)
	v1.POST("/user", userHandler.HandlePostUser)
	v1.PUT("/user/:id", userHandler.HandleUpdateUser)

	v1.GET("/hotel", hotelHandler.HandleGetHotels)

	v1.GET("/hotel/:id", hotelHandler.HandleGetHotel)

	v1.GET("/hotel/:id/rooms", hotelHandler.HandleGetHotelRooms)

	v1.GET("/hotel/search", hotelHandler.HandleHotelSearch)

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: engine,
	}

	chanErrors := make(chan error)
	go func() {
		chanErrors <- runServer(server)
	}()

	chanSignals := make(chan os.Signal, 1)
	signal.Notify(chanSignals, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-chanErrors:
		log.Fatalf("Error while starting server %s", err)
	case s := <-chanSignals:
		log.Printf("Shutting down server in few seconds due to %s", s)
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := Close(ctx, server); err != nil {
			log.Fatal("Server forced to shutdown: ", err)
		}
		log.Print("Server exiting gracefully")
	}

}

func runServer(server *http.Server) error {
	return server.ListenAndServe()
}

func Close(ctx context.Context, server *http.Server) error {
	return server.Shutdown(ctx)
}
