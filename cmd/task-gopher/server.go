package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

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
	// e.PUT("/tasks/:id", handleUpdateTask)
	// e.DELETE("/tasks/:id", handleDeleteTask)

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

// func handleGetTasks(c echo.Context) error {
//     // Task ID from path `users/:id`
//     id := c.Param("id")
//     return c.String(http.StatusOK, id)
// }

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

// /show?team=x-men&member=wolverine

// func show(c echo.Context) error {
//     
// }