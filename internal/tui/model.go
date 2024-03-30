package tui

import (
	"strings"

	"github.com/boejennet/godo/internal/data"
	"github.com/boejennet/godo/internal/todos"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// sessionState is used to track which model is focused
type sessionState uint

const (
	listTodosView sessionState = iota
	newTodoView
)

type TodoModel struct {
	todos     []todos.Todo
	tdb       *data.TodoDB
	textInput textinput.Model
	state     sessionState
	cursor    int
	width     int
	height    int
}

func NewModel(tdb *data.TodoDB, todos []todos.Todo) *TodoModel {
	return &TodoModel{
		tdb:   tdb,
		todos: todos,
		state: listTodosView,
	}
}

func (m TodoModel) Init() tea.Cmd {
	return nil
}

func (m TodoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	// Set the width and height when the window resizes
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	// Is it a key press?
	case tea.KeyMsg:
		switch msg.String() {
		// These keys should exit the program.
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "r":
			// Refresh the list
			m.todos = m.tdb.GetAllTodos()
		case "c":
			// Mark the todo as completed
			m.todos[m.cursor].Completed = !m.todos[m.cursor].Completed
			m.tdb.Update(m.todos[m.cursor])
		case "down", "j":
			if m.cursor < len(m.todos)-1 {
				m.cursor++
			}
		case "ctrl+n":
			// Create a new todo
			if m.state == listTodosView {
				ti := textinput.New()
				ti.Placeholder = "New todo... press enter to save, esc to cancel"
				ti.CharLimit = 256
				ti.Width = m.width
				m.textInput = ti
				m.state = newTodoView
			}
		case "esc":
			// when in the new todo view, cancel
			if m.state == newTodoView {
				m.textInput = textinput.Model{}
				m.state = listTodosView
			}
		case "enter":
			// when in the new todo view, create a new todo
			if m.state == newTodoView {
				m.todos = append(m.todos, todos.NewTodo(m.textInput.Value()))
				m.tdb.Insert(m.todos[len(m.todos)-1])
				m.textInput = textinput.Model{}
				m.state = listTodosView
			}
			// TODO: when in the list view, edit the selected todo name

		}
	}

	if m.state == newTodoView {
		var cmd tea.Cmd
		m.textInput.Focus()
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m TodoModel) View() string {
	doc := strings.Builder{}
	if m.state == listTodosView {
		doc.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, TodoList(m.cursor, m.todos, m.width, m.height)))
	} else if m.state == newTodoView {
		doc.WriteString(m.textInput.View())
	}
	return doc.String()
}
