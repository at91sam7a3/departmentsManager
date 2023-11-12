package main

import (
	"database/sql"
	"html/template"
	"net/http"
)

type DepartmentData struct {
	TableData []string
}

type UsersData struct {
	Department string
	TableData  [][]string
}

func (db DBHolder) EmptyPage(_ http.ResponseWriter, _ *http.Request) {

}

type DBHolder struct {
	db *sql.DB
}

func (db DBHolder) ShowDepartmentsTable(w http.ResponseWriter, r *http.Request) {
	tableData := queryDepartmentData(db.db)
	pageVariables := DepartmentData{
		TableData: tableData,
	}

	// Parse the HTML template
	tmpl, err := template.ParseFiles("templates/departments.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template and write the output to the response writer
	err = tmpl.Execute(w, pageVariables)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (db DBHolder) ShowUsersTable(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "" || r.URL.Path == "/" {
		db.ShowDepartmentsTable(w, r)
		return
	}
	// Sample data for the table
	tableData := queryUserData(db.db, r.URL.Path[1:])

	pageVariables := UsersData{
		Department: r.URL.Path[1:],
		TableData:  tableData,
	}

	// Parse the HTML template
	tmpl, err := template.ParseFiles("templates/users.html")

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Execute the template and write the output to the response writer
	err = tmpl.Execute(w, pageVariables)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (db DBHolder) handleSubmitAddUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	name := r.FormValue("first_name")
	lastName := r.FormValue("last_name")
	department := r.URL.Query().Get("department")
	addUser(db.db, name, lastName, department)
	r.URL.Path = "/" + department
	db.ShowUsersTable(w, r)
}

func (db DBHolder) handleSubmitAddDepartment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	err := r.ParseForm()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	department_name := r.FormValue("department_name")
	addDepartment(db.db, department_name)
	r.URL.Path = "/"
	db.ShowDepartmentsTable(w, r)
}

func (db DBHolder) handleSubmitDeleteDepartment(w http.ResponseWriter, r *http.Request) {
	department := r.URL.Query().Get("department")
	removeDepartment(db.db, department)
	r.URL.Path = "/"
	db.ShowDepartmentsTable(w, r)
}

func (db DBHolder) handleSubmitDeleteUser(w http.ResponseWriter, r *http.Request) {
	first_name := r.URL.Query().Get("first_name")
	last_name := r.URL.Query().Get("last_name")
	department := r.URL.Query().Get("department")
	removeUser(db.db, first_name, last_name)
	r.URL.Path = "/" + department
	db.ShowUsersTable(w, r)
}
