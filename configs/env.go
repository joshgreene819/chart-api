package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Configuration struct {
	MongoURI string
}

func LoadConfiguration() *Configuration {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return &Configuration{MongoURI: os.Getenv("MONGO_URI")}
}
