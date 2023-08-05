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
	userStore  db.UserStore
	ctx        = context.Background()
)

func seedUser(firstName string, lastName string, email string) {
	userParam := domain.CreateUserParams{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  "othmanDev12@",
	}
	user, err := domain.NewCreateUser(userParam)
	if err != nil {
		log.Fatal(err)
	}
	_, err = userStore.CreateUser(context.TODO(), user)
	if err != nil {
		log.Fatal(err)
	}
}

func seedHotel(name string, location string, rating int) {
	hotel := domain.Hotel{
		Name:     name,
		Location: location,
		Rooms:    []primitive.ObjectID{},
		Rating:   rating,
	}

	rooms := []domain.Room{
		{
			Size:  "small",
			Price: 99.9,
		},
		{
			Size:  "medium",
			Price: 122.9,
		},
		{
			Size:  "kingsize",
			Price: 299.9,
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
	seedHotel("Ibisa", "Spain", 4)
	seedHotel("Mogador", "Morocco", 2)
	seedUser("othman", "Test", "othman@gmail.com")

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
	userStore = db.NewMongoUserStore(client)
}
