package main

import (
	"fmt"
	"net"
)

func server() {
	server, err := net.Listen("tcp", ":3302")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(server)
	fmt.Println("Server running on http://localhost:3302")

	for {
		client, err := server.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println(client)
	}
}

func main() {
	go server()

	var execute string
	fmt.Scanln(&execute)
}
