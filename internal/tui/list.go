package tui

import (
	"github.com/boejennet/godo/internal/todos"
	"github.com/charmbracelet/lipgloss"
)

var subtle = lipgloss.AdaptiveColor{Light: "#D9DCCF", Dark: "#383838"}
var special = lipgloss.AdaptiveColor{Light: "#43BF6D", Dark: "#73F59F"}
var listItem = lipgloss.NewStyle().PaddingLeft(2).Render

var listHeader = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderBottom(true).
	BorderForeground(subtle).
	MarginRight(2).
	Render

var checkMark = lipgloss.NewStyle().SetString("✓").
	Foreground(special).
	PaddingRight(1).
	String()

func listDone(s string) string {
	return checkMark + lipgloss.NewStyle().
		Strikethrough(true).
		Foreground(lipgloss.AdaptiveColor{Light: "#969B86", Dark: "#696969"}).
		Render(s)
}

func todoItem(todo todos.Todo) string {
	if todo.Completed {
		return listDone(todo.Name)
	}
	return listItem(todo.Name)
}

func list(width int, height int) lipgloss.Style {
	return lipgloss.NewStyle().
		BorderForeground(subtle).
		MarginRight(2).
		Height(height).
		Width((width / 2) + 1)
}

func TodoList(cursor int, todos []todos.Todo, width int, height int) string {
	todoList := []string{}

	for i := range todos {
		if i == cursor {
			todoList = append(todoList, lipgloss.NewStyle().Foreground(special).Render(todoItem(todos[i]))+" <-")
		} else {
			todoList = append(todoList, (todoItem(todos[i])))
		}
	}

	return list(width, height).Render(
		lipgloss.JoinVertical(
			lipgloss.Left,
			append([]string{listHeader("Todos")}, todoList...)...,
		),
	)
}
