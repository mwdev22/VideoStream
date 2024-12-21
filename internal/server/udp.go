package server

import (
	"fmt"
	"log"
	"net"
)

type serverUDP struct {
	IP        string
	Port      string
	conn      *net.UDPConn
	parentTCP *serverTCP
}

func NewUDP(ip, port string, parentTCP *serverTCP) *serverUDP {
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

	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		return err
	}
	s.conn = conn
	log.Printf("UDP server listening on %s\n", address)

	// 4 bytes of hash, 4096 of video data
	buf := make([]byte, 4100)
	for {
		n, addr, err := conn.ReadFromUDP(buf)
		if err != nil {
			return err
		}
		go s.handlePacket(buf[:n], addr)
	}
}

func (s *serverUDP) handlePacket(data []byte, addr *net.UDPAddr) {
	client := s.parentTCP.retrieveClient(addr.IP.String())
	if client == nil {
		return
	}
	hash := data[:4]
	if string(hash) != string(client.Hash) {
		fmt.Println(hash, "\n", client.Hash)
		return
	}
	fmt.Printf("received %d bytes from %s\n", len(data), addr)
}

func (s *serverUDP) Stop() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
