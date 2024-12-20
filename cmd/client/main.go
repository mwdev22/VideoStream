package main

import (
	"log"
	"net"
	"time"

	"github.com/mwdev22/Custom-Protocol-Server/internal/config"
)

func main() {
	cfg := config.New()
	tcpConn, err := net.Dial("tcp", cfg.IP+":"+cfg.Port)
	if err != nil {
		log.Fatalf("Error connecting to TCP server: %v", err)
	}
	defer tcpConn.Close()

	buf := make([]byte, 1024)
	n, err := tcpConn.Read(buf)
	if err != nil {
		log.Fatalf("error reading from TCP server: %v", err)
	}
	log.Printf("received: %s", string(buf[:n]))
	hash := buf[:n]

	udpAddr := cfg.IP + ":9000"

	udpConn, err := net.Dial("udp", udpAddr)
	if err != nil {
		log.Fatalf("error connecting to UDP server: %v", err)
	}
	defer udpConn.Close()

	for {
		data := hash
		_, err := udpConn.Write([]byte(data))
		if err != nil {
			log.Printf("error sending UDP data: %v", err)
		}
		time.Sleep(1 * time.Second) // Simulate frame rate
	}
}
