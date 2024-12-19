package server

import (
	"log"
	"net"
)

type serverUDP struct {
	IP        string
	Port      string
	conn      *net.UDPConn
	parentTCP *serverTCP
}

func NewServerUDP(ip, port string, parentTCP *serverTCP) *serverUDP {
	return &serverUDP{
		IP:        ip,
		Port:      port,
		parentTCP: parentTCP}
}

func (s *serverUDP) Start() error {
	address := net.JoinHostPort(s.IP, s.Port)
	udpAddr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}

	// make server universal for both tcp/udp combined usage and server interface
	var handleFunc packetHandleFunc
	if s.parentTCP == nil {
		handleFunc = s.handlePacketStandard
	} else {
		handleFunc = s.handlePacketIncludintParentTCP
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
		go handleFunc(buf[:n], addr)
	}
}

func (s *serverUDP) handlePacketIncludintParentTCP(data []byte, addr *net.UDPAddr) {
}

func (s *serverUDP) handlePacketStandard(data []byte, addr *net.UDPAddr) {
	log.Printf("Received from %s: %s\n", addr.String(), string(data))
}

func (s *serverUDP) Stop() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}

type packetHandleFunc func(data []byte, addr *net.UDPAddr)
