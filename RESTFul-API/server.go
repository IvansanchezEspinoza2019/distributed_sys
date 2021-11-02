package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

var Students map[string]map[string]float64
var Subjects map[string]map[string]float64

func GetAllStudents() ([]byte, error) {
	jsonData, err := json.MarshalIndent(Students, "", "    ")
	if err != nil {
		return jsonData, err
	}
	return jsonData, nil
}
func getAllSubjects() {

}

func HandleStudents(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		studentsJson, err := GetAllStudents()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set(
			"Content-Type",
			"application/json",
		)
		res.Write(studentsJson)
	}
}
func main() {

	Students = make(map[string]map[string]float64)
	fmt.Println(Students)
	Students["Pepe"] = make(map[string]float64)

	Students["Pepe"]["mates"] = 65

	http.HandleFunc("/students", HandleStudents)
	http.ListenAndServe(":9999", nil)
}
