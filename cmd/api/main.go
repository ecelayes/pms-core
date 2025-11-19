package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/ecelayes/pms-core/internal/adapter/handler/http"
	"github.com/ecelayes/pms-core/internal/adapter/repository/postgres"
	"github.com/ecelayes/pms-core/internal/infra/database"
	"github.com/ecelayes/pms-core/internal/infra/server"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using system variables")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5432/hotel_pms_db?sslmode=disable"
	}

	dbPool, err := database.NewConnection(dbURL)
	if err != nil {
		log.Fatalf("Error connecting to the DB: %v", err)
	}
	defer dbPool.Close()

	availRepo := postgres.NewAvailabilityRepo(dbPool)
	availHandler := http.NewAvailabilityHandler(availRepo)

	srv := server.New()
	srv.RegisterRoutes(availHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	
	log.Printf("PMS Core running in port %s", port)
	if err := srv.Start(port); err != nil {
		log.Fatal(err)
	}
}
