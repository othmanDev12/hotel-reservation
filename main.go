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

const (
	uriDb  = "mongodb://localhost:27017"
	dbname = "hotel-reservation"
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
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uriDb))
	if err != nil {
		log.Fatal(err)
	}
	userStore := db.NewMongoUserStore(client, dbname)
	userHandler := api.NewUserHandler(userStore)
	appV1.Get("/user/:id", userHandler.HandleGetUser)
	appV1.Get("/users", userHandler.HandleGetUsers)
	appV1.Post("/user", userHandler.HandlePostUser)
	appV1.Delete("/user/:id", userHandler.HandleDeleteUser)
	appV1.Put("/user/:id", userHandler.HandlePutUser)
	err2 := app.Listen(*listenAddress)
	if err2 != nil {
		log.Fatal(err2)
	}
}
