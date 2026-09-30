package main

import (
	"log"
	"os"

	"Event-Explorer/clients"
	"Event-Explorer/routers"
	"Event-Explorer/services"

	beego "github.com/beego/beego/v2/server/web"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("warning: .env file not found")
	}

	googleAPIKey := os.Getenv("GOOGLE_PLACES_API_KEY")
	ticketmasterAPIKey := os.Getenv("TICKETMASTER_API_KEY")

	if googleAPIKey == "" {
		log.Fatal("GOOGLE_PLACES_API_KEY is not configured")
	}

	if ticketmasterAPIKey == "" {
		log.Fatal("TICKETMASTER_API_KEY is not configured")
	}

	googlePlacesClient := clients.NewGooglePlacesClient(
		googleAPIKey,
	)

	ticketmasterClient := clients.NewTicketmasterClient(
		ticketmasterAPIKey,
	)

	locationService := services.NewLocationService(
		googlePlacesClient,
	)

	eventService := services.NewEventService(
		ticketmasterClient,
	)

	routers.Register(
		locationService,
		eventService,
	)

	beego.Run()
}