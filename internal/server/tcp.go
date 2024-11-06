package server

import (
	"log"
	"net"
)

type serverTCP struct {
	IP       string
	Port     string
	listener net.Listener
}

func NewServerTCP(ip, port string) *serverTCP {
	return &serverTCP{IP: ip, Port: port}
}

func (s *serverTCP) Start() error {
	address := net.JoinHostPort(s.IP, s.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	s.listener = listener
	log.Printf("TCP server listening on %s\n", address)

	for {
		conn, err := listener.Accept()
		if err != nil {
			return err
		}
		go s.handleConnection(conn)
	}
}

func (s *serverTCP) handleConnection(conn net.Conn) {
	defer conn.Close()
}

func (s *serverTCP) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}
	return nil
}
