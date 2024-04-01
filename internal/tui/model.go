package tui

import (
	"strings"

	"github.com/boejennet/godo/internal/data"
	"github.com/boejennet/godo/internal/todos"
	"github.com/charmbracelet/bubbles/key"
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
	keys      keyMap
	cursor    int
	width     int
	height    int
}

func NewModel(tdb *data.TodoDB, todos []todos.Todo) *TodoModel {
	return &TodoModel{
		tdb:   tdb,
		todos: todos,
		state: listTodosView,
		keys:  keys,
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
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Up):
			if m.cursor > 0 {
				m.cursor--
			}
		case key.Matches(msg, m.keys.Down):
			if m.cursor < len(m.todos)-1 {
				m.cursor++
			}
		case key.Matches(msg, m.keys.Refresh):
			m.todos = m.tdb.GetAllTodos()
		case key.Matches(msg, m.keys.NewTodo):
			// Create a new todo
			if m.state == listTodosView {
				ti := textinput.New()
				ti.Placeholder = "New todo... press enter to save, esc to cancel"
				ti.CharLimit = 256
				ti.Width = m.width
				m.textInput = ti
				m.state = newTodoView
			}
		case key.Matches(msg, m.keys.ToggleCompleted):
			if len(m.todos) > 0 {
				m.todos[m.cursor].Completed = !m.todos[m.cursor].Completed
				m.tdb.Update(m.todos[m.cursor])
			}
		}
		switch msg.String() {
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
