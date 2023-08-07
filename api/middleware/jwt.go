package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

func JWTAuthentication(ctx *fiber.Ctx) error {
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

	return ctx.Next()
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
