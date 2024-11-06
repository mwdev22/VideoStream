package main

import (
	"log"

	"github.com/mwdev22/Custom-Protocol-Server/internal/config"
	"github.com/mwdev22/Custom-Protocol-Server/internal/server"
)

func main() {
	cfg := config.LoadConfig()
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("failed to initialize server: %v", err)
	}

	if err := srv.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}

}
