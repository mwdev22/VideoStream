package main

import (
	"log"

	"github.com/mwdev22/Custom-Protocol-Server/internal/config"
	"github.com/mwdev22/Custom-Protocol-Server/internal/server"
)

func main() {
	cfg := config.New()
	srv := server.NewTCP(cfg.IP, cfg.Port)

	if err := srv.Start(); err != nil {
		log.Fatalf("server error: %v", err)
	}

	udp := server.NewUDP(cfg.IP, "9000", srv)
	if err := udp.Start(); err != nil {
		log.Fatalf("server error :%v", err)
	}

}
