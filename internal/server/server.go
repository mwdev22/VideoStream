package server

import (
	"errors"

	"github.com/mwdev22/Custom-Protocol-Server/internal/config"
)

type Server interface {
	Start() error
	Stop() error
}

func NewServer(cfg *config.Config, serverType string) (Server, error) {
	switch serverType {
	case "tcp":
		return NewServerTCP(cfg.IP, cfg.Port), nil
	case "udp":
		return NewServerUDP(cfg.IP, cfg.Port, nil), nil
	default:
		return nil, errors.New("unsupported protocol")
	}
}
