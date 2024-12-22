package main

import (
	"encoding/binary"
	"log"
	"net"
	"time"

	cam "github.com/mwdev22/Custom-Protocol-Server/internal/client"
	"github.com/mwdev22/Custom-Protocol-Server/internal/config"
)

var (
	health = make(chan struct{})
)

func main() {
	cfg := config.New()
	tcpConn, err := net.Dial("tcp", cfg.IP+":"+cfg.Port)
	if err != nil {
		log.Fatalf("error connecting to TCP server: %v", err)
	}
	defer tcpConn.Close()

	buf := make([]byte, 32) // SHA-256 produces a 32-byte hash
	n, err := tcpConn.Read(buf)
	if err != nil {
		log.Fatalf("error reading from TCP server: %v", err)
	}
	hash := buf[:n]
	log.Printf("Received hash from server: %x", hash)

	go monitorConnection(tcpConn)

	udpAddr := cfg.IP + ":9000"

	udpConn, err := net.Dial("udp", udpAddr)
	if err != nil {
		log.Fatalf("error connecting to UDP server: %v", err)
	}
	defer udpConn.Close()

	cam, err := cam.New()
	if err != nil {
		log.Fatalf("error setting up camera: %v", err)
	}
	go func() {
		defer cam.Close()
		for {
			frame, err := cam.ReadFrame()
			if err != nil {
				log.Println("error reading frame: ", err)
				continue // skip the current iteration and try reading the next frame
			}

			// calculate number of packets
			numPackets := len(frame) / config.MaxPacketSize
			if len(frame)%config.MaxPacketSize != 0 {
				numPackets++
			}

			for i := 0; i < numPackets; i++ {
				// calculate start and end of packet
				start := i * config.MaxPacketSize
				end := (i + 1) * config.MaxPacketSize
				// last packet is usually shorter
				if end > len(frame) {
					end = len(frame)
				}

				packet := frame[start:end]

				// append hash to the video data
				data := append(hash, packet...)

				// send the packet via UDP
				_, err := udpConn.Write(data)
				if err != nil {
					log.Printf("error sending UDP data: %v", err)
					continue // skip sending this packet and try again
				}
			}

			// sleep to control the frame rate (can be adjusted based on requirements)
			time.Sleep(100 * time.Millisecond)
		}
	}()

	<-health

	tcpConn.Close()
	udpConn.Close()

}

func monitorConnection(conn net.Conn) {
	for {
		buf := make([]byte, 4)
		_, err := conn.Read(buf)
		if err != nil {
			log.Fatalf("error reading from TCP server: %v", err)
		}

		code := config.ErrorCode(binary.BigEndian.Uint32(buf))
		switch code {
		case config.ErrorCodeOK:
			log.Println("--- Connection is OK, received status from TCP server ---")
		case config.ErrorCodeInvalidRequest:
			log.Println(config.ErrorCodeInvalidRequest)
			health <- struct{}{}
		case config.ErrorCodeInvalidHash:
			log.Println(config.ErrorCodeInvalidHash)
			health <- struct{}{}
		}
	}
}
