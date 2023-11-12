package main

import (
	"database/sql"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "users.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()
	//creating tables for first use
	createTables(db)
	myDbHandler := DBHolder{db: db}

	http.HandleFunc("/favicon.ico", myDbHandler.EmptyPage)
	http.HandleFunc("/submitAddUser", myDbHandler.handleSubmitAddUser)
	http.HandleFunc("/submitAddDepartment", myDbHandler.handleSubmitAddDepartment)
	http.HandleFunc("/submitDeleteUser", myDbHandler.handleSubmitDeleteUser)
	http.HandleFunc("/submitDeleteDepartment", myDbHandler.handleSubmitDeleteDepartment)
	http.HandleFunc("/", myDbHandler.ShowUsersTable)

	if err := http.ListenAndServeTLS(":8443", "server.crt", "server.key", nil); err != nil {
		panic(err)
	}
}
