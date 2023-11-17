package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/kancli"
	"github.com/charmbracelet/lipgloss"
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

var addCmd = &cobra.Command {
	Use: "add NAME",
	Short: "Add a new task with an optional description and tag",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db := createDB()
		defer db.Close()
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
		return addTask(db, args[0], description, todo, tag)
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

		newTask := Task{int64(id), name, description, status, time.Now(), tag}
		return editTask(db, newTask)
	},
}

var delCmd = &cobra.Command {
	Use: "del ID",
	Short: "Delete a task by its ID",
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		db := createDB()
		defer db.Close()		
		var id, err = strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		return delTask(db, int64(id))
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all your tasks",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		db := createDB()
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
		{Title: "Name", Width: calculateWidth(LG, w)},
		{Title: "Tag", Width: calculateWidth(MD, w)},
		{Title: "Status", Width: calculateWidth(SM, w)},
		{Title: "Created At", Width: calculateWidth(MD, w)},
	}
	var rows []table.Row
	for _, task := range tasks {
		rows = append(rows, table.Row{
			fmt.Sprintf("%d", task.ID),
			task.Name,
			task.Tag,
			task.Status.String(),
			task.Created.Format("2 Jan 2006, 15:04:05"),
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

var kanbanCmd = &cobra.Command{
	Use:   "kanban",
	Short: "Interact with your tasks in a Kanban board.",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		db := createDB()
		defer db.Close()
		todos, err := getTasksByStatus(db, todo)
		if err != nil {
			return err
		}
		ipr, err := getTasksByStatus(db, inProgress)
		if err != nil {
			return err
		}
		finished, err := getTasksByStatus(db, done)
		if err != nil {
			return err
		}

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
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(listCmd)
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
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(delCmd)
	rootCmd.AddCommand(kanbanCmd)
}