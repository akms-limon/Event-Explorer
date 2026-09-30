package main

import (
	"log"
	"os"

	"Event-Explorer/clients"
	"Event-Explorer/services"

	_ "Event-Explorer/routers"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env file not found")
	}

	googleAPIKey := os.Getenv("GOOGLE_PLACES_API_KEY")

	if googleAPIKey == "" {
		log.Fatal("GOOGLE_PLACES_API_KEY is not configured")
	}

	googlePlacesClient := clients.NewGooglePlacesClient(googleAPIKey)
	locationService := services.NewLocationService(googlePlacesClient)

	_ = locationService

	beego.Run()
}
