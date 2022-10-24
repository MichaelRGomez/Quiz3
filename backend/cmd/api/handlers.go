// File: todoApi/backend/cmd/api/handlers.go
package main

import (
	"errors"
	"fmt"
	"net/http"

	"todo.michaelgomez.net/internal/data"
	"todo.michaelgomez.net/internal/validator"
)

func (app *application) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	//Our target decode destination
	var input struct {
		Title       string `json:"title"`
		Descritpion string `json:"description"`
		Completed   bool   `json:"completed"`
	}

	//Initialize a new json.Decoder instance
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	//coping the valeus from the input struct to the new task struct
	task := &data.Task{
		Title:       input.Title,
		Descritpion: input.Descritpion,
		Completed:   input.Completed,
	}

	//Initialize a new Validator Instance
	v := validator.New()

	//check the map to determine if ther were any validation errors
	if data.ValidateTask(v, task); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//Creating a task
	err = app.models.Tasks.Insert(task)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	//Create a location header for the newly created resource/School
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/toto/%d", task.ID))

	//Writing the JSON response with 201 - created status code with the body
	//being the task data and the header being the headers map
	err = app.writeJSON(w, http.StatusCreated, envelope{"task": task}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// The showentry handler will display an individual task
func (app *application) showTaskHandler(w http.ResponseWriter, r *http.Request) {
	//getting the request data from param function
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundReponse(w, r)
		return
	}

	//Fetching the specific task
	task, err := app.models.Tasks.Get(id)

	//Handling errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundReponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	//Writing the data from the returned get()
	err = app.writeJSON(w, http.StatusOK, envelope{"task": task}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// The updateschool handler will facilitate an update action to the task in the database
func (app *application) updateTaskHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("debug ! 1")

	//This method does a partial replacement
	//Get the id for the task that needs updating
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundReponse(w, r)
		return
	}

	fmt.Println("debug ! 2")

	//Fetch the original record from the database
	task, err := app.models.Tasks.Get(id)

	fmt.Println("debug ! 3")

	//Handling the errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundReponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	fmt.Println("debug ! 4")

	//Creating an input struct to hold data read in from the client
	//Updating the input struct to use pointers because pointers have a default value of nil
	var input struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Completed   *bool   `json:"completed"`
	}

	fmt.Println("debug ! 5")

	//Initilizing a new json.Decoder instance
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	fmt.Println("debug ! 6")

	//checking for any updates
	if input.Title != nil {
		task.Title = *input.Title
	}
	if input.Description != nil {
		task.Descritpion = *input.Description
	}
	if input.Completed != nil {
		task.Completed = *input.Completed
	}

	fmt.Println("debug ! 7")

	//Performing validation on the updated task. If validation fails, then we send a 422 - unprocessable enitiy response to the client
	//Initilize a new Validator Instance
	v := validator.New()

	fmt.Println("debug ! 8")

	//Checking the map to determin if there were any validation errors
	if data.ValidateTask(v, task); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	fmt.Println("debug ! 9")

	//Passing the updated task record to the update() method
	err = app.models.Tasks.Update(task)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	fmt.Println("debug ! 10")

	//Writing the data returned by Get()
	err = app.writeJSON(w, http.StatusOK, envelope{"task": task}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	fmt.Println("debug ! 11")
}

// deletetask handler is to facilitate deletion of a task
func (app *application) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	//geting the id for the task that needs updating
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundReponse(w, r)
		return
	}

	//deleting the school from the database, send a 404 not found status code to the client if there is no matching record
	err = app.models.Tasks.Delete(id)

	//handling errors
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundReponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	//Returning 200 status ok to the client with a success message
	err = app.writeJSON(w, http.StatusOK, envelope{"message": "task sucessfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// The listtask handler allows the client to see a listing of a schools based on a set of criteria
func (app *application) listTasksHandler(w http.ResponseWriter, r *http.Request) {
	//creating an input struct to hold our query parameters
	var input struct {
		Title       string
		Description string
		Completed   bool
		data.Filters
	}

	//Initializing a validator
	v := validator.New()

	//getting the URL values map
	qs := r.URL.Query()

	//Using the helper method to extract the values
	input.Title = app.readString(qs, "title", "")
	input.Description = app.readString(qs, "decription", "")
	input.Completed = app.readBool(qs, "completed", false, v)

	//filering now
	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortList = []string{"id", "title", "completed", "-id", "-description", "-completed"}

	//checking for validation errors
	if data.ValidateFilter(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	//Geting a listing of all tasks
	tasks, metadata, err := app.models.Tasks.GetAll(input.Title, input.Description, input.Completed, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//sending JSON response containing all the schools
	err = app.writeJSON(w, http.StatusOK, envelope{"tasks": tasks, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
