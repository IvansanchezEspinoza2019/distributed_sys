package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func Process(id uint64, c chan string) {
	/* this processes are going to be executed in a concurrently way*/

	/* process id */
	i := uint64(0)
	/* controls if the process still executing*/
	ejecution := true

	for { /* Infinite loop*/
		if ejecution {

			/* get the message of the channel */
			select {
			case ctrl_stm := <-c: /* get the control message from the channel */
				/* split the string into an array*/
				array_stm := strings.Split(ctrl_stm, ";")

				if array_stm[0] == "show" && array_stm[1] == "all" {
					/* show the value of 'i' */
					fmt.Printf("\nID{%d} -> %d", id, i)

				} else if array_stm[0] == "delete" {
					/* delete this process only if the id matches with the id of this process*/
					id_to_delete, _ := strconv.Atoi(array_stm[1])
					if uint64(id_to_delete) == id {
						/* stop the execution of this process*/
						ejecution = false
					}
				}
			}
			i += 1
			time.Sleep(time.Millisecond * 500)
		} else {
			/* end of this process execution*/
			break
		}
	}
}

func signalSender(c chan string, value *string) {
	/* sends all controls messages through the channel*/
	for {
		c <- *value
	}
}

func main() {
	//  process identifiers
	id := uint64(0)

	// channel to communicate all processes
	var c chan string = make(chan string)

	//  control
	var CONTROL string = "show;off"

	/* goroutine that send all control messages to all processes through the channel */
	go signalSender(c, &CONTROL)

	scanner := bufio.NewScanner(os.Stdin)
	var opc string

	for opc != "4" {
		// menu
		menu()
		scanner.Scan()
		opc = scanner.Text()

		switch opc {
		case "1":
			/* create new process with the communication channel as parameter*/
			go Process(id, c)
			fmt.Printf("\nProcess ID{%d} created..\n", id)
			id++

		case "2":
			/* Show all/show off all the processes*/
			if CONTROL == "show;all" {
				CONTROL = "show;off"
			} else if CONTROL == "show;off" {
				CONTROL = "show;all"
			} else {
				CONTROL = "show;all"
			}

		case "3":
			/* id of the process*/
			var deleteID uint64
			fmt.Print("[ID to delete]: ")
			fmt.Scanln(&deleteID)
			/* message control */
			CONTROL = "delete;" + strconv.Itoa(int(deleteID))
		case "4":
			fmt.Println("Bye")
		default:
			fmt.Println("Opción inválida")
		}
	}
}

func menu() {
	fmt.Println("\t\t\t.: Menu :.")
	fmt.Println("\t\t1) Agregar proceso")
	fmt.Println("\t\t2) Mostrar proceso")
	fmt.Println("\t\t3) Terminar proceso")
	fmt.Println("\t\t4) Salir")
	fmt.Print("Opcion: ")
}
