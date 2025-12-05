package main

import (
	"github.com/akansha204/cryptex-secretservice/internal/database"
	"github.com/akansha204/cryptex-secretservice/internal/routes"
	"github.com/gofiber/fiber/v2"
)

func main() {

	database.ConnectDB()

	app := fiber.New()
	routes.SetupRoutes(app)
	app.Listen(":3000")
}
