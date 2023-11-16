package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/mattn/go-sqlite3"
)

const db_fname = "./data/tasks.db"

type Task struct {
	id int64;
	name string;
	description string;
	completed bool;
}

func main() {
	os.Remove(db_fname)

	db, err := sql.Open("sqlite3", db_fname)
	handleErr(err)
	defer db.Close()

	var sqliteVersion, _, _ = sqlite3.Version()
	fmt.Println("SQLite version:", sqliteVersion)

	sqlStatement := `
		CREATE TABLE tasks (id integer not null primary key, name text, description text, completed int);
		DELETE FROM tasks;
	`
	_, err = db.Exec(sqlStatement)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStatement)
		return
	}

	id1 := addTask(db, "test1", "test1", false)
	id2 := addTask(db, "test2", "test2", false)
	_ = addTask(db, "test3", "test3", true)
	_ = addTask(db, "test4", "test4", true)
	delTask(db, id2)
	var task Task
	task.id = id1
	task.name = "edited"
	task.description = "also edited"
	task.completed = true
	var editedRows = editTask(db, task)
	fmt.Println("Edited rows", editedRows)
	var tasks = getTasks(db)

	for _, value := range tasks {
		fmt.Println(value)
	}

	// transaction, err := db.Begin()
	// handleErr(err)


	// statement, err := transaction.Prepare("INSERT INTO tasks(id, name, description, completed) values(?, ?, ?, ?)")
	// handleErr(err)
	
	// defer statement.Close()
	// // add tasks using transcation
	// for i := 0; i < 5; i++ {
	// 	_, err = statement.Exec(1, "test", "some desc", 0)
	// 	handleErr(err)
	// }
	// // commit transaction
	// err = transaction.Commit()
	// handleErr(err)

	// get tasks
	// tasks, err := db.Query(`
	// 	SELECT 
	// 		id, name, description, completed 
	// 	FROM 
	// 		tasks;
	// `)
	// handleErr(err)
	// defer tasks.Close()

	// // print tasks
	// for tasks.Next() {
	// 	var task Task
	// 	err  = tasks.Scan(&task.id, &task.name, &task.description, &task.completed)
	// 	handleErr(err)
	// 	fmt.Println(task)
	// }

	// err  = tasks.Err()
	// handleErr(err)
}

func addTask(db *sql.DB, name string, description string, completed bool) int64 {
	sqlStatement := `
		INSERT INTO 
			tasks(id, name, description, completed) 
			values (
				(SELECT MAX(id) FROM tasks LIMIT 1) + 1, ?, ?, ?
			);`
	res, err := db.Exec(sqlStatement, name, description, completed)
	handleErr(err)
	id, err := res.LastInsertId()
	handleErr(err)
	return id
}

func delTask(db *sql.DB, id int64) {
	sqlStatement := `
		DELETE FROM tasks WHERE id = ?;`
	_, err := db.Exec(sqlStatement, id)
	handleErr(err)
}

func editTask(db *sql.DB, task Task) int64 {
	// // get task first
	// rows, err := db.Query("SELECT FROM tasks WHERE id = ?", id)
	// handleErr(err)
	// var task Task
	// err = rows.Scan(&task.id, &task.name, &task.description, &task.completed)
	// handleErr(err)

	// // only update updated fields

	// update task
	updateStatement := `
		UPDATE tasks
		SET 
			name = ?,
			description = ?,
			completed = ?
		WHERE id = ?;`
	res, err := db.Exec(updateStatement, task.name, task.description, task.completed, task.id)
	handleErr(err)
	numRows, err := res.RowsAffected()
	handleErr(err)
	return numRows
}

func getTasks(db *sql.DB) []Task {
	// get tasks
	rows, err := db.Query(`
		SELECT 
			id, name, description, completed 
		FROM 
			tasks;
	`)
	handleErr(err)
	defer rows.Close()

	var tasks = []Task{}

	// print tasks
	for rows.Next() {
		var task Task
		err = rows.Scan(&task.id, &task.name, &task.description, &task.completed)
		handleErr(err)
		tasks = append(tasks, task)
	}

	err  = rows.Err()
	handleErr(err)
	return tasks
}

func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}