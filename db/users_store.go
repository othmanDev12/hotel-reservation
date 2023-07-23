package db

import (
	"context"
	"github.com/hotel-reservation/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const userCol = "user"

// UserStore create a new user store interface
type UserStore interface {
	GetUserById(ctx context.Context, id string) (*domain.User, error)
	GetUsers(ctx context.Context) ([]*domain.User, error)
	CreateUser(ctx context.Context, user *domain.User) (*domain.User, error)
}

// MongoUserStore MongoStore implementation of this user interface
type MongoUserStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewMongoUserStore NewMongoStore constructor
func NewMongoUserStore(client *mongo.Client) *MongoUserStore {
	return &MongoUserStore{
		client:     client,
		collection: client.Database(Dbname).Collection(userCol),
	}
}

func (m *MongoUserStore) GetUserById(ctx context.Context, id string) (*domain.User, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	var user domain.User
	if err := m.collection.FindOne(ctx, bson.M{"_id": oid}).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (m *MongoUserStore) CreateUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	res, err := m.collection.InsertOne(ctx, user)
	if err != nil {
		return nil, err
	}
	user.Id = res.InsertedID.(primitive.ObjectID)
	return user, nil
}

func (m *MongoUserStore) GetUsers(ctx context.Context) ([]*domain.User, error) {
	cur, err := m.collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	var users []*domain.User
	if err := cur.All(ctx, &users); err != nil {
		return []*domain.User{}, nil
	}
	return users, nil
}
