package main

import (
	"encoding/json"
	"errors"
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
	Student string  `json: "estudiante"`
	Subject string  `json: "materia"`
	Note    float64 `json: "promedio"`
}

var Students map[uint64]Student

/************** Methods **************/
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
	jsonRes := []byte(`{"Msg": "Ok"}`)
	_, exists := Students[data.ID]
	if exists {
		Students[data.ID].Subjects[data.Subject] = data.Note
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
		return nil, errors.New("Student NOT found")
	}

	_, existsSub := Students[id].Subjects[data.Subject]
	if !existsSub {
		return nil, errors.New("Subject NOT found")
	}
	Students[id].Subjects[data.Subject] = data.Note
	return jsonRes, nil
}

/* Get one student*/
func GetStudent(id uint64) ([]byte, error) {

	student, exists := Students[id]
	if !exists {
		return nil, errors.New("Student NOT found")
	}
	jsonData, err := json.MarshalIndent(student, "", "  ")
	if err != nil {
		return nil, errors.New(err.Error())
	}
	return jsonData, nil
}

/* delete a student */
func DeleteStudent(id uint64) ([]byte, error) {
	jsonRes := []byte(`{"Msg": "Ok"}`)
	_, exists := Students[id]
	if !exists {
		return nil, errors.New("Student NOT found")
	}
	delete(Students, id)
	return jsonRes, nil
}

/*********** Handle requests *********/
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
			if errUpdate.Error() == "Subject NOT found" || errUpdate.Error() == "Student NOT found" {
				http.Error(res, errUpdate.Error(), http.StatusNotFound)
				return
			}
			http.Error(res, errUpdate.Error(), http.StatusInternalServerError)
			return
		}

		res.Header().Set(
			"Content-Type",
			"application/json",
		)

		res.Write(jsonRes)
	case "GET":
		studentJson, err := GetStudent(id)
		if err != nil {
			if err.Error() == "Student NOT found" {
				http.Error(res, err.Error(), http.StatusNotFound)
				return
			}
			http.Error(res, err.Error(), http.StatusInternalServerError)
		}

		res.Header().Set(
			"Content-Type",
			"application/json",
		)
		res.Write(studentJson)

	case "DELETE":
		resJson, err := DeleteStudent(id)
		if err != nil {
			http.Error(res, err.Error(), http.StatusNotFound)
			return
		}

		res.Header().Set(
			"Content-Type",
			"application/json",
		)
		res.Write(resJson)
	}
}
func main() {
	Students = make(map[uint64]Student)

	// server
	http.HandleFunc("/students", HandleStudents)
	http.HandleFunc("/students/", HandleStudentId)
	http.ListenAndServe(":9999", nil)
}
