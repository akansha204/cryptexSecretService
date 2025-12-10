package main

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/akansha204/cryptex-secretservice/internal/database"
	"github.com/akansha204/cryptex-secretservice/internal/repository"
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
	startAutoPurgeJob()
	app.Listen(":3000")
}
func startAutoPurgeJob() {
	purgeRepo := repository.NewPurgeRepository()

	daysStr := os.Getenv("PURGE_DAYS")
	if daysStr == "" {
		daysStr = "7" // bydefault 7 days
	}
	days, _ := strconv.Atoi(daysStr)

	go func() {
		ticker := time.NewTicker(24 * time.Hour) // run every 24hrs
		defer ticker.Stop()

		for {
			<-ticker.C

			ctx := context.Background()

			olderThan := time.Duration(days) * 24 * time.Hour
			// olderThan := time.Duration(days) * time.Minute
			err := purgeRepo.PurgeOldData(ctx, olderThan)
			if err != nil {
				fmt.Println("[PURGE ERROR]", err)
			} else {
				fmt.Println("[PURGE] Completed successfully")
			}
		}
	}()
}
