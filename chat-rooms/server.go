package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
)

type ServerInfo struct {
	Temtic     string
	TotalUsers uint64
}

/* Server for the distint tematics*/
type Server struct {
	Title     string
	Listener  net.Listener
	Port      string
	Host      string
	IDCounter uint64
	Clients   map[uint64]net.Conn
}

func (s *Server) Init() error {
	listener, err := net.Listen("tcp", ":"+s.Port)

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

/********************** Microservice server ***********************/
type MicroService struct {
	ChatServers []*Server
}

func (m *MicroService) GetChatRooms(message string, response *[]ServerInfo) error {
	if len(m.ChatServers) == 0 {
		return errors.New("No servers available")
	}
	for _, server := range m.ChatServers {
		si := ServerInfo{Temtic: server.Title, TotalUsers: uint64(len(server.Clients))}
		*response = append(*response, si)
	}
	return nil
}

func main() {
	/* chat rooms*/
	s1 := Server{Title: "Videogames", Host: "", Port: "9997"}
	s2 := Server{Title: "Cooking", Host: "", Port: "9998"}
	s3 := Server{Title: "General", Host: "", Port: "9999"}

	service := &MicroService{ChatServers: []*Server{
		&s1, &s2, &s3,
	}}

	err := InitializeChatRooms(service)

	if err != nil {
		fmt.Println(err)
		return
	}
	// run every chat room
	for _, server := range service.ChatServers {
		go server.Run()
	}
	/** Microservice Server **/

	serviceServer, errServer := net.Listen("tcp", ":8000")
	if errServer != nil {
		fmt.Println(errServer)
		return
	}

	fmt.Println("Chat Room Running!")
	/* Register RPC service*/
	rpc.Register(service)
	/* listen for clients */
	for {
		c, err := serviceServer.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go rpc.ServeConn(c)
	}

	var input string
	fmt.Scanln(&input)
}

func InitializeChatRooms(m *MicroService) error {
	/* Initialize every caht room*/
	for _, server := range m.ChatServers {
		server.Init()
	}
	return nil
}
