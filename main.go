package main

import (
	"log"
	"os"

	_ "Event-Explorer/routers"
	beego "github.com/beego/beego/v2/server/web"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	googleAPIKey := os.Getenv("GOOGLE_PLACES_API_KEY")
	if googleAPIKey == "" {
		log.Fatal("GOOGLE_PLACES_API_KEY is not set in the environment variables")
	}

	_ = googleAPIKey

	beego.Run()
}
