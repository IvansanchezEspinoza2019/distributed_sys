package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/rpc"
	"strings"

	"../common"
)

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

/******* handler functions *******/
func (m *MiddleWareApi) ChatRoomID(res http.ResponseWriter, req *http.Request) {
	// convertimos el {id} de la url a uint64
	id := strings.TrimPrefix(req.URL.Path, "/chatRooms/")

	switch req.Method {
	case "GET":
		{
			IP, err := m.FindIpServer(id)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			res.Header().Set(
				"Content-Type",
				"application/json",
			)
			res.Write(IP)
		}
	}
}

func (m *MiddleWareApi) ChatRooms(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		resJson, errJson := m.GetChatRooms(res)
		if errJson != nil {
			http.Error(res, errJson.Error(), http.StatusInternalServerError)
			return
		}
		/* Response to the client */
		res.Header().Set(
			"Content-Type",
			"application/json",
		)
		res.Write(resJson)
	}
}

func (m *MiddleWareApi) ServersStatus(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		{
			resData, err := m.GetServersStatus()
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			res.Header().Set(
				"Content-Type",
				"application/json",
			)
			res.Write(resData)
		}
	}
}

// run api
func (m *MiddleWareApi) RunApi() {
	/* Run Api*/
	http.HandleFunc("/chatRooms", m.ChatRooms)
	http.HandleFunc("/chatRooms/", m.ChatRoomID)
	http.HandleFunc("/chatRooms/status", m.ServersStatus) // requires adminpermisions

	fmt.Println("RESTful API Running on http://localhost:" + m.Port)
	http.ListenAndServe("localhost:"+m.Port, nil)
}

// functions
func (m *MiddleWareApi) GetChatRooms(res http.ResponseWriter) ([]byte, error) {
	/* Make petition to the chat server*/
	var response []common.ServerInfo
	err := m.RpcClient.Call("MicroService.GetChatRooms", "nil", &response)
	if err != nil {
		return nil, err
	}
	resJson, errJson := json.MarshalIndent(response, "", "  ")
	if errJson != nil {
		return nil, errJson
	}
	return resJson, nil
}

func (m *MiddleWareApi) FindIpServer(serverName string) ([]byte, error) {
	var IPServer string
	errRpc := m.RpcClient.Call("MicroService.GetChatRoomIP", serverName, &IPServer)
	if errRpc != nil {
		return nil, errRpc
	}
	responseJson, errJson := json.MarshalIndent(IPServer, "", "   ")
	if errJson != nil {
		return nil, errJson
	}
	return responseJson, nil
}

func (m *MiddleWareApi) GetServersStatus() ([]byte, error) {
	var serverStatus []common.ServerDetail
	// mget data from server
	errRpc := m.RpcClient.Call("MicroService.ServersStatus", "", &serverStatus)

	if errRpc != nil {
		return nil, errRpc
	}
	// convert to json format
	responseJson, errJson := json.MarshalIndent(serverStatus, "", "   ")
	if errJson != nil {
		return nil, errJson
	}
	return responseJson, nil
}

// main
func main() {
	api := MiddleWareApi{Port: "1001"}
	err := api.Init()
	if err != nil {
		fmt.Println("Error ocurred:", err)
		return
	}
	api.RunApi()
}
