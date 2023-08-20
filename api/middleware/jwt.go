package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/hotel-reservation/db"
	"os"
	"time"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	// decorator pattern
	return func(ctx *fiber.Ctx) error {
		fmt.Println("JWT Authentication ............ ")
		token, ok := ctx.GetReqHeaders()["Authorization"]
		if !ok {
			return fmt.Errorf("unauthorized")
		}
		claims, err := ValidToken(token)
		if err != nil {
			return err
		}
		expiresFloat := claims["expires"].(float64)
		expires := int64(expiresFloat)

		if time.Now().Unix() > expires {
			return fmt.Errorf("token expired")
		}

		userId := claims["id"].(string)
		user, err := userStore.GetUserById(ctx.Context(), userId)
		if err != nil {
			return fmt.Errorf("unautorized")
		}
		// set the current user inside fiber context
		ctx.Context().SetUserValue("user", user)
		return ctx.Next()
	}
}

func ValidToken(tokenString string) (jwt.MapClaims, error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Printf("Unexpected signing method: %v\n", token.Header["alg"])
			return nil, fmt.Errorf("unautorized")
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to jwt parse token: ", err)
		return nil, fmt.Errorf("unautorized")
	}
	if !token.Valid {
		fmt.Println("invalid token: ", err)
		return nil, fmt.Errorf("unautorized")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unautorized")
	}
	return claims, nil
}
