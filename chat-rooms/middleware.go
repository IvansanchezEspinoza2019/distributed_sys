package main

import (
	"fmt"
	"net/rpc"
)

type ServerInfo struct {
	Temtic     string
	TotalUsers uint64
}

func main() {
	fmt.Println("Spy el middleware")
	go client()
	var input string
	fmt.Scanln(&input)
}

func client() {
	c, errConn := rpc.Dial("tcp", "localhost:8000")

	if errConn != nil {
		fmt.Println(errConn)
		return
	}
	var response []ServerInfo

	err2 := c.Call("MicroService.GetChatRooms", "nil", &response)
	if err2 != nil {
		fmt.Println("[ERROR]", err2)
	} else {
		fmt.Println("\nROOMS: [", response, "]")
	}
}
