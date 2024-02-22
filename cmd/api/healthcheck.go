package main

import (
	"net/http"
)

// healthcheckHandler provides information about the application's status.
func (app *Application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {
	health := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}
	err := app.writeJSON(w, http.StatusOK, health, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
