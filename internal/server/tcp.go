package server

import (
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

func NewServerTCP(ip, port string) *serverTCP {
	return &serverTCP{
		IP:      ip,
		Port:    port,
		clients: make(map[string]*Client),
		mutex:   &sync.RWMutex{},
	}
}

type Client struct {
	IP   string
	Port int
	Conn net.Conn
}

func NewClientTCP(ip string, port int, conn net.Conn) *Client {
	return &Client{
		IP:   ip,
		Port: port,
		Conn: conn,
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

	for {
		conn, err := s.Listener.Accept()
		if err != nil {
			return err
		}
		go s.handleConnection(conn)
	}
}

func (s *serverTCP) handleConnection(conn net.Conn) {
	defer conn.Close()

	clientAddr := conn.RemoteAddr().String()
	client := NewClientTCP(
		conn.RemoteAddr().(*net.TCPAddr).IP.String(), conn.RemoteAddr().(*net.TCPAddr).Port, conn,
	)

	s.addClient(clientAddr, client)
	log.Printf("client connected: %s\n", clientAddr)

	defer s.removeClient(clientAddr)
}

func (s *serverTCP) addClient(addr string, client *Client) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.clients[addr] = client
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
