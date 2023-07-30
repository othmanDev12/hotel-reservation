package db

import (
	"context"
	"github.com/hotel-reservation/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RoomStore interface {
	CreateRoom(ctx context.Context, room *domain.Room) (*domain.Room, error)
	GetRooms(ctx context.Context, filter bson.M) ([]*domain.Room, error)
}

type MongoRoomStore struct {
	client     *mongo.Client
	collection *mongo.Collection
	hotelStore HotelStore
}

func NewMongoRoomStore(client *mongo.Client, hotelStore HotelStore) *MongoRoomStore {
	return &MongoRoomStore{
		client:     client,
		collection: client.Database(Dbname).Collection("rooms"),
		hotelStore: hotelStore,
	}
}

func (s *MongoRoomStore) CreateRoom(ctx context.Context, room *domain.Room) (*domain.Room, error) {
	res, err := s.collection.InsertOne(ctx, room)
	if err != nil {
		return nil, err
	}
	room.Id = res.InsertedID.(primitive.ObjectID)
	filter := bson.M{"_id": room.HotelID}
	update := bson.M{"$push": bson.M{"rooms": room.Id}}
	if err = s.hotelStore.UpdateHotel(ctx, filter, update); err != nil {
		return nil, err
	}
	return room, nil
}

func (s *MongoRoomStore) GetRooms(ctx context.Context, filter bson.M) ([]*domain.Room, error) {
	res, err := s.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	var rooms []*domain.Room
	if err := res.All(ctx, &rooms); err != nil {
		return nil, err
	}
	return rooms, nil
}
