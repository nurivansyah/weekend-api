package initializers

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVars() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()

		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}
}
