package api

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/hotel-reservation/db"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	userStore db.UserStore
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func NewAuthHandler(userStore db.UserStore) *AuthHandler {
	return &AuthHandler{userStore: userStore}
}

func (h *AuthHandler) HandleAuthenticate(ctx *fiber.Ctx) error {
	var params AuthParams
	if err := ctx.BodyParser(&params); err != nil {
		return err
	}

	user, err := h.userStore.GetUserByEmail(ctx.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return fmt.Errorf("invalid Credentials")
		}
	}

	if err = bcrypt.CompareHashAndPassword([]byte(user.EncryptedPassword), []byte(params.Password)); err != nil {
		return fmt.Errorf("invalid Credentials")
	}

	fmt.Println(user)

	return nil
}
