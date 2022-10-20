// File: todoApi/backend/cmd/api/handlers.go
package main

import "net/http"

func (app *application) createTaskHandler(w http.ResponseWriter, r *http.Request) {
	//Our target decode destination
	var input struct {
		Title       string `json:"title"`
		Descritpion string `json:"description"`
		Completed   bool   `json:"completed"`
	}

	//Initialize a new json.Decoder instance
	err != app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
}
