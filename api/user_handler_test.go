package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/hotel-reservation/db"
	"github.com/hotel-reservation/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http/httptest"
	"testing"
)

type TestDb struct {
	db.UserStore
}

func (d *TestDb) teardown(t *testing.T) {
	if err := d.UserStore.Drop(context.TODO()); err != nil {
		t.Fatal(err)
	}
}

func Setup(t *testing.T) *TestDb {
	// create mongodb connection
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(db.UriDb))
	if err != nil {
		t.Fatal(err)
	}
	return &TestDb{
		UserStore: db.NewMongoUserStore(client),
	}
}

func TestPostUser(t *testing.T) {
	tdb := Setup(t)
	defer tdb.teardown(t)

	app := fiber.New()
	userHandler := NewUserHandler(tdb.UserStore)
	app.Post("/", userHandler.HandlePostUser)

	params := domain.CreateUserParams{
		FirstName: "JohnLKGFLG",
		LastName:  "SmithJFGJKGF",
		Email:     "johngmail.com",
		Password:  "passwordDb@123",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Error(err)
	}
	var user domain.User
	if err = json.NewDecoder(resp.Body).Decode(&user); err != nil {
		t.Error(err)
	}
	fmt.Println(req.Body)

	if user.FirstName != params.FirstName {
		t.Errorf("expected first name %s but got %s", params.FirstName, user.FirstName)
	}
	if user.LastName != params.LastName {
		t.Errorf("expected first name %s but got %s", params.LastName, user.LastName)
	}
	if user.Email != params.Email {
		t.Errorf("expected first name %s but got %s", params.Email, user.Email)
	}

}
