package main

import (
	"github.com/gofiber/fiber/v2"
	"log"
)

func main() {
	app := fiber.New()
	app.Get("/foo", handleFoo)
	err := app.Listen(":5000")
	if err != nil {
		log.Fatal(err)
		return
	}
}

func handleFoo(ctx *fiber.Ctx) error {
	return ctx.JSON(map[string]string{"msg": "is working fine ................"})
}
