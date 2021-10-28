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
	http.HandleFunc("/general-grade", w.GeneralGrade)
	http.HandleFunc("/subject-grade", w.SubjectGrade)

	/* Run server*/
	http.ListenAndServe(":9999", nil)
}

/* Handler functions (listen por client requests)*/
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

func (w *WebServer) SubjectGrade(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		w.sendFormSubjectGrade(res)
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Println(err)
			return

		}

		w.GetSubjectNote(res, req.FormValue("selected-subject"))
	}
}

func (w *WebServer) newStudentNote(res http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		w.sendAddStudentNoteForm(res)
	case "POST":
		if err := req.ParseForm(); err != nil {
			fmt.Fprintf(res, " ParseForm() $v", err)
			return
		}
		note, _ := strconv.ParseFloat(req.FormValue("note"), 64)
		w.addStudentNote(res, req.FormValue("student"), req.FormValue("subject"), note)
	}

}

func (w *WebServer) GeneralGrade(res http.ResponseWriter, req *http.Request) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	html, _ := readHTML("./html/response.html")

	/* calculates the note aberage of all students */

	if len(w.Students) == 0 {
		fmt.Fprintf(res, "Vacio")
		return
	}

	var students string
	var totalAvg float64

	students = "<h2>Promedio de cada estudiante</h2>"
	for student, subjects := range w.Students {
		var avg float64
		for _, sv := range subjects {
			avg += sv
		}
		avg = (avg / float64(len(subjects)))
		totalAvg += avg
		students += "<li> " + student + " <strong>(" + strconv.FormatFloat(avg, 'f', 2, 64) + ")</strong> </li>"
	}
	generalAvg := "<br><div><h2>Promedio General:  <span><strong>" + strconv.FormatFloat((totalAvg/float64(len(w.Students))), 'f', 2, 64) + "</strong></h2></div>"
	fmt.Fprintf(res, html, "<ul> "+students+"</ul>"+generalAvg)
}

func (w *WebServer) sendFormSubjectGrade(res http.ResponseWriter) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)
	html, _ := readHTML("./html/form-select-subject.html")

	var subjects string
	for subject, _ := range w.Subjects {
		subjects += "<option value='" + subject + "'>" + subject + "</option>"
	}
	fmt.Fprintf(res, html, subjects)
}

func (w *WebServer) GetStudentNote(res http.ResponseWriter, studentName string) {
	/* Returns the student grade*/

	res.Header().Set(
		"Content-Type",
		"text/html",
	)

	html, _ := readHTML("./html/response.html")

	var grade float64 = 0.0
	if v, exists := w.Students[studentName]; exists {
		if exists {
			var avg float64
			for _, sv := range v {
				avg += sv
			}
			grade = (avg) / float64(len(v))
		} else {
			fmt.Fprintf(res, html, "No existe el estudiante")
			return
		}
	} else {
		fmt.Fprintf(res, html, "No existe el estudiante")
		return

	}
	var s string = "EL promedio de '" + studentName + "' es <strong>" + strconv.FormatFloat(grade, 'f', 2, 64) + "</strong>"
	fmt.Fprintf(res, html, s)
}

func (w *WebServer) GetSubjectNote(res http.ResponseWriter, subjectName string) {
	res.Header().Set(
		"Content-Type",
		"text/html",
	)

	html, _ := readHTML("./html/response.html")
	var grade float64 = 0.0
	/* calculates the avdg of a subject */

	if len(w.Subjects) == 0 {
		fmt.Fprintf(res, "Vacio")
		return
	}
	if students, exists := w.Subjects[subjectName]; exists {
		if exists {
			var svg float64
			for _, sv := range students {
				svg += sv
			}
			grade = (svg / float64(len(students)))
		} else {
			fmt.Fprintf(res, "NO existe la materia")
			return
		}
	} else {
		fmt.Fprintf(res, "NO existe la materia")
		return
	}

	fmt.Fprintf(res, html, "La materia '"+subjectName+"' tiene <strong>"+strconv.FormatFloat(grade, 'f', 2, 64)+"</strong>")
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

func (w *WebServer) addStudentNote(res http.ResponseWriter, studentName string, subjectName string, note float64) {

	res.Header().Set(
		"Content-Type",
		"text/html",
	)

	html, _ := readHTML("./html/response.html")
	var response string
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
	response = "<div><h2>STUDENTS</h2>"
	for student, subjets := range w.Students {
		response += student + "<ul>"
		for subject, sv := range subjets {
			response += "<li>" + subject + " <strong>" + strconv.FormatFloat(sv, 'f', 2, 64) + "</strong></li>"
		}
		response += "</ul>"
	}

	response += "<div><h2>SUBJECTS</h2>"
	for subjet, students := range w.Subjects {
		response += subjet + "<ul>"
		for student, sv := range students {
			response += "<li>" + student + " <strong>" + strconv.FormatFloat(sv, 'f', 2, 64) + "</strong></li>"

		}
		response += "</ul>"
	}
	fmt.Fprintf(res, html, response)
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

	html, _ := readHTML("./html/home.html")

	fmt.Fprintf(res, html)

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
