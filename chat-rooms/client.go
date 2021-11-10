package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

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

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(responseData))
}

func main() {
	/* initialize client*/
	c1 := Client{ApiAdd: "http://localhost:1001"}
	go c1.MakePetitionApi()
	var input string
	fmt.Scanln(&input)
}
