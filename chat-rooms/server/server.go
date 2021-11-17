package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"net"
	"net/rpc"

	"../common"
)

/* Server for the distint tematics*/
type Server struct {
	Title       string
	Listener    net.Listener
	Port        string
	Host        string
	Files       []common.File
	GeneralChat []common.MsgMeta
	IDCounter   uint64
	Clients     map[uint64]net.Conn
}

func (s *Server) Init() error {
	listener, err := net.Listen("tcp", s.Host+":"+s.Port)

	if err != nil {
		return err
	}
	s.IDCounter = 0
	s.Listener = listener
	s.Clients = make(map[uint64]net.Conn)
	return nil
}

func (s *Server) Run() {
	/* chat room listenning for clients connections*/
	for {
		client, err := s.Listener.Accept()
		if err != nil {
			continue
		}
		go s.HandleClient(client)
	}
}

func (s *Server) HandleClient(c net.Conn) {
	defer c.Close()
	/* create a new client in the server */
	clientID := s.IDCounter
	s.IDCounter++

	/* Add the new client to the slice */
	s.Clients[clientID] = c

	/* sends its assigned id */
	gob.NewEncoder(c).Encode(clientID)

	var instruction string
	var msg string

	/* Listen for client requests */
	for {
		if c != nil {
			/* receive the instruction (post_msg, etc) */
			receive := gob.NewDecoder(c)
			err := receive.Decode(&instruction)
			if err == nil {
				if instruction == "post_msg" {
					/*  reads the data   */
					receibeMsg := gob.NewDecoder(c)
					errMsg := receibeMsg.Decode(&msg)

					/* error*/
					if errMsg == nil {
						//sends the message through the general chat */
						go s.SendChat(msg, clientID)
					} else {
						fmt.Println("Error ", err)
						return
					}
				} else if instruction == "post_file" {
					f := common.File{}
					receibeFile := gob.NewDecoder(c)
					errFile := receibeFile.Decode(&f)

					if errFile == nil {
						/* send the file to clients*/
						go s.SendFile(&f, clientID)
					} else {
						fmt.Println("Error ", err)
						return
					}
				} else if instruction == "out" {
					/* client disconnect  */
					go s.DisconnectClient(clientID)
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
	meta := common.MsgMeta{MsgBody: msg, CliID: id}
	s.GeneralChat = append(s.GeneralChat, meta)
	for _, c := range s.Clients {
		gob.NewEncoder(c).Encode(&instruction)
		gob.NewEncoder(c).Encode(&meta)
	}
}

func (s *Server) SendFile(file *common.File, id uint64) {
	/* Sends the file to all connected clients. only gnores the file creator*/
	var instruction = "file"
	file.Creator = id
	s.Files = append(s.Files, *file)
	s.GeneralChat = append(s.GeneralChat, common.MsgMeta{CliID: id, MsgBody: file.Filename})

	for ID, c := range s.Clients {
		if ID != id {
			gob.NewEncoder(c).Encode(&instruction)
			gob.NewEncoder(c).Encode(file)
		} else {
			var special_instruccion string = "msg"
			meta := common.MsgMeta{MsgBody: file.Filename, CliID: id}
			gob.NewEncoder(c).Encode(&special_instruccion)
			gob.NewEncoder(c).Encode(&meta)
		}
	}
}

func (s *Server) DisconnectClient(id uint64) {
	/* close the connection and removes the client from the server */
	cli, exists := s.Clients[id]
	if exists {
		cli.Close()
		delete(s.Clients, id)
	}
}

/********************** Microservice server ***********************/
type MicroService struct {
	ChatServers []*Server
}

func (m *MicroService) GetChatRooms(message string, response *[]common.ServerInfo) error {
	if len(m.ChatServers) == 0 {
		return errors.New("No servers available")
	}
	for _, server := range m.ChatServers {
		si := common.ServerInfo{Temtic: server.Title, TotalUsers: uint64(len(server.Clients))}
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
			*response = server.Host + ":" + server.Port
			return nil
		}
	}
	return errors.New("The server was not found")
}

func (m *MicroService) ServersStatus(message string, response *[]common.ServerDetail) error {
	if len(m.ChatServers) == 0 {
		return errors.New("No servers available")
	}
	for _, server := range m.ChatServers {
		*response = append(*response, common.ServerDetail{IP: "http://" + server.Host + ":" + server.Port, Tematic: server.Title, TotalUsers: uint64(len(server.Clients))})
	}
	return nil
}

// main //
func main() {
	/* chat rooms*/
	s1 := Server{Title: "Videogames", Host: "localhost", Port: "9997"}
	s2 := Server{Title: "Cooking", Host: "localhost", Port: "9998"}
	s3 := Server{Title: "General", Host: "localhost", Port: "9999"}

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

	fmt.Println("Chat Rooms Servers Running!")
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
