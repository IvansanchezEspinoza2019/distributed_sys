package main

import (
	"fmt"
	"net"
)

/* Server for the distint tematics*/
type Server struct {
	Title     string
	Listener  net.Listener
	Port      string
	Host      string
	IDCounter uint64
	Clients   map[uint64]net.Conn
}

func (s *Server) Init(host string, port string) error {
	const (
		CONN_HOST = "localhost"
		CONN_PORT = "9999"
		CONN_TYPE = "tcp"
	)
	listener, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)

	if err != nil {
		return err
	}
	s.IDCounter = 0
	s.Listener = listener
	s.Clients = make(map[uint64]net.Conn)
	return nil
}

func (s *Server) Run() {
	for {
		client, err := s.Listener.Accept()
		if err != nil {
			continue
		}
		s.HandleClient(client)
	}
}

func (s *Server) HandleClient(c net.Conn) {
	// adding the new connection
	s.Clients[s.IDCounter] = c
	s.IDCounter++
}

func main() {
	s := Server{Port: "9997", Host: "localhost"}
	err := s.Init("localhost", "997")
	if err != nil {
		fmt.Println(err)
		return
	}
	go s.Run()

	var input string
	fmt.Scanln(&input)
}
