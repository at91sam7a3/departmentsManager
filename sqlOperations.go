package main

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func createTables(db *sql.DB) {
	// Create "departments" table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS departments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL
		)
	`)
	if err != nil {
		panic(err)
	}

	// Create "users" table with foreign key "Department"
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			First_Name TEXT NOT NULL,
			Last_Name TEXT NOT NULL,
			department INTEGER,
			FOREIGN KEY (department) REFERENCES departments(id)
		)
	`)
	if err != nil {
		panic(err)
	}
}

func removeUser(db *sql.DB, user_name string, last_name string) {
	// Prepare the SQL statement for deletion
	stmt, err := db.Prepare("DELETE FROM users WHERE First_Name = ? AND Last_Name = ?")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(user_name, last_name)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func removeDepartment(db *sql.DB, departmente string) {
	departmentID := getDepartmentId(db, departmente)

	// Prepare the SQL statement for deletion
	stmt, err := db.Prepare("DELETE FROM users WHERE department = ?")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(departmentID)
	if err != nil {
		fmt.Println(err)
		return
	}

	stmt, err = db.Prepare("DELETE FROM departments WHERE departments.id = ?")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer stmt.Close()

	// Execute the SQL statement
	_, err = stmt.Exec(departmentID)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func addUser(db *sql.DB, user_name string, last_name string, department string) {
	dep_id := getDepartmentId(db, department)
	_, err := db.Exec("INSERT INTO users (First_Name,Last_Name,department) VALUES ((?),(?),(?))", user_name, last_name, dep_id)
	if err != nil {
		panic(err)
	}
}

func addDepartment(db *sql.DB, department_name string) {
	_, err := db.Exec("INSERT INTO departments (name) VALUES ((?))", department_name)
	if err != nil {
		panic(err)
	}
}

func getDepartmentId(db *sql.DB, department_name string) string {
	rows, err := db.Query("SELECT departments.id FROM departments WHERE departments.name=(?)", department_name)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id string
		err := rows.Scan(&id)
		if err != nil {
			panic(err)
		}

		fmt.Printf("ID: %s\n", id)
		return id
	}
	return ""
}

func queryDepartmentData(db *sql.DB) []string {
	departments := []string{}
	rows, err := db.Query("SELECT  name FROM departments")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		err := rows.Scan(&name)
		if err != nil {
			panic(err)
		}
		departments = append(departments, name)
	}
	return departments
}

func queryUserData(db *sql.DB, departnemtName string) [][]string {
	requestString := fmt.Sprintf("SELECT users.First_Name, users.Last_Name FROM users JOIN departments ON users.department = departments.id WHERE departments.name = '%s';", departnemtName)
	rows, err := db.Query(requestString)
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	result := [][]string{}
	for rows.Next() {
		var name string
		var lastname string
		err := rows.Scan(&name, &lastname)
		if err != nil {
			panic(err)
		}
		row := []string{name, lastname}
		result = append(result, row)
	}
	return result
}
