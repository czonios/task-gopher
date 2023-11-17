/*
A CLI task management tool with remote access capabilities.

Usage:

	task-gopher [flags]
	task-gopher [command]

Available Commands:

	add         Add a new task with an optional description and tag
	completion  Generate the autocompletion script for the specified shell
	del         Delete a task by its ID
	deldb       delete all your tasks
	help        Help about any command
	kanban      Interact with your tasks in a Kanban board
	list        List all your tasks
	serve       create and start a server for the DB
	update      Update an existing task name, description, tag or completion status by its id

Flags:

	-h, --help   help for task-gopher

Use "task-gopher [command] --help" for more information about a command.
*/
package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/mattn/go-sqlite3"
)

const dbFname = "tasks.db"

var homeDir, _ = os.UserHomeDir()
var projectDir = homeDir + "/go/src/github.com/czonios/task-gopher"
var _ = os.Mkdir(projectDir, os.ModePerm)

type status int

// status enum
const (
	todo status = iota
	inProgress
	done
)

func (s status) String() string {
	return [...]string{"todo", "in progress", "done"}[s]
}

// A Task is the representation of a task
type Task struct {
	ID      int64     // unique task ID
	Name    string    // task title
	Desc    string    // optional description
	Status  status    // the status, one of {todo, in progress, done}
	Created time.Time // timestamp of when the task was created
	Tag     string    // optional tag for the task
}

// implement list.Item & list.DefaultItem
func (t Task) FilterValue() string {
	return t.Name
}

func (t Task) Title() string {
	return t.Name
}

func (t Task) Description() string {
	return t.Tag
}

// implement kancli.Status
func (s status) Next() int {
	if s == done {
		return int(todo)
	}
	return int(s + 1)
}

func (s status) Prev() int {
	if s == todo {
		return int(done)
	}
	return int(s - 1)
}

func (s status) Int() int {
	return int(s)
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
			if v, ok := uField.(status); ok && uField != todo {
				oValues.Field(i).SetInt(int64(v))
			}
		}
	}
}

func main() {
	// Find .env file
	err := godotenv.Load(projectDir + "/.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	// execute root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

// createDB returns an opened SQLite database that can be used to run queries
// It creates the directory and the db file, if they don't exist
func createDB(args ...bool) *sql.DB {
	var dataDir = projectDir + "/data/"
	var dbPath = dataDir + dbFname
	_ = os.Mkdir(dataDir, os.ModePerm)

	if len(args) > 0 {
		delete := args[0]
		if delete {
			os.Remove(dbPath)
		}
	}

	// start the database
	var db, err = sql.Open("sqlite3", dbPath)
	handleErr(err)

	if _, err := db.Query("SELECT * FROM tasks"); err == nil {
		// fmt.Println("Found tasks DB!")
	} else {
		sqlStatement := `
            CREATE TABLE "tasks" (
                "id" INTEGER NOT NULL PRIMARY KEY,
                "name" TEXT NOT NULL,
                "description" TEXT,
                "status" INTEGER,
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

// addTask inserts a task into the database
func addTask(db *sql.DB, name string, description string, completed status, tag string) (int64, error) {
	sqlStatement := `
        INSERT INTO 
            tasks(id, name, description, status, tag, created) 
            values ((SELECT MAX(id) FROM tasks LIMIT 1) + 1, ?, ?, ?, ?, ?);`
	res, err := db.Exec(sqlStatement, name, description, completed, tag, time.Now().Format(time.RFC3339))
	handleErr(err)
	id, err := res.LastInsertId()
	handleErr(err)
	return id, err
}

// delTask deletes a task from the database
func delTask(db *sql.DB, id int64) error {
	sqlStatement := `DELETE FROM tasks WHERE id = ?;`
	_, err := db.Exec(sqlStatement, id)
	handleErr(err)
	return err
}

// editTask updates an existing task in the database
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
            status = ?,
            tag = ?,
            created = ?
        WHERE id = ?;`
	res, err := db.Exec(updateStatement, orig.Name, orig.Desc, orig.Status, orig.Tag, orig.Created.Format(time.RFC3339), orig.ID)
	handleErr(err)
	_, err = res.RowsAffected()
	return err
}

// row2Task returns a task scanned from a database row
func row2Task(rows *sql.Rows) Task {
	var task Task
	var timestr string
	var err = rows.Scan(&task.ID, &task.Name, &task.Desc, &task.Status, &task.Tag, &timestr)
	handleErr(err)
	task.Created, err = time.Parse(time.RFC3339, timestr)
	handleErr(err)
	return task
}

// getTask returns the task with a given id
func getTask(db *sql.DB, id int64) Task {
	var row, err = db.Query(`
        SELECT id, name, description, status, tag, created 
        FROM tasks WHERE id = ?
        LIMIT 1
    `, id)
	handleErr(err)
	row.Next()
	defer row.Close()
	return row2Task(row)
}

// getTasks returns all the tasks in the database
func getTasks(db *sql.DB) ([]Task, error) {
	// get tasks
	rows, err := db.Query(`
        SELECT id, name, description, status, tag, created
        FROM tasks
        ORDER BY created ASC;
    `)
	handleErr(err)
	defer rows.Close()

	var tasks = []Task{}

	// print tasks
	for rows.Next() {
		var task = row2Task(rows)
		tasks = append(tasks, task)
	}

	err = rows.Err()
	return tasks, err
}

// handleErr logs a Fatal error if given a non-nil error
func handleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
