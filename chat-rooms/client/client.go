/*
	The client connects to the Middleware by RESTful API.
	When client chooses the chat room, then the connection is made
	by TCP protocol for sending and receiving msg's from other clients
*/

package main

import (
	"bufio"
	"encoding/gob"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"../common"
)

/* strunct client */
type Client struct {
	Connection   net.Conn
	ID           uint64
	ChatName     string
	MsgChan      chan string
	FileChan     chan string
	HasDirectory bool
	GeneralChat  []common.MsgMeta
	DirName      string
	Connected    bool
	ApiAdd       string
	_admi_pass   string
	_admin_email string
	AdminLogged  bool
}

/*********** Connection to MIDDLEWARE *********/
func (c *Client) MakePetitionApi() { // to test the api
	// https://tutorialedge.net/golang/consuming-restful-api-with-go/
	fmt.Println("Making petition to api..")
	response, err := http.Get(c.ApiAdd + "/chatRooms")
	if err != nil {
		fmt.Print(err.Error())
		return
	}

	// get the raw json from the API
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(responseData))

	// convert it to a Go datatype
	var servers []common.ServerInfo
	json.Unmarshal([]byte(responseData), &servers)

	fmt.Println("New", servers)

}

func (c *Client) FirstMenu() {
	fmt.Println("  1) Conectarse a una sala de chat")
	fmt.Println("  2) Salir")
	fmt.Print("Opcion: ")
}

func (c *Client) ApiCallGetRooms() ([]common.ServerInfo, error) {
	// make petition to middleware
	response, err := http.Get(c.ApiAdd + "/chatRooms")
	if err != nil {
		return nil, err
	}
	// get the raw json from the API
	responseData, errRes := ioutil.ReadAll(response.Body)
	if errRes != nil {
		return nil, errRes
	}
	// convert it to a Go datatype
	var servers []common.ServerInfo
	json.Unmarshal([]byte(responseData), &servers)

	return servers, nil
}

func (c *Client) ApiCallGetIPServer(serverName string) (string, error) {
	response, errCall := http.Get(c.ApiAdd + "/chatRooms/" + serverName)
	if errCall != nil {
		return "", errCall
	}

	responseData, errRes := ioutil.ReadAll(response.Body)
	if errRes != nil {
		return "", errCall
	}

	var IP string
	json.Unmarshal([]byte(responseData), &IP)
	return IP, nil
}

func (c *Client) ApiCallGetSeversStatus() ([]common.ServerDetail, error) {
	/* make call to the api*/
	response, err := http.Get(c.ApiAdd + "/chatRooms/status")
	if err != nil {
		return nil, err
	}

	// get the response body
	responseData, errRes := ioutil.ReadAll(response.Body)
	if errRes != nil {
		return nil, errRes
	}

	// convert from json to a Go datatype
	var serverStatus []common.ServerDetail
	json.Unmarshal([]byte(responseData), &serverStatus)
	return serverStatus, nil
}

func (c *Client) MakeConnectionToServer(serverName string) int {
	// get server IP
	IP, errIP := c.ApiCallGetIPServer(serverName)
	if errIP != nil {
		fmt.Println(errIP)
		return 1
	}

	// make connection between client and server
	errConection := c.Init(IP)
	if errConection != nil {
		fmt.Println(errConection)
		return 1
	}
	c.ChatName = serverName
	return 0
}

func (c *Client) ServersMenu() int {
	servers, errCall := c.ApiCallGetRooms()
	if errCall != nil {
		fmt.Println(errCall)
		return 1
	}
	fmt.Println("\t.: SELLECCIONA UN CHAT:.\n")
	for index, server := range servers {
		fmt.Println(index, ")Chat: ["+server.Temtic+"]", "\tUsers: (", server.TotalUsers, ")")
	}
	fmt.Print("\nSelecciona un server: ")
	var selectedServer string
	fmt.Scanln(&selectedServer)

	index, errConv := strconv.Atoi(selectedServer)
	if errConv != nil {
		fmt.Println("[Debes seleccionar una opcion v??lida!]")
		return 1
	}

	if index < 0 || index >= len(servers) {
		fmt.Println("[Debes seleccionar una opcion v??lida!]")
		return 1
	}
	// falta hacer l peticion del server seleccionado y obtener la direccion ip de la sala de chat
	return c.MakeConnectionToServer(servers[index].Temtic)
}

