package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type ErrorCode int

const (
	ErrorCodeOK ErrorCode = iota
	ErrorCodeInvalidRequest
	ErrorCodeInvalidHash
)

func (e ErrorCode) String() string {
	switch e {
	case ErrorCodeOK:
		return "health check OK"
	case ErrorCodeInvalidRequest:
		return "Invalid Request"
	case ErrorCodeInvalidHash:
		return "Invalid Hash"
	default:
		return "Unknown"
	}
}

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
