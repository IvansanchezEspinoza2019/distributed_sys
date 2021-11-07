package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"syscall"
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

/* Client struct  */
type Client struct {
	Connection   net.Conn
	ID           uint64
	MsgChan      chan string
	FileChan     chan string
	GeneralChat  []MsgMeta
	HasDirectory bool
	DirName      string
	Connected    bool
}

/*** Client main functions ****/
func (c *Client) Init() error {
	/* connects to the server  */
	con, err := net.Dial("tcp", "localhost:9999")
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

		f := File{Filename: stat.Name(), Content: make([]byte, total)}
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

func (c *Client) CloseConnection() {
	/* closes connection with the server*/
	if c.Connection != nil {
		c.Connection.Close()
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
					meta := MsgMeta{}
					gob.NewDecoder(c.Connection).Decode(&meta)
					c.GeneralChat = append(c.GeneralChat, meta)
				} else if instruction == "file" {
					file := File{}
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

func (c *Client) HandleFile(file *File) {
	c.GeneralChat = append(c.GeneralChat, MsgMeta{CliID: file.Creator, MsgBody: file.Filename})

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

func (c *Client) Disconenct() {
	if c.Connection != nil {
		defer c.Connection.Close()
		var msg string = "out"
		gob.NewEncoder(c.Connection).Encode(&msg)
	}
}

func main() {
	/* Client init */
	var msg_channel chan string = make(chan string) // channel to comunicate all the sessages th the client will send
	var file_channel chan string = make(chan string)
	cli := Client{MsgChan: msg_channel, FileChan: file_channel} // this is the client

	err := cli.Init() // init the connection with the server
	if err == nil {
		/*go SetupCloseHandler(&cli)
		fmt.Printf("Client {%d} running...\n\n", cli.ID)
		defer cli.CloseConnection() // close the connection when this function ends
		go cli.SendMsg()            // function that sends messges to server
		go cli.SendFile()
		go cli.ListenForUpdates()

		/* auxiliar variables */
		//var opc string
		//scanner := bufio.NewScanner(os.Stdin)
		//var text string

		/* Main rutine */
		/*for opc != "4" {
		if !cli.Connected {
			fmt.Println("\n--SERVER DISCONNECTED!!--\n")
		}
		menu(cli.ID)
		scanner.Scan()
		opc = scanner.Text()

		if opc == "1" {
			/*  creates the message */
		/*fmt.Print("Msg: ")
				scanner.Scan()
				text = scanner.Text()
				msg_channel <- text // sends the message through the channel
			} else if opc == "2" {
				fmt.Println("Filename: ")
				scanner.Scan()
				text = scanner.Text()
				file_channel <- text
			} else if opc == "3" {
				fmt.Println("CHAT GENERAL:")
				cli.PrintGlobalChat()
			} else {
				break
			}
		}

		for {
			// just for handling Ctrl+C input
		}*/
	} else {
		fmt.Println("Error: ", err)
	}

}

func menu(id uint64) {
	fmt.Println(".: Client ", id, " :.")
	fmt.Println("\t1) Enviar mensaje")
	fmt.Println("\t2) Enviar archivo")
	fmt.Println("\t3) Mostrar chat")
	fmt.Print("Opcion: ")
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
