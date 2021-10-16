package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
)

type MicroService struct {
}

func (m *MicroService) Hello(name string, response *string) error {
	*response = "Hello " + name + "!!"
	return nil
}

type Server struct {
	Host    string
	Service MicroService
	Listen  net.Listener
}

func (s *Server) Init() error {

	/* Initialize the server */
	if len(s.Host) > 0 {
		server, err := net.Listen("tcp", s.Host)
		if err != nil {
			return err
		}

		s.Listen = server
		fmt.Println("Server running on", s.Host)
		return nil
	} else {
		return errors.New("error, no host provided")
	}

}

func (s *Server) Run() {
	rpc.Register(new(MicroService))
	for {
		c, err := s.Listen.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go rpc.ServeConn(c)
	}
}

func main() {
	server := Server{Host: ":9999"}
	err := server.Init()
	if err != nil {
		fmt.Println(err)
	}

	go server.Run()

	var input string
	fmt.Scanln(&input)
}
