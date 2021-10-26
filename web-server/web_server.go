package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

type WebServer struct {
	Students map[string]map[string]float64
	Subjects map[string]map[string]float64
}

func (w *WebServer) Init() {
	w.Students = make(map[string]map[string]float64)
	w.Subjects = make(map[string]map[string]float64)
}

func (w *WebServer) Run() {
	/* Routes */
	http.HandleFunc("/", w.welcome)
	http.HandleFunc("/add-student-note", w.newStudentNote)
	http.HandleFunc("/student-grade", w.StudentGrade)

	/* Run server*/
	http.ListenAndServe(":9999", nil)
}

func (w *WebServer) StudentGrade(res http.ResponseWriter, req *http.Request) {

	switch req.Method {
	case "GET":
		w.SendFormStudentGrade(res)

	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Println(err)
			return
		}

		w.GetStudentNote(res, req.FormValue("selected-student"))
	}

}

func (w *WebServer) GetStudentNote(res http.ResponseWriter, studentName string) {

	var grade float64 = 0.0
	if v, exists := w.Students[studentName]; exists {
		if exists {
			var avg float64
			for _, sv := range v {
				avg += sv
			}
			grade = (avg) / float64(len(v))
		} else {
			fmt.Println("No existe el estudiante")
		}
	} else {
		fmt.Println("No existe el estudiante")
	}

	res.Header().Set(
		"Content-Type",
		"text/html",
	)

	html, _ := readHTML("./html/response.html")
	s := fmt.Sprintf("%f", grade)
	fmt.Fprintf(res, html, "EL promedio del estudiante "+studentName+" es ", s)
}

func (w *WebServer) SendFormStudentGrade(res http.ResponseWriter) {

	var students string

	for student, _ := range w.Students {
		students += "<option  value='" +
			student + "'>" +
			student + "</option>"
	}
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	html, _ := readHTML("./html/form-select-student.html")

	fmt.Fprintf(res,
		html,
		students,
	)
}

func (w *WebServer) addStudentNote(studentName string, subjectName string, note float64) {
	// add student info
	if _, exists := w.Students[studentName]; exists { //  already exists
		if exists {
			w.Students[studentName][subjectName] = note
		}
	} else {
		subject := make(map[string]float64)
		subject[subjectName] = note
		w.Students[studentName] = subject
	}

	// add subject info
	if _, exists := w.Subjects[subjectName]; exists {
		if exists {
			w.Subjects[subjectName][studentName] = note
		}
	} else {
		// create the first student
		student := make(map[string]float64)
		student[studentName] = note
		// create the subject
		w.Subjects[subjectName] = student
	}
	fmt.Println("\n.: STUDENTS :.")
	for student, subjets := range w.Students {
		fmt.Println(student, "{")
		for subject, sv := range subjets {
			fmt.Println("\t"+subject+": ", sv)
		}
		fmt.Println("}")
	}

	fmt.Println("\n.: SUBJECTS :.")
	for subjet, students := range w.Subjects {
		fmt.Println(subjet, "{")
		for student, sv := range students {
			fmt.Println("\t"+student+": ", sv)
		}
		fmt.Println("}")
	}
}

func main() {
	webServer := WebServer{}
	webServer.Init()
	webServer.Run()

}

func (w *WebServer) welcome(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(res,
		`<DOCTYPE html>
		<html>
			<head>
				<title>Home</title>
			</head>
			<body>
				<h1>Home page!</h1>
			</body>
		</html>`,
	)
}

func (w *WebServer) newStudentNote(res http.ResponseWriter, req *http.Request) {
	/* */
	switch req.Method {
	case "GET":
		w.sendAddStudentNoteForm(res)
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, " ParseForm() $v", err)
			return
		}
		note, _ := strconv.ParseFloat(req.FormValue("note"), 64)
		w.addStudentNote(req.FormValue("student"), req.FormValue("subject"), note)
	}

}

func (w *WebServer) sendAddStudentNoteForm(res http.ResponseWriter) {
	html, err := readHTML("./html/form-add-student-note.html")
	if err != nil {
		fmt.Println(err)
		fmt.Fprintf(res, `
		<DOCTYPE html>
		<html>
			<head>
				<title>Error</title>
			</head>
			<body>
				<p>Un error ha ocurrido</p>
			</body>
		</html>`)
	}
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	fmt.Fprintf(res, html)
}

func readHTML(path string) (string, error) {
	/* reads html files */
	html, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(html), nil
}
