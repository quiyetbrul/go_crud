package main

import (
	"fmt"
	"net/http"
)

// logError logs an error message and records additional information
// about the request including the HTTP method and URl.
func (app *Application) logError(r *http.Request, err error) {
	app.logger.Println(err)
}

// errorResponse sends JSON-formatted err msgs to the client with a status code.
// Note: an interface is used to provide flexibility over the values included in a response
func (app *Application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message interface{}) {
	env := envelope{"error": message}
	err := app.writeJSON(w, status, env, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}
}

// serverErrorResponse logs detailed err msgs
// and sends 500 Internal Server Error status code and JSON response to the client.
func (app *Application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)
	message := "server could not process your request"
	app.errorResponse(w, r, http.StatusInternalServerError, message)
}

// notFoundResponse sends 404 Not Found status code and JSON response to the client.
func (app *Application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	message := "request not found"
	app.errorResponse(w, r, http.StatusNotFound, message)
}

// methodNotAllowedResponse sends 405 Method Not Allowed status code and JSON response to the client.
func (app *Application) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {
	message := fmt.Sprintf("the $s method is not supported for this request", r.Method)
	app.errorResponse(w, r, http.StatusMethodNotAllowed, message)
}

func (app *Application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.errorResponse(w, r, http.StatusBadRequest, err.Error())
}

// Note that the errors parameter here has the type map[string]string, which is exactly
// the same as the errors map contained in the Validator type.
func (app *Application) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, errors)
}
