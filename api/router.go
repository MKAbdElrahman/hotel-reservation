package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/mkabdelrahman/hotel-reservation/api/handlers"
	"github.com/mkabdelrahman/hotel-reservation/business"
	"github.com/mkabdelrahman/hotel-reservation/db"
	"github.com/mkabdelrahman/hotel-reservation/middleware"
	"go.mongodb.org/mongo-driver/mongo"
)

func setupRouter(client *mongo.Client) *gin.Engine {

	userStore := db.NewMongoUserStore(client, dbName, userColl)
	hotelStore := db.NewMongoHotelStore(client, dbName, hotelColl)
	roomStore := db.NewMongoRoomStore(client, dbName, roomColl)
	bookingStore := db.NewMongoBookingStore(client, dbName, bookingColl)

	hotelManager := business.NewManager(userStore, hotelStore, roomStore, bookingStore)

	// logger

	errorLogger := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	authHandler := handlers.NewAuthHandler(hotelManager, errorLogger)
	userHandler := handlers.NewUserHandler(hotelManager, errorLogger)
	hotelHandler := handlers.NewHotelHandler(hotelManager, errorLogger)
	bookingHandler := handlers.NewBookingHandler(hotelManager, errorLogger)

	engine := gin.New()

	engine.Use(middleware.Logger)
	engine.Use(gin.Recovery())

	v1 := engine.Group("/api/v1")

	// v1.Use(middleware.AuthMiddleware())

	adminRoutes := engine.Group("/admin", middleware.AdminOnlyMiddleware(hotelManager))

	{
		adminRoutes.GET("/dashboard", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "Admin dashboard"})
		})
	}

	engine.POST("/api/auth", authHandler.HandleAuthenticate)

	// users
	v1.GET("/user/:id", userHandler.HandleGetUser)
	v1.DELETE("/user/:id", userHandler.HandleDeleteUser)
	v1.GET("/user/:id/bookings", userHandler.HandleGetUserBookings)

	v1.GET("/user", userHandler.HandleGetUsers)
	v1.POST("/user", userHandler.HandlePostUser)
	v1.PUT("/user/:id", userHandler.HandleUpdateUser)

	// hotel
	v1.GET("/hotel", hotelHandler.HandleGetHotels)

	v1.GET("/hotel/:id", hotelHandler.HandleGetHotel)

	v1.GET("/hotel/:id/rooms", hotelHandler.HandleGetHotelRooms)

	v1.GET("/hotel/search", hotelHandler.HandleHotelSearch)

	// booking

	v1.GET("/booking/:id", bookingHandler.HandleGetBooking)
	v1.GET("/booking", bookingHandler.HandleGetBookings)
	v1.POST("/booking", bookingHandler.HandlePostBooking)
	// used to change the booking status
	v1.DELETE("/booking/:id", bookingHandler.HandleCancelBooking)

	return engine
}
