package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

var db *sql.DB

// serve starts an echo server
// It opens the SQLite database and sets up the accepted routes
func serve(port string) {
	// create or open the database
	db = createDB()
	defer db.Close()

	// create the server
	e := echo.New()

	// set up routes
	e.GET("/tasks", handleGetTasks)
	e.POST("/tasks/add", handleAddTask)
	e.PUT("/tasks/:id", handleUpdateTask)
	e.DELETE("/tasks/:id", handleDeleteTask)

	// start on port
	e.Logger.Fatal(e.Start(":" + port))
}

// getJSONRawBody returns the body of a request c in JSON format
func getJSONRawBody(c echo.Context) (map[string]interface{}, error) {

	jsonBody := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		return nil, err
	}

	return jsonBody, nil
}

// handleGetTasks fetches all tasks from the database and returns them in JSON form in the response
func handleGetTasks(c echo.Context) error {
	log.Println(c.Request().RemoteAddr+":", c.Request().Method, c.Request().RequestURI)
	tasks, err := getTasks(db)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, tasks)
}

// handleDeleteTask deletes a task from the database and returns its id
func handleDeleteTask(c echo.Context) error {
	log.Println(c.Request().RemoteAddr+":", c.Request().Method, c.Request().RequestURI)
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid task id")
	}
	err = delTask(db, int64(id))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Could not delete task "+fmt.Sprint(id))
	}
	return c.String(http.StatusOK, fmt.Sprint(id))
}

// handleAddTask adds a task to the database
// It gets the task data from the request body in JSON form
func handleAddTask(c echo.Context) error {
	log.Println(c.Request().RemoteAddr+":", c.Request().Method, c.Request().RequestURI)
	body, err := getJSONRawBody(c)

	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusBadRequest, "You must provide a request body")
	}

	// get task details from JSON
	name := body["Name"].(string)
	desc := body["Desc"].(string)
	tag := body["Tag"].(string)
	completed := body["Status"].(string)
	var status status
	switch completed {
	case inProgress.String():
		status = inProgress
	case done.String():
		status = done
	default:
		status = todo
	}

	if name == "" {
		return c.String(http.StatusBadRequest, "You must provide a task name")
	}

	// create task
	id, err := addTask(db, name, desc, status, tag)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Could not create task")
	}

	// get the task
	task := getTask(db, id)

	return c.JSON(http.StatusOK, task)
}

// handleUpdateTask updates a task in the database by its id
// The id is given in the request parameters and the changed values are given in the request body
func handleUpdateTask(c echo.Context) error {
	log.Println(c.Request().RemoteAddr+":", c.Request().Method, c.Request().RequestURI)
	body, err := getJSONRawBody(c)

	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusBadRequest, "You must provide a request body")
	}
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		fmt.Println(err)
		return c.String(http.StatusBadRequest, "Invalid task id")
	}

	// get task details from JSON
	name := body["Name"].(string)
	desc := body["Desc"].(string)
	tag := body["Tag"].(string)
	completed := body["Status"].(string)
	fmt.Println(completed)

	var status status
	switch completed {
	case inProgress.String():
		status = inProgress
	case done.String():
		status = done
	default:
		status = todo
	}

	// update task
	newTask := Task{int64(id), name, desc, status, time.Now(), tag}
	err = editTask(db, newTask)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Could not update task")
	}
	return c.NoContent(http.StatusOK)
}
