package main

import (
	"github.com/akansha204/cryptex-secretservice/internal/database"
	"github.com/akansha204/cryptex-secretservice/internal/routes"

	"github.com/akansha204/cryptex-secretservice/internal/utils"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load(".env")
	if err := utils.Init(); err != nil {
		panic(err)
	}

	database.ConnectDB()

	app := fiber.New()
	routes.SetupRoutes(app)
	app.Listen(":3000")
}
