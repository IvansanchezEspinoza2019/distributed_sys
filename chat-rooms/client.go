package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)

type ServerInfo struct {
	Temtic     string
	TotalUsers uint64
}

/* strunct client */
type Client struct {
	Connection   net.Conn
	ID           uint64
	MsgChan      chan string
	FileChan     chan string
	HasDirectory bool
	DirName      string
	Connected    bool
	ApiAdd       string
}

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
	var servers []ServerInfo
	json.Unmarshal([]byte(responseData), &servers)

	fmt.Println("New", servers)

}

func (c *Client) FirstMenu() {
	fmt.Println("  1) Conectarse a una sala de chat")
	fmt.Println("  2) Salir")
	fmt.Print("Opcion: ")
}

func (c *Client) ApiCallGetRooms() ([]ServerInfo, error) {
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
	var servers []ServerInfo
	json.Unmarshal([]byte(responseData), &servers)

	return servers, nil
}

func (c *Client) ServersMenu() int {
	servers, errCall := c.ApiCallGetRooms()
	if errCall != nil {
		fmt.Println(errCall)
		return 1
	}
	fmt.Println("\t.: SELLECCIONA UN CHAT:.\n")
	for _, server := range servers {
		fmt.Println("\tChat: ["+server.Temtic+"]", "\tUsers: (", server.TotalUsers, ")")
	}
	fmt.Print("Server: ")
	var selectedServer string
	fmt.Scanln(&selectedServer)

	// falta hacer l peticion del server seleccionado y obtener la direccion ip de la sala de chat
	return 0
}

func main() {
	/* Initialize client*/
	cli := Client{ApiAdd: "http://localhost:1001"}

	var opc string
	for opc != "2" {
		cli.FirstMenu()
		fmt.Scanln(&opc)
		if opc == "1" {
			cli.ServersMenu()
		} else if opc == "2" {
			os.Exit(0)
		}
	}

	//go c1.MakePetitionApi()
	var input string
	fmt.Scanln(&input)
}
