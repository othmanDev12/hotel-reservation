package api

import (
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hotel-reservation/db"
	"github.com/hotel-reservation/domain"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"os"
	"time"
)

type AuthHandler struct {
	userStore db.UserStore
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResp struct {
	User  *domain.User `json:"user"`
	Token string       `json:"access_token"`
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

	if !domain.CompareEncrAndPlainPassword(user.EncryptedPassword, params.Password) {
		return fmt.Errorf("invalid Credentials")
	}

	token, err := CreateJwtFromUser(user)

	if err != nil {
		return fmt.Errorf("invalid Token")
	}
	resp := AuthResp{
		User:  user,
		Token: token,
	}
	return ctx.JSON(resp)
}

func CreateJwtFromUser(user *domain.User) (string, error) {
	now := time.Now()
	expires := now.Add(time.Hour * 4).Unix()

	claims := jwt.MapClaims{
		"id":      user.Id,
		"email":   user.Email,
		"expires": expires,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		log.Fatal("key is of invalid type")
		return "", err
	}

	return tokenString, nil

}
