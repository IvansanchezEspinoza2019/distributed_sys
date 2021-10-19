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
	fmt.Println("\n.: STUDENTS :.")
	for student, subjets := range m.Students {
		fmt.Println(student, "{")
		for subject, sv := range subjets {
			fmt.Println("\t"+subject+": ", sv)
		}
		fmt.Println("}")
	}

	fmt.Println("\n.: SUBJECTS :.")
	for subjet, students := range m.Subjects {
		fmt.Println(subjet, "{")
		for student, sv := range students {
			fmt.Println("\t"+student+": ", sv)
		}
		fmt.Println("}")
	}
	fmt.Println("----------------------------------")
	response.Msg = "EXITO"
	response.Status = 200
	return nil
}

func (m *MicroService) StudentAverege(studentName string, response *float64) error {
	/* clculates the student averega of all its subjects*/

	if v, exists := m.Students[studentName]; exists {
		if exists {
			var avg float64
			for _, sv := range v {
				avg += sv
			}
			*response = (avg) / float64(len(v))
		} else {
			return errors.New("No existe el estudiante")
		}
	} else {
		return errors.New("No existe el estudiante")
	}
	return nil
}

func (m *MicroService) GeneralAverage(typeAVG string, response *float64) error {
	/* calculates the note aberage of all students */

	if len(m.Students) == 0 {
		return errors.New("Vacío")
	}

	var totalAvg float64
	for student, subjects := range m.Students {
		var avg float64
		for _, sv := range subjects {
			avg += sv
		}
		avg = (avg / float64(len(subjects)))
		totalAvg += avg
		fmt.Println(student+": ", avg)
	}
	fmt.Println("\n-----------------")
	*response = (totalAvg / float64(len(m.Students)))
	return nil
}

func (m *MicroService) SubjectAvg(subject string, response *float64) error {
	/* calculates the avdg of a subject */

	if len(m.Subjects) == 0 {
		return errors.New("Vacio")
	}
	if students, exists := m.Subjects[subject]; exists {
		if exists {
			var svg float64
			for _, sv := range students {
				svg += sv
			}
			*response = (svg / float64(len(students)))
		} else {
			return errors.New("NO existe la materia")
		}
	} else {
		return errors.New("NO existe la materia")
	}
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
