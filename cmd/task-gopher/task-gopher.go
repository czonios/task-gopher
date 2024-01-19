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

// var projectDir = "."
var _ = os.Mkdir(projectDir, os.ModePerm)

// status enum
type status int

const (
	todo status = iota
	inProgress
	done
	invalidStatus
)

func (s status) String() string {
	return [...]string{"todo", "in progress", "done", "invalid"}[s]
}

// task type enum
type task_type int

const (
	generic task_type = iota
	daily
	habit
	invalidType
)

func (s task_type) String() string {
	return [...]string{"generic", "daily", "habit", "invalid"}[s]
}

// A Task is the representation of a task
type Task struct {
	ID      int64     // unique task ID
	Name    string    // task title
	Desc    string    // optional description
	Status  status    // the status, one of {todo, in progress, done}
	Type    task_type // the type of the task, one of {generic, daily, habit}
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
			if v, ok := uField.(status); ok && uField != invalidStatus {
				oValues.Field(i).SetInt(int64(v))
			}
			if v, ok := uField.(task_type); ok && uField != invalidType {
				oValues.Field(i).SetInt(int64(v))
			}
		}
	}

}

func main() {
	// Find .env file
	_ = godotenv.Load(projectDir + "/.env")
	port := os.Getenv("PORT")
	addr := os.Getenv("ADDRESS")
	if port == "" || addr == "" {
		log.Fatalf("Environment variables not set!")
	}

	// execute root command
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

// checkDayStart checks if it is a new day (after 6am) and runs the new day logic.
//
// Parameters:
// - db: the database connection (*sql.DB)
//
// Return type: error
func checkDayStart(db *sql.DB) error {
	now := time.Now().UTC()
	prevDay := time.Date(now.Year(), now.Month(), now.Day(), 6, 0, 0, 0, now.Location())
	// runs every 10 seconds
	for range time.Tick(time.Second * 10) {
		// check if it is a new day (after 6am)
		now := time.Now().UTC()
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 6, 0, 0, 0, now.Location())
		if now.After(startOfDay) && startOfDay.After(prevDay) {
			prevDay = startOfDay
			// impl new day logic (daily tasks should reset to todo status)
			err := resetDailyTasks(db)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func resetDailyTasks(db *sql.DB) error {
	// Get all tasks with Type daily
	dailyTasks, err := getTasksByType(db, daily)
	if err != nil {
		return err
	}
	// Reset the status of each daily task to "todo"
	for _, task := range dailyTasks {
		task.Status = todo
		editTask(db, task)
	}
	return nil
}

func getTasksByType(db *sql.DB, taskType task_type) ([]Task, error) {
	// get tasks by type
	rows, err := db.Query(`
        SELECT id, name, description, status, type, tag, created
        FROM tasks
		WHERE type = ?
        ORDER BY created ASC;
    `, taskType)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks = []Task{}

	// print tasks
	for rows.Next() {
		task, err := row2Task(rows)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	err = rows.Err()
	return tasks, err
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
				"type" INTEGER,
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
func addTask(db *sql.DB, name string, description string, completed status, t_type task_type, tag string) (int64, error) {
	sqlStatement := `
        INSERT INTO 
            tasks(id, name, description, status, type, tag, created) 
            values ((SELECT MAX(id) FROM tasks LIMIT 1) + 1, ?, ?, ?, ?, ?, ?);`
	res, err := db.Exec(sqlStatement, name, description, completed, t_type, tag, time.Now().Format(time.RFC3339))
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, err
}

// delTask deletes a task from the database
func delTask(db *sql.DB, id int64) error {
	sqlStatement := `DELETE FROM tasks WHERE id = ?;`
	_, err := db.Exec(sqlStatement, id)
	if err != nil {
		return err
	}
	return err
}

// editTask updates an existing task in the database
func editTask(db *sql.DB, task Task) error {
	// get existing task
	var orig, err = getTask(db, task.ID)

	if err != nil {
		return err
	}
	orig.merge(task)

	// update task
	updateStatement := `
        UPDATE tasks
        SET 
            name = ?,
            description = ?,
            status = ?,
            type = ?,
            tag = ?,
            created = ?
        WHERE id = ?;`
	res, err := db.Exec(updateStatement, orig.Name, orig.Desc, orig.Status, orig.Type, orig.Tag, orig.Created.Format(time.RFC3339), orig.ID)
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	return err
}

// row2Task returns a task scanned from a database row
func row2Task(rows *sql.Rows) (Task, error) {
	var task Task
	var timestr string
	var err = rows.Scan(&task.ID, &task.Name, &task.Desc, &task.Status, &task.Type, &task.Tag, &timestr)
	if err != nil {
		return Task{}, err
	}
	task.Created, err = time.Parse(time.RFC3339, timestr)
	if err != nil {
		return Task{}, err
	}
	return task, nil
}

// getTask returns the task with a given id
func getTask(db *sql.DB, id int64) (Task, error) {
	var row, err = db.Query(`
        SELECT id, name, description, status, type, tag, created 
        FROM tasks WHERE id = ?
        LIMIT 1
    `, id)
	if err != nil {
		return Task{}, err
	}
	row.Next()
	defer row.Close()
	task, err := row2Task(row)
	if err != nil {
		return Task{}, err
	}
	return task, nil
}

// getTasks returns all the tasks in the database
func getTasks(db *sql.DB) ([]Task, error) {
	// get tasks
	rows, err := db.Query(`
        SELECT id, name, description, status, type, tag, created
        FROM tasks
        ORDER BY created ASC;
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks = []Task{}

	// print tasks
	for rows.Next() {
		task, err := row2Task(rows)
		if err != nil {
			return nil, err
		}
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
