package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/kancli"
	"github.com/charmbracelet/lipgloss"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var rootCmd = &cobra.Command {
	Use:   "task-gopher",
	Short: "A CLI task management tool for ~slaying~ your to do list.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var serveCmd = &cobra.Command{
	Use:   "serve",
	Aliases: []string{"server", "start"},
	Short: "create and start a server for the DB",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		
		// Find .env file
		err := godotenv.Load(projectDir + "/.env")
		if err != nil{
			log.Fatalf("Error loading .env file: %s", err)
		}
		port := os.Getenv("PORT")
		serve(port)
	},
}

var addCmd = &cobra.Command {
	Use: "add NAME",
	Short: "Add a new task with an optional description and tag",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// db := createDB()
		// defer db.Close()
		var description, tag string
		var err error
		description, err = cmd.Flags().GetString("description")
		if err != nil {
			return err
		}
		tag, err = cmd.Flags().GetString("tag")
		if err != nil {
			return err
		}
		// JSON body
		body := []byte(fmt.Sprintf(`{
			"Name": "%v",
			"Desc": "%v",
			"Status": "%v",
			"Tag": "%v"
		}`, args[0], description, todo, tag))
		
		addr := os.Getenv("ADDRESS")
		port := os.Getenv("PORT")
		url := addr + ":" + port + "/tasks/add"
		// fmt.Println(url)
		// fmt.Println(string(body))
		// fmt.Println(http.DetectContentType(body))
		res, err := http.Post(url, "application/json; charset=utf-8", bytes.NewBuffer(body))
		if err != nil {
			return err 
		}
		fmt.Println(res.Status)
		jsonBody := make(map[string]interface{})
		err = json.NewDecoder(res.Body).Decode(&jsonBody)
		if err != nil {
			return err 
		}
		newTask, err := json.Marshal(jsonBody)
		fmt.Println(string(newTask))
		return err
	},
}

var updateCmd = &cobra.Command {
	Use: "update ID",
	Short: "Update an existing task name, description, tag or completion status by its id",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db := createDB()
		defer db.Close()

		var id, err = strconv.Atoi(args[0])
		if err != nil {
			return err
		}

		name, err := cmd.Flags().GetString("name")
		if err != nil {
			return err
		}

		description, err := cmd.Flags().GetString("description")
		if err != nil {
			return err
		}
		tag, err := cmd.Flags().GetString("tag")
		if err != nil {
			return err
		}
		completed, err := cmd.Flags().GetInt("status")
		if err != nil {
			return err
		}

		var status status
		switch completed {
		case int(inProgress):
			status = inProgress
		case int(done):
			status = done
		default:
			status = todo
		}

		// JSON body
		body := []byte(fmt.Sprintf(`{
			"Name": "%v",
			"Desc": "%v",
			"Status": "%v",
			"Tag": "%v"
		}`, name, description, status, tag))
		
		addr := os.Getenv("ADDRESS")
		port := os.Getenv("PORT")
		url := addr + ":" + port + "/tasks/" + fmt.Sprint(id)

		// create a new PUT request
		req, err := http.NewRequest("PUT", url, bytes.NewBuffer(body))
		if err != nil {
			return err
		}
	
		// send the request
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		fmt.Println(res.Status)
		if res.StatusCode == 200 {
			fmt.Println("Updated task", id)
		} else {
			fmt.Println("Something went wrong")
		}
		return nil
	},
}

var delCmd = &cobra.Command {
	Use: "del ID",
	Short: "Delete a task by its ID",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// db := createDB()
		// defer db.Close()		
		var id, err = strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		addr := os.Getenv("ADDRESS")
		port := os.Getenv("PORT")
		url := addr + ":" + port + "/tasks/" + fmt.Sprint(id)
		
		// create a new HTTP client
		client := &http.Client{}

		// create a new DELETE request
		req, err := http.NewRequest("DELETE", url, nil)
		if err != nil {
			return err
		}
	
		// send the request
		resp, err := client.Do(req)
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		fmt.Println("Deleted task", id)
		return nil
	},
}

