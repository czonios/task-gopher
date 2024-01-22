package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var db *sql.DB
var (
	upgrader = websocket.Upgrader{}
)
var clients = make(map[*websocket.Conn]bool)

// serve starts an echo server
// It opens the SQLite database and sets up the accepted routes
func serve(port string) {
	// create or open the database
	db = createDB()
	defer db.Close()

	// create the server
	e := echo.New()
	// e.Pre(middleware.HTTPSRedirect())

	// set up middleware
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogRemoteIP: true,
		LogHost:     true,
		LogMethod:   true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			log.Printf("%v %v from %v. Status: %v\n", v.Method, v.URI, v.RemoteIP, v.Status)
			return nil
		},
	}))
	e.Use(middleware.Recover())
	e.Use(middleware.Secure())

	// set up routes
	e.GET("/tasks", handleGetTasks)
	e.GET("/tasks/:id", handleGetTask)
	e.POST("/tasks/add", handleAddTask)
	e.PUT("/tasks/:id", handleUpdateTask)
	e.DELETE("/tasks/:id", handleDeleteTask)
	e.GET("/ws", handleWebsocket)

	// Goroutine for checking new day start
	go checkDayStart(db)

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

// handleWebsocket handles the WebSocket connection.
func handleWebsocket(c echo.Context) error {
	ws, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return err
	}
	// remove connection from list of connections and close it when done
	defer func() {
		delete(clients, ws)
		ws.Close()
	}()

	err = manageConnections(ws) // keep list of connections
	if err != nil {
		c.Logger().Error(err)
	}

	// Write hello message
	err = ws.WriteMessage(websocket.TextMessage, []byte("Websocket connected!"))
	if err != nil {
		c.Logger().Error(err)
	}

	for {
		// Read message
		_, msg, err := ws.ReadMessage()
		if err != nil {
			c.Logger().Error(err)
		}
		err = handleMessage(ws, msg) // handle message
		if err != nil {
			c.Logger().Error(err)
		}
	}
}

// handleMessage handles an incoming (through websocket) message
func handleMessage(ws *websocket.Conn, msg []byte) error {
	log.Printf("Received message from socket: %s\n", msg)
	//TODO implement message handling logic
	// should notify all clients that a change was made?
	return nil
}

// manageConnections manages the list of websocket connections
func manageConnections(ws *websocket.Conn) error {
	log.Println("WS connection from", ws.RemoteAddr().String())
	clients[ws] = true
	return nil
}

// handleGetTasks fetches all tasks from the database and returns them in JSON form in the response
func handleGetTasks(c echo.Context) error {
	// log.Println(c.Request().RemoteAddr+":", c.Request().Method, c.Request().RequestURI)
	tasks, err := getTasks(db)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, tasks)
}

// handleGetTasks fetches a task from the database by ID and returns it in JSON form in the response
func handleGetTask(c echo.Context) error {
	// log.Println(c.Request().RemoteAddr+":", c.Request().Method, c.Request().RequestURI)
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return c.String(http.StatusBadRequest, "Invalid task id")
	}
	tasks, err := getTask(db, id)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Could not fetch task "+fmt.Sprint(id))
	}
	return c.JSON(http.StatusOK, tasks)
}

// handleDeleteTask deletes a task from the database and returns its id
func handleDeleteTask(c echo.Context) error {
	// log.Println(c.Request().RemoteAddr+":", c.Request().Method, c.Request().RequestURI)
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
	// log.Println(c.Request().RemoteAddr+":", c.Request().Method, c.Request().RequestURI)
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
	type_s := body["Type"].(string)

	var status status
	switch completed {
	case todo.String():
		status = todo
	case inProgress.String():
		status = inProgress
	case done.String():
		status = done
	default:
		status = invalidStatus
	}

	var type_t task_type
	switch type_s {
	case generic.String():
		type_t = generic
	case daily.String():
		type_t = daily
	case habit.String():
		type_t = habit
	default:
		type_t = invalidType
	}

	if name == "" {
		return c.String(http.StatusBadRequest, "You must provide a task name")
	}

	// create task
	id, err := addTask(db, name, desc, status, type_t, tag)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Could not create task")
	}

	// get the task
	task, err := getTask(db, id)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Could not create task")
	}

	return c.JSON(http.StatusOK, task)
}

// handleUpdateTask updates a task in the database by its id
// The id is given in the request parameters and the changed values are given in the request body
func handleUpdateTask(c echo.Context) error {
	// log.Println(c.Request().RemoteAddr+":", c.Request().Method, c.Request().RequestURI)
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
	type_s := body["Type"].(string)

	var status status
	switch completed {
	case todo.String():
		status = todo
	case inProgress.String():
		status = inProgress
	case done.String():
		status = done
	default:
		status = invalidStatus
	}

	var type_t task_type
	switch type_s {
	case generic.String():
		type_t = generic
	case daily.String():
		type_t = daily
	case habit.String():
		type_t = habit
	default:
		type_t = invalidType
	}

	// update task
	newTask := Task{int64(id), name, desc, status, type_t, time.Now(), tag}
	err = editTask(db, newTask)
	if err != nil {
		return c.String(http.StatusInternalServerError, "Could not update task")
	}
	return c.NoContent(http.StatusOK)
}
