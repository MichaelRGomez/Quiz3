// File: todoApi/backend/internal/data/task.go
package data

import (
	"database/sql"
	"time"

	"todo.michaelgomez.net/internal/validator"
)

// task struct supports the infromation for the todo task
type Task struct {
	ID          int64     `json:"id"`
	CreatedAt   time.Time `json:"-"`
	Title       string    `json:"title"`
	Descritpion string    `json:"description"`
	Completed   bool      `json:"completed"`
}

func ValidateTask(v *validator.Validator, task *Task) {
	//using check() method to check our validation checks
	v.Check(task.Title != "", "title", "must be provided")
	v.Check(len(task.Title) <= 250, "title", "must not be more than 250 bytes long")

	v.Check(task.Descritpion != "", "description", "must be provided")
	v.Check(len(task.Descritpion) <= 250, "description", "must no be more than 250 bytes long")

	v.Check(task.Completed != false, "completed", "new task must be false")
}

type TaskModel struct {
	DB *sql.DB
}
