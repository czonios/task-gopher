package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	_ "github.com/mattn/go-sqlite3"
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

// Stringer for Task
func (t Task) String() string {
	var completed string
	if t.Completed {
		completed = "x"
	} else {
		completed = " "
	}
	var parsedTime, err = time.Parse("2006-01-02_15:04:05", t.Created)
	handleErr(err)
	var timestr = parsedTime.Format(time.RFC1123)
	return fmt.Sprintf("- [%v] %v\n\tTag: %v\n\tDescription: %v\n\tCreated at: %v\n\tID: %v", completed, t.Name, t.Tag, t.Description, timestr, t.ID)
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

// func createTask(db *sql.DB, name string, description string, completed bool, tag string) Task {
// 	var task Task
// 	task.Name = name
// 	task.Description = description
// 	task.Completed = completed
// 	task.Tag = tag
// 	task.ID = addTask(db, task)
// 	return task
// }

func main() {
	// this is here because the linter deletes the import
	// var _, _, _ = sqlite3.Version()

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//TODO remove and use for testing
	// create the table for tasks
	var db = createDB()
	defer db.Close()
	// var task1 = createTask(db, "test1", "test1", false, "")
	// var task2 = createTask(db, "test2", "test2", false, "")
	// _ = createTask(db, "test3", "test3", true, "")
	// _ = createTask(db, "test4", "test4", true, "")
	// delTask(db, task2.ID)
	// _ = editTask(db, task1)
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

func createDB() *sql.DB {
	//! delete previous database
	// os.Remove(db_fname)

	// start the database
	var db, err = sql.Open("sqlite3", db_fname)
	handleErr(err)

	if _, err := db.Query("SELECT * FROM tasks"); err == nil {
		fmt.Println("Found tasks DB!")
	} else {
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
		_, err = db.Exec(sqlStatement)
		handleErr(err)
	}
	return db
}

func addTask(db *sql.DB, name string, description string, completed bool, tag string) error {
	sqlStatement := `
		INSERT INTO 
			tasks(id, name, description, completed, tag, created) 
			values ((SELECT MAX(id) FROM tasks LIMIT 1) + 1, ?, ?, ?, ?, ?);`
	_, err := db.Exec(sqlStatement, name, description, completed, tag, time.Now().Format("2006-01-02_15:04:05"))
	// handleErr(err)
	// id, err := res.LastInsertId()
	// handleErr(err)
	// return id
	return err
}

func delTask(db *sql.DB, id int64) error {
	sqlStatement := `DELETE FROM tasks WHERE id = ?;`
	_, err := db.Exec(sqlStatement, id)
	handleErr(err)
	return err
}

func editTask(db *sql.DB, task Task) error {
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
	_, err = res.RowsAffected()
	return err
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