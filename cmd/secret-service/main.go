package main

import (
	"github.com/akansha204/cryptex-secretservice/internal/database"
	"github.com/gofiber/fiber/v2"
)

func main() {

	database.ConnectDB()

	app := fiber.New()
	app.Listen(":3000")
}
