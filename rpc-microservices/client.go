package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
)

/* server response */
type Reply struct {
	Msg    string
	Status int
}

/* struct to save client request for saving a student note*/
type StudentSubject struct {
	Name    string
	Subject string
	Note    float64
}

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
			var response Reply
			data := &StudentSubject{}

			fmt.Print("Name: ")
			fmt.Scanln(&data.Name)

			fmt.Print("Subject: ")
			fmt.Scanln(&data.Subject)

			fmt.Print("Note: ")
			fmt.Scanf("%f", &data.Note)
			fmt.Scanln()

			err = c.Call("MicroService.AddNoteStudentSubjetc", data, &response)
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
	fmt.Println("1 ) Agregar nota de estudiante")
	fmt.Print("Opcion: ")
}
