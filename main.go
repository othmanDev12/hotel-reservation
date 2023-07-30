package main

import (
	"context"
	"flag"
	"github.com/gofiber/fiber/v2"
	"github.com/hotel-reservation/api"
	"github.com/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func main() {
	config := fiber.Config{
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			return ctx.JSON(map[string]string{"message": err.Error()})
		},
	}
	listenAddress := flag.String("listenAddr", ":5000", "address to listen")
	app := fiber.New(config)
	appV1 := app.Group("api/v1")
	// create mongodb connection
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.UriDb))
	if err != nil {
		log.Fatal(err)
	}
	userStore := db.NewMongoUserStore(client, db.Dbname)
	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client, hotelStore)
	userHandler := api.NewUserHandler(userStore)
	hotelHandler := api.NewHotelHandler(hotelStore, roomStore)
	appV1.Get("/user/:id", userHandler.HandleGetUser)
	appV1.Get("/users", userHandler.HandleGetUsers)
	appV1.Post("/user", userHandler.HandlePostUser)
	appV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	appV1.Put("/user/:id", userHandler.HandlePutUser)
	appV1.Get("/hotels", hotelHandler.HandleGetHotels)
	err2 := app.Listen(*listenAddress)
	if err2 != nil {
		log.Fatal(err2)
	}
}
