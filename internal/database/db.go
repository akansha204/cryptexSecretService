package database

import (
	"context"
	"database/sql"
	"log"
	"os"
	"time"

	"github.com/akansha204/cryptex-secretservice/internal/models"
	"github.com/joho/godotenv"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

var DB *bun.DB

func ConnectDB() {
	_ = godotenv.Load()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL missing in .env")
	}
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	DB = bun.NewDB(sqldb, pgdialect.New())
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// ctx := context.Background()

	if err := DB.PingContext(ctx); err != nil {
		log.Fatal("DB connection failed:", err)
	}
	log.Println("Connected to PostgreSQL")

	createTables(ctx)

}

func createTables(ctx context.Context) {

	_, err := DB.NewCreateTable().
		Model((*models.Project)(nil)).
		IfNotExists().
		Exec(ctx)

	if err != nil {
		log.Fatal("Error creating projects table:", err)
	}

	_, err = DB.NewCreateTable().
		Model((*models.Secret)(nil)).
		IfNotExists().
		Exec(ctx)

	if err != nil {
		log.Fatal("Error creating secrets table:", err)
	}

	_, err = DB.NewCreateTable().
		Model((*models.AuditLog)(nil)).
		IfNotExists().
		Exec(ctx)

	if err != nil {
		log.Fatal("Error creating audit logs table:", err)
	}

}
