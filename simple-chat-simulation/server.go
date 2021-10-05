package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"net"
)

type MsgMeta struct {
	CliID   uint64 // who sends the message
	MsgBody string
}

type File struct {
	Filename string
	Content  []byte
	Creator  uint64
}

type Cli struct {
	ID  uint64
	Con net.Conn
}

type Server struct {
	Host        string
	Clients     []Cli
	ClientIds   uint64
	GeneralChat []MsgMeta
	Files       []File
}

func (s *Server) Init() {

}

func (s *Server) Run() error {
	if len(s.Host) > 0 {
		server, err := net.Listen("tcp", s.Host)
		if err != nil {
			return err
		}
		fmt.Println("Server running on", s.Host)
		for {
			client, err := server.Accept()
			if err != nil {
				continue
			}
			go s.HandleClient(client)
		}
		return nil
	} else {
		return errors.New("error, no host provided")
	}
}

func (s *Server) SendFile(file *File, id uint64) {
	var instruction = "file"
	file.Creator = id
	s.Files = append(s.Files, *file)
	s.GeneralChat = append(s.GeneralChat, MsgMeta{CliID: id, MsgBody: file.Filename})

	for _, c := range s.Clients {
		gob.NewEncoder(c.Con).Encode(&instruction)
		gob.NewEncoder(c.Con).Encode(file)
	}

}

func (s *Server) HandleClient(c net.Conn) {
	defer c.Close()

	client := Cli{ID: s.ClientIds, Con: c}
	/* Add the new client to the slice */
	s.Clients = append(s.Clients, client)
	s.ClientIds++

	/* sends its assigned id */
	gob.NewEncoder(c).Encode(s.Clients[len(s.Clients)-1].ID)

	fmt.Println(s.Clients)

	var instruction string
	var msg string

	/* Listen for client requests */
	for {
		/*1.- receive the instruction (post_msg, etc) */
		receive := gob.NewDecoder(c)
		err := receive.Decode(&instruction)
		if err == nil {
			fmt.Println("Instruccion recibida: ", instruction)

			if instruction == "post_msg" {
				/* 2.- reads the data   */
				receibeMsg := gob.NewDecoder(c)
				errMsg := receibeMsg.Decode(&msg)

				/* error*/
				if errMsg == nil {
					go s.SendChat(msg, client.ID)
					fmt.Println("Mensaje recibido: ", msg)
				} else {
					fmt.Println("Error 1", err)
					return
				}
			} else if instruction == "post_file" {
				f := File{}
				receibeFile := gob.NewDecoder(c)
				errFile := receibeFile.Decode(&f)

				if errFile == nil {
					go s.SendFile(&f, client.ID)

					fmt.Println("Archivo recibido: ", f)
				} else {
					fmt.Println("Error ", err)
					return
				}
			}

		} else { /* error */
			fmt.Println(instruction)
			fmt.Println("Error 2", err)
			return
		}
	}

}

func (s *Server) SendChat(msg string, id uint64) {
	var instruction string = "msg"
	meta := MsgMeta{MsgBody: msg, CliID: id}
	s.GeneralChat = append(s.GeneralChat, meta)
	for _, c := range s.Clients {
		gob.NewEncoder(c.Con).Encode(&instruction)
		gob.NewEncoder(c.Con).Encode(&meta)
	}
}

func main() {
	/* serever init */
	s := Server{Host: ":5555", ClientIds: 0}

	go s.Run()

	var in string
	fmt.Scanln(&in)
}
