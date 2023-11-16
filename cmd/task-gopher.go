package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/mattn/go-sqlite3"
)

const db_fname = "./data/tasks.db"

type Task struct {
	ID int64
	Name string
	Description string
	Completed bool
	Created string
	Tag string
}

// merge the changed fields to the original task
func (orig *Task) merge(t Task) {
	uValues := reflect.ValueOf(&t).Elem()
	oValues := reflect.ValueOf(orig).Elem()
	for i := 0; i < uValues.NumField(); i++ {
		uField := uValues.Field(i).Interface()
		if oValues.CanSet() {
			if v, ok := uField.(int64); ok && uField != 0 {
				oValues.Field(i).SetInt(v)
			}
			if v, ok := uField.(string); ok && uField != "" {
				oValues.Field(i).SetString(v)
			}
			if v, ok := uField.(bool); ok {
				oValues.Field(i).SetBool(v)
			}
		}
	}
}

func createTask(name string, description string, completed bool, tag string) Task {
	var task Task
	task.Name = name
	task.Description = description
	task.Completed = completed
	task.Tag = tag
	return task
}

func main() {
	// delete previous database
	os.Remove(db_fname)

	// start the database
	db, err := sql.Open("sqlite3", db_fname)
	handleErr(err)
	defer db.Close()

	// this is here because the linter deletes the import
	var _, _, _ = sqlite3.Version()

	// create the table for tasks
	createDB(db)

	//TODO remove and use for testing
	var task1 = createTask("test1", "test1", false, "NULL")
	task1.ID = addTask(db, task1)
	id2 := addTask(db, createTask("test2", "test2", false, "NULL"))
	_ = addTask(db, createTask("test3", "test3", true, "NULL"))
	_ = addTask(db, createTask("test4", "test4", true, "NULL"))
	delTask(db, id2)
	_ = editTask(db, task1)
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
}

func createDB(db *sql.DB) {
	sqlStatement := `
		CREATE TABLE "tasks" (
			"id" INTEGER NOT NULL PRIMARY KEY,
			"name" TEXT NOT NULL,
			"description" TEXT,
			"completed" INTEGER,
			"created" TEXT,
			"tag" TEXT
		);
		DELETE FROM tasks;
	`
	var _, err = db.Exec(sqlStatement)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStatement)
		return
	}
}

func addTask(db *sql.DB, task Task) int64 {
	sqlStatement := `
		INSERT INTO 
			tasks(id, name, description, completed, tag, created) 
			values ((SELECT MAX(id) FROM tasks LIMIT 1) + 1, ?, ?, ?, ?, ?);`
	res, err := db.Exec(sqlStatement, task.Name, task.Description, task.Completed, task.Tag, time.Now().Format("2006-01-02_15:04:05"))
	handleErr(err)
	id, err := res.LastInsertId()
	handleErr(err)
	return id
}

func delTask(db *sql.DB, id int64) {
	sqlStatement := `DELETE FROM tasks WHERE id = ?;`
	_, err := db.Exec(sqlStatement, id)
	handleErr(err)
}

func editTask(db *sql.DB, task Task) int64 {
	// get existing task
	var orig = getTask(db, task.ID)
	orig.merge(task)

	// update task
	updateStatement := `
		UPDATE tasks
		SET 
			name = ?,
			description = ?,
			completed = ?,
			tag = ?,
			created = ?
		WHERE id = ?;`
	res, err := db.Exec(updateStatement, orig.Name, orig.Description, orig.Completed, orig.Tag, orig.Created, orig.ID)
	handleErr(err)
	numRows, err := res.RowsAffected()
	handleErr(err)
	return numRows
}

func row2Task(rows *sql.Rows) Task {
	var task Task
	var err = rows.Scan(&task.ID, &task.Name, &task.Description, &task.Completed, &task.Tag, &task.Created)
	handleErr(err)
	return task
}

func getTask(db *sql.DB, id int64) Task {
	var row, err = db.Query(`
		SELECT id, name, description, completed, tag, created 
		FROM tasks WHERE id = ?
		LIMIT 1
	`, id)
	handleErr(err)
	row.Next()
	defer row.Close()
	return row2Task(row)
}

func getTasks(db *sql.DB) []Task {
	// get tasks
	rows, err := db.Query(`
		SELECT 
			id, name, description, completed, tag, created
		FROM 
			tasks;
	`)
	handleErr(err)
	defer rows.Close()

	var tasks = []Task{}

	// print tasks
	for rows.Next() {
		var task = row2Task(rows)
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