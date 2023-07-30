package db

import (
	"context"
	"github.com/hotel-reservation/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HotelStore interface {
	CreateHotel(ctx context.Context, hotel *domain.Hotel) (*domain.Hotel, error)
	UpdateHotel(ctx context.Context, filter bson.M, update bson.M) error
	GetHotels(ctx context.Context) ([]*domain.Hotel, error)
}

type MongoHotelStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoHotelStore(client *mongo.Client) *MongoHotelStore {
	return &MongoHotelStore{
		client:     client,
		collection: client.Database(Dbname).Collection("hotels"),
	}
}

func (m *MongoHotelStore) CreateHotel(ctx context.Context, hotel *domain.Hotel) (*domain.Hotel, error) {
	resp, err := m.collection.InsertOne(ctx, hotel)
	if err != nil {
		return nil, err
	}
	hotel.Id = resp.InsertedID.(primitive.ObjectID)
	return hotel, nil
}

func (m *MongoHotelStore) UpdateHotel(ctx context.Context, filter bson.M, update bson.M) error {
	_, err := m.collection.UpdateOne(ctx, filter, update)
	return err
}

func (m *MongoHotelStore) GetHotels(ctx context.Context) ([]*domain.Hotel, error) {
	resp, err := m.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var hotels []*domain.Hotel
	if err = resp.All(ctx, &hotels); err != nil {
		return nil, err
	}
	return hotels, nil
}
