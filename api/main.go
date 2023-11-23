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
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	dbName      = "hotel-reservation"
	userColl    = "users"
	hotelColl   = "hotels"
	roomColl    = "rooms"
	bookingColl = "bookings"
)

const serverShutdownTimeout = 5 * time.Second

type config struct {
	MONGODB_URI string `conf:"default:mongodb://localhost:27017,flag:dburi,env:DB_URI"`
	Port        int    `conf:"default:8080,env:PORT"`
}

func main() {

	// CONFIGS

	var cfg config
	help, err := conf.Parse("", &cfg)
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return
		}
		log.Fatalf("Error parsing configuration: %v\n", err)
	}

	// DATABASE
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(cfg.MONGODB_URI))
	if err != nil {
		log.Fatal(err)
	}

	// SERVER
	engine := setupRouter(client)
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Port),
		Handler: engine,
	}

	err = run(cfg, client, server)
	if err != nil {
		log.Fatal(err)
	}
}

func run(cfg config, client *mongo.Client, server *http.Server) error {

	chanErrors := make(chan error)
	go func() {
		chanErrors <- server.ListenAndServe()
	}()

	chanSignals := make(chan os.Signal, 1)
	signal.Notify(chanSignals, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-chanErrors:
		log.Fatalf("Error while starting server %s", err)
		return err
	case s := <-chanSignals:
		log.Printf("Shutting down server in few seconds due to %s", s)
		ctx, cancel := context.WithTimeout(context.Background(), serverShutdownTimeout)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatal("Server forced to shutdown: ", err)
			return err
		}
		log.Print("Server exiting gracefully")
	}
	return nil
}
