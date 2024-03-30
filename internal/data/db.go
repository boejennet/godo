package data

import (
	"fmt"
	"path/filepath"

	"github.com/boejennet/godo/internal/todos"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TodoDB struct {
	db *gorm.DB
}

func NewTDB() *TodoDB {
	return &TodoDB{}
}

func (t *TodoDB) GetAllTodos() []todos.Todo {
	t.getDBConn()
	todos := []todos.Todo{}
	result := t.db.Find(&todos)
	if result.Error != nil {
		fmt.Println("Error getting todos: ", result.Error)
	}
	return todos
}

func (t *TodoDB) Insert(todo todos.Todo) {
	t.getDBConn()
	result := t.db.Create(&todo)
	if result.Error != nil {
		fmt.Println("Error inserting todo: ", result.Error)
	}
}

func (t *TodoDB) Update(todo todos.Todo) {
	t.getDBConn()
	result := t.db.Save(&todo)
	if result.Error != nil {
		fmt.Println("Error updating todo: ", result.Error)
	}
}

func (t *TodoDB) GetCompleted() []todos.Todo {
	t.getDBConn()
	todos := []todos.Todo{}
	result := t.db.Where("completed = ?", true).Find(&todos)
	if result.Error != nil {
		fmt.Println("Error getting completed todos: ", result.Error)
	}
	return todos
}

func (t *TodoDB) GetByID(id int64) todos.Todo {
	t.getDBConn()
	todo := todos.Todo{}
	result := t.db.First(&todo, id)
	if result.Error != nil {
		fmt.Println("Error getting todo by ID: ", result.Error)
	}
	return todo
}

// Set up connection when required
func (t *TodoDB) getDBConn() error {
	if t.db != nil {
		return nil
	}
	dbFullLocation := filepath.Join(SetupPath(), "todos.db")
	db, errorDB := gorm.Open(sqlite.Open(dbFullLocation), &gorm.Config{})
	if errorDB != nil {
		return fmt.Errorf("failed to connect database: %w", errorDB)
	}

	db.AutoMigrate(&todos.Todo{})
	t.db = db

	return nil
}
