package main

import (
	"strconv"

	"github.com/spf13/cobra"
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
		return addTask(db, args[0], description, false, tag)
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

		var description, tag string
		var completed bool
		description, err = cmd.Flags().GetString("description")
		if err != nil {
			return err
		}
		tag, err = cmd.Flags().GetString("tag")
		if err != nil {
			return err
		}
		completed, err = cmd.Flags().GetBool("completed")
		if err != nil {
			return err
		}
		var task Task
		task.ID = int64(id)
		task.Tag = tag
		task.Completed = completed
		task.Description = description
		return editTask(db, task)
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
	// rootCmd.AddCommand(listCmd)
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
	updateCmd.Flags().BoolP(
		"completed",
		"c",
		false,
		"specify a completion status for your task (0/1 or true/false)",
	)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(delCmd)
	// rootCmd.AddCommand(kanbanCmd)
}