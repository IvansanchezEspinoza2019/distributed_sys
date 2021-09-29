package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"

	"./process"
)

type Server struct {
	ProcessQueue []*process.Process
}

func (s *Server) Init(pLen int) {
	/* create initial processes */
	for i := 0; i < pLen; i++ {
		var p *process.Process = new(process.Process)
		p.ID = uint64(i)
		p.ICount = 0
		p.Executing = true
		s.ProcessQueue = append(s.ProcessQueue, p)
	}
}

func (s *Server) Start() {
	/* starts server*/
	go s.Listen()

	go s.UpdateProcesses()
	time.Sleep(time.Millisecond * 100)

	for _, p := range s.ProcessQueue {
		go p.Execute() // execute every single process
	}
}

func (s *Server) HandleClient(c net.Conn) {
	/* handle client requests */
	var msg string
	err := gob.NewDecoder(c).Decode(&msg)

	if err != nil {
		fmt.Println(err)
		return
	}

	if msg == "get_process" {
		/* new process requests */
		if len(s.ProcessQueue) > 0 {
			gob.NewEncoder(c).Encode(s.ProcessQueue[0])
			s.ProcessQueue[0].Stop() // stop process
			s.ProcessQueue = s.ProcessQueue[1:]
		}
	} else if msg == "put_process" {
		/* the client retrieve its process to the server */
		var newProcess process.Process
		gob.NewDecoder(c).Decode(&newProcess)
		s.ProcessQueue = append(s.ProcessQueue, &newProcess)
		newProcess.Executing = true
		go newProcess.Execute()
	}
}

func (s *Server) UpdateProcesses() {
	/* prints all running processes */
	for {
		for _, p := range s.ProcessQueue {
			fmt.Printf("%d : %d \n", p.ID, p.ICount)
		}
		fmt.Println("-------------------------------")
		time.Sleep(time.Millisecond * 500)
	}
}

func (s *Server) Listen() {
	/***** SERVER RUNNING *****/
	server, err := net.Listen("tcp", ":3302")
	if err != nil {
		fmt.Println(err)
	}
	for {
		client, err := server.Accept()

		if err != nil {
			fmt.Println(err)
			continue
		}
		go s.HandleClient(client)
	}
}

func main() {
	var s Server
	s.Init(5) // init sercer processes
	s.Start() // all needed to start the server

	var execute string
	fmt.Scanln(&execute)
}
