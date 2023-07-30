package main

import (
	"context"
	"github.com/hotel-reservation/db"
	"github.com/hotel-reservation/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

var (
	client     *mongo.Client
	roomStore  db.RoomStore
	hotelStore db.HotelStore
	ctx        = context.Background()
)

func seedHotel(name string, location string) {
	hotel := domain.Hotel{
		Name:     name,
		Location: location,
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

	_, err := hotelStore.CreateHotel(ctx, &hotel)
	if err != nil {
		log.Fatal(err)
	}
	for _, room := range rooms {
		room.HotelID = hotel.Id
		_, err := roomStore.CreateRoom(ctx, &room)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	seedHotel("Ibisa", "Spain")
	seedHotel("Mogador", "Morocco")

}

func init() {
	var err error
	client, err = mongo.Connect(ctx, options.Client().ApplyURI(db.UriDb))
	if err != nil {
		log.Fatal(err)
	}

	if err := client.Database(db.Dbname).Drop(context.TODO()); err != nil {
		log.Fatal(err)
	}

	hotelStore = db.NewMongoHotelStore(client)
	roomStore = db.NewMongoRoomStore(client, hotelStore)
}
