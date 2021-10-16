package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
)

func client() {
	c, err := rpc.Dial("tcp", "127.0.1.1:9999")

	if err != nil {
		fmt.Println(err)
		return
	}

	var opc string
	scanner := bufio.NewScanner(os.Stdin)

	for opc != "3" {
		menu()
		scanner.Scan()
		opc = scanner.Text()

		if opc == "1" {
			var name, response string
			fmt.Print("Name: ")
			fmt.Scanln(&name)

			err = c.Call("MicroService.Hello", name, &response)
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println("Res:", response)
			}
		}
	}
}

func main() {
	client()
}

func menu() {
	fmt.Println("\n.: Menu :.\n")
	fmt.Println("1 ) Hello")
	fmt.Print("Opcion: ")
}
