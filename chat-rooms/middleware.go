package main

import (
	"fmt"
	"net/http"
	"net/rpc"
)

type ServerInfo struct {
	Temtic     string
	TotalUsers uint64
}

/* Middleware */
type MiddleWareApi struct {
	Port      string
	RpcClient rpc.Client
}

func (m *MiddleWareApi) Init() error { // conects to chat servers
	client, err := rpc.Dial("tcp", "localhost:8000")
	if err != nil {
		return err
	}
	m.RpcClient = *client
	return nil
}

func (m *MiddleWareApi) RunApi() {
	/* Run Api*/
	http.HandleFunc("/chatRooms", m.ChatRooms)

	fmt.Println("RESTful API Running on http://localhost:" + m.Port)
	http.ListenAndServe(":"+m.Port, nil)
}

func (m *MiddleWareApi) ChatRooms(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		m.GetChatRooms(res)
	}
}

func (m *MiddleWareApi) GetChatRooms(res http.ResponseWriter) {
	/*var response []ServerInfo
	err := m.RpcClient.Call("MicroService.GetChatRooms", "nil", &response)
	if err != nil {
		fmt.Println("[ERROR]", err)
	} else {
		fmt.Println("\nROOMS: [", response, "]")
	}*/
	res.Header().Set(
		"Content-Type",
		"application/json",
	)
	res.Write([]byte(`{"Message":"Api initialized"}`))
}

func main() {
	api := MiddleWareApi{Port: "1001"}
	err := api.Init()
	if err != nil {
		fmt.Println("Error ocurred:", err)
		return
	}
	api.RunApi()
}

func client() {
	c, errConn := rpc.Dial("tcp", "localhost:8000/api/v1")

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
