package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"net"
	"net/rpc"
)

type Cli struct {
	ID  uint64
	Con net.Conn
}

type MsgMeta struct {
	CliID   uint64 // who sends the message
	MsgBody string
}

type ServerInfo struct {
	Temtic     string
	TotalUsers uint64
}

/* Server for the distint tematics*/
type Server struct {
	Title       string
	Listener    net.Listener
	Port        string
	Host        string
	GeneralChat []MsgMeta
	IDCounter   uint64
	Clients     map[uint64]Cli
}

func (s *Server) Init() error {
	listener, err := net.Listen("tcp", ":"+s.Port)

	if err != nil {
		return err
	}
	s.IDCounter = 0
	s.Listener = listener
	s.Clients = make(map[uint64]Cli)
	return nil
}

func (s *Server) Run() {
	fmt.Println("Running ", s.Title)
	for {
		client, err := s.Listener.Accept()
		fmt.Println("Client: ", client, "\tListener", s.Listener)
		if err != nil {
			continue
		}
		go s.HandleClient(client)
	}
}

func (s *Server) HandleClient(c net.Conn) {
	defer c.Close()

	fmt.Println(s.IDCounter)
	/* create a new client in the server */
	client := Cli{ID: s.IDCounter, Con: c}
	s.IDCounter++

	/* Add the new client to the slice */
	s.Clients[client.ID] = client

	newClient := s.Clients[client.ID]
	fmt.Println(s.IDCounter)
	/* sends its assigned id */
	gob.NewEncoder(newClient.Con).Encode(newClient.ID)

	var instruction string
	var msg string
	fmt.Println(client)

	/* Listen for client requests */
	for {
		if newClient.Con != nil {
			/* receive the instruction (post_msg, etc) */
			receive := gob.NewDecoder(newClient.Con)
			err := receive.Decode(&instruction)
			if err == nil {
				if instruction == "post_msg" {
					/*  reads the data   */
					receibeMsg := gob.NewDecoder(newClient.Con)
					errMsg := receibeMsg.Decode(&msg)

					/* error*/
					if errMsg == nil {
						//sends the message through the general chat */
						go s.SendChat(msg, newClient.ID)
					} else {
						fmt.Println("Error ", err)
						return
					}
				} else if instruction == "post_file" {
					//f := File{}
					//receibeFile := gob.NewDecoder(newClient.Con)
					//errFile := receibeFile.Decode(&f)

					/*if errFile == nil {
					/* send the file to clients*/
					//go s.SendFile(&f, newClient.ID)
					/*} else {
						fmt.Println("Error ", err)
						return
					}*/
				} else if instruction == "out" {
					/* client disconnect  */
					//go s.DisconnectClient(newClient.ID)
					break
				}
			} else { /* error */
				fmt.Println("Error pepe ", err)
				return
			}
		} else {
			break
		}
	}
}
func (s *Server) SendChat(msg string, id uint64) {
	/* sends msg though the general chat*/
	var instruction string = "msg"
	meta := MsgMeta{MsgBody: msg, CliID: id}
	s.GeneralChat = append(s.GeneralChat, meta)
	for _, c := range s.Clients {
		gob.NewEncoder(c.Con).Encode(&instruction)
		gob.NewEncoder(c.Con).Encode(&meta)
	}
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

func (m *MicroService) GetChatRoomIP(ID string, response *string) error {
	if len(m.ChatServers) == 0 {
		return errors.New("No servers available")
	}

	for _, server := range m.ChatServers {
		if server.Title == ID {
			*response = server.Port
			return nil
		}
	}
	return errors.New("The server was not found")
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
