package middleware

import (
	"log"
	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from .env file
func LoadEnv() error {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using default environment variables")
		return err
	}

	log.Println("Successfully loaded .env file")
	return nil
}
