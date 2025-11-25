package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/ecelayes/pms-core/internal/adapter/handler/http"
	"github.com/ecelayes/pms-core/internal/adapter/repository/postgres"
	"github.com/ecelayes/pms-core/internal/core/service"
	"github.com/ecelayes/pms-core/internal/infra/database"
	"github.com/ecelayes/pms-core/internal/infra/server"
)

func main() {
	_ = godotenv.Load()

	dbPool, err := database.NewConnection(os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatalf("Database error: %v", err)
	}
	defer dbPool.Close()

	availRepo := postgres.NewAvailabilityRepo(dbPool)
	bookingRepo := postgres.NewBookingRepo(dbPool)
	authRepo := postgres.NewAuthRepo(dbPool)
	hotelRepo := postgres.NewHotelRepo(dbPool)

	bookingService := service.NewBookingService(bookingRepo)
	authService := service.NewAuthService(authRepo)

	availHandler := http.NewAvailabilityHandler(availRepo)
	bookingHandler := http.NewBookingHandler(bookingService)
	authHandler := http.NewAuthHandler(authService)
	hotelHandler := http.NewHotelHandler(hotelRepo)

	srv := server.New()
	
	srv.RegisterRoutes(availHandler, bookingHandler, authHandler, hotelHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("PMS Core running in port %s", port)
	if err := srv.Start(port); err != nil {
		log.Fatal(err)
	}
}
