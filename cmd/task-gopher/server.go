package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

var db *sql.DB

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

	// check which routes are up
	// for _, route := range e.Routes() {
	// 	fmt.Println(route.Path)
	// }
	
	// start on port
	e.Logger.Fatal(e.Start(":" + port))
}

func getJSONRawBody(c echo.Context) (map[string]interface{}, error) {

	jsonBody := make(map[string]interface{})
	err := json.NewDecoder(c.Request().Body).Decode(&jsonBody)
	if err != nil {
		return nil, err 
	}

   return jsonBody, nil
}

func handleGetTasks(c echo.Context) error {
    tasks, err := getTasks(db)
		if err != nil {
			return err
	}
    return c.JSON(http.StatusOK, tasks)
}

func handleDeleteTask(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid task id")
	}
	err = delTask(db, int64(id))
	if err != nil {
		return c.String(http.StatusInternalServerError, "Could not delete task " + fmt.Sprint(id))
	}
    return c.String(http.StatusOK, fmt.Sprint(id))
}

func handleAddTask(c echo.Context) error {

	// fmt.Println(c.Request().Body)
	// s, _ := io.ReadAll(c.Request().Body)
	// fmt.Println(string(s))
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

func handleUpdateTask(c echo.Context) error {

	// fmt.Println(c.Request().Body)
	// s, _ := io.ReadAll(c.Request().Body)
	// fmt.Println(string(s))
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