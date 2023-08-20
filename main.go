package main

import (
	"context"
	"flag"
	"github.com/gofiber/fiber/v2"
	"github.com/hotel-reservation/api"
	"github.com/hotel-reservation/api/middleware"
	"github.com/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func main() {
	if err := os.Setenv("JWT_SECRET", "mySecret"); err != nil {
		log.Fatal(err)
	}
	config := fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.JSON(map[string]string{"message": err.Error()})
		},
	}
	// create mongodb connection
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.UriDb))
	if err != nil {
		log.Fatal(err)
	}
	var (
		listenAddress = flag.String("listenAddr", ":5000", "address to listen")
		app           = fiber.New(config)
		apiv          = app.Group("auth")
		userStore     = db.NewMongoUserStore(client)
		hotelStore    = db.NewMongoHotelStore(client)
		roomStore     = db.NewMongoRoomStore(client, hotelStore)
		bookingStore  = db.NewMongoBookingStore(client)
		appV1         = app.Group("api/v1", middleware.JWTAuthentication(userStore))
		store         = &db.Store{
			HotelStore:   hotelStore,
			RoomStore:    roomStore,
			UserStore:    userStore,
			BookingStore: bookingStore,
		}
		userHandler  = api.NewUserHandler(userStore)
		authHandler  = api.NewAuthHandler(userStore)
		hotelHandler = api.NewHotelHandler(store)
		roomHandler  = api.NewRoomHandler(store)
	)

	// auth
	apiv.Post("/login", authHandler.HandleAuthenticate)

	// users
	appV1.Get("/user/:id", userHandler.HandleGetUser)
	appV1.Get("/users", userHandler.HandleGetUsers)
	appV1.Post("/user", userHandler.HandlePostUser)
	appV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	appV1.Put("/user/:id", userHandler.HandlePutUser)
	// hotels
	appV1.Get("/hotels", hotelHandler.HandleGetHotels)
	appV1.Get("/hotel/:id/rooms", hotelHandler.HandleGetRoomsByHotelId)
	appV1.Get("/hotel/:id", hotelHandler.HandleGetHotelById)
	appV1.Put("/hotel/:id", hotelHandler.HandlePutHotel)
	appV1.Delete("/hotel/:id", hotelHandler.HandleDeleteHotel)

	// booking
	appV1.Post("/room/:id/book", roomHandler.HandleRoomBooking)

	err2 := app.Listen(*listenAddress)
	if err2 != nil {
		log.Fatal(err2)
	}
}
