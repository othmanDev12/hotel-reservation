package db

import (
	"context"
	"github.com/hotel-reservation/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BookingStore interface {
	InsertBooking(ctx context.Context, booking *domain.Book) (*domain.Book, error)
}

type MongoBookingStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

func NewMongoBookingStore(client *mongo.Client) *MongoBookingStore {
	return &MongoBookingStore{
		client:     client,
		collection: client.Database(Dbname).Collection("bookings"),
	}
}

func (s *MongoBookingStore) InsertBooking(ctx context.Context, booking *domain.Book) (*domain.Book, error) {
	inserted, err := s.collection.InsertOne(ctx, booking)
	if err != nil {
		return nil, err
	}
	booking.Id = inserted.InsertedID.(primitive.ObjectID)
	return booking, nil
}
