package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/hotel-reservation/db"
	"github.com/hotel-reservation/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

func insertUserTest(t *testing.T, userStore db.UserStore) *domain.User {
	userParam := domain.CreateUserParams{
		FirstName: "othman",
		LastName:  "test",
		Email:     "othman@gmail.com",
		Password:  "othmanDev12@",
	}
	user, err := domain.NewCreateUser(userParam)
	if err != nil {
		t.Fatal(err)
	}
	_, err = userStore.CreateUser(context.TODO(), user)
	if err != nil {
		t.Fatal(err)
	}
	return user

}

func TestAuthenticateSuccess(t *testing.T) {
	tdb := Setup(t)
	defer tdb.teardown(t)

	insertedUser := insertUserTest(t, tdb.UserStore)
	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/login", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "othman@gmail.com",
		Password: "othmanDev12@",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("Unexpected status code of 200 but got %d", resp.StatusCode)
	}

	var authResponse AuthResp
	if err := json.NewDecoder(resp.Body).Decode(&authResponse); err != nil {
		t.Fatal(err)
	}

	if len(authResponse.Token) == 0 {
		t.Fatal("expected the JWT token to be present in the auth response")
	}

	if insertedUser.Id != authResponse.User.Id {
		t.Fatal("expected the user to be the inserted user")
	}
}

func TestAuthenticateWithWrongPassword(t *testing.T) {
	tdb := Setup(t)
	defer tdb.teardown(t)

	insertUserTest(t, tdb.UserStore)
	app := fiber.New()
	authHandler := NewAuthHandler(tdb.UserStore)
	app.Post("/login", authHandler.HandleAuthenticate)

	params := AuthParams{
		Email:    "othman@gmail.com",
		Password: "othmanDev12@1111",
	}

	b, _ := json.Marshal(params)
	req := httptest.NewRequest("POST", "/login", bytes.NewReader(b))
	req.Header.Add("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatal(err)
	}
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Unexpected status code of 200 but got %d", resp.StatusCode)
	}

	var genericResp GenericResp
	if err := json.NewDecoder(resp.Body).Decode(&genericResp); err != nil {
		t.Fatal(err)
	}

	if genericResp.Type != "error" {
		t.Fatalf("expected generic resp type to be error but u got %s", genericResp.Type)
	}

	if genericResp.Msg != "invalid credentials" {
		t.Fatalf("expected generic resp msg to be <invalid credentials> %s", genericResp.Msg)
	}
}