func getTasksFromServer() ([]Task, error){
	addr := os.Getenv("ADDRESS")
	port := os.Getenv("PORT")
	url := addr + ":" + port + "/tasks"
	
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	// decode response into tasks array
	var tasks []Task
	err = json.NewDecoder(resp.Body).Decode(&tasks)
	if err != nil {
		return nil, err
	}
	return tasks, nil
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all your tasks",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		tasks, err := getTasksFromServer()
		if err != nil {
			return err
		}
		table := setupTable(tasks)
		fmt.Print(table.View())
		return nil
	},
}

var dropDBCmd = &cobra.Command{
	Use:   "deldb",
	Short: "delete all your tasks",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		db := createDB(true)
		defer db.Close()
		tasks, err := getTasks(db)
		if err != nil {
			return err
		}
		table := setupTable(tasks)
		fmt.Print(table.View())
		return nil
	},
}

func calculateWidth(min, width int) int {
	p := width / 10
	switch min {
	case XS:
		if p < XS {
			return XS
		}
		return p / 2

	case SM:
		if p < SM {
			return SM
		}
		return p / 2
	case MD:
		if p < MD {
			return MD
		}
		return p * 2
	case LG:
		if p < LG {
			return LG
		}
		return p * 3
	default:
		return p
	}
}

const (
	XS int = 1
	SM int = 3
	MD int = 5
	LG int = 10
)

func setupTable(tasks []Task) table.Model {
	// get term size
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		// we don't really want to fail it...
		log.Println("unable to calculate height and width of terminal")
	}

	columns := []table.Column{
		{Title: "ID", Width: calculateWidth(XS, w)},
		{Title: "Name", Width: calculateWidth(MD, w)},
		{Title: "Tag", Width: calculateWidth(SM, w)},
		{Title: "Status", Width: calculateWidth(MD, w)},
		{Title: "Description", Width: calculateWidth(MD, w)},
		{Title: "Created At", Width: calculateWidth(MD, w)},
	}
	var rows []table.Row
	for _, task := range tasks {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", task.ID),
			task.Name,
			task.Tag,
			task.Status.String(),
			task.Desc,
			task.Created.Format("2 Jan 2006"),
		})
	}
	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(false),
		table.WithHeight(len(tasks)),
	)
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	t.SetStyles(s)
	return t
}

func filterTasksByStatus(tasks []Task, s status) []Task {
	var filtered []Task
	for _, task := range tasks {
		if task.Status == s {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

var kanbanCmd = &cobra.Command{
	Use:   "kanban",
	Short: "Interact with your tasks in a Kanban board",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {

		tasks, err := getTasksFromServer()
		if err != nil {
			return err
		}

		todos := filterTasksByStatus(tasks, todo)
		ipr := filterTasksByStatus(tasks, inProgress)
		finished := filterTasksByStatus(tasks, done)

		todoCol := kancli.NewColumn(tasksToItems(todos), todo, true)
		iprCol := kancli.NewColumn(tasksToItems(ipr), inProgress, false)
		doneCol := kancli.NewColumn(tasksToItems(finished), done, false)
		board := kancli.NewDefaultBoard([]kancli.Column{todoCol, iprCol, doneCol})
		p := tea.NewProgram(board)
		_, err = p.Run()
		return err
	},
}

// convert tasks to items for a list
func tasksToItems(tasks []Task) []list.Item {
	var items []list.Item
	for _, t := range tasks {
		items = append(items, t)
	}
	return items
}

func init() {
	// add cmd flags
	addCmd.Flags().StringP(
		"tag",
		"t",
		"",
		"specify a tag for your task",
	)
	addCmd.Flags().StringP(
		"description",
		"d",
		"",
		"specify a description for your task",
	)
	// update cmd flags
	updateCmd.Flags().StringP(
		"name",
		"n",
		"",
		"specify a name for your task",
	)
	updateCmd.Flags().StringP(
		"tag",
		"t",
		"",
		"specify a tag for your task",
	)
	updateCmd.Flags().StringP(
		"description",
		"d",
		"",
		"specify a description for your task",
	)
	updateCmd.Flags().IntP(
		"status",
		"s",
		int(todo),
		"specify a completion status for your task (0/1/2 for todo/in progress/done)",
	)
	// add all commands
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(delCmd)
	rootCmd.AddCommand(kanbanCmd)
	rootCmd.AddCommand(dropDBCmd)
	rootCmd.AddCommand(serveCmd)
}