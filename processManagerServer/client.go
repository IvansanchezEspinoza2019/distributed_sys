package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"time"

	"./process"
)

type Client struct {
	p *process.Process
}

func (cli *Client) clientAskProcess() {
	/*  connect to the server */
	c, err := net.Dial("tcp", ":3302")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer c.Close()

	/* request process */
	msg := "get_process"
	gob.NewEncoder(c).Encode(msg)

	cli.p = new(process.Process)
	/* get the  process*/
	gob.NewDecoder(c).Decode(&cli.p)
	//fmt.Printf("Process {%d} received successfully...\n", p.ID)

	/* Handle process execution */
	go cli.executeProcess()
}

func (clie *Client) executeProcess() {
	// execute process
	go clie.p.Execute()
	for {
		fmt.Printf("%d : %d \n", clie.p.ID, clie.p.ICount)
		time.Sleep(time.Millisecond * 500)
	}
}

func main() {

	client := Client{}
	defer client.clientSendProcess()
	go client.clientAskProcess()

	var input string
	fmt.Scanln(&input)
}

func (cli *Client) clientSendProcess() {

	if cli.p != nil {
		cli.p.Stop()
		/*  connect to the server */
		c, err := net.Dial("tcp", ":3302")
		if err != nil {
			fmt.Println(err)
			return
		}
		defer c.Close()

		msg := "put_process"
		gob.NewEncoder(c).Encode(msg)

		//  send the process
		gob.NewEncoder(c).Encode(cli.p)
	} else {
		fmt.Println("Ningún proceso corriendo")
	}
}
