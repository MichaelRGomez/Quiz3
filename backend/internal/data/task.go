// File: todoApi/backend/internal/data/task.go
package data

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	Version     int32     `json:"version"`
}

func ValidateTask(v *validator.Validator, task *Task) {
	//using check() method to check our validation checks
	v.Check(task.Title != "", "title", "must be provided")
	v.Check(len(task.Title) <= 250, "title", "must not be more than 250 bytes long")

	v.Check(task.Descritpion != "", "description", "must be provided")
	v.Check(len(task.Descritpion) <= 250, "description", "must no be more than 250 bytes long")

	//v.Check(task.Completed, "completed", "new task must be false")
}

type TaskModel struct {
	DB *sql.DB
}

// Insert() allows us to create a new task
func (m TaskModel) Insert(task *Task) error {
	query := `
		INSERT INTO task_list (title, description, completed)
		VALUES ($1, $2, $3)
		RETURNING id, created_at, completed, version
	`

	//creating the context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//Cleaning up to prevent memory leaks
	defer cancel()

	//collect the date field into a slice
	args := []interface{}{task.Title, task.Descritpion, task.Completed}

	return m.DB.QueryRowContext(ctx, query, args...).Scan(&task.ID, &task.CreatedAt, &task.Completed, &task.Version)
}

// Get() allows us to retrieve a specific task
func (m TaskModel) Get(id int64) (*Task, error) {
	//Ensure that there is a valid id
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	//Construct our query with the given id
	query := `
		SELECT id, created_at, title, description, completed, version
		FROM task_list
		WHERE id = $1
	`

	//Declaring the Task varaible to hold the returned data
	var task Task

	//Creating the context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//Cleaning up to prevent memory leaks
	defer cancel()

	err := m.DB.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.CreatedAt,
		&task.Title,
		&task.Descritpion,
		&task.Completed,
		&task.Version,
	)

	if err != nil {
		//Check the type of error
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	//Succes
	return &task, nil
}

// Update() allows us to edit/alter a specific task
// Optimistic locking (version number)
func (m TaskModel) Update(task *Task) error {
	//create a query
	query := `
		UPDATE task_list
		SET title = $1, description = $2, completed = $3, version = version + 1
		WHERE id = $4
		AND version = $5
		RETURNING version
	`
	args := []interface{}{task.Title, task.Descritpion, task.Completed, task.ID, task.Version}

	//Creating the context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//Cleaning up to prevent memory leaks
	defer cancel()

	//Check for edit conflicts
	err := m.DB.QueryRowContext(ctx, query, args...).Scan(&task.Version)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil

}

// Delete() removes a specific task
func (m TaskModel) Delete(id int64) error {
	//Ensure that there is a valid id
	if id < 1 {
		return ErrRecordNotFound
	}
	//creating the delete query
	query := `
		DELETE FROM task_list
		WHERE id = $1
	`

	//creating the context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	//clearing up to prevent memory leaks
	defer cancel()

	//Executing the query
	result, err := m.DB.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	//checking how many rows were affected by the delete operation. we call the RowsAffected() method on the result variable
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	//Check if no rows were affected
	if rowsAffected == 0 {
		return ErrRecordNotFound
	}
	return nil
}

// the GetAll() method returns a list of all tasks sorted by id
func (m TaskModel) GetAll(title string, description string, completed bool, filters Filters) ([]*Task, Metadata, error) {
	//constructing the query
	query := fmt.Sprintf(`
		SELECT COUNT(*) OVER(),
		id, created_at, title, description, completed, version
		FROM task_list
		WHERE (to_tsvector('simple', title) @@ plainto_tsquery('simple', $1) OR $1 = '')
		AND (to_tsvector('simple', description) @@ plainto_tsquery('simple', $2) OR $2 = '')
		AND (completed = $3 OR completed = TRUE)
		ORDER BY %s %s, id ASC
		LIMIT $4 OFFSET $5
	`, filters.sortColumn(), filters.sortOrder())

	//creating the 3 second time out context
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	//fmt.Println("Debug ! 2.5")

	//Execute the query
	args := []interface{}{title, description, completed, filters.limit(), filters.offSet()}
	//fmt.Println("Debug ! 2.42")
	rows, err := m.DB.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, Metadata{}, err
	}

	//fmt.Println("Debug ! 2.41")

	//Closing the result set
	defer rows.Close()
	totalRecords := 0

	//Initialize an empty slice to hold the task data
	tasks := []*Task{}

	//fmt.Println("Debug ! 2.4")

	//Iterate over the rows in the result set
	for rows.Next() {
		var task Task

		//Scanning the valus from the row into the task struct
		err := rows.Scan(
			&totalRecords,
			&task.ID,
			&task.CreatedAt,
			&task.Title,
			&task.Descritpion,
			&task.Completed,
			&task.Version,
		)
		if err != nil {
			return nil, Metadata{}, err
		}
		//Add the task to our slice
		tasks = append(tasks, &task)
	}

	//fmt.Println("Debug ! 2.3")

	//checking for errors after looping through the result set
	if err = rows.Err(); err != nil {
		return nil, Metadata{}, err
	}

	//fmt.Println("Debug ! 2.2")

	metadata := calculateMetaData(totalRecords, filters.Page, filters.PageSize)
	//returning the slice of tasks

	//fmt.Println("Debug ! 2.1")

	return tasks, metadata, nil
}
