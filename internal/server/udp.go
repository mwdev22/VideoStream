package server

import (
	"bytes"
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

	// packet format:
	// | Hash (32 bytes) | Video Data (4096 bytes) |
	buf := make([]byte, 4128)
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
	hash := data[:32]

	videoData := data[32:]

	if !bytes.Equal(hash, client.Hash) {
		log.Printf("Hash mismatch for client %s: expected %x, received %x", addr.String(), client.Hash, hash)
		client.Quit <- struct{}{}
		return
	}
	fmt.Printf("receive %d video bytes from %s\n", len(videoData), addr)

}

func (s *serverUDP) Stop() error {
	if s.conn != nil {
		return s.conn.Close()
	}
	return nil
}
