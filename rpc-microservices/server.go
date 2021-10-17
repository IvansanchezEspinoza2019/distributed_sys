package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
)

/* struct to save client request for saving a student note*/
type StudentSubject struct {
	Name    string
	Subject string
	Note    float64
}

/* server response */
type Reply struct {
	Msg    string
	Status int
}

/* struct microservices */
type MicroService struct {
	Students map[string]map[string]float64
	Subjects map[string]map[string]float64
	// ejemplo: https://stackoverflow.com/questions/68281518/golang-rpc-get-wrong-struct-variable
}

/* CALLS*/
func (m *MicroService) AddNoteStudentSubjetc(data *StudentSubject, response *Reply) error {

	// add student info
	if _, exists := m.Students[data.Name]; exists { //  already exists
		if exists {
			m.Students[data.Name][data.Subject] = data.Note
		}
	} else {
		subject := make(map[string]float64)
		subject[data.Subject] = data.Note
		m.Students[data.Name] = subject
	}

	// add subject info
	if _, exists := m.Subjects[data.Subject]; exists {
		if exists {
			m.Subjects[data.Subject][data.Name] = data.Note
		}
	} else {
		// create the first student
		student := make(map[string]float64)
		student[data.Name] = data.Note
		// create the subject
		m.Subjects[data.Subject] = student
	}
	fmt.Println("STUDENTS: ", m.Students)
	fmt.Println("SUBJECTS: ", m.Subjects)
	response.Msg = "EXITO"
	response.Status = 200
	return nil
}

/* Server Struct*/
type Server struct {
	Host    string
	Listen  net.Listener
	Service *MicroService
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
	/* instance of microservices struct*/
	st := make(map[string]map[string]float64)
	sb := make(map[string]map[string]float64)
	s.Service = &MicroService{Students: st, Subjects: sb}

	rpc.Register(s.Service)
	/* listen for clients */
	for {
		c, err := s.Listen.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go rpc.ServeConn(c)
	}
}

func (s *Server) Hello(name string, response *string) error {
	*response = "Hello Server " + name + " Host " + s.Host
	return nil
}

func main() {
	// init server
	server := Server{Host: ":9999"}
	err := server.Init()
	if err != nil {
		fmt.Println(err)
	}

	go server.Run()

	var input string
	fmt.Scanln(&input)
}
