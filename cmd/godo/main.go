package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/boejennet/godo/internal/data"
	todostype "github.com/boejennet/godo/internal/todos"
	"github.com/boejennet/godo/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "godo",
	Version: "1.0",
	Short:   "Go Do is a simple todo app for the CLI",
}

func NewListCmd(tdb *data.TodoDB) *cobra.Command {
	return &cobra.Command{
		Use:   "list",
		Short: "Lists todos",
		Example: heredoc.Doc(`
			$ godo list
		`),
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			completed, err := cmd.Flags().GetBool("completed")
			if err != nil {
				return err
			}

			var todos []todostype.Todo
			if completed {
				todos = tdb.GetCompleted()
			} else {
				todos = tdb.GetAllTodos()
			}

			for _, td := range todos {
				if completed {
					fmt.Printf("ID: %d, Name: %s, Completed: %t, Completed at: %s\n", td.ID, td.Name, td.Completed, td.CompletedAt.Format("Mon Jan 2 15:04"))
				} else {
					fmt.Printf("ID: %d, Name: %s, Completed: %t\n", td.ID, td.Name, td.Completed)

				}
			}

			return nil
		},
	}
}

func NewUpdateCmd(tdb *data.TodoDB) *cobra.Command {
	return &cobra.Command{
		Use:   "update <id>",
		Short: "Update a todo",
		Example: heredoc.Doc(`
			$ godo update 1 --completed true
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.Atoi(args[0])
			if err != nil {
				return err
			}

			taskToUpdate := tdb.GetByID(int64(id))
			if taskToUpdate.ID == 0 {
				fmt.Println("Todo not found")
				return nil
			}

			complete, err := cmd.Flags().GetBool("completed")
			if err != nil {
				return err
			}

			if taskToUpdate.Completed && complete {
				fmt.Println("Todo already completed")
				return nil
			}
			if complete {
				fmt.Println("Marking todo as completed")
				taskToUpdate.Completed = true
				taskToUpdate.UpdatedAt = time.Now()
				taskToUpdate.CompletedAt = time.Now()
			}

			name, err := cmd.Flags().GetString("name")
			if err != nil {
				return err
			}

			if name != "" {
				taskToUpdate.Name = name
				taskToUpdate.UpdatedAt = time.Now()
			}

			tdb.Update(taskToUpdate)
			fmt.Println("Todo updated:", id)

			return nil
		},
	}
}

func NewAddCmd(tdb *data.TodoDB) *cobra.Command {
	return &cobra.Command{
		Use:   "add <todo>",
		Short: "Add a todo",
		Example: heredoc.Doc(`
			$ godo add "Buy milk" "Buy Eggs"
		`),
		Args: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			for _, arg := range args {
				if arg == "" {
					fmt.Println("Please provide a todo")
				}

				todo := todostype.NewTodo(arg)
				tdb.Insert(todo)
				fmt.Println("Todo added:", todo.Name)
			}
		},
	}
}

func NewTUICmd(tdb *data.TodoDB) *cobra.Command {
	return &cobra.Command{
		Use:   "tui",
		Short: "Start the TUI",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			todos := tdb.GetAllTodos()
			tuiModel := tui.NewModel(tdb, todos)

			p := tea.NewProgram(tuiModel, tea.WithAltScreen())
			if _, err := p.Run(); err != nil {
				fmt.Printf("error starting tui: %v", err)
				os.Exit(1)
			}
		},
	}
}

var whereCmd = &cobra.Command{
	Use:   "where",
	Short: "Show where your todo db is stored",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		_, err := fmt.Println(data.SetupPath())
		return err
	},
}

func main() {
	tdb := data.NewTDB()
	listCmd := NewListCmd(tdb)
	listCmd.Flags().BoolP("completed", "c", false, "Show completed todos")

	updateCmd := NewUpdateCmd(tdb)
	updateCmd.Flags().BoolP("completed", "c", false, "Set the todo as completed")
	updateCmd.Flags().StringP("name", "n", "", "Set the name of the todo")

	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(NewAddCmd(tdb))
	rootCmd.AddCommand(whereCmd)
	rootCmd.AddCommand(NewTUICmd(tdb))

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
