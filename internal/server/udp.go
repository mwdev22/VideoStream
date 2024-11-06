package server

import (
	"log"
	"net"
)

type serverUDP struct {
	IP   string
	Port string
	conn *net.UDPConn
}

func NewServerUDP(ip, port string) *serverUDP {
	return &serverUDP{IP: ip, Port: port}
}

func (s *serverUDP) Start() error {
	address := net.JoinHostPort(s.IP, s.Port)
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}
	s.conn = conn
	log.Printf("UDP server listening on %s\n", address)

	buf := make([]byte, 1024)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			return err
		}
		go s.handlePacket(buf[:n], addr)
	}
}

func (s *serverUDP) handlePacket(data []byte, addr *net.UDPAddr) {
}

func (s *serverUDP) Stop() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
