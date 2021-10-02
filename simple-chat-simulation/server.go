package main

import (
	"errors"
	"fmt"
	"net"
)

type Server struct {
	Host string
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
	fmt.Println(c)
}

func main() {

	s := Server{Host: ":5555"}
	//s.Init()
	go s.Run()

	var in string
	fmt.Scanln(&in)

}
