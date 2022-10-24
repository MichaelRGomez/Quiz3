// File: todoApi/backend/cmd/api/routes.go
package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	//security routes
	router.NotFound = http.HandlerFunc(app.notFoundReponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	//actual routes
	router.HandlerFunc(http.MethodGet, "/v1/todo", app.listTasksHandler)
	router.HandlerFunc(http.MethodPost, "/v1/todo", app.createTaskHandler)
	router.HandlerFunc(http.MethodGet, "/v1/todo/:id", app.showTaskHandler)
	router.HandlerFunc(http.MethodPut, "/v1/todo/:id", app.updateTaskHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/todo/:id", app.deleteTaskHandler)

	return router
}
