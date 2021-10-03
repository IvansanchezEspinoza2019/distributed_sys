package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

/* Client struct  */
type Client struct {
	Connection net.Conn
	MsgChan    chan string
}

/*** Client main functions ****/
func (c *Client) Init() error {
	/* connects to the server  */
	con, err := net.Dial("tcp", ":5555")
	if err != nil {
		return err
	}
	c.Connection = con
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

func (c *Client) CloseConnection() {
	/* closes connection with the server*/
	if c.Connection != nil {
		c.Connection.Close()
	}
}

func main() {
	/* Client init */
	var msg_channel chan string = make(chan string) // channel to comunicate all the sessages th the client will send
	cli := Client{MsgChan: msg_channel}             // this is the client

	err := cli.Init() // init the connection with the server
	if err == nil {
		defer cli.CloseConnection() // close the connection when this function ends
		go cli.SendMsg()            // function that sends messges to server

		/* auxiliar variables */
		var opc string
		scanner := bufio.NewScanner(os.Stdin)
		var text string

		/* Main rutine */
		for opc != "2" {
			menu()
			scanner.Scan()
			opc = scanner.Text()

			if opc == "1" {
				/*  creates the message */
				fmt.Print("Msg: ")
				scanner.Scan()
				text = scanner.Text()
				msg_channel <- text // sends the message through the channel
			} else {
				break
			}
		}
	} else {
		fmt.Println("Error: ", err)
	}

}

func menu() {
	fmt.Println("\t1) Enviar mensaje")
	fmt.Print("Opcion: ")
}
