// File: todoApi/backend/cmd/api/errors.go
package main

import (
	"net/http"
)

// reports a logged error to the terminal
func (app *application) logError(r *http.Request, err error) {
	app.logger.Println(err)
}

// to facilitate a json formatted error repsonse
func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	//creating the json response
	env := envelope{"error": message}
	err := app.writeJson(w, status, env, nil)

	if err != nil {
		app.logError(r, err)
		w.WriteHeader(http.StatusInternalServerError)
	}
}
