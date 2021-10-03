package main

import (
	"encoding/gob"
	"errors"
	"fmt"
	"net"
)

type Server struct {
	Host    string
	Clients []net.Conn
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

func (s *Server) HandleClient(c net.Conn) {
	defer c.Close()

	/* Add the new client to the slice */
	s.Clients = append(s.Clients, c)

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
					fmt.Println("Mensaje recibido: ", msg)
				} else {
					fmt.Println("Error ", err)
					return
				}
			}

		} else { /* error */
			fmt.Println("Error ", err)
			return
		}
	}

}

func main() {
	/* serever init */
	s := Server{Host: ":5555"}

	go s.Run()

	var in string
	fmt.Scanln(&in)
}
