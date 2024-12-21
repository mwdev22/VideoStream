package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

const MaxPacketSize = 4096

type Config struct {
	IP   string
	Port string
}

func New() *Config {

	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	return &Config{
		IP:   os.Getenv("SERVER_IP"),
		Port: os.Getenv("SERVER_PORT"),
	}
}