/************* Connection With a specific CHAT ROOM **************/
func (c *Client) Init(IP string) error {
	/* connects to the server */
	con, err := net.Dial("tcp", IP) //":"+port
	if err != nil {
		return err
	}
	c.Connection = con
	c.Connected = true
	gob.NewDecoder(con).Decode(&c.ID)
	return nil
}

func (c *Client) SendMsg() {
	/* this function is in charge of posting messages to the server*/
	var instruction string = "post_msg"
	for {
		if c.Connection != nil {
			select {
			case msg := <-c.MsgChan:
				// first sends the instruction and then the content of the message
				gob.NewEncoder(c.Connection).Encode(&instruction)
				gob.NewEncoder(c.Connection).Encode(&msg)
			}
		}
	}
}
func (c *Client) ProccesReadFile(filename string) {
	var instruction string = "post_file"
	/* reads file*/
	file, err := os.Open(filename)

	if err == nil {
		defer file.Close()
		stat, err := file.Stat()

		if err != nil {
			fmt.Println(err)
			return
		}
		total := stat.Size()

		f := common.File{Filename: stat.Name(), Content: make([]byte, total)}
		count, err := file.Read(f.Content)
		if err != nil {
			fmt.Println(err)
			return
		}
		if c.Connection != nil && count > 0 {

			/* send to the server */
			gob.NewEncoder(c.Connection).Encode(&instruction)
			gob.NewEncoder(c.Connection).Encode(&f)

		}
	}
}

func (c *Client) SendFile() {
	for {
		if c.Connection != nil {
			select {
			case filename := <-c.FileChan:
				go c.ProccesReadFile(filename)
			}
		}
	}
}

func (c *Client) ListenForUpdates() {
	var instruction string
	for {
		if c.Connection != nil {

			receibe := gob.NewDecoder(c.Connection)
			err := receibe.Decode(&instruction)
			if err == nil {
				if instruction == "msg" {
					meta := common.MsgMeta{}
					gob.NewDecoder(c.Connection).Decode(&meta)
					c.GeneralChat = append(c.GeneralChat, meta)
				} else if instruction == "file" {
					file := common.File{}
					gob.NewDecoder(c.Connection).Decode(&file)
					go c.HandleFile(&file)
				} else if instruction == "disconnect" {
					fmt.Println("\nSERVER DISCONNECTED!!\n")
					c.Connection.Close()
					c.Connected = false
					break
				}
			}
		}
	}
}

func (c *Client) HandleFile(file *common.File) {
	c.GeneralChat = append(c.GeneralChat, common.MsgMeta{CliID: file.Creator, MsgBody: file.Filename})

	if !c.HasDirectory {
		_, err := os.Stat("test")

		if os.IsNotExist(err) {
			var dir string = "CLIENT_" + strconv.Itoa(int(c.ID))
			errDir := os.MkdirAll(dir, 0755)
			if errDir != nil {
				fmt.Println(errDir)
				return
			}
			c.HasDirectory = true
			c.DirName = dir
			/* save the file*/
			fileToSave, err := os.Create(c.DirName + "/" + file.Filename)
			if err != nil {
				fmt.Println(err)
				return
			}
			defer fileToSave.Close()
			_, errFile := fileToSave.Write(file.Content)
			if errFile != nil {
				fmt.Println(errFile)
			}
		}
	} else { // this client already has a folder
		/* save the file*/
		fileToSave, err := os.Create(c.DirName + "/" + file.Filename)
		if err != nil {
			fmt.Println(err)
			return
		}
		defer fileToSave.Close()
		_, errFile := fileToSave.Write(file.Content)
		if errFile != nil {
			fmt.Println(errFile)
		}
	}

}

func (c *Client) PrintGlobalChat() {
	fmt.Println("+------------------------+\n")
	for _, msg := range c.GeneralChat {
		if msg.CliID == c.ID {
			fmt.Printf("\t\t[This client]: %s", msg.MsgBody)
		} else {
			fmt.Println(msg.MsgBody)
		}
		fmt.Println("\n*************\n")
	}
	fmt.Println("\n+------------------------")
}

func (c *Client) CloseConnection() {
	/* closes connection with the server*/
	if c.Connection != nil {
		c.Connection.Close()
	}
}

