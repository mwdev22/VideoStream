package server

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/binary"
	"log"
	"net"
	"sync"
	"time"

	"github.com/mwdev22/Custom-Protocol-Server/internal/config"
)

type serverTCP struct {
	IP       string
	Port     string
	Listener net.Listener
	clients  map[string]*Client
	mutex    *sync.RWMutex
}

func NewTCP(ip, port string) *serverTCP {
	return &serverTCP{
		IP:      ip,
		Port:    port,
		clients: make(map[string]*Client),
		mutex:   &sync.RWMutex{},
	}
}

type Client struct {
	Conn net.Conn
	Hash []byte
	Quit chan struct{}
}

func NewClientTCP(conn net.Conn, hash []byte) *Client {
	return &Client{
		Conn: conn,
		Hash: hash,
		Quit: make(chan struct{}),
	}
}

func (s *serverTCP) Start() error {
	address := net.JoinHostPort(s.IP, s.Port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	s.Listener = listener
	log.Printf("TCP server listening on %s\n", address)

	go func() {
		for {
			conn, err := s.Listener.Accept()
			if err != nil {
				log.Println(err)
				continue
			}
			go s.handleConnection(conn)
		}
	}()
	return nil
}

func (s *serverTCP) handleConnection(conn net.Conn) {
	defer conn.Close()
	clientAddr := conn.RemoteAddr().(*net.TCPAddr).IP.String()
	client := NewClientTCP(
		conn,
		generateHash(),
	)

	s.addClient(clientAddr, client)
	defer s.removeClient(clientAddr)

	// write hash to client
	conn.Write(client.Hash)
	log.Printf("client connected: %s\n", clientAddr)

	buf := make([]byte, 4)
	beatChan := make(chan struct{})
	go func() {
		ticker := time.NewTicker(1 * time.Second) // send a heartbeat every second
		defer ticker.Stop()

		for range ticker.C {
			binary.BigEndian.PutUint32(buf, uint32(config.ErrorCodeOK))
			_, err := conn.Write(buf)
			if err != nil {
				log.Printf("error sending heartbeat to client %s: %v", clientAddr, err)
				beatChan <- struct{}{}
				return
			}
		}
	}()

	select {
	case <-client.Quit:
		binary.BigEndian.PutUint32(buf, uint32(config.ErrorCodeInvalidHash))
		_, err := conn.Write(buf)
		if err != nil {
			log.Printf("error sending status to client, maybe he had already disconnected? %s: %v", clientAddr, err)
		}
		log.Printf("client disconnected: %s\n", clientAddr)
	case <-beatChan:
		log.Printf("client disconnected: %s\n", clientAddr)
	}
}

func (s *serverTCP) addClient(addr string, client *Client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.clients[addr] = client
}
func (s *serverTCP) retrieveClient(addr string) *Client {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.clients[addr]
}

func (s *serverTCP) removeClient(addr string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	delete(s.clients, addr)
}

func (s *serverTCP) Stop() error {
	if s.Listener != nil {
		return s.Listener.Close()
	}
	return nil
}

func generateHash() []byte {
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		log.Fatalf("failed to generate random data for hash: %v", err)
	}

	hash := sha256.New()
	hash.Write(randomBytes)
	return hash.Sum(nil)
}
