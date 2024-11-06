package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	IP       string
	Port     string
	Protocol string
}

func LoadConfig() *Config {

	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	return &Config{
		IP:       os.Getenv("SERVER_IP"),
		Port:     os.Getenv("SERVER_PORT"),
		Protocol: os.Getenv("SERVER_PROTOCOL"),
	}
}