func (c *Client) Disconenct() {
	if c.Connection != nil {
		defer c.Connection.Close()
		var msg string = "out"
		gob.NewEncoder(c.Connection).Encode(&msg)
	}
}

/***** Client admin ******/
func (c *Client) LoginAdmin(pass, email string) {
	if pass == c._admi_pass && email == c._admin_email {
		c.AdminLogged = true
	}
}

func (c *Client) IsAdminLogged() bool {
	return c.AdminLogged
}

func (c *Client) ServersStatus() error {
	if !c.IsAdminLogged() {
		return errors.New("not admin")
	}
	data, err := c.ApiCallGetSeversStatus()
	if err != nil {
		return err
	}

	fmt.Println("\n|Server", "\t\tTotal Users", "\t\tIP")
	for _, server := range data {
		fmt.Println(server.Tematic, "\t\t", server.TotalUsers, "\t\t", server.IP)
	}
	fmt.Println()
	return nil
}

func main() {

	/* args */
	var args_email, args_password string
	if len(os.Args) >= 3 {
		args_email = os.Args[1]
		args_password = os.Args[2]
	}

	/* Client init */
	var msg_channel chan string = make(chan string) // channel to comunicate all the sessages th the client will send
	var file_channel chan string = make(chan string)
	/* Initialize client*/
	cli := Client{ApiAdd: "http://localhost:1001", MsgChan: msg_channel, FileChan: file_channel, AdminLogged: false, _admi_pass: "12345", _admin_email: "admin@gmail.com"}
	cli.LoginAdmin(args_password, args_email)

	if !cli.IsAdminLogged() {
		/* First menu */
		var opc string
		for opc != "2" {
			cli.FirstMenu()
			fmt.Scanln(&opc)
			if opc == "1" {
				if cli.ServersMenu() == 0 {
					fmt.Println("Conneccion exitosa")
					break
				}
			} else if opc == "2" {
				os.Exit(0)
			}
		}

		/* When connected to a chat server*/
		go SetupCloseHandler(&cli)
		fmt.Printf("Client {%d} running... on server {%s}\n\n", cli.ID, cli.ChatName)
		// close the connection when this function ends
		defer cli.CloseConnection()
		// function that sends messges to server
		go cli.SendMsg()
		// function that sends files to the server
		go cli.SendFile()
		go cli.ListenForUpdates()

		/* auxiliar variables */
		scanner := bufio.NewScanner(os.Stdin)
		var text string

		/* Main rutine */
		for opc != "4" {
			if !cli.Connected {
				fmt.Println("\n--SERVER DISCONNECTED!!--\n")
			}
			menu(cli.ID, cli.ChatName)
			scanner.Scan()
			opc = scanner.Text()

			if opc == "1" {
				/* creates the message */
				fmt.Print("Msg: ")
				scanner.Scan()
				text = scanner.Text()
				// sends the message through the channel
				msg_channel <- text
			} else if opc == "2" {
				fmt.Println("Filename: ")
				scanner.Scan()
				text = scanner.Text()
				// sends the filename to find the file to be subbmitted to chat room
				file_channel <- text
			} else if opc == "3" {
				fmt.Println("CHAT GENERAL:")
				cli.PrintGlobalChat()
			} else {
				break
			}
		}

		fmt.Println("\nType [Ctrl+C] command to exit..")
		for {
			// just for handling Ctrl+C input
		}
	} else {
		// admin
		fmt.Println("ADMIN")
		/* First menu */
		var opc string
		for opc != "2" {
			fmt.Println(".: ADMIN :.")
			fmt.Println("  1) Estado salas de chat")
			fmt.Println("  2) Salir")
			fmt.Print("Opcion: ")
			fmt.Scanln(&opc)
			if opc == "1" {
				err := cli.ServersStatus()
				if err != nil {
					fmt.Println(err)
				}
			} else if opc == "2" {
				os.Exit(0)
			}
		}
	}

}

func SetupCloseHandler(cli *Client) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		cli.Disconenct()
		os.Exit(0)
	}()
}
func menu(id uint64, serverName string) {
	fmt.Println(".: CLIENT ", id, "| SERVER ", serverName, ":.")
	fmt.Println("\t1) Enviar mensaje")
	fmt.Println("\t2) Enviar archivo")
	fmt.Println("\t3) Mostrar chat")
	fmt.Print("Opcion: ")
}
