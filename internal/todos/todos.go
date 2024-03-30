package todos

import (
	"time"

	"gorm.io/gorm"
)

type Todo struct {
	gorm.Model

	ID          int64
	Name        string
	Completed   bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
	CompletedAt time.Time
	ParentID    int64
}

func NewTodo(name string) Todo {
	return Todo{
		Name:      name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
