package configs

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type EnvConfiguration struct {
	MongoURI string
}

func LoadConfiguration() *EnvConfiguration {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return &EnvConfiguration{MongoURI: os.Getenv("MONGO_URI")}
}
