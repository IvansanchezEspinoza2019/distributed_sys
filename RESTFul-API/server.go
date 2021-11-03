package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type Student struct {
	ID       uint64
	Student  string
	Subjects map[string]float64
}

type StudentData struct {
	ID      uint64  `json: "id"`
	Student string  `json:  "estudiante"`
	Subject string  `json: "materia"`
	Note    float64 `json: "promedio"`
}

var Students map[uint64]Student
var Subjects map[string]map[string]float64

/***** Methods ******/
/*All students*/
func GetAllStudents() ([]byte, error) {
	jsonData, err := json.MarshalIndent(Students, "", "    ")
	if err != nil {
		return jsonData, err
	}
	return jsonData, nil
}

/* add student data*/
func AddStudentData(data StudentData) ([]byte, error) {
	fmt.Println(data)
	jsonRes := []byte(`{"Msg": "Ok"}`)
	_, exists := Students[data.ID]
	if exists {
		//Students[data.Student][data.Subject] = data.Note
		return []byte(`{"Msg":"Ya añadiste a este estudiante, si queires actualizar su informacion us PUT"}`), errors.New("Estudiante repetido")
	} else {
		student := Student{ID: data.ID, Subjects: make(map[string]float64), Student: data.Student}
		student.Subjects[data.Subject] = data.Note
		Students[data.ID] = student
	}
	return jsonRes, nil
}

/* update student note*/
func UpdateNote(id uint64, data StudentData) ([]byte, error) {
	jsonRes := []byte(`{"Msg": "Ok"}`)

	_, exists := Students[id]
	if !exists {
		return []byte(`{"Msg":"NO existe el estudiante"}`), errors.New("Student NOT found")
	}

	_, existsSub := Students[id].Subjects[data.Subject]
	if !existsSub {
		return []byte(`{"Msg":"NO existe la materia"}`), errors.New("Subject NOT found")
	}
	Students[id].Subjects[data.Subject] = data.Note
	fmt.Println("Aqui esta", Students[id])
	return jsonRes, nil
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
	case "POST":
		var data StudentData
		err := json.NewDecoder(req.Body).Decode(&data)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		json, err := AddStudentData(data)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		res.Header().Set(
			"Content-Type",
			"application/json",
		)
		res.Write(json)
	}
}

func HandleStudentId(res http.ResponseWriter, req *http.Request) {

	/*  get the id from the url*/
	id, err := strconv.ParseUint(strings.TrimPrefix(req.URL.Path, "/students/"), 10, 64) // uint64

	if err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}

	switch req.Method {
	case "PUT":
		var updateS StudentData
		err = json.NewDecoder(req.Body).Decode(&updateS)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonRes, errUpdate := UpdateNote(id, updateS)
		if errUpdate != nil {
			http.Error(res, errUpdate.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set(
			"Content-Type",
			"application/json",
		)

		res.Write(jsonRes)

	}
}
func main() {
	Students = make(map[uint64]Student)

	// server
	http.HandleFunc("/students", HandleStudents)
	http.HandleFunc("/students/", HandleStudentId)
	http.ListenAndServe(":9999", nil)
}
