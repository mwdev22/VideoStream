package main

import (
	"flag"
	"log"
	"net"
	"os"
	"time"

	"github.com/blackjack/webcam"
	"github.com/mwdev22/Custom-Protocol-Server/internal/config"
)

var (
	dev = flag.String("d", "/dev/video0", "video device to use")
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

	cam, err := webcam.Open(*dev)
	if err != nil {
		log.Fatalf("error opening camera: %v", err)
		os.Exit(1)
	}

	defer cam.Close()

	cam.SetFramerate(30)
	format_desc := cam.GetSupportedFormats()
	var format webcam.PixelFormat

	for k, s := range format_desc {
		if s == "Motion-JPEG" {
			format = k
		}
	}
	_, _, _, err = cam.SetImageFormat(format, 640, 480)
	if err != nil {
		log.Fatalf("error setting image format: %v", err)
	}

	err = cam.StartStreaming()
	if err != nil {
		log.Fatalf("error starting streaming: %v", err)
	}
	for {
		frame, err := cam.ReadFrame()
		if err != nil {
			log.Println("error reading frame: ", err)
		}
		numPackets := len(frame) / config.MaxPacketSize
		if len(frame)%config.MaxPacketSize != 0 {
			numPackets++
		}

		for i := 0; i < numPackets; i++ {
			start := i * config.MaxPacketSize
			end := (i + 1) * config.MaxPacketSize
			if end > len(frame) {
				end = len(frame)
			}

			packet := frame[start:end]

			data := append(hash, packet...)

			_, err := udpConn.Write([]byte(data))
			if err != nil {
				log.Printf("Error sending UDP data: %v", err)
			}
			log.Printf("sent %d bytes to %s\n", len(data), udpAddr)
		}

		time.Sleep(100 * time.Millisecond)
	}
}
