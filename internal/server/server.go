package server

import (
	"errors"

	"github.com/mwdev22/Custom-Protocol-Server/internal/config"
)

type Server interface {
	Start() error
	Stop() error
}

func NewServer(cfg *config.Config) (Server, error) {
	switch cfg.Protocol {
	case "TCP":
		return NewServerTCP(cfg.IP, cfg.Port), nil
	case "UDP":
		return NewServerUDP(cfg.IP, cfg.Port), nil
	default:
		return nil, errors.New("unsupported protocol")
	}
}
