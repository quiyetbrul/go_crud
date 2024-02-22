package data

import (
	"database/sql"
	"errors"
)

type Item struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"` // optional
	Completed   bool   `json:"completed"`
}

type ItemModel struct {
	DB *sql.DB
}

// Insert creates a new record in the todolist table.
func (m ItemModel) Insert(i *Item) error {
	query := `INSERT INTO todolist (title, description, completed)
	VALUES ($1, $2, $3)
	RETURNING id`

	args := []interface{}{i.Title, i.Description, i.Completed}

	return m.DB.QueryRow(query, args...).Scan(&i.ID)
}

// Get retrieves a record from the todolist table.
func (m ItemModel) Get(id int64) (*Item, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `SELECT id, title, description, completed
	FROM todolist
	WHERE id = $1`

	var item Item

	err := m.DB.QueryRow(query, id).Scan(&item.ID, &item.Title, &item.Description, &item.Completed)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &item, nil
}

// GetAll retrieves all records from the todolist table.
func (m ItemModel) GetAll() ([]*Item, error) {
	query := `
	SELECT * 
	FROM todolist
	ORDER BY id`
	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	items := []*Item{}
	for rows.Next() {
		var item Item
		err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Description,
			&item.Completed,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, &item)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

// Update updates a record in the todolist table.
func (m ItemModel) Update(i *Item) error {
	query := `UPDATE todolist
	SET title = $1, description = $2, completed = $3
	WHERE id = $4
	RETURNING title`

	args := []interface{}{
		i.Title,
		i.Description,
		i.Completed,
		i.ID,
	}

	return m.DB.QueryRow(query, args...).Scan(&i.Title)
}

// Delete deletes a record in the todolist table.
func (m ItemModel) Delete(id int64) error {
	query := `DELETE FROM todolist
	WHERE id = $1`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
