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

	for opc != "5" {
		menu()
		scanner.Scan()
		opc = scanner.Text()

		if opc == "1" {
			var response Reply
			data := &StudentSubject{}

			fmt.Print("Student name: ")
			fmt.Scanln(&data.Name)

			fmt.Print("Subject: ")
			fmt.Scanln(&data.Subject)

			fmt.Print("Student Note: ")
			fmt.Scanf("%f", &data.Note)
			fmt.Scanln()

			err = c.Call("MicroService.AddNoteStudentSubjetc", data, &response)
			if err != nil {
				fmt.Println("[ERROR]", err)
			} else {
				fmt.Println("Res:", response)
			}
		} else if opc == "2" {
			var (
				name     string
				response float64
			)

			fmt.Print("Student: ")
			fmt.Scanln(&name)

			err = c.Call("MicroService.StudentAverege", name, &response)
			if err != nil {
				fmt.Println("[ERROR]", err)
			} else {
				fmt.Println("\tNOTE AVERAGE: [", response, "]")
			}

		} else if opc == "3" {
			var (
				response float64
			)

			err = c.Call("MicroService.GeneralAverage", "General", &response)
			if err != nil {
				fmt.Println("[ERROR]", err)
			} else {
				fmt.Println("\nGENERAL AVERAGE: [", response, "]")
			}
		} else if opc == "4" {
			var (
				subject  string
				response float64
			)
			fmt.Print("Subject: ")
			fmt.Scanln(&subject)

			err = c.Call("MicroService.SubjectAvg", subject, &response)
			if err != nil {
				fmt.Println("[ERROR]", err)
			} else {
				fmt.Println("\tSUBJECT AVERAGE: [", response, "]")
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
	fmt.Println("2 ) Promedio de un Alumno")
	fmt.Println("3 ) Promedio General")
	fmt.Println("4 ) Promedio Materia")
	fmt.Print("Opcion: ")
}
