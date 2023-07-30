package main

import (
	"context"
	"fmt"
	"github.com/hotel-reservation/db"
	"github.com/hotel-reservation/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

func main() {

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.UriDb))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.Dbname).Drop(context.TODO()); err != nil {
		log.Fatal(err)
	}

	hotelStore := db.NewMongoHotelStore(client)
	roomStore := db.NewMongoRoomStore(client, hotelStore)

	hotel := domain.Hotel{
		Name:     "ibisa",
		Location: "spain",
		Rooms:    []primitive.ObjectID{},
	}

	rooms := []domain.Room{
		{
			Type:  domain.SingleRoomType,
			Price: 99.9,
		},
		{
			Type:  domain.SeaSideRoomType,
			Price: 122.9,
		},
		{
			Type:  domain.DoubleRoomType,
			Price: 125.9,
		},
	}

	_, err = hotelStore.CreateHotel(context.Background(), &hotel)
	if err != nil {
		log.Fatal(err)
	}
	for _, room := range rooms {
		room.HotelID = hotel.Id
		insertedRoom, err := roomStore.CreateRoom(context.Background(), &room)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("inserted room value is: ", insertedRoom)

	}

}
