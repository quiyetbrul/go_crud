package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
)

type Models struct {
	Items interface {
		Insert(item *Item) error
		Get(id int64) (*Item, error)
		GetAll() ([]*Item, error)
		Update(item *Item) error
		Delete(id int64) error
	}
}

// NewModels returns a Models struct containing an initialized ItemModel.
func NewModels(db *sql.DB) Models {
	return Models{
		Items: ItemModel{DB: db},
	}
}
