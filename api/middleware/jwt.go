package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

func JWTAuthentication(ctx *fiber.Ctx) error {
	fmt.Println("JWT Authentication ............ ")
	token, ok := ctx.GetReqHeaders()["Authorization"]
	if !ok {
		return fmt.Errorf("unauthorized")
	}
	fmt.Println("token", token)
	return nil
}

func parseJWTToken(tokenString string) error {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Printf("Unexpected signing method: %v\n", token.Header["alg"])
			return nil, fmt.Errorf("unautorized")
		}
		secret := os.Getenv("JWT_SECRET")
		return secret, nil
	})
	if err != nil {
		fmt.Println("failed to jwt parse token: ", err)
		return fmt.Errorf("unautorized")
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		fmt.Println(claims)
	}
	return fmt.Errorf("unautorized")
}
