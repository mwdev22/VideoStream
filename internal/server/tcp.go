package server

import (
	"crypto/rand"
	"log"
	"net"
	"sync"
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
}

func NewClientTCP(conn net.Conn, hash []byte) *Client {
	return &Client{
		Conn: conn,
		Hash: hash,
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
	// defer conn.Close()

	clientAddr := conn.RemoteAddr().(*net.TCPAddr).IP.String()
	client := NewClientTCP(
		conn,
		generateHash(),
	)

	s.addClient(clientAddr, client)
	conn.Write(client.Hash)
	log.Printf("client connected: %s\n", clientAddr)

	// defer s.removeClient(clientAddr)
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
	log.Printf("client disconnected: %s\n", addr)
}

func (s *serverTCP) Stop() error {
	if s.Listener != nil {
		return s.Listener.Close()
	}
	return nil
}

func generateHash() []byte {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Failed to generate hash: %v", err)
	}
	return bytes
}
