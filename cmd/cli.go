package main

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command {
	Use:   "tasks",
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
		createTask(db, args[0], description, false, tag)
		return nil
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
	// updateCmd.Flags().StringP(
	// 	"name",
	// 	"n",
	// 	"",
	// 	"specify a name for your task",
	// )
	// updateCmd.Flags().StringP(
	// 	"project",
	// 	"p",
	// 	"",
	// 	"specify a project for your task",
	// )
	// updateCmd.Flags().IntP(
	// 	"status",
	// 	"s",
	// 	int(todo),
	// 	"specify a status for your task",
	// )
	// rootCmd.AddCommand(updateCmd)
	// rootCmd.AddCommand(deleteCmd)
	// rootCmd.AddCommand(kanbanCmd)
}