package main

import (
	"bufio"
	"encoding/gob"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
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
	Listen      net.Listener
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
	for {
		client, err := s.Listen.Accept()
		if err != nil {
			continue
		}
		go s.HandleClient(client)
	}
}

func (s *Server) SendFile(file *File, id uint64) {
	/* Sends the file to all connected clients. only gnores the file creator*/
	var instruction = "file"
	file.Creator = id
	s.Files = append(s.Files, *file)
	s.GeneralChat = append(s.GeneralChat, MsgMeta{CliID: id, MsgBody: file.Filename})

	for _, c := range s.Clients {
		if c.ID != id {
			gob.NewEncoder(c.Con).Encode(&instruction)
			gob.NewEncoder(c.Con).Encode(file)
		} else {
			var special_instruccion string = "msg"
			meta := MsgMeta{MsgBody: file.Filename, CliID: id}
			gob.NewEncoder(c.Con).Encode(&special_instruccion)
			gob.NewEncoder(c.Con).Encode(&meta)
		}
	}
}

func (s *Server) HandleClient(c net.Conn) {
	defer c.Close()

	/* create a new client in the server */
	client := Cli{ID: s.ClientIds, Con: c}

	/* Add the new client to the slice */
	s.Clients = append(s.Clients, client)
	s.ClientIds++

	newClient := s.Clients[len(s.Clients)-1]

	/* sends its assigned id */
	gob.NewEncoder(newClient.Con).Encode(newClient.ID)

	var instruction string
	var msg string

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
						/*sends the message through the general chat */
						go s.SendChat(msg, newClient.ID)
					} else {
						fmt.Println("Error ", err)
						return
					}
				} else if instruction == "post_file" {
					f := File{}
					receibeFile := gob.NewDecoder(newClient.Con)
					errFile := receibeFile.Decode(&f)

					if errFile == nil {
						/* send the file to clients*/
						go s.SendFile(&f, newClient.ID)
					} else {
						fmt.Println("Error ", err)
						return
					}
				} else if instruction == "out" {
					/* client disconnect  */
					go s.DisconnectClient(newClient.ID)
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

func (s *Server) DisconnectClient(id uint64) {
	/* close the connection and removes the client from the server */
	var index int = -1
	for i, cli := range s.Clients {
		if cli.ID == id {
			cli.Con.Close()
			index = i
			break
		}
	}
	if index != -1 {
		s.Clients = append(s.Clients[:index], s.Clients[index+1:]...)
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

func (s *Server) ShowMessages() {
	/* prints the general chat */
	fmt.Println("\n---------------\n")
	for _, meta := range s.GeneralChat {
		fmt.Printf("[Client-{%d}] %s\n", meta.CliID, meta.MsgBody)
	}
	fmt.Println("\n---------------\n")
}

func (s *Server) ShowConnectedClients() {
	/* prints the connected clients */
	fmt.Println("\n---------------\n")
	for _, cli := range s.Clients {
		fmt.Printf("Client {%d}\n", cli.ID)
	}
	fmt.Println("\n---------------\n")
}

func (s *Server) Backup() {
	/*  makes a backup of the generl chat (incudes filenames)*/
	file, err := os.Create("ServerBackup.txt")

	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	for _, data := range s.GeneralChat {
		var resp string
		resp = strconv.Itoa(int(data.CliID)) + ", " + data.MsgBody
		_, errfile := file.WriteString(resp + "\n")
		if errfile != nil {
			fmt.Println(errfile)
		}
	}
}

func (s *Server) ShootDownServer() {
	/* disconnect all clients */
	for _, cli := range s.Clients {
		if cli.Con != nil {
			var msg string = "disconnect"
			gob.NewEncoder(cli.Con).Encode(msg)
		}
	}
	s.Clients = nil
	fmt.Println("Server shooted down...")
}

func main() {
	/* serever init */
	s := Server{Host: ":5555", ClientIds: 0}
	err := s.Init()
	if err != nil {
		fmt.Println("Error ", err)
		return
	}

	/* runs the server*/
	go s.Run()
	/* auxiliar variables*/
	scanner := bufio.NewScanner(os.Stdin)
	var opc string

	/* menu */
	for opc != "4" {
		menu()
		scanner.Scan()
		opc = scanner.Text()
		if opc == "1" { // print messages
			fmt.Println("CHAT GENERAL:")
			s.ShowMessages()
		} else if opc == "2" { // backup
			go s.Backup()
		} else if opc == "3" { // connected clients
			fmt.Println("CLIENTES:")
			s.ShowConnectedClients()
		} else if opc == "4" {
			s.ShootDownServer()
			break
		}
	}
}

func menu() {
	fmt.Println(".: SERVIDOR :.\n")
	fmt.Println("1) Mostrar mensajes/archivos")
	fmt.Println("2) Respaldar mensajes/archivos")
	fmt.Println("3) Clientes connectdos")
	fmt.Println("4) Terminar servidor")
	fmt.Print("Opcion: ")
}
