// File: todoApi/backend/cmd/api/helpers.go
package main

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// Defining envelope type for json formatting
type envelope map[string]interface {
}

func (app *application) readIDParam(r *http.Request) (int64, error) {
	//getting request from slice
	params := httprouter.ParamsFromContext(r.Context())

	//getting the id
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	if err != nil || id < 1 {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}
