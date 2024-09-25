package main

import (
	"fmt"
	"net/http"

	"github.com/quiyetbrul/go_crud/internal/data"
	"github.com/quiyetbrul/go_crud/internal/validator"
)

// createItemHandler handles "POST /api/v1/items" endpoint.
func (app *Application) createItemHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Completed   bool   `json:"completed"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	item := &data.Item{
		Title:       input.Title,
		Description: input.Description,
		Completed:   input.Completed,
	}

	v := validator.New()
	if !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Items.Insert(item)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/api/v1/items/%d", item.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"item": item}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// getAllItemsHandler handles "GET /api/v1/items" endpoint and retrieves all items from the todolist.
func (app *Application) getAllItemsHandler(w http.ResponseWriter, r *http.Request) {
	items, err := app.models.Items.GetAll()
	err = app.writeJSON(w, http.StatusOK, envelope{"todolist": items}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// getItemHandler handles "GET /api/v1/item/{:id}" endpoint and retrieves an item from the todolist.
func (app *Application) getItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	item, err := app.models.Items.Get(id)
	if err != nil {
		app.switchHelper(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"item": item}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// updateItemHandler handles "PUT /api/v1/item/{:id}" endpoint and updates an item from the todolist.
func (app *Application) updateItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	item, err := app.models.Items.Get(id)
	if err != nil {
		app.switchHelper(w, r, err)
		return
	}

	var input struct {
		Title       *string `json:"title"`
		Description *string `json:"description"`
		Completed   *bool   `json:"completed"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Title != nil {
		item.Title = *input.Title
	}
	if input.Description != nil {
		item.Description = *input.Description
	}
	if input.Completed != nil {
		item.Completed = *input.Completed
	}

	v := validator.New()
	if ValidateItem(v, item); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Items.Update(item)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"updated item": item}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

// deleteItemHandler handles "DELETE /api/v1/item/{:id}" endpoint and deletes an item from the todolist.
func (app *Application) deleteItemHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	err = app.models.Items.Delete(id)
	if err != nil {
		app.switchHelper(w, r, err)
		return
	}
	fmt.Fprintf(w, "Delete item: %v\n", id)
}

func ValidateItem(v *validator.Validator, item *data.Item) {
	v.Check(item.Title != "", "title", "must be provided")
	v.Check(item.Description != "", "description", "must be provided")
	v.Check(item.Completed || !item.Completed, "completed", "must be provided")
}
