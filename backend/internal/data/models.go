// File: todoApi/backend/internal/data/models.go
package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// A wrapper for out data models
type Models struct {
	Tasks TaskModel
}

// NewModels() allows us to create a new model
func NewModels(db *sql.DB) Models {
	return Models{
		Tasks: TaskModel{DB: db},
	}
}
