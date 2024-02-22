package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// stable, fast, and well-tested compared to http.ServeMux
// provides validation for request methods as well as HTTP inputs

// The httprouter package provides a few configuration options that you can use to
// customize the behavior of your application further, including enabling trailing slash
// redirects and enabling automatic URL path cleaning.
// https://pkg.go.dev/github.com/julienschmidt/httprouter?tab=doc#Router
func (app *Application) routes() *httprouter.Router {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/api/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/api/v1/items", app.createItemHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/items", app.getAllItemsHandler)
	router.HandlerFunc(http.MethodGet, "/api/v1/item/:id", app.getItemHandler)
	router.HandlerFunc(http.MethodPut, "/api/v1/item/:id", app.updateItemHandler)
	router.HandlerFunc(http.MethodDelete, "/api/v1/item/:id", app.deleteItemHandler)

	return router
}
